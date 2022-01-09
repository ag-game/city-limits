package world

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"code.rocketnine.space/tslocum/citylimits/asset"
	"code.rocketnine.space/tslocum/citylimits/component"
	. "code.rocketnine.space/tslocum/citylimits/ecs"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

const TileSize = 128

const (
	StructureHouse1 = iota + 1
	StructurePoliceStation
)

var World = &GameWorld{
	CamScale:     1,
	CamMoving:    true,
	PlayerWidth:  8,
	PlayerHeight: 32,
	TileImages:   make(map[uint32]*ebiten.Image),
	ResetGame:    true,
	Level:        NewLevel(256),
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

	MessageVisible  bool
	MessageTicks    int
	MessageDuration int
	MessageUpdated  bool
	MessageText     string

	PlayerX, PlayerY float64

	CamX, CamY float64
	CamScale   float64
	CamMoving  bool

	PlayerWidth  float64
	PlayerHeight float64

	HoverStructure         int
	HoverX, HoverY         int
	HoverLastX, HoverLastY int

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

	TileImages map[uint32]*ebiten.Image

	ResetGame bool

	GotCursorPosition bool

	tilesets []*ebiten.Image

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

	World.ObjectGroups = nil
	World.HazardRects = nil
	World.CreepRects = nil
	World.CreepEntities = nil
	World.TriggerEntities = nil
	World.TriggerRects = nil
	World.TriggerNames = nil

	World.MessageVisible = false
}

func BuildStructure(structureType int, placeX int, placeY int) (*Structure, error) {
	tt := time.Now()
	defer func() {
		log.Println(time.Since(tt))
	}()

	loader := tiled.Loader{
		FileSystem: asset.FS,
	}

	var filePath string
	switch structureType {
	case StructureHouse1:
		filePath = "map/house1.tmx"
	case StructurePoliceStation:
		filePath = "map/policestation.tmx"
	default:
		panic(fmt.Sprintf("unknown structure %d", structureType))
	}

	// Parse .tmx file.
	m, err := loader.LoadFromFile(filepath.FromSlash(filePath))
	if err != nil {
		log.Fatalf("error parsing world: %+v", err)
	}

	if placeX < 0 || placeY < 0 || placeX+m.Width > 256 || placeY+m.Height > 256 {
		return nil, errors.New("invalid location: building does not fit")
	}

	// Load tileset.

	tileset := m.Tilesets[0]
	if len(World.tilesets) == 0 {
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
	}

	// Load tiles.

	for i := uint32(0); i < uint32(tileset.TileCount); i++ {
		rect := tileset.GetTileRect(i)
		World.TileImages[i+tileset.FirstGID] = World.tilesets[0].SubImage(rect).(*ebiten.Image)
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

	structure := &Structure{
		Type: structureType,
		X:    placeX,
		Y:    placeY,
	}

	// TODO Add entity

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

				for i > len(World.Level.Tiles)-2 {
					World.Level.AddLayer()
				}

				if World.Level.Tiles[i][x+placeX][y+placeY] == nil {
					World.Level.Tiles[i][x+placeX][y+placeY] = &Tile{}
				}
				World.Level.Tiles[i][x+placeX][y+placeY].Sprite = World.TileImages[t.Tileset.FirstGID+t.ID]

				// TODO handle flipping

				//tileX, tileY := TileToGameCoords(x, y)
				//e := createTileEntity(t, tileX+float64(layer.OffsetY*2), tileY+float64(layer.OffsetY*2))
				//tileEntities = append(tileEntities, e)
			}
		}
	}

	// Load ObjectGroups.

	var objects []*tiled.ObjectGroup
	var loadObjects func(grp *tiled.Group)
	loadObjects = func(grp *tiled.Group) {
		for _, subGrp := range grp.Groups {
			loadObjects(subGrp)
		}
		for _, objGrp := range grp.ObjectGroups {
			objects = append(objects, objGrp)
		}
	}
	for _, grp := range m.Groups {
		loadObjects(grp)
	}
	for _, objGrp := range m.ObjectGroups {
		objects = append(objects, objGrp)
	}

	World.Map = m
	World.ObjectGroups = objects

	for _, grp := range World.ObjectGroups {
		if grp.Name == "TRIGGERS" {
			for _, obj := range grp.Objects {
				mapTile := ECS.NewEntity()
				ECS.AddComponent(mapTile, &component.PositionComponent{
					X: obj.X,
					Y: obj.Y - 32,
				})
				ECS.AddComponent(mapTile, &component.SpriteComponent{
					Image: World.TileImages[obj.GID],
				})

				World.TriggerNames = append(World.TriggerNames, obj.Name)
				World.TriggerEntities = append(World.TriggerEntities, mapTile)
				World.TriggerRects = append(World.TriggerRects, ObjectToRect(obj))
			}
		} else if grp.Name == "HAZARDS" {
			for _, obj := range grp.Objects {
				r := ObjectToRect(obj)
				r.Min.Y += 32
				r.Max.Y += 32
				World.HazardRects = append(World.HazardRects, r)
			}
		} else if grp.Name == "CREEPS" {
			for _, obj := range grp.Objects {
				creepType := component.CreepSnowblower
				switch obj.GID {
				case 9:
					creepType = component.CreepSmallRock
				case 18:
					creepType = component.CreepMediumRock
				case 23:
					creepType = component.CreepLargeRock
				}
				r := ObjectToRect(obj)
				c := NewActor(creepType, int64(obj.ID), float64(r.Min.X), float64(r.Min.Y))
				World.CreepRects = append(World.CreepRects, r)
				World.CreepEntities = append(World.CreepEntities, c)
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
}

func SetMessage(message string, duration int) {
	World.MessageText = message
	World.MessageVisible = true
	World.MessageUpdated = true
	World.MessageDuration = duration
	World.MessageTicks = 0
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
	return cx, cy
}

func IsoToScreen(x, y float64) (float64, float64) {
	cx, cy := float64(World.ScreenW/2), float64(World.ScreenH/2)
	return ((x - World.CamX) * World.CamScale) + cx, ((y - World.CamY) * World.CamScale) + cy
}

func ScreenToIso(x, y int) (float64, float64) {
	cx, cy := float64(World.ScreenW/2), float64(World.ScreenH/2)
	return ((float64(x) - cx) / World.CamScale) + World.CamX, ((float64(y) - cy) / World.CamScale) + World.CamY
}
