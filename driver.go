package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"bench_dispatch/datamodels"
	"bench_dispatch/geoloc"
	"bench_dispatch/tools/clog"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/xid"
)

// UserState : Etat d'un Driver
type UserState int

const (
	idle UserState = iota
	ready
	waitOK
	moving
	onRide
	err
)

// Driver : Représente une connexion avec une voiture / taxi
// Cett structure contient tous les infos de communication
type Driver struct {
	io          sync.Mutex
	mu          sync.RWMutex
	conn        io.ReadWriteCloser
	hub         *Hub
	id          int
	name        string
	driverState UserState
	coord       datamodels.Coordinates
	dice        int
	ride        datamodels.Ride
	toDest      float64
}

////////////////
// Lectures
////////////////

// Receive : Lit le message en attente.
func (d *Driver) Receive() error {
	req, _ := d.readResponse()

	if req == nil {
		// Message vide de contôle.
		// clog.Trace("main", "Driver", "Empty request")
		clog.File("main", "Driver", "Empty request")
		req = &datamodels.Response{}
	} else {
		switch req.Method {
		case "NewRide":
			d.acceptRide(req.Params)
		case "AcceptRideResponse":
			d.computeRideResponse(req.Status.ID, req.Params)
		default:
			clog.File("main", "Driver", "RECV: [%d] %s", req.Status.ID, req.Status.Message)
		}
	}
	return nil
}

func (d *Driver) readResponse() (*datamodels.Response, error) {
	d.io.Lock()
	defer d.io.Unlock()

	h, r, err := wsutil.NextReader(d.conn, ws.StateClientSide)
	if err != nil {
		// clog.Error("Driver", "readRequest::NextReader", "%s", err)
		return nil, err
	}
	if h.OpCode.IsControl() {
		// clog.Info("Driver", "readRequest::IsControl", "Opcode Control : %v", h)
		return nil, wsutil.ControlFrameHandler(d.conn, ws.StateClientSide)(h, r)
	}

	req := &datamodels.Response{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(req); err != nil {
		// clog.Error("Driver", "readRequest::Decode", "%s", err)
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
		Latitude:  d.coord.Latitude,
		Longitude: d.coord.Longitude,
	}

	d.writeRequest(d.id, "UpdateDriverLocation", params)
}

func dice(nb int) int {
	return rand.Intn(nb) + 1
}

func (d *Driver) updateRide(state datamodels.RideState) {
	params := datamodels.RideUpdate{
		ID:    d.ride.ID,
		State: state,
	}

	clog.File("Driver", "updateRide", "Driver: %d ask Ride: %s -> %d", d.id, d.ride.ID, state)
	d.writeRequest(d.id, "ChangeRideState", params)
}

func (d *Driver) acceptRide(params datamodels.DataParams) {
	var ride datamodels.Ride
	mapstructure.Decode(params, &ride)

	d.mu.Lock()
	if d.driverState == ready {
		clog.File("Driver", "AcceptRide", "Driver: %d -> Ride: %s", d.id, ride.ID)
		d.writeRequest(d.id, "AcceptRide", datamodels.AcceptRide{RideID: ride.ID})
		d.driverState = waitOK
	}
	d.mu.Unlock()
}

func (d *Driver) computeRideResponse(responseCode int, params datamodels.DataParams) {
	var ride datamodels.Ride
	mapstructure.Decode(params, &ride)

	defer d.mu.Unlock()
	d.mu.Lock()

	if responseCode != 0 {
		clog.File("Driver", "RideRejected", "Driver: %d -> Ride: %s", d.id, ride.ID)
		d.driverState = ready
		return
	}

	if d.driverState == waitOK {
		clog.File("Driver", "Ride OK", "%d -> %s", d.id, ride.ID)
		d.ride = ride
		d.updateRide(datamodels.Approach)
		d.toDest = geoloc.DistanceAccurate(d.coord.Latitude, d.coord.Longitude, ride.FromAddress.Coord.Latitude, ride.FromAddress.Coord.Longitude) / 1000
		d.driverState = moving
		return
	}

	d.driverState = ready
	clog.File("Driver", "Ride OK ERROR", "%d -> %s", d.id, ride.ID)
}

func (d *Driver) createCourse() {
	now := time.Now()

	ride := datamodels.Ride{
		Origin:      datamodels.Defaut,
		ID:          xid.New().String(),
		Date:        now.Format(time.RFC3339),
		ValidUntil:  now.Format(time.RFC3339),
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
		ID:     d.id,
		Method: "NewRide",
		Params: ride,
	}

	d.write(req)
}

func (d *Driver) login() {
	login := datamodels.Login{
		ID:   d.id,
		Name: d.name,
	}
	d.writeRequest(d.id, "Login", login)
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
			switch d.driverState {
			case idle:
				if idleCount == 0 {
					d.driverState = ready
				} else {
					idleCount--
				}
			case ready:
				if dice(100) < conf.Bench.PercentForIdle {
					d.driverState = idle
					idleCount = conf.Bench.IdleDuration
					sendPosCount = 0
					if conf.Bench.IdleCreateRide {
						d.createCourse()
					}
				}
			case moving:
				d.toDest -= float64(conf.Bench.KmByBT)
				if d.toDest <= 0 {
					d.updateRide(datamodels.PickUpPassenger)
					d.coord = d.ride.FromAddress.Coord
					d.toDest = geoloc.DistanceAccurate(d.coord.Latitude, d.coord.Longitude, d.ride.ToAddress.Coord.Latitude, d.ride.ToAddress.Coord.Longitude) / 1000
					d.driverState = onRide
				}
			case onRide:
				d.toDest -= float64(conf.Bench.KmByBT)
				if d.toDest <= 0 {
					d.toDest = 0
					d.updateRide(datamodels.Ended)
					d.coord = d.ride.ToAddress.Coord
					d.ride = datamodels.Ride{}
					d.driverState = ready
				}
			}
			d.mu.Unlock()

			if sendPosCount == 0 {
				d.sendCoord()
				sendPosCount = conf.Bench.SendPos
			}
			sendPosCount--
			d.dice = dice(100)

			// case s := <-d.in:
			// 	d.mu.Lock()
			// 	d.driverState = s
			// 	clog.File("Driver", "updateDriver", "%d -> %d", d.id, s)
			// 	d.mu.Unlock()
		}
	}
}
