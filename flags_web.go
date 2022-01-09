//go:build js && wasm
// +build js,wasm

package main

import (
	"code.rocketnine.space/tslocum/citylimits/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func parseFlags() {
	world.World.DisableEsc = true

	ebiten.SetFullscreen(true)
}
