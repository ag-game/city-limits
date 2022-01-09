package component

import (
	"time"

	. "code.rocketnine.space/tslocum/citylimits/ecs"

	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteComponent struct {
	Image          *ebiten.Image
	HorizontalFlip bool
	VerticalFlip   bool
	DiagonalFlip   bool // TODO unimplemented

	Angle float64

	Overlay            *ebiten.Image
	OverlayX, OverlayY float64 // Overlay offset

	Frame     int
	Frames    []*ebiten.Image
	FrameTime time.Duration
	LastFrame time.Time
	NumFrames int

	DamageTicks int

	OverrideColorScale bool
	ColorScale         float64
}

var SpriteComponentID = ECS.NewComponentID()

func (p *SpriteComponent) ComponentID() gohan.ComponentID {
	return SpriteComponentID
}

func Sprite(ctx *gohan.Context) *SpriteComponent {
	c, ok := ctx.Component(SpriteComponentID).(*SpriteComponent)
	if !ok {
		return nil
	}
	return c
}
