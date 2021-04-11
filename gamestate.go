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
	RoomTransitionManager *RoomTransitionManager
	SystemsManager        *SystemsManager
	GameStateManager      *GameStateManager
	UI                    UI
	LocaleMessages        LocaleMessagesMap
	CollisionSystem       *SystemCollision
	InputSystem           *SystemInput
	Spritesheet           map[int]*pixel.Sprite
	EntitiesMap           EntityByEntityID
	AllMapDrawData        map[string]MapData
	RoomWarps             map[EntityID]EntityConfig
	Entities              Entities
	CurrentState          State
	RoomData              *RoomData
	MapConfig             MapConfig
	WindowConfig          WindowConfig
	FrameRate             int
}

func NewGameStateManager(
	roomTransitionManager *RoomTransitionManager,
	systemsManager *SystemsManager,
	ui UI,
	localeMessages LocaleMessagesMap,
	collisionSystem *SystemCollision,
	inputSystem *SystemInput,
	spritesheet map[int]*pixel.Sprite,
	entitiesMap EntityByEntityID,
	allMapDrawData map[string]MapData,
	roomWarps map[EntityID]EntityConfig,
	entities Entities,
	roomData *RoomData,
	mapConfig MapConfig,
	windowConfig WindowConfig,
	frameRate int,
) GameStateManager {
	return GameStateManager{
		RoomTransitionManager: roomTransitionManager,
		SystemsManager:        systemsManager,
		UI:                    ui,
		LocaleMessages:        localeMessages,
		CollisionSystem:       collisionSystem,
		InputSystem:           inputSystem,
		Spritesheet:           spritesheet,
		EntitiesMap:           entitiesMap,
		AllMapDrawData:        allMapDrawData,
		RoomWarps:             roomWarps,
		Entities:              entities,
		CurrentState:          StateStart,
		RoomData:              roomData,
		MapConfig:             mapConfig,
		WindowConfig:          windowConfig,
		FrameRate:             frameRate,
	}
}

func (g *GameStateManager) Update() {
	switch g.CurrentState {
	case StateStart:
		GameStateStart(g.UI, g.LocaleMessages, g.GameStateManager, g.MapConfig)
	case StateGame:
		GameStateGame(
			g.UI,
			g.MapConfig,
			g.WindowConfig,
			g.Spritesheet,
			g.RoomData,
			g.Entities,
			g.RoomWarps,
			g.EntitiesMap,
			g.AllMapDrawData,
			g.InputSystem,
			RoomsMap,
			g.SystemsManager,
			g.GameStateManager,
			g.FrameRate,
		)
	case StatePause:
		GameStatePause(g.UI, g.LocaleMessages, g.GameStateManager, g.MapConfig)
	case StateOver:
		GameStateOver(g.UI, g.LocaleMessages, g.GameStateManager, g.MapConfig)
	case StateMapTransition:
		GameStateMapTransition(
			g.UI,
			g.Spritesheet,
			g.Entities,
			g.AllMapDrawData,
			g.InputSystem,
			g.SystemsManager,
			RoomsMap,
			g.CollisionSystem,
			g.GameStateManager,
			g.RoomData,
			g.RoomTransitionManager,
			g.WindowConfig,
			g.MapConfig,
		)
	}
}
