package system

import (
	"image/color"
	"strings"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RenderMessageSystem struct {
	op      *ebiten.DrawImageOptions
	logoImg *ebiten.Image
	msgImg  *ebiten.Image
	tmpImg  *ebiten.Image
}

func NewRenderMessageSystem() *RenderMessageSystem {
	s := &RenderMessageSystem{
		op:      &ebiten.DrawImageOptions{},
		logoImg: ebiten.NewImage(1, 1),
		msgImg:  ebiten.NewImage(1, 1),
		tmpImg:  ebiten.NewImage(200, 200),
	}

	return s
}

func (s *RenderMessageSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *RenderMessageSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *RenderMessageSystem) Update(_ *gohan.Context) error {
	if !world.World.GameStarted || world.World.GameOver || !world.World.MessageVisible {
		return nil
	}

	world.World.MessageTicks++
	if world.World.MessageTicks == world.World.MessageDuration {
		world.World.MessageVisible = false
		return nil
	}
	return nil
}

func (s *RenderMessageSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	if !world.World.GameStarted || !world.World.MessageVisible {
		return nil
	}

	// Draw message.
	if world.World.MessageUpdated {
		s.drawMessage()
	}
	bounds := s.msgImg.Bounds()
	x := (float64(world.World.ScreenW) / 2) - (float64(bounds.Dx()) / 2)
	y := (float64(world.World.ScreenH) / 2) - (float64(bounds.Dy()) / 2)
	s.op.GeoM.Reset()
	s.op.GeoM.Translate(x, y)
	screen.DrawImage(s.msgImg, s.op)
	return nil
}

func (s *RenderMessageSystem) drawMessage() {
	split := strings.Split(world.World.MessageText, "\n")
	width := 0
	for _, line := range split {
		lineSize := len(line) * 12
		if lineSize > width {
			width = lineSize
		}
	}
	height := len(split) * 32

	const padding = 8
	width, height = width+padding*2, height+padding*2

	s.msgImg = ebiten.NewImage(width, height)
	s.msgImg.Fill(color.RGBA{17, 17, 17, 255})

	s.tmpImg.Clear()
	s.tmpImg = ebiten.NewImage(width*2, height*2)
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(2, 2)
	s.op.GeoM.Translate(float64(padding), float64(padding))
	ebitenutil.DebugPrint(s.tmpImg, world.World.MessageText)
	s.msgImg.DrawImage(s.tmpImg, s.op)
	s.op.ColorM.Reset()
}
