package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type TaxSystem struct {
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Weapon
}

func NewTaxSystem() *TaxSystem {
	s := &TaxSystem{}

	return s
}

func (s *TaxSystem) Update(_ gohan.Entity) error {
	if world.World.Paused {
		return nil
	}

	if world.World.Ticks%world.YearTicks != 0 {
		return nil
	}

	taxCollectionAmount := 27.77
	for _, zone := range world.World.Zones {
		if zone.Population == 0 {
			continue
		}

		taxRate := world.World.TaxR
		if zone.Type == world.StructureCommercialZone {
			taxRate = world.World.TaxC
		} else if zone.Type == world.StructureIndustrialZone {
			taxRate = world.World.TaxI
		}

		world.World.Funds += int(taxCollectionAmount * taxRate * float64(zone.Population))
	}
	return nil
}

func (s *TaxSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}
