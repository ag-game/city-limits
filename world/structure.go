package world

import "code.rocketnine.space/tslocum/gohan"

type Structure struct {
	Type int
	X, Y int

	Entity   gohan.Entity
	Children []gohan.Entity
}
