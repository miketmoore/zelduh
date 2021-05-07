package zelduh

func buildTestLevel01(
	roomFactory *RoomFactory,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
) Level {

	room1 := roomFactory.NewRoom("overworldFourWallsDoorBottomRight")

	rooms := []*Room{
		room1,
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
			{room1.ID},
		},
		// This is mutated
		roomByIDMap,
	)

	return Level{
		RoomByIDMap: roomByIDMap,
	}
}
