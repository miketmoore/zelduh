package zelduh

func buildTestLevel01(
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
) Level {
	// Build a map of RoomIDs to Room structs
	roomByIDMap := BuildRooms01(entityConfigPresetFnManager, tileSize)

	BuildMapRoomIDToRoom(
		// Overworld is a multi-dimensional array representing the overworld
		// Each room ID should be unique
		[][]RoomID{
			{1},
		},
		// This is mutated
		roomByIDMap,
	)

	return Level{
		RoomByIDMap: roomByIDMap,
	}
}

// TODO move to zelduh cmd file since it is configuration
// Map of RoomID to a Room configuration
func BuildRooms01(entityConfigPresetFnManager *EntityConfigPresetFnManager, tileSize float64) RoomByIDMap {

	// buildWarpStone := buildWarpStoneFnFactory(entityConfigPresetFnManager)

	// buildWarp := buildWarpFnFactory(tileSize, Dimensions{
	// 	Width:  tileSize,
	// 	Height: tileSize,
	// })

	return RoomByIDMap{
		1: NewRoom("overworldFourWallsDoorBottomRight"), // entityConfigPresetFnManager.GetPreset(PresetNameDialogCorner)(Coordinates{X: 3, Y: 11}),
		// entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(Coordinates{X: 10, Y: 10}),
		// entityConfigPresetFnManager.GetPreset(PresetNamePuzzleBox)(Coordinates{X: 5, Y: 5}),
		// entityConfigPresetFnManager.GetPreset(PresetNameFloorSwitch)(Coordinates{X: 5, Y: 6}),
		// entityConfigPresetFnManager.GetPreset(PresetNameToggleObstacle)(Coordinates{X: 10, Y: 7}),

		// 2: NewRoom("overworldFourWallsDoorTopBottom",
		// 	entityConfigPresetFnManager.GetPreset(PresetNameEnemySkull)(NewCoordinates(6, 6)),
		// 	// entityConfigPresetFnManager.GetPreset(PresetNameEnemySkeleton)(Coordinates{X: 11, Y: 9}),
		// 	// entityConfigPresetFnManager.GetPreset(PresetNameEnemySpinner)(Coordinates{X: 7, Y: 9}),
		// 	// entityConfigPresetFnManager.GetPreset(PresetNameEnemyEyeBurrower)(Coordinates{X: 8, Y: 9}),
		// ),
		// 3: NewRoom("overworldFourWallsDoorRightTopBottom",
		// 	buildWarpStone(6, Coordinates{X: 3, Y: 7}, 5),
		// ),
		// 5: NewRoom("rockWithCaveEntrance",
		// 	buildWarp(
		// 		11,
		// 		Coordinates{
		// 			X: (tileSize * 7) + tileSize/2,
		// 			Y: (tileSize * 9) + tileSize/2,
		// 		},
		// 		30,
		// 	),
		// 	buildWarp(
		// 		11,
		// 		Coordinates{
		// 			X: (tileSize * 8) + tileSize/2,
		// 			Y: (tileSize * 9) + tileSize/2,
		// 		},
		// 		30,
		// 	),
		// ),
		// 6:  NewRoom("rockPathLeftRightEntrance"),
		// 7:  NewRoom("overworldFourWallsDoorLeftTop"),
		// 8:  NewRoom("overworldFourWallsDoorBottom"),
		// 9:  NewRoom("overworldFourWallsDoorTop"),
		// 10: NewRoom("overworldFourWallsDoorLeft"),
		// 11: NewRoom("dungeonFourDoors",
		// 	// South door of cave - warp to cave entrance
		// 	buildWarp(
		// 		5,
		// 		Coordinates{
		// 			X: (tileSize * 6) + tileSize + (tileSize / 2.5),
		// 			Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
		// 		},
		// 		15,
		// 	),
		// 	buildWarp(
		// 		5,
		// 		Coordinates{
		// 			X: (tileSize * 7) + tileSize + (tileSize / 2.5),
		// 			Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
		// 		},
		// 		15,
		// 	),
		// ),
	}
}
