package main

import (
	"bench_dispatch/clog"
	"bench_dispatch/datamodels"
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func output() {
	if err := termbox.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	termbox.HideCursor()

	go func() {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 'q' {
				termbox.Close()
				hub.disconnectAll()
				os.Exit(0)
			}
		}
	}()

	for {
		displayHub()
		time.Sleep(time.Second)
	}
}

// DisplayHub : Affiche l'etat du Hub
func displayHub() {
	hub.mu.RLock()
	nbDrivers := len(hub.drivers)
	hub.mu.RUnlock()

	var i int
	termbox.SetCursor(1, 1)
	for i = 1; i <= nbDrivers; i++ {
		hub.mu.RLock()
		zeDriver := hub.drivers[i]
		hub.mu.RUnlock()

		zeDriver.mu.Lock()
		tbprintf(1, i, termbox.ColorDefault, termbox.ColorDefault, "%s", zeDriver.Name)
		switch zeDriver.DriverState {
		case datamodels.Offline:
			tbprintf(15, i, termbox.ColorDefault, termbox.ColorDefault, "Offline ")
		case datamodels.Free:
			tbprintf(15, i, termbox.ColorGreen, termbox.ColorDefault, " Ready  ")
		case datamodels.WaitOK:
			tbprintf(15, i, termbox.ColorBlack, termbox.ColorCyan, "Wait OK ")
		case datamodels.Moving:
			tbprintf(15, i, termbox.ColorBlack, termbox.ColorGreen, "Approch ")
		case datamodels.Occupied:
			tbprintf(15, i, termbox.ColorBlack, termbox.ColorYellow, "Occupied")
		case datamodels.Billing:
			tbprintf(15, i, termbox.ColorYellow, termbox.ColorBlack, "Payement")
		case datamodels.Err:
			tbprintf(15, i, termbox.ColorRed, termbox.ColorDefault, "-Error- ")
		case datamodels.WaitACK:
			tbprintf(15, i, termbox.ColorRed, termbox.ColorBlack, "Wait ACK")
		default:
		}
		tbprintf(26, i, termbox.ColorDefault, termbox.ColorDefault, "%f %f", zeDriver.Coord.Latitude, zeDriver.Coord.Longitude)
		tbprintf(46, i, termbox.ColorDefault, termbox.ColorDefault, "%.1f Km ", zeDriver.ToDest)
		if zeDriver.Ride.ToAddress.Name == "" {
			tbprintf(55, i, termbox.ColorDefault, termbox.ColorDefault, "                                                    ")
		} else {
			tbprintf(55, i, termbox.ColorDefault, termbox.ColorDefault, "%s", zeDriver.Ride.ToAddress.Name)
		}
		zeDriver.mu.Unlock()
	}
	t := time.Now()
	tbprintf(28, i, termbox.ColorDefault, termbox.ColorDefault, "%s", t.Format("15:04:05"))

	err := termbox.Flush()
	if err != nil {
		clog.Fatal("display", "flush", err)
	}
}

// This function is often useful:
func tbprintf(x, y int, fg, bg termbox.Attribute, format string, vars ...interface{}) {
	for _, c := range fmt.Sprintf(format, vars...) {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
