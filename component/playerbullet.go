package component

import (
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type PlayerBulletComponent struct {
}

var PlayerBulletComponentID = ECS.NewComponentID()

func (p *PlayerBulletComponent) ComponentID() gohan.ComponentID {
	return PlayerBulletComponentID
}

func PlayerBullet(ctx *gohan.Context) *PlayerBulletComponent {
	c, ok := ctx.Component(PlayerBulletComponentID).(*PlayerBulletComponent)
	if !ok {
		return nil
	}
	return c
}
