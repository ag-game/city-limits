package component

import (
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type CreepBulletComponent struct {
	Invulnerable bool // Invulnerable to hazards
}

var CreepBulletComponentID = ECS.NewComponentID()

func (p *CreepBulletComponent) ComponentID() gohan.ComponentID {
	return CreepBulletComponentID
}

func CreepBullet(ctx *gohan.Context) *CreepBulletComponent {
	c, ok := ctx.Component(CreepBulletComponentID).(*CreepBulletComponent)
	if !ok {
		return nil
	}
	return c
}
