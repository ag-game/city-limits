package world

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"path/filepath"
	"strconv"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"code.rocketnine.space/tslocum/citylimits/asset"
	"code.rocketnine.space/tslocum/citylimits/component"
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

const startingYear = 1950

const maxPopulation = 100000

const (
	MonthTicks = 144 * 3
	YearTicks  = MonthTicks * 12
)

const TileSize = 64

var DirtTile = uint32(9*32 + (0))

const startingFunds = 10000

const startingZoom = 1.0

const SidebarWidth = 198

type HUDButton struct {
	Sprite                       *ebiten.Image
	SpriteOffsetX, SpriteOffsetY float64
	Label                        string
	StructureType                int
}

var HUDButtons []*HUDButton

var CameraMinZoom = 0.4
var CameraMaxZoom = 1.0

var World = &GameWorld{
	CamScale:       startingZoom,
	CamScaleTarget: startingZoom,
	CamMoving:      true,
	PlayerWidth:    8,
	PlayerHeight:   32,
	TileImages:     make(map[uint32]*ebiten.Image),
	ResetGame:      true,
	Level:          NewLevel(256),
	Printer:        message.NewPrinter(language.English),
	Power:          newPowerMap(),
	PowerOuts:      newPowerOuts(),
}

type Zone struct {
	Type       int // StructureResidentialZone, StructureCommercialZone or StructureIndustrialZone
	X, Y       int
	Population int
	Powered    bool
}

type PowerPlant struct {
	Type int
	X, Y int
}

type GameWorld struct {
	Level *GameLevel

	Player gohan.Entity

	ScreenW, ScreenH int

	DisableEsc bool

	Debug  int
	NoClip bool

	GameStarted      bool
	GameStartedTicks int
	GameOver         bool

	PlayerX, PlayerY float64

	CamX, CamY     float64
	CamScale       float64
	CamScaleTarget float64
	CamMoving      bool

	PlayerWidth  float64
	PlayerHeight float64

	HoverStructure         int
	HoverX, HoverY         int
	HoverLastX, HoverLastY int
	HoverValid             bool

	Map             *tiled.Map
	ObjectGroups    []*tiled.ObjectGroup
	HazardRects     []image.Rectangle
	CreepRects      []image.Rectangle
	CreepEntities   []gohan.Entity
	TriggerEntities []gohan.Entity
	TriggerRects    []image.Rectangle
	TriggerNames    []string

	NativeResolution bool

	BrokenPieceA, BrokenPieceB gohan.Entity

	TileImages         map[uint32]*ebiten.Image
	TileImagesFirstGID uint32

	ResetGame bool

	MuteMusic        bool
	MuteSoundEffects bool // TODO

	GotCursorPosition bool

	tilesets []*ebiten.Image

	EnvironmentSprites int

	SelectedStructure *Structure

	HUDUpdated     bool
	HUDButtonRects []image.Rectangle

	PowerPlants []*PowerPlant
	Zones       []*Zone

	PowerOuts [][]bool

	Ticks int

	Paused bool

	Funds int

	Printer *message.Printer

	TransparentStructures bool

	Messages      []string
	MessagesTicks []int

	Power          PowerMap
	PowerUpdated   bool
	PowerAvailable int
	PowerNeeded    int

	resetTipShown bool
}

func TileToGameCoords(x, y int) (float64, float64) {
	//return float64(x) * 32, float64(g.currentMap.Height*32) - float64(y)*32 - 32
	return float64(x) * TileSize, float64(y) * TileSize
}

func Reset() {
	for _, e := range ECS.Entities() {
		ECS.RemoveEntity(e)
	}
	World.Player = 0

	World.Funds = startingFunds

	World.ObjectGroups = nil
	World.HazardRects = nil
	World.CreepRects = nil
	World.CreepEntities = nil
	World.TriggerEntities = nil
	World.TriggerRects = nil
	World.TriggerNames = nil
}

func LoadMap(structureType int) (*tiled.Map, error) {
	filePath := StructureFilePaths[structureType]
	if filePath == "" {
		panic(fmt.Sprintf("unknown structure %d", structureType))
	}

	// Parse .tmx file.
	m, err := tiled.LoadFile(filepath.FromSlash(filePath), tiled.WithFileSystem(asset.FS))
	if err != nil {
		log.Fatalf("error parsing world: %+v", err)
	}

	return m, err
}

func DrawMap(structureType int) *ebiten.Image {
	img := ebiten.NewImage(128, 128)

	m, err := LoadMap(structureType)
	if err != nil {
		panic(err)
	}

	var t *tiled.LayerTile
	for i, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}

				tileImg := World.TileImages[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}

				xi, yi := CartesianToIso(float64(x), float64(y))

				scale := 0.9 / float64(m.Width)
				if m.Width < 2 {
					scale = 0.6
				}

				paddingX := 64.0
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(xi+(paddingX*(float64(m.Width)-1)), (yi+float64(i*-40))+92)
				op.GeoM.Scale(scale, scale)
				img.DrawImage(tileImg, op)
			}
		}
	}

	return img
}

func LoadTileset() error {
	m, err := LoadMap(StructureResidentialLow)
	if err != nil {
		return err
	}

	// Load tileset.

	if len(World.tilesets) != 0 {
		return nil // Already loaded.
	}

	tileset := m.Tilesets[0]
	imgPath := filepath.Join("./image/tileset/", tileset.Image.Source)
	f, err := asset.FS.Open(filepath.FromSlash(imgPath))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	World.tilesets = append(World.tilesets, ebiten.NewImageFromImage(img))

	// Load tiles.

	for i := uint32(0); i < uint32(tileset.TileCount); i++ {
		rect := tileset.GetTileRect(i)
		World.TileImages[i+tileset.FirstGID] = World.tilesets[0].SubImage(rect).(*ebiten.Image)
	}

	World.TileImagesFirstGID = tileset.FirstGID
	return nil
}

func BuildStructure(structureType int, hover bool, placeX int, placeY int) (*Structure, error) {
	// For previewing buildings
	/*v := rand.Intn(3)
	if structureType == StructureResidentialZone {
		switch v {
		case 0:
			structureType = StructureResidentialLow
		case 1:
			structureType = StructureResidentialMedium
		case 2:
			structureType = StructureResidentialHigh
		}
	} else if structureType == StructureCommercialZone {
		switch v {
		case 0:
			structureType = StructureCommercialLow
		case 1:
			structureType = StructureCommercialMedium
		case 2:
			structureType = StructureCommercialHigh
		}
	} else if structureType == StructureIndustrialZone {
		switch v {
		case 0:
			structureType = StructureIndustrialLow
		case 1:
			structureType = StructureIndustrialMedium
		case 2:
			structureType = StructureIndustrialHigh
		}
	}*/

	m, err := LoadMap(structureType)
	if err != nil {
		return nil, err
	}

	if m.Width != 1 || m.Height != 1 {
		if placeX == 0 {
			placeX = 1
		}
		if placeY == 0 {
			placeY = 1
		}
	}

	w := m.Width - 1
	h := m.Height - 1

	if placeX-w < 0 || placeY-h < 0 || placeX >= 256 || placeY >= 256 {
		return nil, errors.New("invalid location: building does not fit")
	}

	structure := &Structure{
		Type: structureType,
		X:    placeX,
		Y:    placeY,
	}

	if structureType == StructureBulldozer && !hover {
		// TODO bulldoze entire structure, remove from zones
		var bulldozed bool
		for i := range World.Level.Tiles {
			if World.Level.Tiles[i][placeX][placeY].Sprite != nil {
				World.Level.Tiles[i][placeX][placeY].Sprite = nil
				bulldozed = true
			}

			var img *ebiten.Image
			if i == 0 {
				img = World.TileImages[DirtTile+World.TileImagesFirstGID]
			}
			if World.Level.Tiles[i][placeX][placeY].EnvironmentSprite != img {
				World.Level.Tiles[i][placeX][placeY].EnvironmentSprite = img
				bulldozed = true
			}
		}
		if !bulldozed {
			return nil, errors.New("nothing to bulldoze")
		}
		World.Power.SetTile(placeX, placeY, false)
		return structure, nil
	}

	createTileEntity := func(t *tiled.LayerTile, x float64, y float64) gohan.Entity {
		mapTile := ECS.NewEntity()
		ECS.AddComponent(mapTile, &component.PositionComponent{
			X: x,
			Y: y,
		})

		sprite := &component.SpriteComponent{
			Image:          World.TileImages[t.Tileset.FirstGID+t.ID],
			HorizontalFlip: t.HorizontalFlip,
			VerticalFlip:   t.VerticalFlip,
			DiagonalFlip:   t.DiagonalFlip,
		}
		ECS.AddComponent(mapTile, sprite)

		return mapTile
	}
	_ = createTileEntity

	// TODO Add entity

	tileOccupied := func(tx int, ty int) bool {
		return World.Level.Tiles[1][tx][ty].Sprite != nil || (World.Level.Tiles[0][tx][ty].Sprite != nil && (structureType != StructureRoad || World.Level.Tiles[0][tx][ty].Sprite != World.TileImages[World.TileImagesFirstGID]))
	}

	valid := true
	var existingRoadTiles int
VALIDBUILD:
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tx, ty := (x+placeX)-w, (y+placeY)-h
			if structureType == StructureRoad && World.Level.Tiles[0][tx][ty].Sprite == World.TileImages[World.TileImagesFirstGID] {
				existingRoadTiles++
			}
			if tileOccupied(tx, ty) && structureType != StructureBulldozer {
				valid = false
				break VALIDBUILD
			}
		}
	}
	if structureType == StructureRoad && existingRoadTiles == 4 {
		valid = false
	}
	if hover {
		if structureType == StructureBulldozer {
			World.HoverValid = true
		} else {
			World.HoverValid = valid
		}
	} else if !valid {
		return nil, errors.New("invalid location: space already occupied")
	}

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			tx, ty := (x+placeX)-w, (y+placeY)-h
			if hover {
				if !tileOccupied(tx, ty) || structureType == StructureBulldozer {
					if structureType != StructureBulldozer {
						World.Level.Tiles[0][tx][ty].HoverSprite = World.TileImages[World.TileImagesFirstGID]
					}
					// Hide environment sprites temporarily.
					for i := 1; i < len(World.Level.Tiles); i++ {
						World.Level.Tiles[i][tx][ty].HoverSprite = asset.ImgBlank
					}
				}
			} else {
				World.Level.Tiles[0][tx][ty].Sprite = World.TileImages[World.TileImagesFirstGID]
				World.Level.Tiles[0][tx][ty].EnvironmentSprite = nil
				World.Level.Tiles[1][tx][ty].EnvironmentSprite = nil
			}
		}
	}

	var t *tiled.LayerTile
	for i, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}

				tileImg := World.TileImages[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}

				layerNum := i
				if structureType != StructureRoad {
					layerNum++
				}

				for layerNum > len(World.Level.Tiles)-1 {
					World.Level.AddLayer()
				}

				tx, ty := (x+placeX)-w, (y+placeY)-h
				if hover {
					if !tileOccupied(tx, ty) || structureType == StructureBulldozer {
						World.Level.Tiles[layerNum][tx][ty].HoverSprite = World.TileImages[t.Tileset.FirstGID+t.ID]
					}
				} else {
					World.Level.Tiles[layerNum][tx][ty].Sprite = World.TileImages[t.Tileset.FirstGID+t.ID]

					if structureType == StructureRoad {
						World.Power.SetTile(tx, ty, true)
					}
					World.PowerUpdated = true
				}

				// TODO handle flipping
			}
		}
	}

	return structure, nil
}

func ObjectToRect(o *tiled.Object) image.Rectangle {
	x, y, w, h := int(o.X), int(o.Y), int(o.Width), int(o.Height)
	y -= 32
	return image.Rect(x, y, x+w, y+h)
}

func LevelCoordinatesToScreen(x, y float64) (float64, float64) {
	return (x - World.CamX) * World.CamScale, (y - World.CamY) * World.CamScale
}

func (w *GameWorld) SetGameOver(vx, vy float64) {
	if w.GameOver {
		return
	}

	w.GameOver = true

	// TODO
}

// TODO move
func NewActor(creepType int, creepID int64, x float64, y float64) gohan.Entity {
	actor := ECS.NewEntity()

	ECS.AddComponent(actor, &component.PositionComponent{
		X: x,
		Y: y,
	})

	ECS.AddComponent(actor, &component.ActorComponent{
		Type:       creepType,
		Health:     64,
		FireAmount: 8,
		FireRate:   144 / 4,
		Rand:       rand.New(rand.NewSource(creepID)),
	})

	return actor
}

func StartGame() {
	if World.GameStarted {
		return
	}
	World.GameStarted = true

	if !World.MuteMusic {
		asset.SoundMusic.Play()
	}
}

// CartesianToIso transforms cartesian coordinates into isometric coordinates.
func CartesianToIso(x, y float64) (float64, float64) {
	ix := (x - y) * float64(TileSize/2)
	iy := (x + y) * float64(TileSize/4)
	return ix, iy
}

// CartesianToIso transforms cartesian coordinates into isometric coordinates.
func IsoToCartesian(x, y float64) (float64, float64) {
	cx := (x/float64(TileSize/2) + y/float64(TileSize/4)) / 2
	cy := (y/float64(TileSize/4) - (x / float64(TileSize/2))) / 2
	cx-- // TODO Why is this necessary?
	return cx, cy
}

func IsoToScreen(x, y float64) (float64, float64) {
	cx, cy := float64(World.ScreenW/2), float64(World.ScreenH/2)
	return ((x - World.CamX) * World.CamScale) + cx, ((y - World.CamY) * World.CamScale) + cy
}

func ScreenToIso(x, y int) (float64, float64) {
	// Offset cursor to first above ground layer.
	y += int(float64(16) * World.CamScale)

	cx, cy := float64(World.ScreenW/2), float64(World.ScreenH/2)
	return ((float64(x) - cx) / World.CamScale) + World.CamX, ((float64(y) - cy) / World.CamScale) + World.CamY
}

func ScreenToCartesian(x, y int) (float64, float64) {
	xi, yi := ScreenToIso(x, y)
	return IsoToCartesian(xi, yi)
}

func HUDButtonAt(x, y int) *HUDButton {
	for i, colorRect := range World.HUDButtonRects {
		point := image.Point{x, y}
		if point.In(colorRect) {
			return HUDButtons[i]
		}
	}
	return nil
}

func SetHoverStructure(structureType int) {
	World.HoverStructure = structureType
	World.HUDUpdated = true
}

func Satisfaction() (r, c, i float64) {
	return 0.01, 0.0, 0.02
}

func TargetPopulation() (r, c, i int) {
	currentMax := maxPopulation * ((1 + float64(World.Ticks/(MonthTicks*7))) / 108)

	satisfactionR, satisfactionC, satisfactionI := Satisfaction()
	return int(satisfactionR * currentMax), int(satisfactionC * currentMax), int(satisfactionI * currentMax)
}

func Demand() (r, c, i float64) {
	targetR, targetC, targetI := TargetPopulation()

	populationR, populationC, populationI := Population()
	r, c, i = float64(targetR)-float64(populationR), float64(targetC)-float64(populationC), float64(targetI)-float64(populationI)
	max := r
	if c > max {
		max = c
	}
	if i > max {
		max = i
	}
	barPeak := 100.0
	r, c, i = r/barPeak, c/barPeak, i/barPeak
	clamp := func(v float64) float64 {
		if math.IsNaN(v) {
			return 0
		}
		if v < -1 {
			v = -1
		} else if v > 1 {
			v = 1
		}
		return v
	}
	return clamp(r), clamp(c), clamp(i)
}

var StructureTooltips = map[int]string{
	StructureToggleTransparentStructures: "Transparent buildings",
	StructureBulldozer:                   "Bulldozer",
	StructureRoad:                        "Road",
	StructurePoliceStation:               "Police station",
	StructurePowerPlantCoal:              "Coal power plant",
	StructureResidentialZone:             "Residential zone",
	StructureCommercialZone:              "Commercial zone",
	StructureIndustrialZone:              "Industrial zone",
}

var StructureCosts = map[int]int{
	StructureBulldozer:       5,
	StructureRoad:            25,
	StructurePoliceStation:   1000,
	StructurePowerPlantCoal:  4000,
	StructureResidentialZone: 100,
	StructureCommercialZone:  200,
	StructureIndustrialZone:  100,
}

func Tooltip() string {
	tooltipText := StructureTooltips[World.HoverStructure]
	cost := StructureCosts[World.HoverStructure]
	if cost > 0 {
		tooltipText += World.Printer.Sprintf("\n$%d", cost)
	}
	return tooltipText
}

var monthNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func Date() (month string, year string) {
	y, m := World.Ticks/YearTicks, (World.Ticks%YearTicks)/MonthTicks
	return monthNames[m], strconv.Itoa(startingYear + y)
}

func Population() (r, c, i int) {
	for _, zone := range World.Zones {
		switch zone.Type {
		case StructureResidentialZone:
			r += zone.Population
		case StructureCommercialZone:
			c += zone.Population
		case StructureIndustrialZone:
			i += zone.Population
		}
	}
	return r, c, i
}

var messageLock = &sync.Mutex{}

const messageDuration = 144 * 3

func TickMessages() {
	messageLock.Lock()
	defer messageLock.Unlock()

	var removed int
	for j := 0; j < len(World.MessagesTicks); j++ {
		i := j - removed
		if World.MessagesTicks[i] == 0 {
			World.Messages = append(World.Messages[:i], World.Messages[i+1:]...)
			World.MessagesTicks = append(World.MessagesTicks[:i], World.MessagesTicks[i+1:]...)
			removed++

			World.HUDUpdated = true
		} else if World.MessagesTicks[i] > 0 {
			World.MessagesTicks[i]--
		}
	}
}

func ShowMessage(message string, duration int) {
	messageLock.Lock()
	defer messageLock.Unlock()

	World.Messages = append(World.Messages, message)
	World.MessagesTicks = append(World.MessagesTicks, duration)

	World.HUDUpdated = true
}

func ValidXY(x, y int) bool {
	return x >= 0 && y >= 0 && x < 256 && y < 256
}

var PowerPlantCapacities = map[int]int{
	StructurePowerPlantCoal: 60,
}

var ZonePowerRequirement = map[int]int{
	StructureResidentialZone: 1,
	StructureCommercialZone:  1,
	StructureIndustrialZone:  1,
}
