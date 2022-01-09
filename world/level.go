package world

import "github.com/hajimehoshi/ebiten/v2"

type Tile struct {
	Sprite *ebiten.Image
}

type GameLevel struct {
	Tiles [][][]*Tile

	size int
}

func NewLevel(size int) *GameLevel {
	l := &GameLevel{
		size: size,
	}
	const numLayers = 3
	for i := 0; i < numLayers; i++ {
		l.AddLayer()
	}
	return l
}

func (l *GameLevel) AddLayer() {
	tileMap := make([][]*Tile, l.size)
	for x := 0; x < l.size; x++ {
		tileMap[x] = make([]*Tile, l.size)
	}
	l.Tiles = append(l.Tiles, tileMap)
}
