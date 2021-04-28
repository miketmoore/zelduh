package main

import "github.com/miketmoore/zelduh"

func buildTestLevel(
	entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager,
	tileSize float64,
) zelduh.Level {
	// Build a map of RoomIDs to Room structs
	roomByIDMap := BuildRooms(entityConfigPresetFnManager, tileSize)

	zelduh.BuildMapRoomIDToRoom(
		// Overworld is a multi-dimensional array representing the overworld
		// Each room ID should be unique
		[][]zelduh.RoomID{
			{1, 10},
			{2, 0, 0, 8},
			{3, 5, 6, 7},
			{9},
			{11},
		},
		// This is mutated
		roomByIDMap,
	)

	return zelduh.Level{
		RoomByIDMap: roomByIDMap,
	}
}

// TODO move to zelduh cmd file since it is configuration
// Map of RoomID to a Room configuration
func BuildRooms(entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager, tileSize float64) zelduh.RoomByIDMap {

	buildWarpStone := buildWarpStoneFnFactory(entityConfigPresetFnManager)

	buildWarp := buildWarpFnFactory(tileSize, zelduh.Dimensions{
		Width:  tileSize,
		Height: tileSize,
	})

	return zelduh.RoomByIDMap{
		1: zelduh.NewRoom("overworldFourWallsDoorBottomRight"), // entityConfigPresetFnManager.GetPreset(PresetNameDialogCorner)(zelduh.Coordinates{X: 3, Y: 11}),
		// entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(zelduh.Coordinates{X: 10, Y: 10}),
		// entityConfigPresetFnManager.GetPreset(PresetNamePuzzleBox)(zelduh.Coordinates{X: 5, Y: 5}),
		// entityConfigPresetFnManager.GetPreset(PresetNameFloorSwitch)(zelduh.Coordinates{X: 5, Y: 6}),
		// entityConfigPresetFnManager.GetPreset(PresetNameToggleObstacle)(zelduh.Coordinates{X: 10, Y: 7}),

		2: zelduh.NewRoom("overworldFourWallsDoorTopBottom",
			entityConfigPresetFnManager.GetPreset(PresetNameEnemySkull)(zelduh.NewCoordinates(1, 1)),
			// entityConfigPresetFnManager.GetPreset(PresetNameEnemySkeleton)(zelduh.Coordinates{X: 11, Y: 9}),
			// entityConfigPresetFnManager.GetPreset(PresetNameEnemySpinner)(zelduh.Coordinates{X: 7, Y: 9}),
			// entityConfigPresetFnManager.GetPreset(PresetNameEnemyEyeBurrower)(zelduh.Coordinates{X: 8, Y: 9}),
		),
		3: zelduh.NewRoom("overworldFourWallsDoorRightTopBottom",
			buildWarpStone(6, zelduh.Coordinates{X: 3, Y: 7}, 5),
		),
		5: zelduh.NewRoom("rockWithCaveEntrance",
			buildWarp(
				11,
				zelduh.Coordinates{
					X: (tileSize * 7) + tileSize/2,
					Y: (tileSize * 9) + tileSize/2,
				},
				30,
			),
			buildWarp(
				11,
				zelduh.Coordinates{
					X: (tileSize * 8) + tileSize/2,
					Y: (tileSize * 9) + tileSize/2,
				},
				30,
			),
		),
		6:  zelduh.NewRoom("rockPathLeftRightEntrance"),
		7:  zelduh.NewRoom("overworldFourWallsDoorLeftTop"),
		8:  zelduh.NewRoom("overworldFourWallsDoorBottom"),
		9:  zelduh.NewRoom("overworldFourWallsDoorTop"),
		10: zelduh.NewRoom("overworldFourWallsDoorLeft"),
		11: zelduh.NewRoom("dungeonFourDoors",
			// South door of cave - warp to cave entrance
			buildWarp(
				5,
				zelduh.Coordinates{
					X: (tileSize * 6) + tileSize + (tileSize / 2.5),
					Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
				},
				15,
			),
			buildWarp(
				5,
				zelduh.Coordinates{
					X: (tileSize * 7) + tileSize + (tileSize / 2.5),
					Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
				},
				15,
			),
		),
	}
}
