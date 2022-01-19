package system

import (
	"code.rocketnine.space/tslocum/citylimits/component"
	"code.rocketnine.space/tslocum/citylimits/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type PopulateSystem struct {
}

func NewPopulateSystem() *PopulateSystem {
	s := &PopulateSystem{}

	return s
}

func (s *PopulateSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
	}
}

func (s *PopulateSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *PopulateSystem) Update(_ *gohan.Context) error {
	if world.World.Paused {
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
			case population <= lowDensity:
				return world.StructureResidentialLow
			case population <= mediumDensity:
				return world.StructureResidentialMedium
			default:
				return world.StructureResidentialHigh
			}
		case world.StructureCommercialZone:
			switch {
			case population <= lowDensity:
				return world.StructureCommercialLow
			case population <= mediumDensity:
				return world.StructureCommercialMedium
			default:
				return world.StructureCommercialHigh
			}
		case world.StructureIndustrialZone:
			switch {
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
	if world.World.Ticks%144 == 0 {
		for _, zone := range world.World.Zones {
			if zone.Population < maxPopulation {
				zone.Population++
			}
			newType := buildStructureType(zone.Type, zone.Population)
			for offsetX := 0; offsetX < 2; offsetX++ {
				for offsetY := 0; offsetY < 2; offsetY++ {
					world.BuildStructure(world.StructureBulldozer, false, zone.X-offsetX, zone.Y-offsetY)
				}
			}
			world.BuildStructure(newType, false, zone.X, zone.Y)
		}
	}

	// TODO populate and de-populate zones by target population
	// for zone in zones
	return nil
}

func (s *PopulateSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
