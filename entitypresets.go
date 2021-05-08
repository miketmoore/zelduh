package zelduh

const (
	PresetNameArrow            PresetName = "arrow"
	PresetNameBomb             PresetName = "bomb"
	PresetNameCoin             PresetName = "coin"
	PresetNameExplosion        PresetName = "explosion"
	PresetNameObstacle         PresetName = "obstacle"
	PresetNamePlayer           PresetName = "player"
	PresetNameFloorSwitch      PresetName = "floorSwitch"
	PresetNameToggleObstacle   PresetName = "toggleObstacle"
	PresetNamePuzzleBox        PresetName = "puzzleBox"
	PresetNameWarpStone        PresetName = "warpStone"
	PresetNameUICoin           PresetName = "uiCoin"
	PresetNameEnemySpinner     PresetName = "spinner"
	PresetNameEnemySkull       PresetName = "skull"
	PresetNameEnemySkeleton    PresetName = "skeleton"
	PresetNameHeart            PresetName = "heart"
	PresetNameEnemyEyeBurrower PresetName = "eyeBurrower"
	PresetNameSword            PresetName = "sword"
	PresetNameDialogCorner     PresetName = "dialogCorner"
	PresetNameDialogSide       PresetName = "dialogSide"
)

type BuildWarpFn func(
	warpToRoomID RoomID,
	coordinates Coordinates,
	hitboxRadius float64,
) EntityConfig

// // TODO move this to a higher level configuration location
// func BuildEntityConfigPresetFnsMap(tileSize float64) map[PresetName]EntityConfigPresetFn {

// 	dimensions := Dimensions{
// 		Width:  tileSize,
// 		Height: tileSize,
// 	}

// 	buildCoordinates := func(coordinates Coordinates) Coordinates {
// 		return Coordinates{
// 			X: tileSize * coordinates.X,
// 			Y: tileSize * coordinates.Y,
// 		}
// 	}

// 	return map[PresetName]EntityConfigPresetFn{
// 		// "square": func(coordinates Coordinates) EntityConfig {
// 		// 	return EntityConfig{
// 		// 		Category:    CategoryIgnore,
// 		// 		Dimensions:  dimensions,
// 		// 		Coordinates: buildCoordinates(coordinates),
// 		// 		// TODO draw a square!
// 		// 	}
// 		// },,,
// 		// this is an impassable obstacle that can be toggled "remotely"
// 		// it has two visual states that coincide with each toggle state,
// 	}
// }
