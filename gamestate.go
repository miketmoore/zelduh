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
	GameModel                 *GameModel
	SystemsManager            *SystemsManager
	UI                        UI
	LocaleMessages            LocaleMessagesMap
	CollisionSystem           *CollisionSystem
	InputSystem               *InputSystem
	ShouldAddEntities         *bool
	CurrentRoomID, NextRoomID *RoomID
	CurrentState              *State
	Spritesheet               Spritesheet
	MapDrawData               MapDrawData
	RoomTransition            *RoomTransition
	EntitiesMap               EntitiesMap
	Player                    *Entity
	Hearts                    []Entity
}

func NewGameStateManager(
	gameModel *GameModel,
	systemsManager *SystemsManager,
	ui UI,
	localeMessages LocaleMessagesMap,
	collisionSystem *CollisionSystem,
	inputSystem *InputSystem,
	shouldAddEntities *bool,
	currentRoomID *RoomID,
	nextRoomID *RoomID,
	currentState *State,
	spritesheet Spritesheet,
	mapDrawData MapDrawData,
	roomTransition *RoomTransition,
	entitiesMap EntitiesMap,
	player *Entity,
	hearts []Entity,
) GameStateManager {
	return GameStateManager{
		GameModel:         gameModel,
		SystemsManager:    systemsManager,
		UI:                ui,
		LocaleMessages:    localeMessages,
		CollisionSystem:   collisionSystem,
		InputSystem:       inputSystem,
		ShouldAddEntities: shouldAddEntities,
		CurrentRoomID:     currentRoomID,
		NextRoomID:        nextRoomID,
		CurrentState:      currentState,
		Spritesheet:       spritesheet,
		MapDrawData:       mapDrawData,
		RoomTransition:    roomTransition,
		EntitiesMap:       entitiesMap,
		Player:            player,
		Hearts:            hearts,
	}
}

func (g *GameStateManager) Update() {
	switch *g.CurrentState {
	case StateStart:
		GameStateStart(g.UI, g.LocaleMessages, g.GameModel, g.CurrentState)
	case StateGame:
		GameStateGame(
			g.UI,
			g.GameModel,
			RoomsMap,
			g.SystemsManager,
			g.InputSystem,
			g.ShouldAddEntities,
			g.CurrentRoomID,
			g.CurrentState,
			g.Spritesheet,
			g.MapDrawData,
			g.EntitiesMap,
			g.Player,
			g.Hearts,
		)
	case StatePause:
		GameStatePause(g.UI, g.LocaleMessages, g.GameModel, g.CurrentState)
	case StateOver:
		GameStateOver(g.UI, g.LocaleMessages, g.GameModel, g.CurrentState)
	case StateMapTransition:
		GameStateMapTransition(
			g.UI,
			g.SystemsManager,
			RoomsMap,
			g.CollisionSystem,
			g.GameModel,
			g.InputSystem,
			g.CurrentRoomID,
			g.NextRoomID,
			g.CurrentState,
			g.Spritesheet,
			g.MapDrawData,
			g.RoomTransition,
			g.Player,
		)
	}
}
