package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const CameraMoveSpeed = 0.132

type CameraSystem struct {
}

func NewCameraSystem() *CameraSystem {
	s := &CameraSystem{}

	return s
}
func (_ *CameraSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.WeaponComponentID,
		component.PositionComponentID,
	}
}

func (_ *CameraSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *CameraSystem) Update(ctx *gohan.Context) error {
	if !world.World.GameStarted || world.World.GameOver {
		return nil
	}

	world.World.CamMoving = true
	// TODO
	return nil
}

func (_ *CameraSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
