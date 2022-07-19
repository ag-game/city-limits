package component

type Weapon struct {
	Equipped bool

	Damage int

	// In ticks
	FireRate int
	NextFire int

	BulletSpeed float64
}
