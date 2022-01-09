package entity

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayer() gohan.Entity {
	player := ECS.NewEntity()

	ECS.AddComponent(player, &component.PositionComponent{})

	ECS.AddComponent(player, &component.VelocityComponent{})

	weapon := &component.WeaponComponent{
		Damage:      1,
		FireRate:    144 / 16,
		BulletSpeed: 8,
	}
	ECS.AddComponent(player, weapon)

	return player
}
