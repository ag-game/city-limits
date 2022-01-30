package system

import (
	"code.rocketnine.space/tslocum/citylimits/asset"
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type TickSystem struct {
}

func NewTickSystem() *TickSystem {
	s := &TickSystem{}

	return s
}

func (s *TickSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *TickSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *TickSystem) Update(_ *gohan.Context) error {
	if world.World.Paused {
		return nil
	}

	// Update date display.
	if world.World.Ticks%world.MonthTicks == 0 {
		world.World.HUDUpdated = true
	}
	if world.World.Ticks%144 == 0 {
		world.TickMessages()

		if !world.World.MuteMusic && !asset.SoundMusic1.IsPlaying() && !asset.SoundMusic2.IsPlaying() && !asset.SoundMusic3.IsPlaying() {
			world.PlayNextSong()
		}
	}
	world.World.Ticks++
	return nil
}

func (s *TickSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
