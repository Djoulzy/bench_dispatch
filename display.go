package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func output() {
	err := termbox.Init()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		displayHub()

		time.Sleep(time.Second)
	}
}

// DisplayHub : Affiche l'etat du Hub
func displayHub() {
	termbox.SetCursor(1, 1)
	for i := 0; i < len(hub.drivers); i++ {
		tbprintf(1, i, termbox.ColorDefault, termbox.ColorDefault, "%d", hub.drivers[i].id)
		switch hub.drivers[i].driverState {
		case idle:
			tbprintf(4, i, termbox.ColorDefault, termbox.ColorDefault, "I")
		case ready:
			tbprintf(4, i, termbox.ColorGreen, termbox.ColorDefault, "R")
		case onRide:
			tbprintf(4, i, termbox.ColorYellow, termbox.ColorDefault, "O")
		case err:
			tbprintf(4, i, termbox.ColorRed, termbox.ColorDefault, "E")
		}
		tbprintf(6, i, termbox.ColorDefault, termbox.ColorDefault, "%f %f", hub.drivers[i].coord.Latitude, hub.drivers[i].coord.Longitude)
	}
	termbox.Flush()
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		if ev.Ch == 'q' {
			os.Exit(0)
		}
	}
}

// This function is often useful:
func tbprintf(x, y int, fg, bg termbox.Attribute, format string, vars ...interface{}) {
	for _, c := range fmt.Sprintf(format, vars...) {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
