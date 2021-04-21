package zelduh

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
	Name           RoomName
	connectedRooms *ConnectedRooms
	EntityConfigs  []EntityConfig
}

// ConnectedRooms returns the room's map name
func (r Room) ConnectedRooms() *ConnectedRooms {
	return r.connectedRooms
}

// SetConnectedRoom sets the connected room IDs
func (r Room) SetConnectedRoom(direction Direction, id RoomID) {
	switch direction {
	case DirectionUp:
		r.connectedRooms.Top = id
	case DirectionRight:
		r.connectedRooms.Right = id
	case DirectionDown:
		r.connectedRooms.Bottom = id
	case DirectionLeft:
		r.connectedRooms.Left = id
	}
}

// NewRoom builds a new Room
func NewRoom(name RoomName, entityConfigs ...EntityConfig) *Room {
	return &Room{
		Name:           name,
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
					connectRooms(roomByIDMap, roomAID, roomBID, DirectionUp)
				}
			}
			// Right
			if len(layout[row]) > col+1 {
				roomBID := layout[row][col+1]
				if roomBID > 0 {
					connectRooms(roomByIDMap, roomAID, roomBID, DirectionRight)
				}
			}
			// Bottom
			if (len(layout) > row+1) && (len(layout[row+1]) > col) {
				roomBID := layout[row+1][col]
				if roomBID > 0 {
					connectRooms(roomByIDMap, roomAID, roomBID, DirectionDown)
				}
			}
			// Left
			if col > 0 {
				roomBID := layout[row][col-1]
				if roomBID > 0 {
					connectRooms(roomByIDMap, roomAID, roomBID, DirectionLeft)
				}
			}
		}
	}
}

func connectRooms(roomByIDMap RoomByIDMap, roomAID, roomBID RoomID, dir Direction) {

	roomA, roomAOK := roomByIDMap[roomAID]
	roomB, roomBOK := roomByIDMap[roomBID]

	if !roomAOK || !roomBOK {
		return
	}

	var directionAToB Direction
	var directionBToA Direction

	switch dir {
	case DirectionUp:
		// b is above a
		directionAToB = DirectionUp
		directionBToA = DirectionDown
	case DirectionRight:
		// b is right of a
		directionAToB = DirectionRight
		directionBToA = DirectionLeft
	case DirectionDown:
		// b is below a
		directionAToB = DirectionDown
		directionBToA = DirectionUp
	case DirectionLeft:
		// b is left of a
		directionAToB = DirectionLeft
		directionBToA = DirectionRight
	}

	// Cache room references by their IDs
	roomByIDMap[roomAID] = roomA
	roomByIDMap[roomBID] = roomB

	// Connect rooms
	roomA.SetConnectedRoom(directionAToB, roomBID)
	roomB.SetConnectedRoom(directionBToA, roomAID)

}
