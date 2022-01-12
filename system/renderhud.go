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

	sidebarShade := uint8(111)
	sidebarColor := color.RGBA{sidebarShade, sidebarShade, sidebarShade, 255}
	s.tmpImg.Fill(sidebarColor)

	paddingSize := 1
	borderSize := 2
	columns := 3

	buttonShade := uint8(142)
	colorButton := color.RGBA{buttonShade, buttonShade, buttonShade, 255}

	lightBorderShade := uint8(216)
	colorLightBorder := color.RGBA{lightBorderShade, lightBorderShade, lightBorderShade, 255}

	mediumBorderShade := uint8(56)
	colorMediumBorder := color.RGBA{mediumBorderShade, mediumBorderShade, mediumBorderShade, 255}

	darkBorderShade := uint8(42)
	colorDarkBorder := color.RGBA{darkBorderShade, darkBorderShade, darkBorderShade, 255}

	swatchWidth := world.SidebarWidth / columns
	swatchHeight := swatchWidth
	world.World.HUDButtonRects = make([]image.Rectangle, len(world.HUDButtons))
	for i, button := range world.HUDButtons {
		row := i / columns
		x, y := (i%columns)*swatchWidth, row*swatchHeight
		r := image.Rect(x+paddingSize, y+paddingSize, x+swatchWidth-paddingSize, y+swatchHeight-paddingSize)

		bgColor := colorButton
		topLeftBorder := colorLightBorder
		bottomRightBorder := colorMediumBorder
		if world.World.HoverStructure == button.StructureType {
			bgColor = sidebarColor
			topLeftBorder = colorDarkBorder
			bottomRightBorder = colorLightBorder
		}

		// Draw background.
		s.tmpImg.SubImage(r).(*ebiten.Image).Fill(bgColor)

		// Draw sprite.
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+paddingSize)+button.SpriteOffsetX, float64(y+paddingSize)+button.SpriteOffsetY)
		s.tmpImg.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).DrawImage(button.Sprite, op)

		// Draw top and left border.
		s.tmpImg.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+borderSize)).(*ebiten.Image).Fill(topLeftBorder)
		s.tmpImg.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+borderSize, r.Max.Y)).(*ebiten.Image).Fill(topLeftBorder)

		// Draw bottom and right border.
		s.tmpImg.SubImage(image.Rect(r.Min.X, r.Max.Y-borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(bottomRightBorder)
		s.tmpImg.SubImage(image.Rect(r.Max.X-borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(bottomRightBorder)

		world.World.HUDButtonRects[i] = r
	}

	s.hudImg.DrawImage(s.tmpImg, nil)
}
