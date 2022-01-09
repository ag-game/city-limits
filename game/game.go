package game

import (
	"image/color"
	"math/rand"
	"os"
	"sync"
	"time"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/entity"

	"code.rocketnine.space/tslocum/citylimits/asset"
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/citylimits/system"
	"code.rocketnine.space/tslocum/citylimits/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

// game is an isometric demo game.
type game struct {
	w, h int

	audioContext *audio.Context

	op *ebiten.DrawImageOptions

	disableEsc bool

	debugMode  bool
	cpuProfile *os.File

	movementSystem *system.MovementSystem
	renderSystem   *system.RenderSystem

	addedSystems bool

	sync.Mutex
}

// NewGame returns a new isometric demo game.
func NewGame() (*game, error) {
	g := &game{
		audioContext: audio.NewContext(sampleRate),
		op:           &ebiten.DrawImageOptions{},
	}

	err := g.loadAssets()
	if err != nil {
		panic(err)
	}

	const numEntities = 30000
	ECS.Preallocate(numEntities)

	return g, nil
}

func (g *game) tileToGameCoords(x, y int) (float64, float64) {
	return float64(x) * 32, float64(y) * 32
}

// Layout is called when the game's layout changes.
func (g *game) Layout(w, h int) (int, int) {
	if w != g.w || h != g.h {
		world.World.ScreenW, world.World.ScreenH = w, h
		g.w, g.h = w, h
	}
	return g.w, g.h
}

func (g *game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.Exit()
		return nil
	}

	if world.World.ResetGame {
		world.Reset()

		world.BuildStructure(world.StructureHouse1, 0, 0)

		world.BuildStructure(world.StructureHouse1, 8, 12)

		world.BuildStructure(world.StructurePoliceStation, 12, 12)

		// TODO

		world.World.HoverStructure = world.StructurePoliceStation

		if world.World.Player == 0 {
			world.World.Player = entity.NewPlayer()
		}

		const playerStartOffset = 128

		w := float64(world.World.Map.Width * world.World.Map.TileWidth)
		h := float64(world.World.Map.Height * world.World.Map.TileHeight)

		position := ECS.Component(world.World.Player, component.PositionComponentID).(*component.PositionComponent)
		position.X, position.Y = w/2, h-playerStartOffset

		if !g.addedSystems {
			g.addSystems()

			g.addedSystems = true // TODO
		}

		rand.Seed(time.Now().UnixNano())

		world.World.ResetGame = false
		world.World.GameOver = false
	}

	err := ECS.Update()
	if err != nil {
		return err
	}
	return nil
}

// renderSprite renders a sprite on the screen.
func (g *game) renderSprite(x float64, y float64, offsetx float64, offsety float64, angle float64, geoScale float64, colorScale float64, alpha float64, hFlip bool, vFlip bool, sprite *ebiten.Image, target *ebiten.Image) int {
	if alpha < .01 || colorScale < .01 {
		return 0
	}

	xi, yi := world.CartesianToIso(float64(x), float64(y))

	padding := float64(world.TileSize) * world.World.CamScale
	cx, cy := float64(world.World.ScreenW/2), float64(world.World.ScreenH/2)

	// Skip drawing tiles that are out of the screen.
	drawX, drawY := world.IsoToScreen(xi, yi)
	if drawX+padding < 0 || drawY+padding < 0 || drawX > float64(world.World.ScreenW) || drawY > float64(world.World.ScreenH) {
		//log.Println("SKIP", drawX, drawY, world.World.ScreenW, world.World.ScreenH)
		return 0
	}

	g.op.GeoM.Reset()

	/*if hFlip {
		s.op.GeoM.Scale(-1, 1)
		s.op.GeoM.Translate(TileWidth, 0)
	}
	if vFlip {
		s.op.GeoM.Scale(1, -1)
		s.op.GeoM.Translate(0, TileWidth)
	}*/

	// Move to current isometric position.
	g.op.GeoM.Translate(xi, yi+offsety)
	// Translate camera position.
	g.op.GeoM.Translate(-world.World.CamX, -world.World.CamY)
	// Zoom.
	g.op.GeoM.Scale(world.World.CamScale, world.World.CamScale)
	// Center.
	g.op.GeoM.Translate(cx, cy)

	target.DrawImage(sprite, g.op)

	/*s.op.GeoM.Scale(geoScale, geoScale)
	// Rotate
	s.op.GeoM.Translate(offsetx, offsety)
	s.op.GeoM.Rotate(angle)
	// Move to current isometric position.
	s.op.GeoM.Translate(x, y)
	// Translate camera position.
	s.op.GeoM.Translate(-world.World.CamX, -world.World.CamY)
	// Zoom.
	s.op.GeoM.Scale(s.camScale, s.camScale)
	// Center.
	//s.op.GeoM.Translate(float64(s.ScreenW/2.0), float64(s.ScreenH/2.0))

	s.op.ColorM.Scale(colorScale, colorScale, colorScale, alpha)

	target.DrawImage(sprite, s.op)

	s.op.ColorM.Reset()*/

	return 1
}

func (g *game) Draw(screen *ebiten.Image) {
	// Handle background rendering separately to simplify design.
	for i := range world.World.Level.Tiles {
		for x := range world.World.Level.Tiles[i] {
			for y, tile := range world.World.Level.Tiles[i][x] {
				if tile == nil || tile.Sprite == nil {
					continue
				}
				g.renderSprite(float64(x), float64(y), 0, float64(i*-80), 0, 1, 1, 1, false, false, tile.Sprite, screen)
			}
		}
	}

	err := ECS.Draw(screen)
	if err != nil {
		panic(err)
	}
}

func (g *game) addSystems() {
	ecs := ECS

	g.movementSystem = system.NewMovementSystem()
	ecs.AddSystem(system.NewPlayerMoveSystem(world.World.Player, g.movementSystem))
	ecs.AddSystem(system.NewplayerFireSystem())
	ecs.AddSystem(g.movementSystem)
	ecs.AddSystem(system.NewCreepSystem())
	ecs.AddSystem(system.NewCameraSystem())
	g.renderSystem = system.NewRenderSystem()
	ecs.AddSystem(g.renderSystem)
	ecs.AddSystem(system.NewRenderMessageSystem())
	ecs.AddSystem(system.NewRenderDebugTextSystem(world.World.Player))
	ecs.AddSystem(system.NewProfileSystem(world.World.Player))
}

func (g *game) loadAssets() error {
	asset.ImgWhiteSquare.Fill(color.White)
	asset.LoadSounds(g.audioContext)
	return nil
}

func (g *game) Exit() {
	os.Exit(0)
}
