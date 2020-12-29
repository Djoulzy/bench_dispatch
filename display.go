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

	termbox.SetCursor(1, 1)
	for i := 0; i < len(driverList); i++ {
		tbprintf(1, i, termbox.ColorDefault, termbox.ColorDefault, "%d", driverList[i].id)
		switch driverList[i].driverState {
		case idle:
			tbprintf(4, i, termbox.ColorDefault, termbox.ColorDefault, "I")
		case ready:
			tbprintf(4, i, termbox.ColorGreen, termbox.ColorDefault, "R")
		case moving:
			tbprintf(4, i, termbox.ColorBlack, termbox.ColorGreen, "M")
		case onRide:
			tbprintf(4, i, termbox.ColorYellow, termbox.ColorDefault, "O")
		case err:
			tbprintf(4, i, termbox.ColorRed, termbox.ColorDefault, "E")
		default:
		}
		tbprintf(6, i, termbox.ColorDefault, termbox.ColorDefault, "%f %f", driverList[i].coord.Latitude, driverList[i].coord.Longitude)
	}

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
