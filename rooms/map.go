package rooms

import (
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/direction"
)

// RoomID is a room ID
type RoomID int

// Rooms is a type of map that indexes rooms by their ID
type Rooms map[RoomID]Room

// ConnectedRooms is used to configure adjacent rooms
type ConnectedRooms struct {
	Top    RoomID
	Right  RoomID
	Bottom RoomID
	Left   RoomID
}

// EntityConfig is used to simplify building entities
type EntityConfig struct {
	Category                 categories.Category
	X, Y, W, H, HitBoxRadius float64
	SpriteFrames             []int
	WarpToRoomID             RoomID
	Invincible               bool
	PatternName              string
	Direction                direction.Name
}

// Room represents one map section
type Room struct {
	MapName        string
	ConnectedRooms ConnectedRooms
	EntityConfigs  []EntityConfig
}

// NewRoom builds a new Room
func NewRoom(name string, entityConfigs ...EntityConfig) Room {
	return Room{
		MapName:       name,
		EntityConfigs: entityConfigs,
	}
}
