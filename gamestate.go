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
	SystemsManager              *SystemsManager
	UI                          UI
	LocaleMessages              LocaleMessagesMap
	CollisionSystem             *CollisionSystem
	InputSystem                 *InputSystem
	ShouldAddEntities           *bool
	CurrentRoomID, NextRoomID   *RoomID
	CurrentState                *State
	Spritesheet                 Spritesheet
	MapDrawData                 MapDrawData
	RoomTransition              *RoomTransition
	EntitiesMap                 EntitiesMap
	Player                      *Entity
	Hearts                      []Entity
	RoomWarps                   RoomWarps
	LevelManager                *LevelManager
	EntityConfigPresetFnManager *EntityConfigPresetFnManager
	TileSize                    float64
	FrameRate                   int
	NonObstacleSprites          map[int]bool
	WindowConfig                WindowConfig
	ActiveSpaceRectangle        ActiveSpaceRectangle
}

func NewGameStateManager(
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
	roomWarps RoomWarps,
	levelManager *LevelManager,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
	frameRate int,
	nonObstacleSprites map[int]bool,
	windowConfig WindowConfig,
	activeSpaceRectangle ActiveSpaceRectangle,
) GameStateManager {
	return GameStateManager{
		SystemsManager:              systemsManager,
		UI:                          ui,
		LocaleMessages:              localeMessages,
		CollisionSystem:             collisionSystem,
		InputSystem:                 inputSystem,
		ShouldAddEntities:           shouldAddEntities,
		CurrentRoomID:               currentRoomID,
		NextRoomID:                  nextRoomID,
		CurrentState:                currentState,
		Spritesheet:                 spritesheet,
		MapDrawData:                 mapDrawData,
		RoomTransition:              roomTransition,
		EntitiesMap:                 entitiesMap,
		Player:                      player,
		Hearts:                      hearts,
		RoomWarps:                   roomWarps,
		LevelManager:                levelManager,
		EntityConfigPresetFnManager: entityConfigPresetFnManager,
		TileSize:                    tileSize,
		FrameRate:                   frameRate,
		NonObstacleSprites:          nonObstacleSprites,
		WindowConfig:                windowConfig,
		ActiveSpaceRectangle:        activeSpaceRectangle,
	}
}

func (g *GameStateManager) Update() {
	switch *g.CurrentState {
	case StateStart:
		GameStateStart(g.UI, g.LocaleMessages, g.CurrentState, g.ActiveSpaceRectangle)
	case StateGame:
		GameStateGame(
			g.UI,
			g.LevelManager.CurrentLevel.RoomByIDMap,
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
			g.RoomWarps,
			g.EntityConfigPresetFnManager,
			g.TileSize,
			g.FrameRate,
			g.NonObstacleSprites,
			g.WindowConfig,
			g.ActiveSpaceRectangle,
		)
	case StatePause:
		GameStatePause(g.UI, g.LocaleMessages, g.CurrentState, g.ActiveSpaceRectangle)
	case StateOver:
		GameStateOver(g.UI, g.LocaleMessages, g.CurrentState, g.ActiveSpaceRectangle)
	case StateMapTransition:
		GameStateMapTransition(
			g.UI,
			g.SystemsManager,
			g.LevelManager.CurrentLevel.RoomByIDMap,
			g.CollisionSystem,
			g.InputSystem,
			g.CurrentRoomID,
			g.NextRoomID,
			g.CurrentState,
			g.Spritesheet,
			g.MapDrawData,
			g.RoomTransition,
			g.Player,
			g.TileSize,
			g.WindowConfig,
			g.ActiveSpaceRectangle,
		)
	}
}
