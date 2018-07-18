package config

import "github.com/miketmoore/terraform2d"

const TranslationFile = "i18n/zelduh/active.en.toml"

const FrameRate int = 5

// TileSize defines the width and height of a tile
const TileSize float64 = 48

const (
	WinX float64 = 0
	WinY float64 = 0
	WinW float64 = 800
	WinH float64 = 800
)

const (
	MapW float64 = 672 // 48 * 14
	MapH float64 = 576 // 48 * 12
	MapX         = (WinW - MapW) / 2
	MapY         = (WinH - MapH) / 2
)

const SpritesheetPath string = "assets/spritesheet.png"

const TilemapDir = "assets/tilemaps/"

var TilemapFiles = []string{
	"overworldOpen",
	"overworldOpenCircleOfTrees",
	"overworldFourWallsDoorBottom",
	"overworldFourWallsDoorLeftTop",
	"overworldFourWallsDoorRightTop",
	"overworldFourWallsDoorTopBottom",
	"overworldFourWallsDoorRightTopBottom",
	"overworldFourWallsDoorBottomRight",
	"overworldFourWallsDoorTop",
	"overworldFourWallsDoorRight",
	"overworldFourWallsDoorLeft",
	"overworldTreeClusterTopRight",
	"overworldFourWallsClusterTrees",
	"overworldFourWallsDoorsAllSides",
	"rockPatternTest",
	"rockPathOpenLeft",
	"rockWithCaveEntrance",
	"rockPathLeftRightEntrance",
	"test",
	"dungeonFourDoors",
}

// Overworld is a multi-dimensional array representing the overworld
// Each room ID should be unique
var Overworld = [][]terraform2d.RoomID{
	[]terraform2d.RoomID{1, 10},
	[]terraform2d.RoomID{2, 0, 0, 8},
	[]terraform2d.RoomID{3, 5, 6, 7},
	[]terraform2d.RoomID{9},
	[]terraform2d.RoomID{11},
}

// NonObstacleSprites defines which sprites are not obstacles
var NonObstacleSprites = map[int]bool{
	8:   true,
	9:   true,
	24:  true,
	37:  true,
	38:  true,
	52:  true,
	53:  true,
	66:  true,
	86:  true,
	136: true,
	137: true,
}
