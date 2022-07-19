package system

import (
	"image"
	"image/color"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const rewindThreshold = 1

type MovementSystem struct {
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Velocity `gohan:"?"`
}

func NewMovementSystem() *MovementSystem {
	s := &MovementSystem{}
	return s
}

func drawDebugRect(r image.Rectangle, c color.Color, overrideColorScale bool) gohan.Entity {
	rectEntity := gohan.NewEntity()

	rectImg := ebiten.NewImage(r.Dx(), r.Dy())
	rectImg.Fill(c)

	rectEntity.AddComponent(&component.Position{
		X: float64(r.Min.X),
		Y: float64(r.Min.Y),
	})

	rectEntity.AddComponent(&component.Sprite{
		Image:              rectImg,
		OverrideColorScale: overrideColorScale,
	})

	return rectEntity
}

func (s *MovementSystem) Update(e gohan.Entity) error {
	if !world.World.GameStarted {
		return nil
	}

	if world.World.GameOver && e == world.World.Player {
		return nil
	}

	position := s.Position
	velocity := s.Velocity

	vx, vy := velocity.X, velocity.Y
	if e == world.World.Player && (world.World.NoClip || world.World.Debug != 0) && ebiten.IsKeyPressed(ebiten.KeyShift) {
		vx, vy = vx*2, vy*2
	}

	position.X, position.Y = position.X+vx, position.Y+vy

	// Force player to remain within the screen bounds.
	// TODO same for bullets
	if e == world.World.Player {
		screenX, screenY := s.levelCoordinatesToScreen(position.X, position.Y)
		if screenX < 0 {
			diff := screenX / world.World.CamScale
			position.X -= diff
		} else if screenX > float64(world.World.ScreenW)-world.World.PlayerWidth {
			diff := (float64(world.World.ScreenW) - world.World.PlayerWidth - screenX) / world.World.CamScale
			position.X += diff
		}
		if screenY < 0 {
			diff := screenY / world.World.CamScale
			position.Y -= diff
		} else if screenY > float64(world.World.ScreenH)-world.World.PlayerHeight {
			diff := (float64(world.World.ScreenH) - world.World.PlayerHeight - screenY) / world.World.CamScale
			position.Y += diff
		}

		world.World.PlayerX, world.World.PlayerY = position.X, position.Y

		// Check player hazard collision.
		if world.World.NoClip {
			return nil
		}
		playerRect := image.Rect(int(position.X), int(position.Y), int(position.X+world.World.PlayerWidth), int(position.Y+world.World.PlayerHeight))
		for _, r := range world.World.HazardRects {
			if playerRect.Overlaps(r) {
				world.World.SetGameOver(0, 0)
				return nil
			}
		}
	}
	return nil
}

func (s *MovementSystem) levelCoordinatesToScreen(x, y float64) (float64, float64) {
	return (x - world.World.CamX) * world.World.CamScale, (y - world.World.CamY) * world.World.CamScale
}

func (_ *MovementSystem) Draw(_ gohan.Entity, screen *ebiten.Image) error {
	return gohan.ErrUnregister
}
