package system

import (
	"log"

	"github.com/beefsack/go-astar"

	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type PowerScanSystem struct {
}

func NewPowerScanSystem() *PowerScanSystem {
	s := &PowerScanSystem{}

	return s
}

func (s *PowerScanSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *PowerScanSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *PowerScanSystem) Update(_ *gohan.Context) error {
	if world.World.Paused {
		return nil
	}

	const scanTicks = 144 * 2
	if world.World.Ticks%scanTicks != 0 {
		return nil
	}

	if !world.World.PowerUpdated {
		return nil
	}

	var totalPowerAvailable int

	powerRemaining := make([]int, len(world.World.PowerPlants))
	for i, plant := range world.World.PowerPlants {
		powerRemaining[i] = world.PowerPlantCapacities[plant.Type]
		totalPowerAvailable += world.PowerPlantCapacities[plant.Type]
	}

	const (
		plantSize = 5
	)
	powerSourceTiles := make([][]*world.PowerMapTile, len(world.World.PowerPlants))
	log.Println("POWER TILES")
	for i, plant := range world.World.PowerPlants {
		for y := 0; y < plantSize; y++ {
			t := world.World.Power.GetTile(plant.X+1, plant.Y-y)
			if t != nil {
				powerSourceTiles[i] = append(powerSourceTiles[i], t)
			}

			t = world.World.Power.GetTile(plant.X-plantSize, plant.Y-y)
			if t != nil {
				powerSourceTiles[i] = append(powerSourceTiles[i], t)
			}
		}
		for x := 0; x < plantSize; x++ {
			t := world.World.Power.GetTile(plant.X-x, plant.Y+1)
			if t != nil {
				powerSourceTiles[i] = append(powerSourceTiles[i], t)
			}

			t = world.World.Power.GetTile(plant.X, plant.Y-plantSize)
			if t != nil {
				powerSourceTiles[i] = append(powerSourceTiles[i], t)
			}
		}
	}

	var totalPowerRequired int

	var havePowerOut bool

	world.ResetPowerOuts()

	// TODO use a consistent procedure to check each building that needs power
	// as connected via road to a power plant, and power-out buildings without enough power
	// "citizens report brown-outs"
	for _, zone := range world.World.Zones {
		// TODO lock, set powered status on build immediately

		powerRequired := world.ZonePowerRequirement[zone.Type]
		_ = powerRequired

		const zoneSize = 2
		var powerDestinationTiles []*world.PowerMapTile
		for y := 0; y < zoneSize; y++ {
			t := world.World.Power.GetTile(zone.X+1, zone.Y-y)
			if t != nil {
				powerDestinationTiles = append(powerDestinationTiles, t)
			}

			t = world.World.Power.GetTile(zone.X-zoneSize, zone.Y-y)
			if t != nil {
				powerDestinationTiles = append(powerDestinationTiles, t)
			}
		}
		for x := 0; x < zoneSize; x++ {
			t := world.World.Power.GetTile(zone.X-x, zone.Y+1)
			if t != nil {
				powerDestinationTiles = append(powerDestinationTiles, t)
			}

			t = world.World.Power.GetTile(zone.X, zone.Y-zoneSize)
			if t != nil {
				powerDestinationTiles = append(powerDestinationTiles, t)
			}
		}

		var powered bool
	FINDPOWERPATH:
		for j := range powerRemaining {
			if powerRemaining[j] < powerRequired {
				continue
			}

			for _, powerSource := range powerSourceTiles[j] {
				for _, to := range powerDestinationTiles {
					if to == nil {
						continue
					}

					from := world.World.Power.GetTile(powerSource.X, powerSource.Y)

					log.Println("SEARCH", from.X, from.Y, "TO", to.X, to.Y)

					/*for _, n := range powerSource.PathNeighbors() {
						t := n.(*world.PowerMapTile)
						log.Println("NEIGHBOR", t.X, t.Y, t.CarriesPower)
					}*/

					p, dist, found := astar.Path(from, to)
					if found {
						log.Printf("Resulting path\n%+v %f", p, dist)
						powerRemaining[j] -= powerRequired
						powered = true
						break FINDPOWERPATH
					}
				}
			}
		}
		zone.Powered = powered
		log.Println("ZONE", zone, zone.Powered)

		if !powered {
			havePowerOut = true
			world.World.PowerOuts[zone.X][zone.Y] = true
		}

		totalPowerRequired += powerRequired
	}

	if !havePowerOut {
		world.World.PowerUpdated = false
	}

	world.World.PowerAvailable, world.World.PowerNeeded = totalPowerAvailable, totalPowerRequired

	return nil
}

func (s *PowerScanSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
