package zelduh

// GameModel contains data used throughout the game
type GameModel struct {
	AddEntities               bool
	CurrentRoomID, NextRoomID RoomID
	RoomTransition            *RoomTransition
	CurrentState              State
}
