package zelduh

import (
	"fmt"
	"os"
	"strings"

	"github.com/miketmoore/mazegen"
	"github.com/miketmoore/zelduh/core/tmx"
)

func buildLevelMaze(
	roomFactory *RoomFactory,
	entityFactory *EntityFactory,
	tileSize float64,
) Level {
	fmt.Println("building level...")

	fmt.Println("building maze data...")
	rows := 3
	cols := 3
	random := mazegen.NewRandom()
	grid, err := mazegen.BuildMaze(rows, cols, random)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println("maze data built")

	tmxNameMap := map[string]tmx.TMXFileName{
		"nes_": "overworldFourWallsDoorLeft",
		"ne_w": "overworldFourWallsDoorBottom",
		"n_sw": "overworldFourWallsDoorRight",
		"_esw": "overworldFourWallsDoorTop",
		"ne__": "overworldFourWallsDoorBottomLeft",
		"n__w": "overworldFourWallsDoorBottomRight",
		"_es_": "overworldFourWallsDoorLeftTop",
		"__sw": "overworldFourWallsDoorRightTop",
		"_e__": "overworldFourWallsDoorTopBottomLeft",
		"n_s_": "overworldFourWallsDoorRightLeft",
		"n___": "overworldFourWallsDoorRightBottomLeft",
		"_e_w": "overworldFourWallsDoorTopBottom",
		"__s_": "overworldFourWallsDoorTopRightLeft",
	}

	roomByIdMap := RoomByIDMap{}

	buildKeyChar := func(cell *mazegen.Cell, direction mazegen.DirectionValue) string {
		if cell.Walls[direction] {
			return strings.Split(string(direction), "")[0]
		}
		return "_"
	}

	// layout := [][]RoomID{}

	var layout [][]RoomID = make([][]RoomID, grid.Rows)

	for rowIndex, row := range grid.Cells {
		// fmt.Printf("building row index=%d\n", rowIndex)

		// Make row of cells slice
		layout[rowIndex] = make([]RoomID, len(row))

		for columnIndex, cell := range row {
			// fmt.Println(rowIndex, columnIndex, cell)

			// Determine room TMX name
			key := fmt.Sprintf("%s%s%s%s",
				buildKeyChar(cell, mazegen.North),
				buildKeyChar(cell, mazegen.East),
				buildKeyChar(cell, mazegen.South),
				buildKeyChar(cell, mazegen.West),
			)
			fmt.Printf("key=%s\n", key)
			roomName := tmxNameMap[key]

			fmt.Printf("roomName=%s\n", roomName)

			// Build room
			room := roomFactory.NewRoom(roomName)

			// Index room by ID
			roomByIdMap[room.ID] = room

			// Insert room ID in multi-dimensional slice representing actual map layout
			layout[rowIndex][columnIndex] = room.ID
		}
	}

	// // top left
	// room1 := roomFactory.NewRoom("overworldFourWallsDoorBottomRight",
	// 	entityFactory.PresetSkeleton()(NewCoordinates(4, 4)),
	// 	entityFactory.PresetSkull()(NewCoordinates(5, 4)),
	// 	entityFactory.PresetEyeBurrower()(NewCoordinates(6, 4)),
	// )

	// // top right
	// room2 := roomFactory.NewRoom("overworldFourWallsDoorBottomLeft")

	// // bottom left
	// room3 := roomFactory.NewRoom("overworldFourWallsDoorRightTop")

	// // bottom right
	// room4 := roomFactory.NewRoom("overworldFourWallsDoorLeftTop")

	// rooms := []*Room{
	// 	room1,
	// 	room2,
	// 	room3,
	// 	room4,
	// }

	// // Build a map of RoomIDs to Room structs
	// roomByIDMap := RoomByIDMap{}
	// for _, room := range rooms {
	// 	roomByIDMap[room.ID] = room
	// }

	// fmt.Printf("room by ID map: %\n", roomByIdMap)
	// fmt.Println(roomByIdMap)
	for rowIndex, row := range grid.Cells {
		for columnIndex := range row {
			fmt.Printf("%d ", layout[rowIndex][columnIndex])
		}
		fmt.Println("")
	}

	connectRooms(
		// Overworld is a multi-dimensional array representing the overworld
		// Each room ID should be unique
		layout,
		// This is mutated
		roomByIdMap,
	)

	return Level{
		RoomByIDMap:  roomByIdMap,
		RoomIdLayout: layout,
	}
}
