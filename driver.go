package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"bench_dispatch/datamodels"
	"bench_dispatch/tools/clog"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)


// UserState : Etat d'un Driver
type UserState int

const (
	idle UserState = iota
	ready
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
}

////////////////
// Lectures
////////////////

// Receive : Lit le message en attente.
func (d *Driver) Receive() error {
	req, _ := d.readRequest()

	if req == nil {
		// Message vide de contôle.
		clog.Trace("main", "Driver", "Empty request")
		req = &datamodels.Request{}
	} else {
		switch req.Method {
		// case "NewRide":
		// 	tmpRide := rm.NewRide(req.Params)
		// 	// dist := geoloc.DistanceFromHome(tmpRide.ToAddress.Coord.Latitude, tmpRide.ToAddress.Coord.Longitude)
		// 	clog.Trace("Driver", "NewRide", "ID : %s", tmpRide.ID)
		// case "AcceptRide":
		// 	ride, err := rm.AcceptRide(d, req.Params)
		// 	r := datamodels.Request{ID: req.ID, Method: "AcceptRideResponse", Params: ride, Status: err}
		// 	return d.write(r)
		// case "UpdateDriverLocation":
		// 	mapstructure.Decode(req.Params, &d.coord)
		// 	return nil
		default:
			clog.Trace("main", "Driver", "RECV: %v", req.Params)
		}
	}
	return nil
}

func (d *Driver) readRequest() (*datamodels.Request, error) {
	d.io.Lock()
	defer d.io.Unlock()

	h, r, err := wsutil.NextReader(d.conn, ws.StateClientSide)
	if err != nil {
		clog.Error("Driver", "readRequest::NextReader", "%s", err)
		return nil, err
	}
	if h.OpCode.IsControl() {
		clog.Info("Driver", "readRequest::IsControl", "Opcode Control : %v", h)
		return nil, wsutil.ControlFrameHandler(d.conn, ws.StateClientSide)(h, r)
	}

	req := &datamodels.Request{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(req); err != nil {
		clog.Error("Driver", "readRequest::Decode", "%s", err)
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

// Life : Simulation des actions d'un Driver
func (d *Driver) Life() {
	posTimer, _ := time.ParseDuration(fmt.Sprintf("%ds", conf.Bench.SendPosInterval))
	ticker := time.NewTicker(posTimer)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			d.sendCoord()
		}
	}
}
