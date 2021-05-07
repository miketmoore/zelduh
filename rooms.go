package zelduh

import "github.com/miketmoore/zelduh/core/direction"

// RoomID is a room ID
type RoomID int

type RoomName string

// ConnectedRooms is used to configure adjacent rooms
type ConnectedRooms struct {
	Top    RoomID
	Right  RoomID
	Bottom RoomID
	Left   RoomID
}

// RoomByIDMap is a type of map that indexes rooms by their ID
type RoomByIDMap map[RoomID]*Room

// Room represents one map section
type Room struct {
	TMXFileName    RoomName
	connectedRooms *ConnectedRooms
	EntityConfigs  []EntityConfig
}

// ConnectedRooms returns the room's map name
func (r Room) ConnectedRooms() *ConnectedRooms {
	return r.connectedRooms
}

// SetConnectedRoom sets the connected room IDs
func (r Room) SetConnectedRoom(dir direction.Direction, id RoomID) {
	switch dir {
	case direction.DirectionUp:
		r.connectedRooms.Top = id
	case direction.DirectionRight:
		r.connectedRooms.Right = id
	case direction.DirectionDown:
		r.connectedRooms.Bottom = id
	case direction.DirectionLeft:
		r.connectedRooms.Left = id
	}
}

// NewRoom builds a new Room
func NewRoom(tmxFileName RoomName, entityConfigs ...EntityConfig) *Room {
	return &Room{
		TMXFileName:    tmxFileName,
		connectedRooms: &ConnectedRooms{},
		EntityConfigs:  entityConfigs,
	}
}

// BuildMapRoomIDToRoom transforms a multi-dimensional array of RoomID values into a map of Room structs, indexed by RoomID
func BuildMapRoomIDToRoom(layout [][]RoomID, roomByIDMap RoomByIDMap) {

	for row := 0; row < len(layout); row++ {
		for col := 0; col < len(layout[row]); col++ {
			roomAID := layout[row][col]
			// Top
			if (row > 0) && (len(layout[row-1]) > col) {
				roomBID := layout[row-1][col]
				if roomBID > 0 {
					connectRooms(roomByIDMap, roomAID, roomBID, direction.DirectionUp)
				}
			}
			// Right
			if len(layout[row]) > col+1 {
				roomBID := layout[row][col+1]
				if roomBID > 0 {
					connectRooms(roomByIDMap, roomAID, roomBID, direction.DirectionRight)
				}
			}
			// Bottom
			if (len(layout) > row+1) && (len(layout[row+1]) > col) {
				roomBID := layout[row+1][col]
				if roomBID > 0 {
					connectRooms(roomByIDMap, roomAID, roomBID, direction.DirectionDown)
				}
			}
			// Left
			if col > 0 {
				roomBID := layout[row][col-1]
				if roomBID > 0 {
					connectRooms(roomByIDMap, roomAID, roomBID, direction.DirectionLeft)
				}
			}
		}
	}
}

func connectRooms(roomByIDMap RoomByIDMap, roomAID, roomBID RoomID, dir direction.Direction) {

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

	// Cache room references by their IDs
	roomByIDMap[roomAID] = roomA
	roomByIDMap[roomBID] = roomB

	// Connect rooms
	roomA.SetConnectedRoom(directionAToB, roomBID)
	roomB.SetConnectedRoom(directionBToA, roomAID)

}
