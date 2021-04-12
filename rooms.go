package zelduh

// RoomID is a room ID
type RoomID int

// ConnectedRooms is used to configure adjacent rooms
type ConnectedRooms struct {
	Top    RoomID
	Right  RoomID
	Bottom RoomID
	Left   RoomID
}

// Room defines an API for room implementations
type Roomer interface {
	MapName() string
	ConnectedRooms() *ConnectedRooms
	SetConnectedRoom(Direction, RoomID)
}

// Rooms is a type of map that indexes rooms by their ID
type Rooms map[RoomID]Roomer

// Room represents one map section
type Room struct {
	mapName        string
	connectedRooms *ConnectedRooms
	EntityConfigs  []EntityConfig
}

// MapName returns the room's map name
func (r Room) MapName() string {
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
func NewRoom(name string, entityConfigs ...EntityConfig) *Room {
	return &Room{
		mapName:        name,
		connectedRooms: &ConnectedRooms{},
		EntityConfigs:  entityConfigs,
	}
}

func indexRoom(roomsMap Rooms, a, b RoomID, dir Direction) {
	// fmt.Printf("indexRoom a:%d b:%d dir:%s\n", a, b, dir)
	roomA, okA := roomsMap[a]
	roomB, okB := roomsMap[b]
	if okA && okB {
		switch dir {
		case DirectionUp:
			// b is above a
			roomA.SetConnectedRoom(DirectionUp, b)
			roomsMap[a] = roomA
			roomB.SetConnectedRoom(DirectionDown, a)
			roomsMap[b] = roomB
		case DirectionRight:
			// b is right of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.SetConnectedRoom(DirectionRight, b)
				roomsMap[a] = roomA
				roomB.SetConnectedRoom(DirectionLeft, a)
				roomsMap[b] = roomB
			}
		case DirectionDown:
			// b is below a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.SetConnectedRoom(DirectionDown, b)
				roomsMap[a] = roomA
				roomB.SetConnectedRoom(DirectionUp, a)
				roomsMap[b] = roomB
			}
		case DirectionLeft:
			// b is left of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.SetConnectedRoom(DirectionLeft, b)
				roomsMap[a] = roomA
				roomB.SetConnectedRoom(DirectionRight, a)
				roomsMap[b] = roomB
			}
		}
	}

}

// BuildMapRoomIDToRoom transforms a multi-dimensional array of RoomID values into a map of Room structs, indexed by RoomID
func BuildMapRoomIDToRoom(layout [][]RoomID, roomsMap Rooms) {
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
						indexRoom(roomsMap, roomID, n, DirectionUp)
					}
				}
			}
			// Right
			if len(layout[row]) > col+1 {
				n := layout[row][col+1]
				if n > 0 {
					// fmt.Printf("\t%d is left of %d\n", roomID, n)
					indexRoom(roomsMap, roomID, n, DirectionRight)
				}
			}
			// Bottom
			if len(layout) > row+1 {
				if len(layout[row+1]) > col {
					n := layout[row+1][col]
					if n > 0 {
						// fmt.Printf("\t%d is above %d\n", roomID, n)
						indexRoom(roomsMap, roomID, n, DirectionDown)
					}
				}
			}
			// Left
			if col > 0 {
				n := layout[row][col-1]
				if n > 0 {
					// fmt.Printf("\t%d is right of %d\n", roomID, n)
					indexRoom(roomsMap, roomID, n, DirectionLeft)
				}
			}
		}
	}
}
