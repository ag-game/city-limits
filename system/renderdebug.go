package system

import (
	"fmt"
	"image/color"
	_ "image/png"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RenderDebugTextSystem struct {
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Weapon

	player   gohan.Entity
	op       *ebiten.DrawImageOptions
	debugImg *ebiten.Image
}

func NewRenderDebugTextSystem(player gohan.Entity) *RenderDebugTextSystem {
	s := &RenderDebugTextSystem{
		player:   player,
		op:       &ebiten.DrawImageOptions{},
		debugImg: ebiten.NewImage(70, 98),
	}

	return s
}

func (s *RenderDebugTextSystem) Update(_ gohan.Entity) error {
	return gohan.ErrUnregister
}

func (s *RenderDebugTextSystem) Draw(e gohan.Entity, screen *ebiten.Image) error {
	if world.World.Debug <= 0 {
		return nil
	}

	s.debugImg.Fill(color.RGBA{0, 0, 0, 80})
	ebitenutil.DebugPrintAt(s.debugImg, fmt.Sprintf("ENV %d\nENT %d\nUPD %d\nDRA %d\nTPS %0.0f\nFPS %0.0f", world.World.EnvironmentSprites, gohan.CurrentEntities(), gohan.CurrentUpdates(), gohan.CurrentDraws(), ebiten.CurrentTPS(), ebiten.CurrentFPS()), 2, 0)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(world.SidebarWidth, 0)
	screen.DrawImage(s.debugImg, op)
	return nil
}
