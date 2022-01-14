package game

import (
	"image/color"
	"math/rand"
	"os"
	"sync"
	"time"

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

	updateTicks int

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

// Layout is called when the game's layout changes.
func (g *game) Layout(w, h int) (int, int) {
	if w != g.w || h != g.h {
		world.World.ScreenW, world.World.ScreenH = w, h
		g.w, g.h = w, h
		world.World.HUDUpdated = true
	}
	return g.w, g.h
}

func (g *game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.Exit()
		return nil
	}

	const updateSidebarDelay = 144 * 3
	g.updateTicks++
	if g.updateTicks == updateSidebarDelay {
		world.World.HUDUpdated = true
		g.updateTicks = 0
	}

	if world.World.ResetGame {
		world.Reset()

		rand.Seed(time.Now().UnixNano())

		err := world.LoadTileset()
		if err != nil {
			return err
		}

		// Fill below ground layer.
		grassTile := uint32(11*32 + (0))
		treeTileA := uint32(5*32 + (25))
		treeTileB := uint32(5*32 + (27))
		treeTileA = treeTileB // todo
		var img uint32
		for x := range world.World.Level.Tiles[0] {
			for y := range world.World.Level.Tiles[0][x] {
				img = world.DirtTile
				if rand.Intn(128) == 0 {
					img = grassTile
					world.World.Level.Tiles[0][x][y].EnvironmentSprite = world.World.TileImages[img+world.World.TileImagesFirstGID]
					for offsetX := -2 - rand.Intn(7); offsetX < 2+rand.Intn(7); offsetX++ {
						for offsetY := -2 - rand.Intn(7); offsetY < 2+rand.Intn(7); offsetY++ {
							if x+offsetX >= 0 && y+offsetY >= 0 && x+offsetX < 256 && y+offsetY < 256 {
								world.World.Level.Tiles[0][x+offsetX][y+offsetY].EnvironmentSprite = world.World.TileImages[img+world.World.TileImagesFirstGID]
								if rand.Intn(2) == 0 {
									if rand.Intn(3) == 0 {
										world.World.Level.Tiles[1][x+offsetX][y+offsetY].EnvironmentSprite = world.World.TileImages[treeTileA+world.World.TileImagesFirstGID]
									} else {
										world.World.Level.Tiles[1][x+offsetX][y+offsetY].EnvironmentSprite = world.World.TileImages[treeTileB+world.World.TileImagesFirstGID]
									}
								}
							}
						}
					}
				} else {
					if world.World.Level.Tiles[0][x][y].EnvironmentSprite != nil {
						continue
					}
					world.World.Level.Tiles[0][x][y].EnvironmentSprite = world.World.TileImages[img+world.World.TileImagesFirstGID]
				}
			}
		}

		// Load HUD sprites.

		world.HUDButtons = []*world.HUDButton{
			{
				Label:         "Bulldoze",
				Sprite:        world.DrawMap(world.StructureBulldozer),
				SpriteOffsetX: 12,
				SpriteOffsetY: -48,
				StructureType: world.StructureBulldozer,
			}, {
				Label:         "Road",
				Sprite:        world.DrawMap(world.StructureRoad),
				SpriteOffsetX: -12,
				SpriteOffsetY: -28,
				StructureType: world.StructureRoad,
			}, {
				Label:         "House",
				Sprite:        world.DrawMap(world.StructureHouse1),
				SpriteOffsetX: -12,
				SpriteOffsetY: -28,
				StructureType: world.StructureHouse1,
			}, {
				Label:         "Business",
				Sprite:        world.DrawMap(world.StructureBusiness1),
				SpriteOffsetX: -10,
				SpriteOffsetY: -28,
				StructureType: world.StructureBusiness1,
			}, {
				Label:         "Police Station",
				SpriteOffsetX: -19,
				Sprite:        world.DrawMap(world.StructurePoliceStation),
				StructureType: world.StructurePoliceStation,
			},
		}

		// TODO

		if world.World.Player == 0 {
			world.World.Player = entity.NewPlayer()
		}

		if !g.addedSystems {
			g.addSystems()

			g.addedSystems = true // TODO
		}

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
	if alpha < .01 {
		return 0
	}

	xi, yi := world.CartesianToIso(float64(x), float64(y))

	padding := float64(world.TileSize) * world.World.CamScale
	cx, cy := float64(world.World.ScreenW/2), float64(world.World.ScreenH/2)

	// Skip drawing tiles that are out of the screen.
	drawX, drawY := world.IsoToScreen(xi, yi)
	if drawX+padding < 0 || drawY+padding < 0 || drawX-padding > float64(world.World.ScreenW) || drawY-padding > float64(world.World.ScreenH) {
		return 0
	}

	g.op.GeoM.Reset()

	if hFlip {
		g.op.GeoM.Scale(-1, 1)
		g.op.GeoM.Translate(world.TileSize, 0)
	}
	if vFlip {
		g.op.GeoM.Scale(1, -1)
		g.op.GeoM.Translate(0, world.TileSize)
	}

	// Move to current isometric position.
	g.op.GeoM.Translate(xi, yi+offsety)
	// Translate camera position.
	g.op.GeoM.Translate(-world.World.CamX, -world.World.CamY)
	// Zoom.
	g.op.GeoM.Scale(world.World.CamScale, world.World.CamScale)
	// Center.
	g.op.GeoM.Translate(cx, cy)

	g.op.ColorM.Reset()
	g.op.ColorM.Scale(colorScale, colorScale, colorScale, alpha)

	target.DrawImage(sprite, g.op)
	g.op.ColorM.Reset()

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
	var drawn int
	for i := range world.World.Level.Tiles {
		for x := range world.World.Level.Tiles[i] {
			for y, tile := range world.World.Level.Tiles[i][x] {
				if tile == nil {
					continue
				}
				var sprite *ebiten.Image
				colorScale := 1.0
				alpha := 1.0
				if tile.HoverSprite != nil {
					sprite = tile.HoverSprite
					colorScale = 0.6
					if !world.World.HoverValid {
						colorScale = 0.2
					}
				} else if tile.Sprite != nil {
					sprite = tile.Sprite
					// TODO
					/*if i > 0 && world.World.HoverStructure == world.StructureRoad {
						alpha = 0.6
					}*/
				} else if tile.EnvironmentSprite != nil {
					sprite = tile.EnvironmentSprite
				} else {
					continue
				}
				drawn += g.renderSprite(float64(x), float64(y), 0, float64(i*-40), 0, 1, colorScale, alpha, false, false, sprite, screen)
			}
		}
	}
	world.World.EnvironmentSprites = drawn

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
	ecs.AddSystem(system.NewRenderHudSystem())
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
