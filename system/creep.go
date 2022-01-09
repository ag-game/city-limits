package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/citylimits/entity"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

// pause time, screen X, screen Y
type CreepSystem struct {
}

func NewCreepSystem() *CreepSystem {
	s := &CreepSystem{}

	return s
}
func (_ *CreepSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.ActorComponentID,
		component.PositionComponentID,
	}
}

func (_ *CreepSystem) Uses() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.WeaponComponentID,
	}
}

func (s *CreepSystem) Update(ctx *gohan.Context) error {
	if !world.World.GameStarted {
		return nil
	}

	creep := component.Creep(ctx)
	position := component.Position(ctx)

	if creep.Health <= 0 {
		for i, e := range world.World.CreepEntities {
			if e == ctx.Entity {
				world.World.CreepRects = append(world.World.CreepRects[:i], world.World.CreepRects[i+1:]...)
				world.World.CreepEntities = append(world.World.CreepEntities[:i], world.World.CreepEntities[i+1:]...)
				ctx.RemoveEntity()
				return nil
			}
		}
	}

	// Skip inactive creeps.
	sx, sy := world.LevelCoordinatesToScreen(position.X, position.Y)
	inactive := sx < 0 || sy < 0 || sx > 640 || sy > 480
	if creep.Active != !inactive {
		creep.Active = !inactive
	}
	if inactive {
		return nil
	}

	l := len(creep.Movements)
	if l > creep.Movement {
		if creep.MovementTicks == 0 {
			m := creep.Movements[creep.Movement]
			position.X, position.Y = m[0], m[1]
			creep.Movement++

			creep.MovementTicks = int(m[2])
		}
		creep.MovementTicks--
	}

	randVelocity := func() (float64, float64) {
		for {
			vx := creep.Rand.Float64()*0.5 + (0.5 - creep.Rand.Float64())
			vy := creep.Rand.Float64()*0.5 + (0.5 - creep.Rand.Float64())
			if vx > 0.5 || vx < -0.5 || vy > 0.5 || vy < -0.5 {
				return vx, vy
			}
		}
	}

	if creep.FireTicks == 0 {
		for i := 0; i < creep.FireAmount; i++ {
			vx, vy := randVelocity()

			if creep.Rand.Intn(2) == 0 {
				vx *= -1
			}
			if creep.Rand.Intn(2) == 0 {
				vy *= -1
			}
			entity.NewCreepBullet(position.X, position.Y, vx, vy)
		}
		creep.FireTicks = creep.FireRate
	}
	creep.FireTicks--

	if creep.DamageTicks > 0 {
		creep.DamageTicks--

		sprite := ECS.Component(ctx.Entity, component.SpriteComponentID)
		if sprite != nil {
			sp := sprite.(*component.SpriteComponent)
			if creep.DamageTicks > 0 {
				if creep.DamageTicks%2 == 0 {
					sp.ColorScale = 100
				} else {
					sp.ColorScale = .01
				}
				sp.OverrideColorScale = true
			} else {
				sp.OverrideColorScale = false
			}
		}
	}

	return nil
}

func (_ *CreepSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
