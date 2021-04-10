package zelduh

// GameModel contains data used throughout the game
type GameModel struct {
	CurrentRoomID, NextRoomID RoomID
	RoomTransition            *RoomTransition
}
