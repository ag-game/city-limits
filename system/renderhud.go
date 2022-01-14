package system

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type RenderHudSystem struct {
	op           *ebiten.DrawImageOptions
	hudImg       *ebiten.Image
	tmpImg       *ebiten.Image
	sidebarColor color.RGBA
}

func NewRenderHudSystem() *RenderHudSystem {
	s := &RenderHudSystem{
		op:     &ebiten.DrawImageOptions{},
		hudImg: ebiten.NewImage(1, 1),
		tmpImg: ebiten.NewImage(1, 1),
	}

	sidebarShade := uint8(111)
	s.sidebarColor = color.RGBA{sidebarShade, sidebarShade, sidebarShade, 255}

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

	// Fill background.
	s.tmpImg.Fill(s.sidebarColor)

	// Draw buttons.

	paddingSize := 1
	columns := 3

	buttonWidth := world.SidebarWidth / columns
	buttonHeight := buttonWidth
	world.World.HUDButtonRects = make([]image.Rectangle, len(world.HUDButtons))
	var lastButtonY int
	for i, button := range world.HUDButtons {
		row := i / columns
		x, y := (i%columns)*buttonWidth, row*buttonHeight
		r := image.Rect(x+paddingSize, y+paddingSize, x+buttonWidth-paddingSize, y+buttonHeight-paddingSize)

		selected := world.World.HoverStructure == button.StructureType

		// Draw background.
		s.drawButtonBackground(s.tmpImg, r, selected)

		// Draw sprite.
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+paddingSize)+button.SpriteOffsetX, float64(y+paddingSize)+button.SpriteOffsetY)
		s.tmpImg.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).DrawImage(button.Sprite, op)

		s.drawButtonBorder(s.tmpImg, r, selected)

		world.World.HUDButtonRects[i] = r
		lastButtonY = y
	}

	// Draw RCI indicator.
	rciPadding := buttonWidth / 2
	const rciSize = 100
	rciX := buttonWidth
	rciY := lastButtonY + buttonHeight + rciPadding

	// Draw RCI bars.
	colorR := color.RGBA{0, 255, 0, 255}
	colorC := color.RGBA{0, 0, 255, 255}
	colorI := color.RGBA{231, 231, 72, 255}
	demandR, demandC, demandI := world.Demand()
	drawDemandBar := func(demand float64, clr color.RGBA, i int) {
		barOffsetSize := 12
		barOffset := -barOffsetSize + (i * barOffsetSize)
		barWidth := 7
		barX := rciX + buttonWidth/2 - barWidth/2 + barOffset
		barHeight := int((float64(rciSize) / 2) * demand)
		s.tmpImg.SubImage(image.Rect(barX, rciY+(rciSize/2), barX+barWidth, rciY+(rciSize/2)-barHeight)).(*ebiten.Image).Fill(clr)
	}
	drawDemandBar(demandR, colorR, 0)
	drawDemandBar(demandC, colorC, 1)
	drawDemandBar(demandI, colorI, 2)

	// Draw RCI button.
	const rciButtonPadding = 12
	const rciButtonHeight = 20
	const rciButtonLabelPaddingX = 6
	const rciButtonLabelPaddingY = 1
	rciButtonY := rciY + (rciSize / 2) - (rciButtonHeight / 2)
	rciButtonRect := image.Rect(rciX+rciButtonPadding, rciButtonY, rciX+buttonWidth-rciButtonPadding, rciButtonY+rciButtonHeight)

	s.drawButtonBackground(s.tmpImg, rciButtonRect, false) // TODO

	// Draw RCI label.
	ebitenutil.DebugPrintAt(s.tmpImg, "R C I", rciX+rciButtonPadding+rciButtonLabelPaddingX, rciButtonY+rciButtonLabelPaddingY)

	s.drawButtonBorder(s.tmpImg, rciButtonRect, false) // TODO

	s.hudImg.DrawImage(s.tmpImg, nil)
}

func (s *RenderHudSystem) drawButtonBackground(img *ebiten.Image, r image.Rectangle, selected bool) {
	buttonShade := uint8(142)
	colorButton := color.RGBA{buttonShade, buttonShade, buttonShade, 255}

	bgColor := colorButton
	if selected {
		bgColor = s.sidebarColor
	}

	img.SubImage(r).(*ebiten.Image).Fill(bgColor)
}

func (s *RenderHudSystem) drawButtonBorder(img *ebiten.Image, r image.Rectangle, selected bool) {
	borderSize := 2

	lightBorderShade := uint8(216)
	colorLightBorder := color.RGBA{lightBorderShade, lightBorderShade, lightBorderShade, 255}

	mediumBorderShade := uint8(56)
	colorMediumBorder := color.RGBA{mediumBorderShade, mediumBorderShade, mediumBorderShade, 255}

	darkBorderShade := uint8(42)
	colorDarkBorder := color.RGBA{darkBorderShade, darkBorderShade, darkBorderShade, 255}

	topLeftBorder := colorLightBorder
	bottomRightBorder := colorMediumBorder
	if selected {
		topLeftBorder = colorDarkBorder
		bottomRightBorder = colorLightBorder
	}

	// Draw top and left border.
	img.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+borderSize)).(*ebiten.Image).Fill(topLeftBorder)
	img.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Min.X+borderSize, r.Max.Y)).(*ebiten.Image).Fill(topLeftBorder)

	// Draw bottom and right border.
	img.SubImage(image.Rect(r.Min.X, r.Max.Y-borderSize, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(bottomRightBorder)
	img.SubImage(image.Rect(r.Max.X-borderSize, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).Fill(bottomRightBorder)
}
