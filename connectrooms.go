package zelduh

import "github.com/miketmoore/zelduh/core/direction"

// connectRooms traverses a multi-dimensional slice of room IDs and "connects" adjacent rooms
// by calling Room.SetConnectedRoom(...) on each room
func connectRooms(layout [][]RoomID, roomByIDMap RoomByIDMap) {

	for row := 0; row < len(layout); row++ {
		for col := 0; col < len(layout[row]); col++ {
			roomAID := layout[row][col]
			// Top
			if (row > 0) && (len(layout[row-1]) > col) {
				roomBID := layout[row-1][col]
				if roomBID > 0 {
					connectTwoRooms(roomByIDMap, roomAID, roomBID, direction.DirectionUp)
				}
			}
			// Right
			if len(layout[row]) > col+1 {
				roomBID := layout[row][col+1]
				if roomBID > 0 {
					connectTwoRooms(roomByIDMap, roomAID, roomBID, direction.DirectionRight)
				}
			}
			// Bottom
			if (len(layout) > row+1) && (len(layout[row+1]) > col) {
				roomBID := layout[row+1][col]
				if roomBID > 0 {
					connectTwoRooms(roomByIDMap, roomAID, roomBID, direction.DirectionDown)
				}
			}
			// Left
			if col > 0 {
				roomBID := layout[row][col-1]
				if roomBID > 0 {
					connectTwoRooms(roomByIDMap, roomAID, roomBID, direction.DirectionLeft)
				}
			}
		}
	}
}

func connectTwoRooms(roomByIDMap RoomByIDMap, roomAID, roomBID RoomID, dir direction.Direction) {

	roomA, roomAOK := roomByIDMap[roomAID]
	roomB, roomBOK := roomByIDMap[roomBID]

	if !roomAOK || !roomBOK {
		return
	}

	var directionAToB direction.Direction
	var directionBToA direction.Direction

	switch dir {
	case direction.DirectionUp:
		// b is above a
		directionAToB = direction.DirectionUp
		directionBToA = direction.DirectionDown
	case direction.DirectionRight:
		// b is right of a
		directionAToB = direction.DirectionRight
		directionBToA = direction.DirectionLeft
	case direction.DirectionDown:
		// b is below a
		directionAToB = direction.DirectionDown
		directionBToA = direction.DirectionUp
	case direction.DirectionLeft:
		// b is left of a
		directionAToB = direction.DirectionLeft
		directionBToA = direction.DirectionRight
	}

	// Connect rooms
	roomA.SetConnectedRoom(directionAToB, roomBID)
	roomB.SetConnectedRoom(directionBToA, roomAID)

}
