package entity

import (
	"code.rocketnine.space/tslocum/citylimits/asset"
	"code.rocketnine.space/tslocum/citylimits/component"
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayerBullet(x, y, xSpeed, ySpeed float64) gohan.Entity {
	bullet := ECS.NewEntity()

	ECS.AddComponent(bullet, &component.PositionComponent{
		X: x,
		Y: y,
	})

	ECS.AddComponent(bullet, &component.VelocityComponent{
		X: xSpeed,
		Y: ySpeed,
	})

	ECS.AddComponent(bullet, &component.SpriteComponent{
		Image: asset.ImgBlackSquare,
	})

	ECS.AddComponent(bullet, &component.PlayerBulletComponent{})

	return bullet
}
