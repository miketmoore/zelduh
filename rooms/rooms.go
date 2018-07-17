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
