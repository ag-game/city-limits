package system

import (
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"code.rocketnine.space/tslocum/citylimits/asset"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type playerMoveSystem struct {
	player       gohan.Entity
	movement     *MovementSystem
	lastWalkDirL bool

	rewindTicks    int
	nextRewindTick int

	scrollDragX, scrollDragY         int
	scrollCamStartX, scrollCamStartY float64
}

func NewPlayerMoveSystem(player gohan.Entity, m *MovementSystem) *playerMoveSystem {
	return &playerMoveSystem{
		player:      player,
		movement:    m,
		scrollDragX: -1,
		scrollDragY: -1,
	}
}

func (_ *playerMoveSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (_ *playerMoveSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *playerMoveSystem) Update(ctx *gohan.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) && !world.World.DisableEsc {
		os.Exit(0)
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
		v := 1
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			v = 2
		}
		if world.World.Debug == v {
			world.World.Debug = 0
		} else {
			world.World.Debug = v
		}
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		world.World.NoClip = !world.World.NoClip
		return nil
	}

	if !world.World.GameStarted {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			world.StartGame()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		world.World.MuteMusic = !world.World.MuteMusic
		if world.World.MuteMusic {
			asset.SoundMusic.Pause()
		} else {
			asset.SoundMusic.Play()
		}
	}

	if world.World.GameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			world.World.ResetGame = true
		}
		return nil
	}

	// Update target zoom level.
	var scrollY float64
	if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		scrollY = -0.25
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		scrollY = .25
	} else {
		_, scrollY = ebiten.Wheel()
		if scrollY < -1 {
			scrollY = -1
		} else if scrollY > 1 {
			scrollY = 1
		}
	}
	world.World.CamScaleTarget += scrollY * (world.World.CamScaleTarget / 7)
	const minZoom = .4
	const maxZoom = 1
	if world.World.CamScaleTarget < minZoom {
		world.World.CamScaleTarget = minZoom
	} else if world.World.CamScaleTarget > maxZoom {
		world.World.CamScaleTarget = maxZoom
	}

	// Smooth zoom transition.
	div := 10.0
	if world.World.CamScaleTarget > world.World.CamScale {
		world.World.CamScale += (world.World.CamScaleTarget - world.World.CamScale) / div
	} else if world.World.CamScaleTarget < world.World.CamScale {
		world.World.CamScale -= (world.World.CamScale - world.World.CamScaleTarget) / div
	}

	pressLeft := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	pressRight := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	pressUp := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	pressDown := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)

	const camSpeed = 10
	if (pressLeft && !pressRight) ||
		(pressRight && !pressLeft) {
		if pressLeft {
			world.World.CamX -= camSpeed
		} else {
			world.World.CamX += camSpeed
		}
	}

	if (pressUp && !pressDown) ||
		(pressDown && !pressUp) {
		if pressUp {
			world.World.CamY -= camSpeed
		} else {
			world.World.CamY += camSpeed
		}
	}

	const scrollEdgeSize = 1
	x, y := ebiten.CursorPosition()
	if !world.World.GotCursorPosition {
		if x != 0 || y != 0 {
			world.World.GotCursorPosition = true
		} else {
			return nil
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		if s.scrollDragX == -1 && s.scrollDragY == -1 {
			// TODO disabled due to possible ebiten bug
			//ebiten.SetCursorMode(ebiten.CursorModeCaptured)

			s.scrollDragX, s.scrollDragY = x, y
			s.scrollCamStartX, s.scrollCamStartY = world.World.CamX, world.World.CamY
		} else {
			dx, dy := float64(x-s.scrollDragX)/world.World.CamScale, float64(y-s.scrollDragY)/world.World.CamScale
			world.World.CamX, world.World.CamY = s.scrollCamStartX-dx, s.scrollCamStartY-dy
		}
	} else {
		if s.scrollDragX != -1 && s.scrollDragY != -1 {
			s.scrollDragX, s.scrollDragY = -1, -1
			//ebiten.SetCursorMode(ebiten.CursorModeVisible)
		} else if x >= -2 && y >= -2 && x < world.World.ScreenW+2 && y < world.World.ScreenH+2 {
			// Pan via screen edge.
			if x <= scrollEdgeSize {
				world.World.CamX -= camSpeed
			} else if x >= world.World.ScreenW-scrollEdgeSize-1 {
				world.World.CamX += camSpeed
			}
			if y <= scrollEdgeSize {
				world.World.CamY -= camSpeed
			} else if y >= world.World.ScreenH-scrollEdgeSize-1 {
				world.World.CamY += camSpeed
			}
		}
	}

	if x < world.SidebarWidth {
		world.World.Level.ClearHoverSprites()
		world.World.HoverX, world.World.HoverY = 0, 0
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			button := world.HUDButtonAt(x, y)
			if button != nil {
				if button.StructureType != 0 {
					if world.World.HoverStructure == button.StructureType {
						world.SetHoverStructure(0) // Deselect.
					} else {
						world.SetHoverStructure(button.StructureType)
					}
					asset.SoundSelect.Rewind()
					asset.SoundSelect.Play()
				}
			}
		}
	} else if world.World.HoverStructure != 0 {
		tileX, tileY := world.ScreenToCartesian(x, y)
		if tileX >= 0 && tileY >= 0 && tileX < 256 && tileY < 256 {
			multiUseStructure := world.World.HoverStructure == world.StructureBulldozer || world.World.HoverStructure == world.StructureRoad
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || (multiUseStructure && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
				if world.World.HoverStructure == world.StructureBulldozer {
					for i := range world.World.Level.Tiles {
						world.World.Level.Tiles[i][int(tileX)][int(tileY)].Sprite = nil

						var img *ebiten.Image
						if i == 0 {
							img = world.World.TileImages[world.DirtTile+world.World.TileImagesFirstGID]
						}
						world.World.Level.Tiles[i][int(tileX)][int(tileY)].EnvironmentSprite = img
					}
				} else {
					_, err := world.BuildStructure(world.World.HoverStructure, false, int(tileX), int(tileY))
					if err == nil {
						sounds := []*audio.Player{
							asset.SoundPop1,
							asset.SoundPop2,
							asset.SoundPop3,
							asset.SoundPop4,
							asset.SoundPop5,
						}
						sound := sounds[rand.Intn(len(sounds))]
						sound.Rewind()
						sound.Play()
					}
				}
				world.BuildStructure(world.World.HoverStructure, true, int(tileX), int(tileY))
			} else if int(tileX) != world.World.HoverX || int(tileY) != world.World.HoverY {
				world.BuildStructure(world.World.HoverStructure, true, int(tileX), int(tileY))
			}
			world.World.HoverX, world.World.HoverY = int(tileX), int(tileY)
		}
	}

	return nil
}

func (s *playerMoveSystem) Draw(_ *gohan.Context, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}

func deltaXY(x1, y1, x2, y2 float64) (dx float64, dy float64) {
	dx, dy = x1-x2, y1-y2
	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}
	return dx, dy
}
