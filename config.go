package zelduh

type WindowConfig struct {
	X, Y, Width, Height float64
}

type MapConfig struct {
	X, Y, Width, Height float64
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
