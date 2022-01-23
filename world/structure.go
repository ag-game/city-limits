package world

import "code.rocketnine.space/tslocum/gohan"

const (
	StructureToggleHelp = iota + 1
	StructureToggleTransparentStructures
	StructureBulldozer
	StructureRoad
	StructureResidentialZone
	StructureResidentialLow
	StructureResidentialMedium
	StructureResidentialHigh
	StructureCommercialZone
	StructureCommercialLow
	StructureCommercialMedium
	StructureCommercialHigh
	StructureIndustrialZone
	StructureIndustrialLow
	StructureIndustrialMedium
	StructureIndustrialHigh
	StructurePoliceStation
	StructurePowerPlantCoal
	StructurePowerPlantSolar
	StructurePowerPlantNuclear
)

var StructureFilePaths = map[int]string{
	StructureBulldozer:         "map/bulldozer.tmx",
	StructureRoad:              "map/road.tmx",
	StructureResidentialZone:   "map/residential_zone.tmx",
	StructureResidentialLow:    "map/residential_low1.tmx",
	StructureResidentialMedium: "map/residential_med1.tmx",
	StructureResidentialHigh:   "map/residential_high1.tmx",
	StructureCommercialZone:    "map/commercial_zone.tmx",
	StructureCommercialLow:     "map/commercial_low1.tmx",
	StructureCommercialMedium:  "map/commercial_med1.tmx",
	StructureCommercialHigh:    "map/commercial_high1.tmx",
	StructureIndustrialZone:    "map/industrial_zone.tmx",
	StructureIndustrialLow:     "map/industrial_low1.tmx",
	StructureIndustrialMedium:  "map/industrial_med1.tmx",
	StructureIndustrialHigh:    "map/industrial_high1.tmx",
	StructurePoliceStation:     "map/policestation.tmx",
	StructurePowerPlantCoal:    "map/power_coal.tmx",
	StructurePowerPlantSolar:   "map/power_solar.tmx",
	StructurePowerPlantNuclear: "map/power_nuclear.tmx",
}

type Structure struct {
	Type int
	X, Y int

	Entity   gohan.Entity
	Children []gohan.Entity
}
