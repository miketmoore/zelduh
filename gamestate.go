package zelduh

import "github.com/faiface/pixel"

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
	SystemsManager  *SystemsManager
	UI              UI
	LocaleMessages  LocaleMessagesMap
	CollisionSystem *SystemCollision
	InputSystem     *SystemInput
	Spritesheet     map[int]*pixel.Sprite
	EntitiesMap     EntityByEntityID
	AllMapDrawData  map[string]MapData
	RoomWarps       map[EntityID]Config
}

func NewGameStateManager(
	gameModel *GameModel,
	systemsManager *SystemsManager,
	ui UI,
	localeMessages LocaleMessagesMap,
	collisionSystem *SystemCollision,
	inputSystem *SystemInput,
	spritesheet map[int]*pixel.Sprite,
	entitiesMap EntityByEntityID,
	allMapDrawData map[string]MapData,
	roomWarps map[EntityID]Config,
) GameStateManager {
	return GameStateManager{
		GameModel:       gameModel,
		SystemsManager:  systemsManager,
		UI:              ui,
		LocaleMessages:  localeMessages,
		CollisionSystem: collisionSystem,
		InputSystem:     inputSystem,
		Spritesheet:     spritesheet,
		EntitiesMap:     entitiesMap,
		AllMapDrawData:  allMapDrawData,
		RoomWarps:       roomWarps,
	}
}

func (g *GameStateManager) Update() {
	switch g.GameModel.CurrentState {
	case StateStart:
		GameStateStart(g.UI, g.LocaleMessages, g.GameModel)
	case StateGame:
		GameStateGame(g.UI, g.Spritesheet, g.GameModel, g.RoomWarps, g.EntitiesMap, g.AllMapDrawData, g.InputSystem, RoomsMap, g.SystemsManager)
	case StatePause:
		GameStatePause(g.UI, g.LocaleMessages, g.GameModel)
	case StateOver:
		GameStateOver(g.UI, g.LocaleMessages, g.GameModel)
	case StateMapTransition:
		GameStateMapTransition(g.UI, g.Spritesheet, g.AllMapDrawData, g.InputSystem, g.SystemsManager, RoomsMap, g.CollisionSystem, g.GameModel)
	}
}
