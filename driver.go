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
	DriverState datamodels.TaximeterState
	Coord       datamodels.Coordinates
	Ride        datamodels.RideData
	ToDest      float64

	reqID int
}

////////////////
// Lectures
////////////////

// Receive : Lit le message en attente.
func (d *Driver) Receive() error {
	d.io.Lock()
	defer d.io.Unlock()

	header, err := ws.ReadHeader(d.conn)
	if err != nil {
		// handle error
		clog.Error("Driver", "readRequest | ReadHeader", "%s", err)
		clog.File("R-ERR", d.Name, "Error  -> %s", err)
	}

	payload := make([]byte, header.Length)
	_, err = io.ReadFull(d.conn, payload)

	if err != nil {
		// handle error
		clog.Error("Driver", "readRequest | ReadFull", "%s", err)
		clog.File("R-ERR", d.Name, "ReadFull  -> %s", err)
		return err
	}

	pool.Schedule(func() {
		d.HandleProtocol(header, payload)
	})
	return nil
}

func (d *Driver) HandleProtocol(header ws.Header, payload []byte) error {
	var req *datamodels.Response

	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	if header.OpCode.IsControl() {
		if header.OpCode == ws.OpClose {
			req = &datamodels.Response{ID: d.ID, Method: "close"}
		}
		clog.Error("Driver", "readHeader", "OpCode : %v", header.OpCode)
		clog.File("R-ERR", d.Name, "OpCode  -> %v", header.OpCode)
		// return nil, wsutil.ControlFrameHandler(d.conn, ws.StateServerSide)(h, r)
	} else {
		req = &datamodels.Response{}
		if err := json.Unmarshal(payload, &req); err != nil {
			clog.Error("Driver", "readHeader", "%s", err)
			clog.File("R-ERR", d.Name, "Erreur Decode  -> %s", payload)
			return err
		}
	}

	if req == nil {
		// Message vide de contôle.
		clog.File("R-ERR", d.Name, "Empty request")
		return errors.New("empty request")
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

////////////////
// Ecritures
////////////////

// writeResultTo : Retourne le resultat de la méthode à l'appelant
func (d *Driver) writeRequest(method string, req datamodels.DataParams) {
	d.reqID++
	request := datamodels.Request{
		ID:     d.reqID,
		Method: method,
		Params: req,
	}

	go d.write(request, d.reqID, method)
}

func (d *Driver) write(x interface{}, id int, met string) error {
	w := wsutil.NewWriter(d.conn, ws.StateClientSide, ws.OpText)
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(x); err != nil {
		return err
	}

	d.io.Lock()
	time.Sleep(time.Millisecond * 1000)
	err := w.Flush()
	d.io.Unlock()

	if err != nil {
		clog.File("S-ERR", d.Name, "%d | %s", id, met)
		return err
	}
	clog.File("SEND", d.Name, "%d | %s", id, met)
	return nil
}

/////////////////////////////////
// Protocole
/////////////////////////////////

func (d *Driver) sendCoord() {
	updateDriverLocation := datamodels.UpdateDriverLocation{
		Coord: datamodels.Coordinates{
			Latitude:  d.Coord.Latitude,
			Longitude: d.Coord.Longitude,
		},
		VehicleType: datamodels.Berline,
	}
	updateDriverLocation.VehicleOptions = append(updateDriverLocation.VehicleOptions, datamodels.CovidShield)
	d.writeRequest("UpdateDriverLocation", updateDriverLocation)
}

func dice(nb int) int {
	return rand.Intn(nb) + 1
}

/////////////////////////////////
// ChangeRideState
/////////////////////////////////

func (d *Driver) updateRide(state datamodels.RideState) {
	params := datamodels.ChangeRideState{
		ID:    d.Ride.ID,
		State: state,
	}

	d.writeRequest("ChangeRideState", params)
}

func (d *Driver) computeChangeRideStateResponse(responseCode int, params datamodels.DataParams) {
	var rideState datamodels.ChangeRideState
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
	var newRide datamodels.CreateRide
	mapstructure.Decode(params, &newRide)

	d.mu.Lock()
	if d.DriverState == datamodels.Free {
		d.DriverState = datamodels.WaitOK
		d.writeRequest("AcceptRide", datamodels.AcceptRide{ID: newRide.Ride.ID})
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

func (d *Driver) requestChangeTaximeterStateReponse(newState datamodels.TaximeterState) {
	state := datamodels.ChangeTaximeterState{
		State: newState,
	}
	d.mu.Lock()
	d.DriverState = datamodels.WaitACK
	d.mu.Unlock()
	d.writeRequest("ChangeTaximeterState", state)
}

func (d *Driver) computeChangeTaximeterStateReponse(responseCode int, params datamodels.DataParams) {
	var newState datamodels.ChangeTaximeterState
	mapstructure.Decode(params, &newState)

	if responseCode != 0 {
		d.DriverState = datamodels.WaitACK
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
	var rideState datamodels.PendingPaymentResponse
	mapstructure.Decode(params, &rideState)

	if responseCode != 0 {
		return
	}
	d.mu.Lock()
	d.Ride.State = rideState.Ride.State
	d.DriverState = datamodels.Billing
	d.mu.Unlock()
}

/////////////////////////////////
// NewCourse
/////////////////////////////////

func (d *Driver) createRide() {
	var options []datamodels.VehicleOption
	optionsList := []datamodels.VehicleOption{
		datamodels.CPAM,
		datamodels.CovidShield,
		datamodels.EnglishSpoken,
		datamodels.Mkids1,
		datamodels.Mkids2,
		datamodels.Mkids3,
		datamodels.Mkids4,
		datamodels.Pets,
		datamodels.Access,
	}
	nbOptions := rand.Intn(3) + 1

	ride := datamodels.RideData{
		ExternalID:  xid.New().String(),
		Origin:      datamodels.Defaut,
		StartDate:   datamodels.FormatDateForIOS(time.Now()),
		State:       datamodels.Pending,
		IsImmediate: true,
		FromAddress: getNewAdress(),
		ToAddress:   getNewAdress(),
		// Luggages:    0,
		// Passengers:  1,
		// Vehicle:     datamodels.Other,
	}

	for i := 0; i < nbOptions; i++ {
		choice := rand.Intn(len(optionsList))
		options = append(options, optionsList[choice])
	}

	createRide := datamodels.CreateRide{
		Ride:      ride,
		Passenger: datamodels.Passenger{},
		SearchOptions: datamodels.SearchOptions{
			Memo:           "Pas de retard PLEASE!",
			Reference:      "Pote du Maire",
			VehicleOptions: options,
			VehicleType:    datamodels.Berline,
		},
		Proposal: datamodels.Proposal{},
	}

	req := datamodels.Request{
		ID:     d.ID,
		Method: "CreateRide",
		Params: createRide,
	}

	d.write(req, d.ID, "CreateRide")
}

func (d *Driver) login() {
	login := datamodels.Login{
		ID:    d.ID,
		Name:  d.Name,
		State: d.DriverState,
		Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjQwMTU2MjMsImlhdCI6MTYyNDAwODQyMywidXNlclV1aWQiOiIxNWIzZDZlNC02NDYzLTEwM2ItOTEzMy01OTA4N2Q3ZjkwZmMiLCJ1c2VySWQiOjMxN30.KC3Pd2zpK-Y64HUOwI0JOx48trxRLBZdy5vuuo2Z7RtVzlRAL2zT8gdlhZpTgfDJ9PUwMz4RG3r9QcL_dmFQRRseWVbWX2_KYdxktYSvtOb-Zr9MsEdJrxZ4_t4p5brldEkRRwGr_wKK8EiqozVPgRo5uBSb4bVf7zrLfgb-anBR7GNcicu1e455f_UBvMbef0Oa5rxR144fPmEuVLzWgTjVsEBI1jUBaHrhjOlOe-VPQ64Gt74Rm5Q0mtxtTqTx2rLEPu_v6DmrEKKkF9i8dUA6gJqmOcLISI6jree7eBmOMRPIQo0soVewCxdW5NADBg-KJwecsHpivzi49FY0Hg",
	}
	d.writeRequest("Login", login)
}

func (d *Driver) closeConnection() {
	ws.WriteFrame(d.conn, ws.NewCloseFrame(ws.NewCloseFrameBody(ws.StatusNormalClosure, "Client close connection")))
	clog.Trace("Driver", "Close", "%s (%d) is diconnected", d.Name, d.ID)
	d.conn.Close()
}

func (d *Driver) computeLoginResponse(responseCode int, params datamodels.DataParams) {
	d.requestChangeTaximeterStateReponse(datamodels.Free)
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
		<-ticker.C
		switch d.DriverState {
		case datamodels.WaitACK:
		case datamodels.WaitOK:
		case datamodels.Offline:
			if idleCount == 0 {
				d.requestChangeTaximeterStateReponse(datamodels.Free)
			} else {
				idleCount--
			}
		case datamodels.Free:
			if dice(100) < conf.Bench.PercentForIdle {
				if conf.Bench.IdleCreateRide {
					d.createRide()
				}
				d.requestChangeTaximeterStateReponse(datamodels.Offline)
				idleCount = conf.Bench.IdleDuration
				// sendPosCount = 0
			}
		case datamodels.Moving:
			d.ToDest -= float64(conf.Bench.KmByBT)
			if d.ToDest <= 0 {
				d.updateRide(datamodels.PickUpPassenger)
				d.requestChangeTaximeterStateReponse(datamodels.Occupied)

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
		case datamodels.Billing:
			d.updateRide(datamodels.Ended)
			d.requestChangeTaximeterStateReponse(datamodels.Free)

			d.mu.Lock()
			d.Coord = d.Ride.ToAddress.Coord
			d.Ride = datamodels.RideData{}
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
