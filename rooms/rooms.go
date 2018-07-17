package rooms

import (
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh/entityconfig"
)

// Rooms is a type of map that indexes rooms by their ID
type Rooms map[terraform2d.RoomID]terraform2d.Room

// Room represents one map section
type Room struct {
	mapName        string
	connectedRooms *terraform2d.ConnectedRooms
	EntityConfigs  []entityconfig.Config
}

// MapName returns the room's map name
func (r Room) MapName() string {
	return r.mapName
}

// ConnectedRooms returns the room's map name
func (r Room) ConnectedRooms() *terraform2d.ConnectedRooms {
	return r.connectedRooms
}

// SetConnectedRoom sets the connected room IDs
func (r Room) SetConnectedRoom(direction terraform2d.Direction, id terraform2d.RoomID) {
	switch direction {
	case terraform2d.DirectionUp:
		r.connectedRooms.Top = id
	case terraform2d.DirectionRight:
		r.connectedRooms.Right = id
	case terraform2d.DirectionDown:
		r.connectedRooms.Bottom = id
	case terraform2d.DirectionLeft:
		r.connectedRooms.Left = id
	}
}

// NewRoom builds a new Room
func NewRoom(name string, entityConfigs ...entityconfig.Config) *Room {
	return &Room{
		mapName:        name,
		connectedRooms: &terraform2d.ConnectedRooms{},
		EntityConfigs:  entityConfigs,
	}
}

func indexRoom(roomsMap Rooms, a, b terraform2d.RoomID, dir terraform2d.Direction) {
	// fmt.Printf("indexRoom a:%d b:%d dir:%s\n", a, b, dir)
	roomA, okA := roomsMap[a]
	roomB, okB := roomsMap[b]
	if okA && okB {
		switch dir {
		case terraform2d.DirectionUp:
			// b is above a
			roomA.SetConnectedRoom(terraform2d.DirectionUp, b)
			roomsMap[a] = roomA
			roomB.SetConnectedRoom(terraform2d.DirectionDown, a)
			roomsMap[b] = roomB
		case terraform2d.DirectionRight:
			// b is right of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.SetConnectedRoom(terraform2d.DirectionRight, b)
				roomsMap[a] = roomA
				roomA.SetConnectedRoom(terraform2d.DirectionLeft, a)
				roomsMap[b] = roomB
			}
		case terraform2d.DirectionDown:
			// b is below a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.SetConnectedRoom(terraform2d.DirectionDown, b)
				roomsMap[a] = roomA
				roomA.SetConnectedRoom(terraform2d.DirectionUp, a)
				roomsMap[b] = roomB
			}
		case terraform2d.DirectionLeft:
			// b is left of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.SetConnectedRoom(terraform2d.DirectionLeft, b)
				roomsMap[a] = roomA
				roomA.SetConnectedRoom(terraform2d.DirectionRight, a)
				roomsMap[b] = roomB
			}
		}
	}

}

// ProcessMapLayout processes the game maps
func ProcessMapLayout(layout [][]terraform2d.RoomID, roomsMap Rooms) {
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
