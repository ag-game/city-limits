package system

import (
	"image"
	"image/color"

	"code.rocketnine.space/tslocum/citylimits/component"
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const rewindThreshold = 1

type MovementSystem struct {
	ScreenW, ScreenH float64
}

func NewMovementSystem() *MovementSystem {
	s := &MovementSystem{
		ScreenW: 640,
		ScreenH: 480,
	}

	return s
}

func drawDebugRect(r image.Rectangle, c color.Color, overrideColorScale bool) gohan.Entity {
	rectEntity := ECS.NewEntity()

	rectImg := ebiten.NewImage(r.Dx(), r.Dy())
	rectImg.Fill(c)

	ECS.AddComponent(rectEntity, &component.PositionComponent{
		X: float64(r.Min.X),
		Y: float64(r.Min.Y),
	})

	ECS.AddComponent(rectEntity, &component.SpriteComponent{
		Image:              rectImg,
		OverrideColorScale: overrideColorScale,
	})

	return rectEntity
}

func (_ *MovementSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
	}
}

func (_ *MovementSystem) Uses() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.WeaponComponentID,
	}
}

func (s *MovementSystem) Update(ctx *gohan.Context) error {
	if !world.World.GameStarted {
		return nil
	}

	if world.World.GameOver && ctx.Entity == world.World.Player {
		return nil
	}

	position := component.Position(ctx)
	velocity := component.Velocity(ctx)

	vx, vy := velocity.X, velocity.Y
	if ctx.Entity == world.World.Player && (world.World.NoClip || world.World.Debug != 0) && ebiten.IsKeyPressed(ebiten.KeyShift) {
		vx, vy = vx*2, vy*2
	}

	position.X, position.Y = position.X+vx, position.Y+vy

	// Force player to remain within the screen bounds.
	// TODO same for bullets
	if ctx.Entity == world.World.Player {
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
	} else if ctx.Entity == world.World.BrokenPieceA || ctx.Entity == world.World.BrokenPieceB {
		sprite := ECS.Component(ctx.Entity, component.SpriteComponentID).(*component.SpriteComponent)
		if ctx.Entity == world.World.BrokenPieceA {
			sprite.Angle -= 0.05
		} else {
			sprite.Angle += 0.05
		}
	}

	// Check creepBullet collision.
	if world.World.NoClip {
		return nil
	}
	bulletSize := 8.0
	bulletRect := image.Rect(int(position.X), int(position.Y), int(position.X+bulletSize), int(position.Y+bulletSize))

	creepBullet := ECS.Component(ctx.Entity, component.CreepBulletComponentID)
	playerBullet := ECS.Component(ctx.Entity, component.PlayerBulletComponentID)

	// Check hazard collisions.
	if creepBullet != nil || playerBullet != nil {
		var invulnerable bool
		if creepBullet != nil {
			b := creepBullet.(*component.CreepBulletComponent)
			invulnerable = b.Invulnerable
		}
		if !invulnerable {
			for _, hazardRect := range world.World.HazardRects {
				if bulletRect.Overlaps(hazardRect) {
					ctx.RemoveEntity()
					return nil
				}
			}
		}
	}

	if creepBullet != nil {
		playerRect := image.Rect(int(world.World.PlayerX), int(world.World.PlayerY), int(world.World.PlayerX+world.World.PlayerWidth), int(world.World.PlayerY+world.World.PlayerHeight))

		if bulletRect.Overlaps(playerRect) {
			world.World.SetGameOver(velocity.X, velocity.Y)
			return nil
		}
		return nil
	}

	if playerBullet != nil {
		for i, creepRect := range world.World.CreepRects {
			if bulletRect.Overlaps(creepRect) {
				creep := ECS.Component(world.World.CreepEntities[i], component.ActorComponentID).(*component.ActorComponent)
				if creep.Active {
					creep.Health--
					creep.DamageTicks = 6
					ctx.RemoveEntity()
					return nil
				}
			}
		}
	}

	return nil
}

func (s *MovementSystem) levelCoordinatesToScreen(x, y float64) (float64, float64) {
	return (x - world.World.CamX) * world.World.CamScale, (y - world.World.CamY) * world.World.CamScale
}

func (_ *MovementSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
