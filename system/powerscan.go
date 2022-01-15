package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type PowerScanSystem struct {
}

func NewPowerScanSystem() *PowerScanSystem {
	s := &PowerScanSystem{}

	return s
}

func (s *PowerScanSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *PowerScanSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *PowerScanSystem) Update(_ *gohan.Context) error {
	if world.World.Paused {
		return nil
	}

	// TODO use a consistent procedure to check each building that needs power
	// as connected via road to a power plant, and power-out buildings without enough power
	// "citizens report brown-outs"
	return nil
}

func (s *PowerScanSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
