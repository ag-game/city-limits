package system

import (
	"image"
	"image/color"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type RenderHudSystem struct {
	op     *ebiten.DrawImageOptions
	hudImg *ebiten.Image
	tmpImg *ebiten.Image
}

func NewRenderHudSystem() *RenderHudSystem {
	s := &RenderHudSystem{
		op:     &ebiten.DrawImageOptions{},
		hudImg: ebiten.NewImage(1, 1),
		tmpImg: ebiten.NewImage(1, 1),
	}

	return s
}

func (s *RenderHudSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *RenderHudSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *RenderHudSystem) Update(_ *gohan.Context) error {
	return nil
}

func (s *RenderHudSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	// Draw HUD.
	if world.World.HUDUpdated {
		s.drawHUD()
		world.World.HUDUpdated = false
	}
	screen.DrawImage(s.hudImg, nil)
	return nil
}

func (s *RenderHudSystem) drawHUD() {
	bounds := s.hudImg.Bounds()
	if bounds.Dx() != world.World.ScreenW || bounds.Dy() != world.World.ScreenH {
		s.hudImg = ebiten.NewImage(world.World.ScreenW, world.World.ScreenH)
		s.tmpImg = ebiten.NewImage(world.SidebarWidth, world.World.ScreenH)
	} else {
		s.hudImg.Clear()
		s.tmpImg.Clear()
	}
	w := world.SidebarWidth
	if bounds.Dx() < w {
		w = bounds.Dx()
	}

	sidebarShade := uint8(108)
	sidebarColor := color.RGBA{sidebarShade, sidebarShade, sidebarShade, 255}
	s.tmpImg.Fill(sidebarColor)

	borderSize := 1
	columns := 3

	buttonShade := uint8(142)
	colorButton := color.RGBA{buttonShade, buttonShade, buttonShade, 255}

	swatchWidth := world.SidebarWidth / columns
	swatchHeight := swatchWidth
	world.World.HUDButtonRects = make([]image.Rectangle, len(world.HUDButtons))
	for i, button := range world.HUDButtons {
		row := i / columns
		x, y := (i%columns)*swatchWidth, row*swatchHeight
		world.World.HUDButtonRects[i] = image.Rect(x+borderSize, y+borderSize, x+swatchWidth-(borderSize*2), y+swatchHeight-(borderSize*2))

		s.tmpImg.SubImage(world.World.HUDButtonRects[i]).(*ebiten.Image).Fill(colorButton)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+borderSize)+button.SpriteOffsetX, float64(y+borderSize)+button.SpriteOffsetY)
		s.tmpImg.DrawImage(button.Sprite, op)
	}

	s.hudImg.DrawImage(s.tmpImg, nil)
}
