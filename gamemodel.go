package zelduh

// GameModel contains data used throughout the game
type GameModel struct {
	RoomTransition *RoomTransition
	EntitiesMap    map[EntityID]Entity
	Entities       Entities
	RoomWarps      map[EntityID]EntityConfig
}

type Entities struct {
	Player    Entity
	Bomb      Entity
	Explosion Entity
	Sword     Entity
	Arrow     Entity
	Hearts    []Entity
}
