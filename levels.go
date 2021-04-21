package zelduh

type Level struct {
	RoomByIDMap RoomByIDMap
}

type LevelManager struct {
	CurrentLevel *Level
}

func NewLevelManager(currentLevel *Level) LevelManager {
	return LevelManager{
		CurrentLevel: currentLevel,
	}
}
