package world

import (
	"github.com/beefsack/go-astar"
)

type PowerMapTile struct {
	X            int
	Y            int
	CarriesPower bool // Set to true for roads and all building tiles (even power plants)
}

func (t *PowerMapTile) Up() *PowerMapTile {
	tx, ty := t.X, t.Y-1
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

func (t *PowerMapTile) Down() *PowerMapTile {
	tx, ty := t.X, t.Y+1
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

func (t *PowerMapTile) Left() *PowerMapTile {
	tx, ty := t.X-1, t.Y
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

func (t *PowerMapTile) Right() *PowerMapTile {
	tx, ty := t.X+1, t.Y
	if !ValidXY(tx, ty) {
		return nil
	}
	n := World.Power[tx][ty]
	if !n.CarriesPower {
		return nil
	}
	return n
}

type PowerMap [][]*PowerMapTile

func newPowerMap() PowerMap {
	m := make(PowerMap, 256)
	for x := 0; x < 256; x++ {
		m[x] = make([]*PowerMapTile, 256)
		for y := 0; y < 256; y++ {
			m[x][y] = &PowerMapTile{
				X: x,
				Y: y,
			}
		}
	}
	return m
}

func newPowerOuts() [][]bool {
	m := make([][]bool, 256)
	for x := 0; x < 256; x++ {
		m[x] = make([]bool, 256)
	}
	return m
}

func ResetPowerOuts() {
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			World.PowerOuts[x][y] = false
		}
	}
	World.HavePowerOut = false
}

func (m PowerMap) GetTile(x, y int) *PowerMapTile {
	if !ValidXY(x, y) {
		return nil
	}
	return m[x][y]
}

func (m PowerMap) SetTile(x, y int, carriesPower bool) {
	t := m[x][y]
	if t.CarriesPower == carriesPower {
		return
	}
	t.CarriesPower = carriesPower

	World.PowerUpdated = true
}

func (t *PowerMapTile) PathNeighbors() []astar.Pather {
	var neighbors []astar.Pather
	n := t.Up()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	n = t.Down()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	n = t.Left()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	n = t.Right()
	if n != nil {
		neighbors = append(neighbors, n)
	}
	return neighbors
}

func (t *PowerMapTile) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*PowerMapTile)
	if !toT.CarriesPower {
		return 0
	}
	return 1
}

func (t *PowerMapTile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*PowerMapTile)
	absX := toT.X - t.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}
