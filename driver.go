package main

import (
	"encoding/json"
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
}

////////////////
// Lectures
////////////////

// Receive : Lit le message en attente.
func (d *Driver) Receive() error {
	req, err := d.readResponse()

	if err != nil {
		clog.File("Driver", "Receive", "Driver: %s - ERROR: %v", d.Name, err)
		return err
	}
	if req == nil {
		// Message vide de contôle.
		clog.File("Driver", "Receive", "Driver: %s - Empty request", d.Name)
	} else {
		switch req.Method {
		case "LoginResponse":
			d.computeLoginResponse(req.Status.ID, req.Params)
		case "NewRide":
			d.requestRide(req.Params)
		case "AcceptRideResponse":
			d.computeRideResponse(req.Status.ID, req.Params)
		case "ChangeTaximeterStateReponse":
			d.computeChangeState(req.Status.ID, req.Params)
		default:
			clog.File("Driver", "Receive", "Driver: %s - RECV: Method: %s [%d] %s", d.Name, req.Method, req.Status.ID, req.Status.Message)
		}
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
	}
	if header.OpCode.IsControl() {
		if header.OpCode == ws.OpClose {
			return &datamodels.Response{ID: d.ID, Method: "close"}, nil
		}
		clog.Error("Driver", "readHeader", "OpCode : %v", header.OpCode)
		// return nil, wsutil.ControlFrameHandler(d.conn, ws.StateServerSide)(h, r)
	}

	payload := make([]byte, header.Length)
	_, err = io.ReadFull(d.conn, payload)
	if err != nil {
		// handle error
		clog.Error("Driver", "readRequest | ReadFull", "%s", err)
	}
	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	req := &datamodels.Response{}
	if err := json.Unmarshal(payload, &req); err != nil {
		clog.Error("Driver", "readHeader", "%s", err)
		clog.File("Driver", "readHeader", "Decode - Driver: %s -> %s", d.Name, payload)
		return nil, err
	}

	return req, nil
}

////////////////
// Ecritures
////////////////

// writeResultTo : Retourne le resultat de la méthode à l'appelant
func (d *Driver) writeRequest(reqID int, method string, req datamodels.DataParams) error {
	request := datamodels.Request{
		ID:     reqID,
		Method: method,
		Params: req,
	}

	return d.write(request)
}

func (d *Driver) writeRaw(p []byte) error {
	d.io.Lock()
	defer d.io.Unlock()

	_, err := d.conn.Write(p)

	return err
}

func (d *Driver) write(x interface{}) error {
	w := wsutil.NewWriter(d.conn, ws.StateClientSide, ws.OpText)
	encoder := json.NewEncoder(w)

	d.io.Lock()
	defer d.io.Unlock()

	if err := encoder.Encode(x); err != nil {
		return err
	}

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

	d.writeRequest(d.ID, "UpdateDriverLocation", params)
}

func dice(nb int) int {
	return rand.Intn(nb) + 1
}

func (d *Driver) updateRide(state datamodels.RideState) {
	params := datamodels.RideUpdate{
		ID:    d.Ride.ID,
		State: state,
	}

	clog.File("Driver", "updateRide", "Driver: %s ask Ride: %s -> %d", d.Name, d.Ride.ID, state)
	d.writeRequest(d.ID, "ChangeRideState", params)
}

func (d *Driver) requestRide(params datamodels.DataParams) {
	var ride datamodels.Ride
	mapstructure.Decode(params, &ride)

	d.mu.Lock()
	if d.DriverState == datamodels.Free {
		clog.File("Driver", "AcceptRide", "Driver: %s -> Ride: %s", d.Name, ride.ID)
		d.writeRequest(d.ID, "AcceptRide", datamodels.AcceptRide{ID: ride.ID})
		d.DriverState = datamodels.WaitOK
	}
	d.mu.Unlock()
}

func (d *Driver) computeRideResponse(responseCode int, params datamodels.DataParams) {
	var ride datamodels.Ride
	mapstructure.Decode(params, &ride)

	defer d.mu.Unlock()
	d.mu.Lock()

	if responseCode != 0 {
		clog.File("Driver", "RideRejected", "Driver: %s -> Ride: %s", d.Name, ride.ID)
		d.DriverState = datamodels.Free
		return
	}

	if d.DriverState == datamodels.WaitOK {
		clog.File("Driver", "Ride OK", "%s -> %s", d.Name, ride.ID)
		d.Ride = ride
		d.updateRide(datamodels.Approach)
		d.ToDest = geoloc.DistanceAccurate(d.Coord.Latitude, d.Coord.Longitude, ride.FromAddress.Coord.Latitude, ride.FromAddress.Coord.Longitude) / 1000
		d.DriverState = datamodels.Moving
		return
	}

	d.DriverState = datamodels.Free
	clog.File("Driver", "Ride OK ERROR", "%s -> %s", d.Name, ride.ID)
}

func (d *Driver) requestChangeState(newState datamodels.DriverState) {
	state := datamodels.DriverStateChange{
		State: newState,
	}
	d.writeRequest(d.ID, "ChangeTaximeterState", state)
}

func (d *Driver) computeChangeState(responseCode int, params datamodels.DataParams) {
	var newState datamodels.DriverStateChange
	mapstructure.Decode(params, &newState)

	d.mu.Lock()
	d.DriverState = newState.State
	d.mu.Unlock()
}

func (d *Driver) createCourse() {
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
		Method: "NewRide",
		Params: ride,
	}

	d.write(req)
}

func (d *Driver) login() {
	login := datamodels.Login{
		ID:    d.ID,
		Name:  d.Name,
		State: d.DriverState,
	}
	d.writeRequest(d.ID, "Login", login)
}

func (d *Driver) computeLoginResponse(responseCode int, params datamodels.DataParams) {
	d.requestChangeState(datamodels.Free)
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
			d.mu.Lock()
			switch d.DriverState {
			case datamodels.Offline:
				if idleCount == 0 {
					d.DriverState = datamodels.WaitACK
					d.requestChangeState(datamodels.Free)
				} else {
					idleCount--
				}
			case datamodels.Free:
				if dice(100) < conf.Bench.PercentForIdle {
					d.DriverState = datamodels.WaitACK
					d.requestChangeState(datamodels.Offline)
					idleCount = conf.Bench.IdleDuration
					sendPosCount = 0
					if conf.Bench.IdleCreateRide {
						d.createCourse()
					}
				}
			case datamodels.Moving:
				d.ToDest -= float64(conf.Bench.KmByBT)
				if d.ToDest <= 0 {
					d.DriverState = datamodels.WaitACK
					d.requestChangeState(datamodels.Occupied)
					d.updateRide(datamodels.PickUpPassenger)
					d.Coord = d.Ride.FromAddress.Coord
					d.ToDest = geoloc.DistanceAccurate(d.Coord.Latitude, d.Coord.Longitude, d.Ride.ToAddress.Coord.Latitude, d.Ride.ToAddress.Coord.Longitude) / 1000
				}
			case datamodels.Occupied:
				d.ToDest -= float64(conf.Bench.KmByBT)
				if d.ToDest <= 0 {
					d.DriverState = datamodels.WaitACK
					d.requestChangeState(datamodels.Free)
					d.updateRide(datamodels.Ended)
					d.ToDest = 0
					d.Coord = d.Ride.ToAddress.Coord
					d.Ride = datamodels.Ride{}
				}
			case datamodels.WaitACK:
			}
			d.mu.Unlock()

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
