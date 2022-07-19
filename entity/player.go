package entity

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayer() gohan.Entity {
	player := gohan.NewEntity()

	player.AddComponent(&component.Position{})

	player.AddComponent(&component.Velocity{})

	player.AddComponent(&component.Weapon{
		Damage:      1,
		FireRate:    144 / 16,
		BulletSpeed: 8,
	})

	return player
}
