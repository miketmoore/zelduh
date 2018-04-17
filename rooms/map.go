package rooms

import (
	"github.com/faiface/pixel/imdraw"
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

// AnimationConfig is a map of animation types to sprite index lists
type AnimationConfig map[string][]int

// MovementConfig is used to configure an entity's Movement component
type MovementConfig struct {
	LastDirection  direction.Name
	Direction      direction.Name
	MaxSpeed       float64
	Speed          float64
	MaxMoves       int
	RemainingMoves int
	HitSpeed       float64
	MovingFromHit  bool
	HitBackMoves   int
	PatternName    string
}

// HitboxConfig is used to configure an entity's hitbox
type HitboxConfig struct {
	Box    *imdraw.IMDraw
	Radius float64
}

// EntityConfig is used to simplify building entities
type EntityConfig struct {
	Category                                               categories.Category
	Moveable, Animated, Toggleable, Impassable, Invincible bool
	X, Y, W, H                                             float64
	SpriteFrames                                           []int
	WarpToRoomID                                           RoomID
	PatternName                                            string
	Toggled                                                bool
	Animation                                              AnimationConfig
	Hitbox                                                 *HitboxConfig
	Health                                                 int
	Expiration                                             int
	Ignore                                                 bool

	Movement *MovementConfig
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
