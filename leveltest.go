package zelduh

func buildTestLevel(
	roomFactory *RoomFactory,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
) Level {

	buildWarpStone := buildWarpStoneFnFactory(entityConfigPresetFnManager)

	buildWarp := buildWarpFnFactory(tileSize, Dimensions{
		Width:  tileSize,
		Height: tileSize,
	})

	room1 := roomFactory.NewRoom("overworldFourWallsDoorBottomRight")
	// entityConfigPresetFnManager.GetPreset(PresetNameDialogCorner)(Coordinates{X: 3, Y: 11}),
	// entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(Coordinates{X: 10, Y: 10}),
	// entityConfigPresetFnManager.GetPreset(PresetNamePuzzleBox)(Coordinates{X: 5, Y: 5}),
	// entityConfigPresetFnManager.GetPreset(PresetNameFloorSwitch)(Coordinates{X: 5, Y: 6}),
	// entityConfigPresetFnManager.GetPreset(PresetNameToggleObstacle)(Coordinates{X: 10, Y: 7}),

	room2 := roomFactory.NewRoom("overworldFourWallsDoorTopBottom")
	// entityConfigPresetFnManager.GetPreset(PresetNameEnemySkeleton)(Coordinates{X: 11, Y: 9}),
	// entityConfigPresetFnManager.GetPreset(PresetNameEnemySpinner)(Coordinates{X: 7, Y: 9}),
	// entityConfigPresetFnManager.GetPreset(PresetNameEnemyEyeBurrower)(Coordinates{X: 8, Y: 9}),

	room3 := roomFactory.NewRoom("overworldFourWallsDoorRightTopBottom",
		buildWarpStone(6, Coordinates{X: 3, Y: 7}, 5),
	)

	room5 := roomFactory.NewRoom("rockWithCaveEntrance",
		buildWarp(
			11,
			Coordinates{
				X: (tileSize * 7) + tileSize/2,
				Y: (tileSize * 9) + tileSize/2,
			},
			30,
		),
		buildWarp(
			11,
			Coordinates{
				X: (tileSize * 8) + tileSize/2,
				Y: (tileSize * 9) + tileSize/2,
			},
			30,
		),
	)

	room6 := roomFactory.NewRoom("rockPathLeftRightEntrance")
	room7 := roomFactory.NewRoom("overworldFourWallsDoorLeftTop")
	room8 := roomFactory.NewRoom("overworldFourWallsDoorBottom")
	room9 := roomFactory.NewRoom("overworldFourWallsDoorTop")
	room10 := roomFactory.NewRoom("overworldFourWallsDoorLeft")

	room11 := roomFactory.NewRoom("dungeonFourDoors",
		// South door of cave - warp to cave entrance
		buildWarp(
			5,
			Coordinates{
				X: (tileSize * 6) + tileSize + (tileSize / 2.5),
				Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
			},
			15,
		),
		buildWarp(
			5,
			Coordinates{
				X: (tileSize * 7) + tileSize + (tileSize / 2.5),
				Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
			},
			15,
		),
	)

	rooms := []*Room{
		room1,
		room2,
		room3,
		room5,
		room6,
		room7,
		room8,
		room9,
		room10,
		room11,
	}

	// Build a map of RoomIDs to Room structs
	roomByIDMap := RoomByIDMap{}
	for _, room := range rooms {
		roomByIDMap[room.ID] = room
	}

	BuildMapRoomIDToRoom(
		// Overworld is a multi-dimensional array representing the overworld
		// Each room ID should be unique
		[][]RoomID{
			{room1.ID, room10.ID},
			{room2.ID, 0, 0, room8.ID},
			{room3.ID, room5.ID, room6.ID, room7.ID},
			{room9.ID},
			{room11.ID},
		},
		// This is mutated
		roomByIDMap,
	)

	return Level{
		RoomByIDMap: roomByIDMap,
	}
}
