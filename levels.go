package zelduh

type Level struct {
	Map Rooms
}

type LevelManager struct {
	CurrentLevel *Level
}

func NewLevelManager(currentLevel *Level) LevelManager {
	return LevelManager{
		CurrentLevel: currentLevel,
	}
}
