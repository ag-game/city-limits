package component

import (
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type RailComponent struct {
}

var RailComponentID = ECS.NewComponentID()

func (p *RailComponent) ComponentID() gohan.ComponentID {
	return RailComponentID
}

func Rail(ctx *gohan.Context) *RailComponent {
	c, ok := ctx.Component(RailComponentID).(*RailComponent)
	if !ok {
		return nil
	}
	return c
}
