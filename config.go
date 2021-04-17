package zelduh

// FrameRate is used to determine which sprite to use for animations
const FrameRate int = 5

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

// TODO move to zelduh cmd file since it is configuration
// Map of RoomID to a Room configuration
func BuildRooms(entityConfigPresetFnManager *EntityConfigPresetFnManager, tileSize float64) Rooms {
	return Rooms{
		1: NewRoom("overworldFourWallsDoorBottomRight",
			entityConfigPresetFnManager.GetPreset("puzzleBox")(5, 5),
			entityConfigPresetFnManager.GetPreset("floorSwitch")(5, 6),
			entityConfigPresetFnManager.GetPreset("toggleObstacle")(10, 7),
		),
		2: NewRoom("overworldFourWallsDoorTopBottom",
			entityConfigPresetFnManager.GetPreset("skull")(5, 5),
			entityConfigPresetFnManager.GetPreset("skeleton")(11, 9),
			entityConfigPresetFnManager.GetPreset("spinner")(7, 9),
			entityConfigPresetFnManager.GetPreset("eyeburrower")(8, 9),
		),
		3: NewRoom("overworldFourWallsDoorRightTopBottom",
			WarpStone(entityConfigPresetFnManager, 3, 7, 6, 5),
		),
		5: NewRoom("rockWithCaveEntrance",
			EntityConfig{
				Category:     CategoryWarp,
				WarpToRoomID: 11,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 7) + tileSize/2,
				Y:            (tileSize * 9) + tileSize/2,
				Hitbox: &HitboxConfig{
					Radius: 30,
				},
			},
			EntityConfig{
				Category:     CategoryWarp,
				WarpToRoomID: 11,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 8) + tileSize/2,
				Y:            (tileSize * 9) + tileSize/2,
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
			EntityConfig{
				Category:     CategoryWarp,
				WarpToRoomID: 5,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 6) + tileSize + (tileSize / 2.5),
				Y:            (tileSize * 1) + tileSize + (tileSize / 2.5),
				Hitbox: &HitboxConfig{
					Radius: 15,
				},
			},
			EntityConfig{
				Category:     CategoryWarp,
				WarpToRoomID: 5,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 7) + tileSize + (tileSize / 2.5),
				Y:            (tileSize * 1) + tileSize + (tileSize / 2.5),
				Hitbox: &HitboxConfig{
					Radius: 15,
				},
			},
		),
	}
}
