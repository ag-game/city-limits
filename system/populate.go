package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type PopulateSystem struct {
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Weapon
}

func NewPopulateSystem() *PopulateSystem {
	s := &PopulateSystem{}

	return s
}

func (s *PopulateSystem) Update(_ gohan.Entity) error {
	if world.World.Paused {
		return nil
	}

	const popDuration = 144 * 7
	if world.World.Ticks%popDuration != 0 {
		return nil
	}

	// Thresholds.
	const (
		lowDensity    = 3
		mediumDensity = 7
	)
	buildStructureType := func(structureType int, population int) int {
		switch structureType {
		case world.StructureResidentialZone:
			switch {
			case population == 0:
				return world.StructureResidentialZone
			case population <= lowDensity:
				return world.StructureResidentialLow
			case population <= mediumDensity:
				return world.StructureResidentialMedium
			default:
				return world.StructureResidentialHigh
			}
		case world.StructureCommercialZone:
			switch {
			case population == 0:
				return world.StructureCommercialZone
			case population <= lowDensity:
				return world.StructureCommercialLow
			case population <= mediumDensity:
				return world.StructureCommercialMedium
			default:
				return world.StructureCommercialHigh
			}
		case world.StructureIndustrialZone:
			switch {
			case population == 0:
				return world.StructureIndustrialZone
			case population <= lowDensity:
				return world.StructureIndustrialLow
			case population <= mediumDensity:
				return world.StructureIndustrialMedium
			default:
				return world.StructureIndustrialHigh
			}
		default:
			return structureType
		}
	}

	const maxPopulation = 10
	popR, popC, popI := world.Population()
	targetR, targetC, targetI := world.TargetPopulation()
	for _, zone := range world.World.Zones {
		var offset int
		if zone.Type == world.StructureResidentialZone {
			if popR < targetR {
				offset = 1
			} else if popR > targetR {
				offset = -1
			}
		} else if zone.Type == world.StructureCommercialZone {
			if popC < targetC {
				offset = 1
			} else if popC > targetC {
				offset = -1
			}
		} else { // Industrial
			if popI < targetI {
				offset = 1
			} else if popI > targetI {
				offset = -1
			}
		}
		if offset == -1 && zone.Population > 0 {
			zone.Population--
			if zone.Type == world.StructureResidentialZone {
				popR--
			} else if zone.Type == world.StructureCommercialZone {
				popC--
			} else { // Industrial
				popI--
			}
		} else if offset == 1 && zone.Population < maxPopulation && zone.Powered {
			zone.Population++
			if zone.Type == world.StructureResidentialZone {
				popR++
			} else if zone.Type == world.StructureCommercialZone {
				popC++
			} else { // Industrial
				popI++
			}
		}
		newType := buildStructureType(zone.Type, zone.Population)
		// TODO only bulldoze when changed
		for offsetX := 0; offsetX < 2; offsetX++ {
			for offsetY := 0; offsetY < 2; offsetY++ {
				world.BuildStructure(world.StructureBulldozer, false, zone.X-offsetX, zone.Y-offsetY, true)
			}
		}
		world.BuildStructure(newType, false, zone.X, zone.Y, true)
	}

	// TODO populate and de-populate zones by target population
	// for zone in zones
	return nil
}

func (s *PopulateSystem) Draw(_ gohan.Entity, screen *ebiten.Image) error {
	return gohan.ErrUnregister
}
