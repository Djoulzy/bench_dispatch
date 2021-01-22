package main

import (
	"bench_dispatch/clog"
	"bench_dispatch/datamodels"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/copier"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type cpyDriver struct {
	ID          int
	Name        string
	DriverState UserState
	Coord       datamodels.Coordinates
	Dice        int
	Ride        datamodels.Ride
	ToDest      float64
}

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
	var zeDriver cpyDriver
	termbox.SetCursor(1, 1)
	for i = 0; i < nbDrivers; i++ {
		hub.mu.RLock()
		if err := copier.Copy(&zeDriver, hub.drivers[i]); err != nil {
			panic("Can't copy driver")
		}
		// fmt.Printf("%v\n", zeDriver)
		hub.mu.RUnlock()

		tbprintf(1, i, termbox.ColorDefault, termbox.ColorDefault, "%s", zeDriver.Name)
		switch zeDriver.DriverState {
		case idle:
			tbprintf(15, i, termbox.ColorDefault, termbox.ColorDefault, "I d l e")
		case ready:
			tbprintf(15, i, termbox.ColorGreen, termbox.ColorDefault, " Ready ")
		case waitOK:
			tbprintf(15, i, termbox.ColorBlack, termbox.ColorCyan, "Wait OK")
		case moving:
			tbprintf(15, i, termbox.ColorBlack, termbox.ColorGreen, "Approch")
		case onRide:
			tbprintf(15, i, termbox.ColorBlack, termbox.ColorYellow, "On Ride")
		case err:
			tbprintf(15, i, termbox.ColorRed, termbox.ColorDefault, "-Error-")
		default:
		}
		tbprintf(25, i, termbox.ColorDefault, termbox.ColorDefault, "%f %f", zeDriver.Coord.Latitude, zeDriver.Coord.Longitude)
		tbprintf(44, i, termbox.ColorDefault, termbox.ColorDefault, "%.1f Km ", zeDriver.ToDest)
		if zeDriver.Ride.ToAddress.Name == "" {
			tbprintf(54, i, termbox.ColorDefault, termbox.ColorDefault, "                                                    ")
		} else {
			tbprintf(54, i, termbox.ColorDefault, termbox.ColorDefault, "%s", zeDriver.Ride.ToAddress.Name)
		}
	}
	t := time.Now()
	tbprintf(25, i, termbox.ColorDefault, termbox.ColorDefault, "%s", t.Format("15:04:05"))

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
