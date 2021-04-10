package zelduh

import (
	"math/rand"
)

// GameModel contains data used throughout the game
type GameModel struct {
	AddEntities               bool
	CurrentRoomID, NextRoomID RoomID
	RoomTransition            *RoomTransition
	CurrentState              State
	Rand                      *rand.Rand
	Entities                  Entities
	RoomWarps                 map[EntityID]Config
	AllMapDrawData            map[string]MapData
}

type Entities struct {
	Player    Entity
	Bomb      Entity
	Explosion Entity
	Sword     Entity
	Arrow     Entity
	Hearts    []Entity
}
