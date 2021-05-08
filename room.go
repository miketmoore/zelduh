package zelduh

import (
	"fmt"

	"github.com/miketmoore/zelduh/core/direction"
	"github.com/miketmoore/zelduh/core/tmx"
)

// RoomID is a room ID
type RoomID int

// ConnectedRooms is used to configure adjacent rooms
type ConnectedRooms struct {
	Top    RoomID
	Right  RoomID
	Bottom RoomID
	Left   RoomID
}

func (c ConnectedRooms) String() string {
	return fmt.Sprintf("ConnectedRooms Top=%d Right=%d Bottom=%d Left=%d", c.Top, c.Right, c.Bottom, c.Left)
}

// RoomByIDMap is a type of map that indexes rooms by their ID
type RoomByIDMap map[RoomID]*Room

// Room represents one map section
type Room struct {
	ID             RoomID
	TMXFileName    tmx.TMXFileName
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
