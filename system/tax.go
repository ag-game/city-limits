package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type TaxSystem struct {
}

func NewTaxSystem() *TaxSystem {
	s := &TaxSystem{}

	return s
}

func (s *TaxSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *TaxSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *TaxSystem) Update(_ *gohan.Context) error {
	if world.World.Paused {
		return nil
	}

	if world.World.Ticks%world.YearTicks != 0 {
		return nil
	}

	for _, zone := range world.World.Zones {
		if !zone.Powered {
			continue
		}
		_ = zone
	}
	return nil
}

func (s *TaxSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
