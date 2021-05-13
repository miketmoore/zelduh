package zelduh

type Level struct {
	RoomByIDMap  RoomByIDMap
	RoomIdLayout [][]RoomID
}

type LevelManager struct {
	CurrentLevel *Level
}

func NewLevelManager(currentLevel *Level) *LevelManager {
	return &LevelManager{
		CurrentLevel: currentLevel,
	}
}

func (l *LevelManager) SetCurrentLevel(level *Level) {
	l.CurrentLevel = level
}
