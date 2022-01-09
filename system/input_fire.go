package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/entity"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	fireSpeed = 1.5
)

type playerFireSystem struct {
}

func NewplayerFireSystem() *playerFireSystem {
	return &playerFireSystem{}
}

func (_ *playerFireSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
		component.SpriteComponentID,
	}
}

func (_ *playerFireSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *playerFireSystem) Update(ctx *gohan.Context) error {
	if !world.World.GameStarted || world.World.GameOver {
		return nil
	}

	position := component.Position(ctx)
	weapon := component.Weapon(ctx)
	if ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if weapon.NextFire == 0 {
			entity.NewPlayerBullet(position.X-8, position.Y-8, 0, -weapon.BulletSpeed)
			entity.NewPlayerBullet(position.X+8, position.Y-8, 0, -weapon.BulletSpeed)
			weapon.NextFire = weapon.FireRate
		}
	}
	if weapon.NextFire > 0 {
		weapon.NextFire--
	}
	return nil
}

func (s *playerFireSystem) Draw(_ *gohan.Context, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
