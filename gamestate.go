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
	GameModel        *GameModel
	SystemsManager   *SystemsManager
	GameStateManager *GameStateManager
	UI               UI
	LocaleMessages   LocaleMessagesMap
	CollisionSystem  *SystemCollision
	InputSystem      *SystemInput
	Spritesheet      map[int]*pixel.Sprite
	EntitiesMap      EntityByEntityID
	AllMapDrawData   map[string]MapData
	RoomWarps        map[EntityID]Config
	Entities         Entities
	CurrentState     State
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
	entities Entities,
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
		Entities:        entities,
		CurrentState:    StateStart,
	}
}

func (g *GameStateManager) Update() {
	switch g.CurrentState {
	case StateStart:
		GameStateStart(g.UI, g.LocaleMessages, g.GameModel, g.GameStateManager)
	case StateGame:
		GameStateGame(g.UI, g.Spritesheet, g.GameModel, g.Entities, g.RoomWarps, g.EntitiesMap, g.AllMapDrawData, g.InputSystem, RoomsMap, g.SystemsManager, g.GameStateManager)
	case StatePause:
		GameStatePause(g.UI, g.LocaleMessages, g.GameModel, g.GameStateManager)
	case StateOver:
		GameStateOver(g.UI, g.LocaleMessages, g.GameModel, g.GameStateManager)
	case StateMapTransition:
		GameStateMapTransition(g.UI, g.Spritesheet, g.Entities, g.AllMapDrawData, g.InputSystem, g.SystemsManager, RoomsMap, g.CollisionSystem, g.GameModel, g.GameStateManager)
	}
}
