package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"code.rocketnine.space/tslocum/citylimits/world"

	"code.rocketnine.space/tslocum/citylimits/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowTitle("City Limits")
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetMaxTPS(144)
	ebiten.SetRunnableOnUnfocused(true) // Note - this currently does nothing in ebiten
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}

	parseFlags()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-sigc

		g.Exit()
	}()

	world.StartGame()

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
