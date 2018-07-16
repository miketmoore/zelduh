package gamemap

import (
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh/config"
	"github.com/miketmoore/zelduh/rooms"
)

func indexRoom(roomsMap rooms.Rooms, a, b terraform2d.RoomID, dir terraform2d.Direction) {
	// fmt.Printf("indexRoom a:%d b:%d dir:%s\n", a, b, dir)
	roomA, okA := roomsMap[a]
	roomB, okB := roomsMap[b]
	if okA && okB {
		switch dir {
		case terraform2d.DirectionUp:
			// b is above a
			roomA.ConnectedRooms.Top = b
			roomsMap[a] = roomA
			roomB.ConnectedRooms.Bottom = a
			roomsMap[b] = roomB
		case terraform2d.DirectionRight:
			// b is right of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.ConnectedRooms.Right = b
				roomsMap[a] = roomA
				roomB.ConnectedRooms.Left = a
				roomsMap[b] = roomB
			}
		case terraform2d.DirectionDown:
			// b is below a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.ConnectedRooms.Bottom = b
				roomsMap[a] = roomA
				roomB.ConnectedRooms.Top = a
				roomsMap[b] = roomB
			}
		case terraform2d.DirectionLeft:
			// b is left of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.ConnectedRooms.Left = b
				roomsMap[a] = roomA
				roomB.ConnectedRooms.Right = a
				roomsMap[b] = roomB
			}
		}
	}

}

// ProcessMapLayout processes the game maps
func ProcessMapLayout(roomsMap rooms.Rooms) {
	layout := config.Overworld
	// transform multi-dimensional array into map of Room structs, indexed by ID
	for row := 0; row < len(layout); row++ {
		for col := 0; col < len(layout[row]); col++ {
			roomID := layout[row][col]
			// fmt.Printf("Room ID: %d\n", roomID)
			// Top
			if row > 0 {
				if len(layout[row-1]) > col {
					n := layout[row-1][col]
					if n > 0 {
						// fmt.Printf("\t%d is below %d\n", roomID, n)
						indexRoom(roomsMap, roomID, n, terraform2d.DirectionUp)
					}
				}
			}
			// Right
			if len(layout[row]) > col+1 {
				n := layout[row][col+1]
				if n > 0 {
					// fmt.Printf("\t%d is left of %d\n", roomID, n)
					indexRoom(roomsMap, roomID, n, terraform2d.DirectionRight)
				}
			}
			// Bottom
			if len(layout) > row+1 {
				if len(layout[row+1]) > col {
					n := layout[row+1][col]
					if n > 0 {
						// fmt.Printf("\t%d is above %d\n", roomID, n)
						indexRoom(roomsMap, roomID, n, terraform2d.DirectionDown)
					}
				}
			}
			// Left
			if col > 0 {
				n := layout[row][col-1]
				if n > 0 {
					// fmt.Printf("\t%d is right of %d\n", roomID, n)
					indexRoom(roomsMap, roomID, n, terraform2d.DirectionLeft)
				}
			}
		}
	}
}
