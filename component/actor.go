package component

import (
	"math/rand"
)

type Actor struct {
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
