package main

import (
	"context"
	"encoding/csv"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"bench_dispatch/datamodels"
	"bench_dispatch/gopool"
	"bench_dispatch/tools/clog"
	"bench_dispatch/tools/confload"

	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
)

var (
	ioTimeout = time.Millisecond * 100
	conf      = &datamodels.ConfigData{}
	pool      *gopool.Pool
	hub       *Hub
	address   []datamodels.AddressRide
	nbAdress  int
)

// Deadliner : Wrapper de connection pour ajouter un timer avant chaque lecture / ecriture
type Deadliner struct {
	net.Conn
	t time.Duration
}

func (d Deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

func (d Deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}

func connect(i int, u url.URL) net.Conn {
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), u.String())
	if err != nil {
		clog.Fatal("main", "Connect", err)
	}

	return conn
}

func loadCSV() int {
	count := -2
	csvfile, err := os.Open("Marseille.csv")
	if err != nil {
		clog.Fatal("main", "CSV", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	r.Comma = ';'
	r.FieldsPerRecord = 3
	//r := csv.NewReader(bufio.NewReader(csvfile))

	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			clog.Fatal("main", "CSV", err)
		}
		long, _ := strconv.ParseFloat(record[1], 64)
		lat, _ := strconv.ParseFloat(record[2], 64)
		tmp := datamodels.AddressRide{Address: record[0], Coord: datamodels.Coordinates{Longitude: long, Latitude: lat}}
		address = append(address, tmp)
		count++
	}
	clog.Trace("main", "Address CSV", "Loaded %d address ...", count)
	return count
}

func main() {
	var exit = make(chan struct{})
	confload.Load("config.ini", conf)

	clog.LogLevel = 5
	clog.StartLogging = true

	nbAdress = loadCSV()

	pool = gopool.NewPool(conf.Workers, conf.QueueSize, 10)
	hub = NewHub(pool)

	u := url.URL{Scheme: "ws", Host: conf.WSserver.Addr, Path: "/ws"}

	poller, err := netpoll.New(nil)
	if err != nil {
		clog.Fatal("server", "WebSocket", err)
	}

	for i := 0; i < conf.Bench.NbDrivers; i++ {
		newCon := connect(i, u)
		safeConn := Deadliner{newCon, ioTimeout}
		driver := hub.Register(safeConn, i)

		desc := netpoll.Must(netpoll.HandleRead(newCon))

		// On ajoute un listener sur la connection
		poller.Start(desc, func(ev netpoll.Event) {
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				// Connexion perdue ou terminée par le client
				poller.Stop(desc)
				hub.Remove(driver)
				return
			}
			// Nouveau message entrant
			pool.Schedule(func() {
				if err := driver.Receive(); err != nil {
					// Pb de reception, la connexion est rompue
					poller.Stop(desc)
					hub.Remove(driver)
				}
			})
		})
	}

	go output()

	<-exit
}
