package system

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	helpW = 480
	helpH = 220
)

type RenderHudSystem struct {
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Weapon

	op           *ebiten.DrawImageOptions
	hudImg       *ebiten.Image
	tmpImg       *ebiten.Image
	tmpImg2      *ebiten.Image
	helpImg      *ebiten.Image
	sidebarColor color.RGBA
}

func NewRenderHudSystem() *RenderHudSystem {
	s := &RenderHudSystem{
		op:      &ebiten.DrawImageOptions{},
		hudImg:  ebiten.NewImage(1, 1),
		tmpImg:  ebiten.NewImage(1, 1),
		tmpImg2: ebiten.NewImage(1, 1),
		helpImg: ebiten.NewImage(helpW, helpH),
	}

	sidebarShade := uint8(108)
	s.sidebarColor = color.RGBA{sidebarShade, sidebarShade, sidebarShade, 255}

	return s
}

func (s *RenderHudSystem) Update(_ gohan.Entity) error {
	return nil
}

func (s *RenderHudSystem) Draw(_ gohan.Entity, screen *ebiten.Image) error {
	// Draw HUD.
	if world.World.HUDUpdated {
		s.hudImg.Clear()
		s.drawSidebar()
		s.drawMessages()
		s.drawTooltip()
		s.drawRCIWindow()
		s.drawHelp()
		world.World.HUDUpdated = false
	}
	screen.DrawImage(s.hudImg, nil)
	return nil
}

const columns = 3
const buttonWidth = world.SidebarWidth / columns

func (s *RenderHudSystem) drawSidebar() {
	bounds := s.hudImg.Bounds()
	if bounds.Dx() != world.World.ScreenW || bounds.Dy() != world.World.ScreenH {
		s.hudImg = ebiten.NewImage(world.World.ScreenW, world.World.ScreenH)
		s.tmpImg = ebiten.NewImage(world.World.ScreenW, world.World.ScreenH)
		s.tmpImg2 = ebiten.NewImage(world.SidebarWidth, world.World.ScreenH)
	} else {
		s.hudImg.Clear()
		s.tmpImg.Clear()
		s.tmpImg2.Clear()
	}

	// Fill background.
	s.hudImg.SubImage(image.Rect(0, 0, world.SidebarWidth, world.World.ScreenH)).(*ebiten.Image).Fill(s.sidebarColor)

	// Draw buttons.

	const paddingSize = 1
	const buttonHeight = buttonWidth
	world.World.HUDButtonRects = make([]image.Rectangle, len(world.HUDButtons))
	var lastButtonY int
	for i, button := range world.HUDButtons {
		row := i / columns
		x, y := (i%columns)*buttonWidth, row*buttonHeight
		r := image.Rect(x+paddingSize, y+paddingSize, x+buttonWidth-paddingSize, y+buttonHeight-paddingSize)

		if button != nil {
			selected := world.World.HoverStructure == button.StructureType
			if button.StructureType == world.StructureToggleHelp {
				selected = world.World.HelpPage != -1
			} else if button.StructureType == world.StructureToggleTransparentStructures {
				selected = world.World.TransparentStructures
			}

			// Draw background.
			s.drawButtonBackground(s.tmpImg, r, selected)

			// Draw sprite.
			colorScale := 1.0
			if selected {
				colorScale = 0.9
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x+paddingSize)+button.SpriteOffsetX, float64(y+paddingSize)+button.SpriteOffsetY)
			op.ColorM.Scale(colorScale, colorScale, colorScale, 1)
			s.tmpImg.SubImage(image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)).(*ebiten.Image).DrawImage(button.Sprite, op)

			s.drawButtonBorder(s.tmpImg, r, selected)
		}

		world.World.HUDButtonRects[i] = r
		if button != nil {
			nonHUDButton := button.StructureType == world.StructureToggleHelp || button.StructureType == world.StructureToggleTransparentStructures
			if !nonHUDButton {
				lastButtonY = y
			}
		}
	}

	dateY := lastButtonY + buttonHeight*2 - buttonHeight/2 - 16
	s.drawDate(dateY)
	s.drawFunds(dateY + 50)

	indicatorY := dateY + 179
	// Draw RCI indicator.
	s.drawDemand(buttonWidth/2, indicatorY)

	// Draw PWR indicator.
	s.drawPower(buttonWidth/2+buttonWidth, indicatorY)

	s.drawPopulation(world.World.ScreenH - 45)

	s.hudImg.DrawImage(s.tmpImg, nil)

	s.hudImg.SubImage(image.Rect(world.SidebarWidth-1, 0, world.SidebarWidth, world.World.ScreenH)).(*ebiten.Image).Fill(color.Black)
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

func (s *RenderHudSystem) drawTooltip() {
	label := world.Tooltip()
	if label == "" {
		return
	}

	lines := 1 + strings.Count(label, "\n")
	max := maxLen(strings.Split(label, "\n"))

	scale := 3.0
	x, y := world.SidebarWidth, 0
	w, h := (max*6+10)*int(scale), 16*(int(scale))*lines+10
	r := image.Rect(x, y, x+w, y+h)
	s.hudImg.SubImage(r).(*ebiten.Image).Fill(color.RGBA{0, 0, 0, 120})

	s.tmpImg.Clear()
	ebitenutil.DebugPrint(s.tmpImg, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(world.SidebarWidth+(4*scale), 2)
	s.hudImg.DrawImage(s.tmpImg, op)
}

func maxLen(v []string) int {
	max := 0
	for _, line := range v {
		l := len(line)
		if l > max {
			max = l
		}
	}
	return max
}

func (s *RenderHudSystem) drawDemand(x, y int) {
	const rciSize = 100
	rciX := x
	rciY := y

	const rciButtonHeight = 20

	colorR := color.RGBA{0, 255, 0, 255}
	colorC := color.RGBA{0, 0, 255, 255}
	colorI := color.RGBA{231, 231, 72, 255}
	demandR, demandC, demandI := world.Demand()
	drawDemandBar := func(demand float64, clr color.RGBA, i int) {
		barOffsetSize := 12
		barOffset := -barOffsetSize + (i * barOffsetSize)
		barWidth := 7
		barX := rciX + buttonWidth/2 - barWidth/2 + barOffset
		barY := rciY + (rciSize / 2)
		if demand < 0 {
			barY += rciButtonHeight / 2
		} else {
			barY -= rciButtonHeight / 2
		}
		barHeight := int((float64(rciSize) / 2) * demand)
		s.tmpImg.SubImage(image.Rect(barX, barY, barX+barWidth, barY-barHeight)).(*ebiten.Image).Fill(clr)
	}
	drawDemandBar(demandR, colorR, 0)
	drawDemandBar(demandC, colorC, 1)
	drawDemandBar(demandI, colorI, 2)

	// Draw button.
	const rciButtonPadding = 12
	const rciButtonLabelPaddingX = 6
	const rciButtonLabelPaddingY = 1
	rciButtonY := rciY + (rciSize / 2) - (rciButtonHeight / 2)
	rciButtonRect := image.Rect(rciX+rciButtonPadding, rciButtonY, rciX+buttonWidth-rciButtonPadding, rciButtonY+rciButtonHeight)

	s.drawButtonBackground(s.tmpImg, rciButtonRect, world.World.ShowRCIWindow)

	// Draw label.
	ebitenutil.DebugPrintAt(s.tmpImg, "R C I", rciX+rciButtonPadding+rciButtonLabelPaddingX, rciButtonY+rciButtonLabelPaddingY)

	s.drawButtonBorder(s.tmpImg, rciButtonRect, world.World.ShowRCIWindow)

	world.World.RCIButtonRect = rciButtonRect
}

func (s *RenderHudSystem) drawPower(x, y int) {
	const rciSize = 100
	rciX := x
	rciY := y

	const rciButtonHeight = 20

	colorPowerNormal := color.RGBA{0, 255, 0, 255}
	colorPowerOut := color.RGBA{255, 0, 0, 255}
	colorPowerCapacity := color.RGBA{16, 16, 16, 255}
	drawPowerBar := func(demand float64, clr color.RGBA, i int) {
		barOffsetSize := 7
		barOffset := -barOffsetSize + (i * barOffsetSize)
		barWidth := 7
		barX := rciX + buttonWidth/2 - barWidth/2 + barOffset + 4
		barY := rciY + (rciSize / 2)
		if demand < 0 {
			barY += rciButtonHeight / 2
		} else {
			barY -= rciButtonHeight / 2
		}
		barHeight := int((float64(rciSize) / 2) * demand)
		s.tmpImg.SubImage(image.Rect(barX, barY, barX+barWidth, barY-barHeight)).(*ebiten.Image).Fill(clr)
	}

	powerColor := colorPowerNormal
	if world.World.HavePowerOut || world.World.PowerNeeded > world.World.PowerAvailable {
		powerColor = colorPowerOut
	}

	max := world.World.PowerNeeded
	if world.World.PowerAvailable > max {
		max = world.World.PowerAvailable
	}

	pctUsage, pctCapacity := float64(world.World.PowerNeeded)/float64(max), float64(world.World.PowerAvailable)/float64(max)
	clamp := func(v float64) float64 {
		if math.IsNaN(v) {
			return 0
		}
		if v < -1 {
			v = -1
		} else if v > 1 {
			v = 1
		}
		return v
	}

	drawPowerBar(clamp(pctUsage), powerColor, 0)
	drawPowerBar(clamp(pctCapacity), colorPowerCapacity, 1)

	// Draw button.
	const rciButtonPadding = 12
	const rciButtonLabelPaddingX = 6
	const rciButtonLabelPaddingY = 1
	rciButtonY := rciY + (rciSize / 2) - (rciButtonHeight / 2)
	rciButtonRect := image.Rect(rciX+rciButtonPadding, rciButtonY, rciX+buttonWidth-rciButtonPadding, rciButtonY+rciButtonHeight)

	s.drawButtonBackground(s.tmpImg, rciButtonRect, false) // TODO

	// Draw label.
	ebitenutil.DebugPrintAt(s.tmpImg, "POWER", rciX+rciButtonPadding+rciButtonLabelPaddingX, rciButtonY+rciButtonLabelPaddingY)

	s.drawButtonBorder(s.tmpImg, rciButtonRect, false) // TODO
}

func (s *RenderHudSystem) drawMessages() {
	lines := len(world.World.Messages)
	if lines == 0 {
		return
	}
	/*var label string
	max := maxLen(world.World.Messages)
	for i := lines - 1; i >= 0; i-- {
		if i != lines-1 {
			label += "\n"
		}
		for j := max - len(world.World.Messages[i]); j > 0; j-- {
			label += " "
		}
		label += world.World.Messages[i]
	}*/

	label := world.World.Messages[len(world.World.Messages)-1]
	max := len(label)
	lines = 1

	const padding = 10

	scale := 2.0
	w, h := (max*6+10)*int(scale), 16*(int(scale))*lines+10
	x, y := world.World.ScreenW-w, 0
	r := image.Rect(x, y, x+w, y+h)
	s.hudImg.SubImage(r).(*ebiten.Image).Fill(color.RGBA{0, 0, 0, 120})

	s.tmpImg.Clear()
	ebitenutil.DebugPrint(s.tmpImg, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x)+padding, 3)
	s.hudImg.DrawImage(s.tmpImg, op)
}

func (s *RenderHudSystem) drawDate(y int) {
	const datePadding = 10
	month, year := world.Date()
	label := month

	scale := 2.0
	x, y := datePadding, y

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)

	label = year

	x = world.SidebarWidth - 1 - datePadding - (len(label) * 6 * int(scale))

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op.GeoM.Reset()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)
}

func (s *RenderHudSystem) drawFunds(y int) {
	label := world.World.Printer.Sprintf("$%d", world.World.Funds)

	scale := 2.0
	x, y := world.SidebarWidth/2-(len(label)*12)/2, y

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)
}

func (s *RenderHudSystem) drawPopulation(y int) {
	var population int
	for _, zone := range world.World.Zones {
		population += zone.Population
	}

	const datePadding = 10
	label := "Pop"

	scale := 2.0
	x, y := datePadding, y

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)

	label = world.World.Printer.Sprintf("%d", population)

	x = world.SidebarWidth - 1 - datePadding - (len(label) * 6 * int(scale))

	s.tmpImg2.Clear()
	ebitenutil.DebugPrint(s.tmpImg2, label)
	op.GeoM.Reset()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	s.hudImg.DrawImage(s.tmpImg2, op)
}

func (s *RenderHudSystem) drawHelp() {
	if world.World.HelpPage < 0 {
		return
	}

	if world.World.HelpUpdated {
		s.helpImg.Fill(s.sidebarColor)

		label := strings.TrimSpace(world.HelpText[world.World.HelpPage])

		s.tmpImg.Clear()
		ebitenutil.DebugPrint(s.tmpImg, label)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(5, 0)
		s.helpImg.DrawImage(s.tmpImg, op)

		s.helpImg.SubImage(image.Rect(0, 0, helpW, 1)).(*ebiten.Image).Fill(color.Black)
		s.helpImg.SubImage(image.Rect(0, 0, 1, helpH)).(*ebiten.Image).Fill(color.Black)

		// Draw prev/next buttons.
		buttonSize := 32
		buttonPadding := 4
		prevRect := image.Rect(buttonPadding+2, helpH-buttonSize-buttonPadding+1, buttonSize+buttonPadding+2, helpH-buttonPadding+1)
		closeRect := image.Rect(helpW/2-buttonSize/2, helpH-buttonSize-buttonPadding+1, helpW/2+buttonSize/2, helpH-buttonPadding+1)
		nextRect := image.Rect(helpW-buttonPadding, helpH-buttonSize-buttonPadding+1, helpW-buttonSize-buttonPadding, helpH-buttonPadding+1)

		drawButton := func(r image.Rectangle, l string) {
			s.drawButtonBackground(s.helpImg, r, false)
			ebitenutil.DebugPrintAt(s.helpImg, l, r.Min.X+buttonSize/2-4, r.Min.Y+buttonSize/2-10)
			s.drawButtonBorder(s.helpImg, r, false)
		}

		if world.World.HelpPage > 0 {
			drawButton(prevRect, "<")
		}
		drawButton(closeRect, "X")
		if world.World.HelpPage < len(world.HelpText)-1 {
			drawButton(nextRect, ">")
		}

		world.World.HelpButtonRects = []image.Rectangle{
			prevRect,
			closeRect,
			nextRect,
		}

		world.World.HelpUpdated = false
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(world.World.ScreenW)-helpW, float64(world.World.ScreenH)-helpH)
	s.hudImg.DrawImage(s.helpImg, op)
}

func (s *RenderHudSystem) drawRCIWindow() {
	if !world.World.ShowRCIWindow {
		return
	}

	const paddingX = 8

	const (
		rciWindowW = 425
		rciWindowH = 100
	)

	rciWindowRect := image.Rect(world.World.ScreenW/2-rciWindowW/2, world.World.ScreenH/2-rciWindowH/2, world.World.ScreenW/2+rciWindowW, world.World.ScreenH/2+rciWindowH)
	s.hudImg.SubImage(rciWindowRect).(*ebiten.Image).Fill(s.sidebarColor)

	percentBar := func(tax float64) string {
		if tax >= 1.0 {
			tax = .99
		}
		bar := "----------"
		bar = bar[:int(tax*10)] + "%" + bar[int(tax*10)+1:]
		return bar
	}

	label := fmt.Sprintf(`
Residential %3s%%  - |%s| +
Commercial  %3s%%  - |%s| +
Industrial  %3s%%  - |%s| +
`,
		strconv.Itoa(int(world.World.TaxR*100)), percentBar(world.World.TaxR),
		strconv.Itoa(int(world.World.TaxC*100)), percentBar(world.World.TaxC),
		strconv.Itoa(int(world.World.TaxI*100)), percentBar(world.World.TaxI))

	s.tmpImg.Clear()
	ebitenutil.DebugPrint(s.tmpImg, strings.TrimSpace(label))

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(float64(rciWindowRect.Min.X)+paddingX, float64(rciWindowRect.Min.Y))
	s.hudImg.DrawImage(s.tmpImg, op)

	s.hudImg.SubImage(image.Rect(rciWindowRect.Min.X, rciWindowRect.Min.Y, rciWindowRect.Max.X, rciWindowRect.Min.Y+1)).(*ebiten.Image).Fill(color.Black)
	s.hudImg.SubImage(image.Rect(rciWindowRect.Min.X, rciWindowRect.Max.Y-1, rciWindowRect.Max.X, rciWindowRect.Max.Y)).(*ebiten.Image).Fill(color.Black)
	s.hudImg.SubImage(image.Rect(rciWindowRect.Min.X, rciWindowRect.Min.Y, rciWindowRect.Min.X+1, rciWindowRect.Max.Y)).(*ebiten.Image).Fill(color.Black)
	s.hudImg.SubImage(image.Rect(rciWindowRect.Max.X-1, rciWindowRect.Min.Y, rciWindowRect.Max.X, rciWindowRect.Max.Y)).(*ebiten.Image).Fill(color.Black)

	world.World.RCIWindowRect = rciWindowRect
}
