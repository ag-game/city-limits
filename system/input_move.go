package system

import (
	"os"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	moveSpeed = 1.5
)

type playerMoveSystem struct {
	player       gohan.Entity
	movement     *MovementSystem
	lastWalkDirL bool

	rewindTicks    int
	nextRewindTick int
}

func NewPlayerMoveSystem(player gohan.Entity, m *MovementSystem) *playerMoveSystem {
	return &playerMoveSystem{
		player:   player,
		movement: m,
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
		// TODO mute sound
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

	// Smooth zoom transition.
	div := 10.0
	if world.World.CamScaleTarget > world.World.CamScale {
		world.World.CamScale += (world.World.CamScaleTarget - world.World.CamScale) / div
	} else if world.World.CamScaleTarget < world.World.CamScale {
		world.World.CamScale -= (world.World.CamScale - world.World.CamScaleTarget) / div
	}

	pressLeft := ebiten.IsKeyPressed(ebiten.KeyLeft)
	pressRight := ebiten.IsKeyPressed(ebiten.KeyRight)
	pressUp := ebiten.IsKeyPressed(ebiten.KeyUp)
	pressDown := ebiten.IsKeyPressed(ebiten.KeyDown)

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

	x, y := ebiten.CursorPosition()
	if !world.World.GotCursorPosition {
		if x != 0 || y != 0 {
			world.World.GotCursorPosition = true
		} else {
			return nil
		}
	}
	if x == 0 {
		world.World.CamX -= camSpeed
	} else if x == world.World.ScreenW-1 {
		world.World.CamX += camSpeed
	}
	if y == 0 {
		world.World.CamY -= camSpeed
	} else if y == world.World.ScreenH-1 {
		world.World.CamY += camSpeed
	}

	if world.World.HoverStructure != 0 {
		xx, yy := world.ScreenToIso(x, y)
		tileX, tileY := world.IsoToCartesian(xx, yy)
		if tileX >= 0 && tileY >= 0 && tileX < 256 && tileY < 256 {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				world.BuildStructure(world.World.HoverStructure, false, int(tileX), int(tileY))
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
