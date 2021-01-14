package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"bench_dispatch/datamodels"
	"bench_dispatch/gopool"
	"bench_dispatch/clog"
	"bench_dispatch/confload"

	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
)

var (
	ioTimeout = time.Millisecond * 100
	conf      = &datamodels.ConfigData{}
	pool      *gopool.Pool
	hub       *Hub
	address   []datamodels.Address
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
		tmp := datamodels.Address{Name: record[0], Coord: datamodels.Coordinates{Longitude: long, Latitude: lat}}
		address = append(address, tmp)
		count++
	}
	clog.Trace("main", "Address CSV", "Loaded %d address ...", count)
	return count
}

func getName(nb int) string {
	f, err := os.Open("name.txt")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 0; i <= nb; i++ {
		scanner.Scan()
	}
	return scanner.Text()
}

func getNewAdress() datamodels.Address {
	tmp := rand.Intn(nbAdress) + 1
	return address[tmp]
}

func main() {
	var exit = make(chan struct{})
	confload.Load("config.ini", conf)

	clog.LogLevel = 5
	clog.StartLogging = true
	if conf.FileLog != "" {
		clog.EnableFileLog(conf.FileLog)
	}

	nbAdress = loadCSV()

	pool = gopool.NewPool(conf.Workers, conf.QueueSize, 10)
	hub = NewHub(pool)

	u := url.URL{Scheme: "ws", Host: conf.WSserver.Addr, Path: "/ws"}

	poller, err := netpoll.New(nil)
	if err != nil {
		clog.Fatal("server", "WebSocket", err)
	}

	go output()

	for i := 0; i < conf.Bench.NbDrivers; i++ {
		newCon := connect(i, u)
		safeConn := Deadliner{newCon, ioTimeout}
		driver := hub.Register(safeConn, i, getName(i))

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

		time.Sleep(time.Second)
	}

	<-exit
}
