package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"bench_dispatch/clog"
	"bench_dispatch/datamodels"
	"bench_dispatch/geoloc"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
)

// Driver : Représente une connexion avec une voiture / taxi
// Cett structure contient tous les infos de communication
type Driver struct {
	io          sync.Mutex
	mu          sync.RWMutex
	conn        io.ReadWriteCloser
	hub         *Hub
	ID          int
	Name        string
	DriverState datamodels.DriverState
	Coord       datamodels.Coordinates
	Ride        datamodels.Ride
	ToDest      float64

	reqID int
}

////////////////
// Lectures
////////////////

// Receive : Lit le message en attente.
func (d *Driver) Receive() error {
	req, err := d.readResponse()

	if err != nil {
		clog.File("R-ERR", d.Name, "ERROR: %v", err)
		return err
	}
	if req == nil {
		// Message vide de contôle.
		clog.File("R-ERR", d.Name, "Empty request")
		return errors.New("Empty request")
	}

	clog.File("RECV", d.Name, "%d | %s | %s", req.ID, req.Method, req.Status.Message)
	switch req.Method {
	case "LoginResponse":
		d.computeLoginResponse(req.Status.ID, req.Params)
	case "NewRide":
		d.requestRide(req.Params)
	case "AcceptRideResponse":
		d.computeAcceptRideResponse(req.Status.ID, req.Params)
	case "ChangeRideStateResponse":
		d.computeChangeRideStateResponse(req.Status.ID, req.Params)
	case "ChangeTaximeterStateReponse":
		d.computeChangeTaximeterStateReponse(req.Status.ID, req.Params)
	case "PendingPaymentResponse":
		d.computePaymentResponse(req.Status.ID, req.Params)
	default:
		clog.File("R-ERR", d.Name, "Erreur Method: %s [code: %d] %s", req.Method, req.Status.ID, req.Status.Message)
	}

	return nil
}

func (d *Driver) readResponse() (*datamodels.Response, error) {
	d.io.Lock()
	defer d.io.Unlock()

	header, err := ws.ReadHeader(d.conn)
	if err != nil {
		// handle error
		clog.Error("Driver", "readRequest | ReadHeader", "%s", err)
		clog.File("R-ERR", d.Name, "Error  -> %s", err)
	}
	if header.OpCode.IsControl() {
		if header.OpCode == ws.OpClose {
			return &datamodels.Response{ID: d.ID, Method: "close"}, nil
		}
		clog.Error("Driver", "readHeader", "OpCode : %v", header.OpCode)
		clog.File("R-ERR", d.Name, "OpCode  -> %v", header.OpCode)
		// return nil, wsutil.ControlFrameHandler(d.conn, ws.StateServerSide)(h, r)
	}

	payload := make([]byte, header.Length)
	_, err = io.ReadFull(d.conn, payload)
	if err != nil {
		// handle error
		clog.Error("Driver", "readRequest | ReadFull", "%s", err)
		clog.File("R-ERR", d.Name, "ReadFull  -> %s", err)
	}
	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	req := &datamodels.Response{}
	if err := json.Unmarshal(payload, &req); err != nil {
		clog.Error("Driver", "readHeader", "%s", err)
		clog.File("R-ERR", d.Name, "Erreur Decode  -> %s", payload)
		return nil, err
	}

	return req, nil
}

////////////////
// Ecritures
////////////////

// writeResultTo : Retourne le resultat de la méthode à l'appelant
func (d *Driver) writeRequest(method string, req datamodels.DataParams) error {
	d.reqID++
	request := datamodels.Request{
		ID:     d.reqID,
		Method: method,
		Params: req,
	}

	return d.write(request, d.reqID, method)
}

func (d *Driver) writeRaw(p []byte) error {
	d.io.Lock()
	defer d.io.Unlock()

	_, err := d.conn.Write(p)

	return err
}

func (d *Driver) write(x interface{}, id int, met string) error {
	w := wsutil.NewWriter(d.conn, ws.StateClientSide, ws.OpText)
	encoder := json.NewEncoder(w)

	d.io.Lock()
	defer d.io.Unlock()

	if err := encoder.Encode(x); err != nil {
		return err
	}
	clog.File("SEND", d.Name, "%d | %s", id, met)
	return w.Flush()
}

/////////////////////////////////
// Protocole
/////////////////////////////////

func (d *Driver) sendCoord() {
	params := datamodels.Coordinates{
		Latitude:  d.Coord.Latitude,
		Longitude: d.Coord.Longitude,
	}

	d.writeRequest("UpdateDriverLocation", params)
}

func dice(nb int) int {
	return rand.Intn(nb) + 1
}

/////////////////////////////////
// ChangeRideState
/////////////////////////////////

func (d *Driver) updateRide(state datamodels.RideState) {
	params := datamodels.RideUpdate{
		ID:    d.Ride.ID,
		State: state,
	}

	d.writeRequest("ChangeRideState", params)
}

func (d *Driver) computeChangeRideStateResponse(responseCode int, params datamodels.DataParams) {
	var rideState datamodels.RideUpdate
	mapstructure.Decode(params, &rideState)

	if responseCode != 0 {
		return
	}
	d.mu.Lock()
	d.Ride.State = rideState.State
	d.mu.Unlock()
}

/////////////////////////////////
// AcceptRide
/////////////////////////////////
func (d *Driver) requestRide(params datamodels.DataParams) {
	var ride datamodels.Ride
	mapstructure.Decode(params, &ride)

	d.mu.Lock()
	if d.DriverState == datamodels.Free {
		d.DriverState = datamodels.WaitOK
		d.writeRequest("AcceptRide", datamodels.AcceptRide{ID: ride.ID})
	}
	d.mu.Unlock()
}

func (d *Driver) computeAcceptRideResponse(responseCode int, params datamodels.DataParams) {
	var rideResp datamodels.AcceptRideResponse
	mapstructure.Decode(params, &rideResp)

	defer d.mu.Unlock()
	d.mu.Lock()

	if responseCode != 0 {
		d.DriverState = datamodels.Free
		return
	}

	if d.DriverState == datamodels.WaitOK {
		d.Ride = rideResp.Ride
		d.updateRide(datamodels.Approach)
		d.ToDest = geoloc.DistanceAccurate(d.Coord.Latitude, d.Coord.Longitude, rideResp.Ride.FromAddress.Coord.Latitude, rideResp.Ride.FromAddress.Coord.Longitude) / 1000
		d.DriverState = datamodels.Moving
		return
	}

	d.DriverState = datamodels.Free
}

/////////////////////////////////
// ChangeTaximeterState
/////////////////////////////////

func (d *Driver) requestDriverChangeState(newState datamodels.DriverState) {
	state := datamodels.DriverStateChange{
		State: newState,
	}
	d.mu.Lock()
	d.DriverState = datamodels.WaitACK
	d.mu.Unlock()
	d.writeRequest("ChangeTaximeterState", state)
}

func (d *Driver) computeChangeTaximeterStateReponse(responseCode int, params datamodels.DataParams) {
	var newState datamodels.DriverStateChange
	mapstructure.Decode(params, &newState)

	if responseCode != 0 {
		d.DriverState = datamodels.Free
		return
	}

	d.mu.Lock()
	d.DriverState = newState.State
	d.mu.Unlock()
}

/////////////////////////////////
// PendingPayment
/////////////////////////////////

func (d *Driver) computePaymentResponse(responseCode int, params datamodels.DataParams) {
	var rideState datamodels.PaymentResponse
	mapstructure.Decode(params, &rideState)

	if responseCode != 0 {
		return
	}
	d.mu.Lock()
	d.Ride.State = rideState.Ride.State
	d.DriverState = datamodels.Payment
	d.mu.Unlock()
}

/////////////////////////////////
// NewCourse
/////////////////////////////////

func (d *Driver) createRide() {
	ride := datamodels.Ride{
		ExternalID:  xid.New().String(),
		Origin:      datamodels.Defaut,
		Date:        time.Now().Format(time.RFC3339),
		ValidUntil:  time.Now().Format(time.RFC3339),
		State:       datamodels.Pending,
		IsImmediate: true,
		FromAddress: getNewAdress(),
		ToAddress:   getNewAdress(),
		Options: datamodels.OptionsRide{
			Luggages:   0,
			Passengers: 1,
			Vehicle:    datamodels.Other,
		},
	}

	req := datamodels.Request{
		ID:     d.ID,
		Method: "CreateRide",
		Params: ride,
	}

	d.write(req, d.ID, "CreateRide")
}

func (d *Driver) login() {
	login := datamodels.Login{
		ID:    d.ID,
		Name:  d.Name,
		State: d.DriverState,
	}
	d.writeRequest("Login", login)
}

func (d *Driver) computeLoginResponse(responseCode int, params datamodels.DataParams) {
	d.requestDriverChangeState(datamodels.Free)
}

// Life : Simulation des actions d'un Driver
func (d *Driver) Life() {
	baseTimer, _ := time.ParseDuration(fmt.Sprintf("%ds", conf.Bench.BaseTimer))
	ticker := time.NewTicker(baseTimer)
	defer func() {
		ticker.Stop()
	}()

	d.login()

	idleCount := 0
	sendPosCount := 0

	for {
		select {
		case <-ticker.C:
			switch d.DriverState {
			case datamodels.WaitACK:
			case datamodels.WaitOK:
			case datamodels.Offline:
				if idleCount == 0 {
					d.requestDriverChangeState(datamodels.Free)
				} else {
					idleCount--
				}
			case datamodels.Free:
				if dice(100) < conf.Bench.PercentForIdle {
					if conf.Bench.IdleCreateRide {
						d.createRide()
					}
					d.requestDriverChangeState(datamodels.Offline)
					idleCount = conf.Bench.IdleDuration
					sendPosCount = 0
				}
			case datamodels.Moving:
				d.ToDest -= float64(conf.Bench.KmByBT)
				if d.ToDest <= 0 {
					d.requestDriverChangeState(datamodels.Occupied)
					d.updateRide(datamodels.PickUpPassenger)
					d.mu.Lock()
					d.Coord = d.Ride.FromAddress.Coord
					d.ToDest = geoloc.DistanceAccurate(d.Coord.Latitude, d.Coord.Longitude, d.Ride.ToAddress.Coord.Latitude, d.Ride.ToAddress.Coord.Longitude) / 1000
					d.mu.Unlock()
				}
			case datamodels.Occupied:
				d.ToDest -= float64(conf.Bench.KmByBT)
				if d.ToDest <= 0 {
					d.mu.Lock()
					d.DriverState = datamodels.WaitACK
					d.mu.Unlock()
					d.updateRide(datamodels.PendingPayment)
					d.ToDest = 0
				}
			case datamodels.Payment:
				d.requestDriverChangeState(datamodels.Free)
				d.updateRide(datamodels.Ended)
				d.mu.Lock()
				d.Coord = d.Ride.ToAddress.Coord
				d.Ride = datamodels.Ride{}
				d.mu.Unlock()
			}

			if sendPosCount == 0 {
				d.sendCoord()
				sendPosCount = conf.Bench.SendPos
			}
			sendPosCount--

			// case s := <-d.in:
			// 	d.mu.Lock()
			// 	d.driverState = s
			// 	clog.File("Driver", "updateDriver", "%d -> %d", d.id, s)
			// 	d.mu.Unlock()
		}
	}
}
