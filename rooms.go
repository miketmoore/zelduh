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

// Rooms is a type of map that indexes rooms by their ID
type Rooms map[RoomID]*Room

// Room represents one map section
type Room struct {
	mapName        RoomName
	connectedRooms *ConnectedRooms
	EntityConfigs  []EntityConfig
}

// RoomName returns the room's map name
func (r Room) RoomName() RoomName {
	return r.mapName
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
		mapName:        name,
		connectedRooms: &ConnectedRooms{},
		EntityConfigs:  entityConfigs,
	}
}

// BuildMapRoomIDToRoom transforms a multi-dimensional array of RoomID values into a map of Room structs, indexed by RoomID
func BuildMapRoomIDToRoom(layout [][]RoomID, roomsMap Rooms) {
	for row := 0; row < len(layout); row++ {
		for col := 0; col < len(layout[row]); col++ {
			roomID := layout[row][col]
			// Top
			if row > 0 {
				if len(layout[row-1]) > col {
					n := layout[row-1][col]
					if n > 0 {
						connectRooms(roomsMap, roomID, n, DirectionUp)
					}
				}
			}
			// Right
			if len(layout[row]) > col+1 {
				n := layout[row][col+1]
				if n > 0 {
					connectRooms(roomsMap, roomID, n, DirectionRight)
				}
			}
			// Bottom
			if len(layout) > row+1 {
				if len(layout[row+1]) > col {
					n := layout[row+1][col]
					if n > 0 {
						connectRooms(roomsMap, roomID, n, DirectionDown)
					}
				}
			}
			// Left
			if col > 0 {
				n := layout[row][col-1]
				if n > 0 {
					connectRooms(roomsMap, roomID, n, DirectionLeft)
				}
			}
		}
	}
}

func connectRooms(roomsMap Rooms, a, b RoomID, dir Direction) {
	roomA, okA := roomsMap[a]
	roomB, okB := roomsMap[b]

	doConnect := false
	var directionAToB Direction
	var directionBToA Direction

	if okA && okB {
		doConnect = true
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
	}

	if doConnect {
		roomA.SetConnectedRoom(directionAToB, b)
		roomsMap[a] = roomA
		roomB.SetConnectedRoom(directionBToA, a)
		roomsMap[b] = roomB
	}

}
