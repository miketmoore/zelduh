package zelduh

// FrameRate is used to determine which sprite to use for animations
const FrameRate int = 5

// TileSize defines the width and height of a tile
const TileSize float64 = 48

const (
	// WinX is the x coordinate of the window
	WinX float64 = 0

	// WinY is the y coordinate of the window
	WinY float64 = 0

	// WinW is the width of the window
	WinW float64 = 800

	// WinH is the height of the window
	WinH float64 = 800
)

const (
	// MapW is the width of the game map
	MapW float64 = 672 // 48 * 14

	// MapH is the height of the game map
	MapH float64 = 576 // 48 * 12

	// MapX is the x coordinate of the game map
	MapX = (WinW - MapW) / 2

	// MapY is the y coordinate of the game map
	MapY = (WinH - MapH) / 2
)

// SpritesheetPath is the file path for the spritesheet
const SpritesheetPath string = "assets/spritesheet.png"

// TilemapDir is the directory where the tilemap files are located
const TilemapDir = "assets/tilemaps/"

// TilemapFiles is a list of tilemap filenames
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
var Overworld = [][]RoomID{
	{1, 10},
	{2, 0, 0, 8},
	{3, 5, 6, 7},
	{9},
	{11},
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

// Map of RoomID to a Room configuration
var RoomsMap = Rooms{
	1: NewRoom("overworldFourWallsDoorBottomRight",
		GetPreset("puzzleBox")(5, 5),
		GetPreset("floorSwitch")(5, 6),
		GetPreset("toggleObstacle")(10, 7),
	),
	2: NewRoom("overworldFourWallsDoorTopBottom",
		GetPreset("skull")(5, 5),
		GetPreset("skeleton")(11, 9),
		GetPreset("spinner")(7, 9),
		GetPreset("eyeburrower")(8, 9),
	),
	3: NewRoom("overworldFourWallsDoorRightTopBottom",
		WarpStone(3, 7, 6, 5),
	),
	5: NewRoom("rockWithCaveEntrance",
		Config{
			Category:     CategoryWarp,
			WarpToRoomID: 11,
			W:            TileSize,
			H:            TileSize,
			X:            (TileSize * 7) + TileSize/2,
			Y:            (TileSize * 9) + TileSize/2,
			Hitbox: &HitboxConfig{
				Radius: 30,
			},
		},
		Config{
			Category:     CategoryWarp,
			WarpToRoomID: 11,
			W:            TileSize,
			H:            TileSize,
			X:            (TileSize * 8) + TileSize/2,
			Y:            (TileSize * 9) + TileSize/2,
			Hitbox: &HitboxConfig{
				Radius: 30,
			},
		},
	),
	6:  NewRoom("rockPathLeftRightEntrance"),
	7:  NewRoom("overworldFourWallsDoorLeftTop"),
	8:  NewRoom("overworldFourWallsDoorBottom"),
	9:  NewRoom("overworldFourWallsDoorTop"),
	10: NewRoom("overworldFourWallsDoorLeft"),
	11: NewRoom("dungeonFourDoors",
		// South door of cave - warp to cave entrance
		Config{
			Category:     CategoryWarp,
			WarpToRoomID: 5,
			W:            TileSize,
			H:            TileSize,
			X:            (TileSize * 6) + TileSize + (TileSize / 2.5),
			Y:            (TileSize * 1) + TileSize + (TileSize / 2.5),
			Hitbox: &HitboxConfig{
				Radius: 15,
			},
		},
		Config{
			Category:     CategoryWarp,
			WarpToRoomID: 5,
			W:            TileSize,
			H:            TileSize,
			X:            (TileSize * 7) + TileSize + (TileSize / 2.5),
			Y:            (TileSize * 1) + TileSize + (TileSize / 2.5),
			Hitbox: &HitboxConfig{
				Radius: 15,
			},
		},
	),
}

// Just a stub for now since English is the only language supported at this time
var LocaleMessages = map[string]map[string]string{
	"en": {
		"gameTitle":             "Zelduh",
		"pauseScreenMessage":    "Paused",
		"gameOverScreenMessage": "Game Over",
	},
	"es": {
		"gameTitle":             "Zelduh",
		"pauseScreenMessage":    "Paused",
		"gameOverScreenMessage": "Game Over",
	},
}
