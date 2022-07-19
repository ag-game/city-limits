package system

import (
	_ "image/png"
	"time"

	"golang.org/x/image/colornames"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileWidth = 64

	logoText      = "POWERED BY EBITEN"
	logoTextScale = 4.75
	logoTextWidth = 6.0 * float64(len(logoText)) * logoTextScale
	logoTime      = 144 * 3.5

	fadeInTime = 144 * 0.75
)

type RenderSystem struct {
	Position *component.Position
	Sprite   *component.Sprite

	img *ebiten.Image
	op  *ebiten.DrawImageOptions

	camScale float64

	renderer gohan.Entity
}

func NewRenderSystem() *RenderSystem {
	s := &RenderSystem{
		renderer: gohan.NewEntity(),
		img:      ebiten.NewImage(320, 100),
		op:       &ebiten.DrawImageOptions{},
		camScale: 1,
	}

	return s
}

func (s *RenderSystem) Update(_ gohan.Entity) error {
	return gohan.ErrUnregister
}

func (s *RenderSystem) levelCoordinatesToScreen(x, y float64) (float64, float64) {
	px, py := world.World.CamX, world.World.CamY
	py *= -1
	return ((x - px) * s.camScale), ((y + py) * s.camScale)
}

// renderSprite renders a sprite on the screen.
func (s *RenderSystem) renderSprite(x float64, y float64, offsetx float64, offsety float64, angle float64, geoScale float64, colorScale float64, alpha float64, hFlip bool, vFlip bool, sprite *ebiten.Image, target *ebiten.Image) int {
	if alpha < .01 || colorScale < .01 {
		return 0
	}

	xi, yi := world.CartesianToIso(float64(x), float64(y))

	padding := float64(world.TileSize) * world.World.CamScale
	cx, cy := float64(world.World.ScreenW/2), float64(world.World.ScreenH/2)

	// Skip drawing tiles that are out of the screen.
	drawX, drawY := world.IsoToScreen(xi, yi)
	if drawX+padding < 0 || drawY+padding < 0 || drawX-padding > float64(world.World.ScreenW) || drawY-padding > float64(world.World.ScreenH) {
		return 0
	}

	s.op.GeoM.Reset()

	if hFlip {
		s.op.GeoM.Scale(-1, 1)
		s.op.GeoM.Translate(TileWidth, 0)
	}
	if vFlip {
		s.op.GeoM.Scale(1, -1)
		s.op.GeoM.Translate(0, TileWidth)
	}

	// Move to current isometric position.
	s.op.GeoM.Translate(xi, yi)
	// Translate camera position.
	s.op.GeoM.Translate(-world.World.CamX, -world.World.CamY)
	// Zoom.
	s.op.GeoM.Scale(world.World.CamScale, world.World.CamScale)
	// Center.
	s.op.GeoM.Translate(cx, cy)

	target.DrawImage(sprite, s.op)

	/*s.op.GeoM.Scale(geoScale, geoScale)
	// Rotate
	s.op.GeoM.Translate(offsetx, offsety)
	s.op.GeoM.Rotate(angle)
	// Move to current isometric position.
	s.op.GeoM.Translate(x, y)
	// Translate camera position.
	s.op.GeoM.Translate(-world.World.CamX, -world.World.CamY)
	// Zoom.
	s.op.GeoM.Scale(s.camScale, s.camScale)
	// Center.
	//s.op.GeoM.Translate(float64(s.ScreenW/2.0), float64(s.ScreenH/2.0))

	s.op.ColorM.Scale(colorScale, colorScale, colorScale, alpha)

	target.DrawImage(sprite, s.op)

	s.op.ColorM.Reset()*/

	return 1
}

func (s *RenderSystem) Draw(e gohan.Entity, screen *ebiten.Image) error {
	if !world.World.GameStarted {
		// TODO
		if e == world.World.Player {
			screen.Fill(colornames.Purple)
		}
		return nil
	}

	position := s.Position
	sprite := s.Sprite

	if sprite.NumFrames > 0 && time.Since(sprite.LastFrame) > sprite.FrameTime {
		sprite.Frame++
		if sprite.Frame >= sprite.NumFrames {
			sprite.Frame = 0
		}
		sprite.Image = sprite.Frames[sprite.Frame]
		sprite.LastFrame = time.Now()
	}

	colorScale := 1.0
	if sprite.OverrideColorScale {
		colorScale = sprite.ColorScale
	}

	s.renderSprite(position.X, position.Y, 0, 0, sprite.Angle, 1.0, colorScale, 1.0, sprite.HorizontalFlip, sprite.VerticalFlip, sprite.Image, screen)
	return nil
}
