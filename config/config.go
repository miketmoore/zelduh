package config

const (
	TranslationFile = "i18n/zelduh/en-US.all.json"
	Lang            = "en-US"
)

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
