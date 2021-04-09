package zelduh

import (
	"math/rand"

	"github.com/faiface/pixel"
)

// GameModel contains data used throughout the game
type GameModel struct {
	AddEntities               bool
	CurrentRoomID, NextRoomID RoomID
	RoomTransition            *RoomTransition
	CurrentState              State
	Rand                      *rand.Rand
	EntitiesMap               map[EntityID]Entity
	Spritesheet               map[int]*pixel.Sprite
	Entities                  Entities
	RoomWarps                 map[EntityID]Config
	AllMapDrawData            map[string]MapData
	HealthSystem              *SystemHealth
	InputSystem               *SystemInput
	SpatialSystem             *SystemSpatial
}

type Entities struct {
	Player    Entity
	Bomb      Entity
	Explosion Entity
	Sword     Entity
	Arrow     Entity
	Hearts    []Entity
}
