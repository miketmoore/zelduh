package zelduh

import (
	"github.com/faiface/pixel"
)

// GameModel contains data used throughout the game
type GameModel struct {
	CurrentRoomID, NextRoomID RoomID
	RoomTransition            *RoomTransition
	CurrentState              State
	EntitiesMap               map[EntityID]Entity
	Spritesheet               map[int]*pixel.Sprite
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
