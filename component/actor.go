package component

import (
	"math/rand"

	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type ActorComponent struct {
	Type   int
	Active bool

	Health     int
	FireAmount int
	FireRate   int // In ticks
	FireTicks  int //Ticks until next action

	Movement      int
	Movements     [][3]float64 // X, Y, pre-delay in ticks
	MovementTicks int          // Ticks until next action

	DamageTicks int

	Rand *rand.Rand
}

const (
	CreepSnowblower = iota + 1
	CreepSmallRock
	CreepMediumRock
	CreepLargeRock
)

var ActorComponentID = ECS.NewComponentID()

func (p *ActorComponent) ComponentID() gohan.ComponentID {
	return ActorComponentID
}

func Creep(ctx *gohan.Context) *ActorComponent {
	c, ok := ctx.Component(ActorComponentID).(*ActorComponent)
	if !ok {
		return nil
	}
	return c
}
