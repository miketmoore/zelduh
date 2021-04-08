package zelduh

// State is a type of game state
type State string

const (
	StateStart         State = "start"
	StateGame          State = "game"
	StatePause         State = "pause"
	StateOver          State = "over"
	StateMapTransition State = "mapTransition"
)

type GameStateManager struct {
	GameModel       *GameModel
	GameWorld       *World
	UI              UI
	LocaleMessages  map[string]string
	CollisionSystem *SystemCollision
}

func NewGameStateManager(
	gameModel *GameModel,
	gameWorld *World,
	ui UI,
	localeMessages map[string]string,
	collisionSystem *SystemCollision,
) GameStateManager {
	return GameStateManager{
		GameModel:       gameModel,
		GameWorld:       gameWorld,
		UI:              ui,
		LocaleMessages:  localeMessages,
		CollisionSystem: collisionSystem,
	}
}

func (g *GameStateManager) Update() {
	switch g.GameModel.CurrentState {
	case StateStart:
		GameStateStart(g.UI, g.LocaleMessages, g.GameModel)
	case StateGame:
		GameStateGame(g.UI, g.GameModel, RoomsMap, g.GameWorld)
	case StatePause:
		GameStatePause(g.UI, g.LocaleMessages, g.GameModel)
	case StateOver:
		GameStateOver(g.UI, g.LocaleMessages, g.GameModel)
	case StateMapTransition:
		GameStateMapTransition(g.UI, g.GameWorld, RoomsMap, g.CollisionSystem, g.GameModel)
	}
}
