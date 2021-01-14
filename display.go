package main

import (
	"bench_dispatch/tools/clog"
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func output() {
	err := termbox.Init()
	termbox.HideCursor()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
	driverList := hub.drivers
	hub.mu.RUnlock()

	var i int
	termbox.SetCursor(1, 1)
	for i = 0; i < len(driverList); i++ {
		tbprintf(1, i, termbox.ColorDefault, termbox.ColorDefault, "%s", driverList[i].name)
		switch driverList[i].driverState {
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
		tbprintf(25, i, termbox.ColorDefault, termbox.ColorDefault, "%f %f", driverList[i].coord.Latitude, driverList[i].coord.Longitude)
		tbprintf(44, i, termbox.ColorDefault, termbox.ColorDefault, "%.1f Km ", driverList[i].toDest)
		if driverList[i].ride.ToAddress.Name == "" {
			tbprintf(54, i, termbox.ColorDefault, termbox.ColorDefault, "                                                    ")
		} else {
			tbprintf(54, i, termbox.ColorDefault, termbox.ColorDefault, "%s", driverList[i].ride.ToAddress.Name)
		}
	}
	t := time.Now()
	tbprintf(25, i, termbox.ColorDefault, termbox.ColorDefault, "%d:%d:%d", t.Hour(), t.Minute(), t.Second())

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
