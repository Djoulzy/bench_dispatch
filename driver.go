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
	conn        io.ReadWriteCloser
	hub         *Hub
	id          int
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
	req, _ := d.readRequest()

	if req == nil {
		// Message vide de contôle.
		// clog.Trace("main", "Driver", "Empty request")
		req = &datamodels.Request{}
	} else {
		switch req.Method {
		case "NewRide":
			if d.driverState == ready {
				d.acceptRide(req.Params)
				d.driverState = waitOK
			}
		case "AcceptRideResponse":
			var tmpRide datamodels.Ride
			mapstructure.Decode(req.Params, &tmpRide)

			if d.driverState == waitOK && req.Status.ID == 0 {
				d.ride = tmpRide
				d.driverState = moving
				d.updateRide(datamodels.Approach)
				d.toDest = geoloc.DistanceAccurate(d.coord.Latitude, d.coord.Longitude, tmpRide.FromAddress.Coord.Latitude, tmpRide.FromAddress.Coord.Longitude) / 1000
			} else {
				d.driverState = ready
			}
		// case "AcceptRide":
		// 	ride, err := rm.AcceptRide(d, req.Params)
		// 	r := datamodels.Request{ID: req.ID, Method: "AcceptRideResponse", Params: ride, Status: err}
		// 	return d.write(r)
		// case "UpdateDriverLocation":
		// 	mapstructure.Decode(req.Params, &d.coord)
		// 	return nil
		default:
			// clog.Trace("main", "Driver", "RECV: %v", req.Params)
		}
	}
	return nil
}

func (d *Driver) readRequest() (*datamodels.Request, error) {
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

	req := &datamodels.Request{}
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
func (d *Driver) writeResultTo(req *datamodels.Request, result datamodels.DataParams) error {
	return d.write(datamodels.Response{
		ID:     req.ID,
		Result: result,
	})
}

func (d *Driver) writeErrorTo(req *datamodels.Request, id int, err string) error {
	req.Status = datamodels.Error{
		ID:      id,
		Message: err,
	}
	return d.write(req)
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

func (d *Driver) sendCoord() {
	req := datamodels.Request{
		ID:     d.id,
		Method: "UpdateDriverLocation",
		Params: datamodels.Coordinates{
			Latitude:  d.coord.Latitude,
			Longitude: d.coord.Longitude,
		},
	}

	d.write(req)
}

func dice(nb int) int {
	return rand.Intn(nb) + 1
}

func (d *Driver) updateRide(state datamodels.RideState) {
	req := datamodels.Request{
		ID:     d.id,
		Method: "ChangeRideState",
		Params: datamodels.RideUpdate{
			ID:    d.ride.ID,
			State: state,
		},
	}

	d.write(req)
}

func (d *Driver) acceptRide(params datamodels.DataParams) {
	var tmpRide datamodels.Ride
	mapstructure.Decode(params, &tmpRide)

	req := datamodels.Request{
		ID:     d.id,
		Method: "AcceptRide",
		Params: datamodels.AcceptRide{RideID: tmpRide.ID},
	}

	d.write(req)
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

// Life : Simulation des actions d'un Driver
func (d *Driver) Life() {
	baseTimer, _ := time.ParseDuration(fmt.Sprintf("%ds", conf.Bench.BaseTimer))
	ticker := time.NewTicker(baseTimer)
	defer func() {
		ticker.Stop()
	}()

	idleCount := 0
	sendPosCount := 0

	for {
		select {
		case <-ticker.C:
			switch d.driverState {
			case idle:
				if idleCount == 0 {
					d.driverState = ready
				} else {
					idleCount--
				}
			case ready:
				if dice(100) > 95 {
					d.driverState = idle
					idleCount = conf.Bench.IdleDuration
					sendPosCount = 0
					d.createCourse()
				}
			case moving:
				d.toDest -= float64(conf.Bench.KmByBT)
				if d.toDest <= 0 {
					d.driverState = onRide
					d.updateRide(datamodels.PickUpPassenger)
					d.coord = d.ride.FromAddress.Coord
					d.toDest = geoloc.DistanceAccurate(d.coord.Latitude, d.coord.Longitude, d.ride.ToAddress.Coord.Latitude, d.ride.ToAddress.Coord.Longitude) / 1000
				}
			case onRide:
				d.toDest -= float64(conf.Bench.KmByBT)
				if d.toDest <= 0 {
					d.driverState = ready
					d.updateRide(datamodels.Ended)
					d.coord = d.ride.ToAddress.Coord
				}
			}

			if sendPosCount == 0 {
				d.sendCoord()
				sendPosCount = conf.Bench.SendPos
			}
			sendPosCount--
			d.dice = dice(100)
		}
	}
}
