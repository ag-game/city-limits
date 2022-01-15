//go:build js && wasm
// +build js,wasm

package main

import (
	"code.rocketnine.space/tslocum/citylimits/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func parseFlags() {
	world.World.DisableEsc = true

	// Adjust minimum zoom level due to performance decrease when targeting WASM.
	world.CameraMinZoom = 0.6

	ebiten.SetFullscreen(true)
}
