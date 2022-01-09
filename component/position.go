package component

import (
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type PositionComponent struct {
	X, Y float64
}

var PositionComponentID = ECS.NewComponentID()

func (p *PositionComponent) ComponentID() gohan.ComponentID {
	return PositionComponentID
}

func Position(ctx *gohan.Context) *PositionComponent {
	c, ok := ctx.Component(PositionComponentID).(*PositionComponent)
	if !ok {
		return nil
	}
	return c
}
