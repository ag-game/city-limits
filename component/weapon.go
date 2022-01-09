package component

import (
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type WeaponComponent struct {
	Equipped bool

	Damage int

	// In ticks
	FireRate int
	NextFire int

	BulletSpeed float64
}

var WeaponComponentID = ECS.NewComponentID()

func (p *WeaponComponent) ComponentID() gohan.ComponentID {
	return WeaponComponentID
}

func Weapon(ctx *gohan.Context) *WeaponComponent {
	c, ok := ctx.Component(WeaponComponentID).(*WeaponComponent)
	if !ok {
		return nil
	}
	return c
}
