package main

import (
	"net"
	"sync"

	"bench_dispatch/gopool"
	"bench_dispatch/clog"
)

type message struct {
	userType int
	content  []byte
}

// Hub :
type Hub struct {
	mu      sync.RWMutex
	drivers map[int]*Driver

	pool *gopool.Pool
	out  chan message // Channel de sortie
}

// NewHub : Creation du Hub de Driver
func NewHub(pool *gopool.Pool) *Hub {
	hub := &Hub{
		pool:    pool,
		drivers: make(map[int]*Driver),
	}

	clog.Info("main", "Hub", "Driver Hub initialized.")

	return hub
}

// Register : registers new connection as a User.
func (h *Hub) Register(conn net.Conn, id int, name string) *Driver {
	loc := getNewAdress()
	driver := &Driver{
		hub:         h,
		conn:        conn,
		DriverState: ready,
		Coord:       loc.Coord,
	}
	// driver.in = make(chan UserState, 1)

	h.mu.Lock()
	{
		driver.ID = id
		driver.Name = name
		h.drivers[driver.ID] = driver
	}
	h.mu.Unlock()

	h.pool.Schedule(func() {
		driver.Life()
	})

	return driver
}

func (h *Hub) remove(driver *Driver) bool {
	if _, has := h.drivers[driver.ID]; !has {
		return false
	}

	delete(h.drivers, driver.ID)

	return true
}

// Remove : Supprime un driver / fin de comm
func (h *Hub) Remove(driver *Driver) {
	h.mu.Lock()
	removed := h.remove(driver)
	h.mu.Unlock()

	if !removed {
		return
	}
}
