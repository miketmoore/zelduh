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
	SystemsManager            *SystemsManager
	UI                        UI
	LocaleMessages            LocaleMessagesMap
	CollisionSystem           *CollisionSystem
	InputSystem               *InputSystem
	ShouldAddEntities         *bool
	CurrentRoomID, NextRoomID *RoomID
	CurrentState              *State
	SpriteMap                 SpriteMap
	MapDrawData               MapDrawData
	RoomTransition            *RoomTransition
	EntitiesMap               EntitiesMap
	Player                    *Entity
	// Hearts                      []Entity
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
	spriteMap SpriteMap,
	mapDrawData MapDrawData,
	roomTransition *RoomTransition,
	entitiesMap EntitiesMap,
	player *Entity,
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
		SpriteMap:                   spriteMap,
		MapDrawData:                 mapDrawData,
		RoomTransition:              roomTransition,
		EntitiesMap:                 entitiesMap,
		Player:                      player,
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

func (g *GameStateManager) Update() error {
	var err error
	switch *g.CurrentState {
	case StateStart:
		err = GameStateStart(g.UI, g.LocaleMessages, g.CurrentState, g.ActiveSpaceRectangle)
	case StateGame:
		err = GameStateGame(
			g.UI,
			g.LevelManager.CurrentLevel.RoomByIDMap,
			g.SystemsManager,
			g.InputSystem,
			g.ShouldAddEntities,
			g.CurrentRoomID,
			g.CurrentState,
			g.SpriteMap,
			g.MapDrawData,
			g.EntitiesMap,
			g.Player,
			g.RoomWarps,
			g.EntityConfigPresetFnManager,
			g.TileSize,
			g.FrameRate,
			g.NonObstacleSprites,
			g.WindowConfig,
			g.ActiveSpaceRectangle,
		)
	case StatePause:
		err = GameStatePause(g.UI, g.LocaleMessages, g.CurrentState, g.ActiveSpaceRectangle)
	case StateOver:
		err = GameStateOver(g.UI, g.LocaleMessages, g.CurrentState, g.ActiveSpaceRectangle)
	case StateMapTransition:
		err = GameStateMapTransition(
			g.UI,
			g.SystemsManager,
			g.LevelManager.CurrentLevel.RoomByIDMap,
			g.CollisionSystem,
			g.InputSystem,
			g.CurrentRoomID,
			g.NextRoomID,
			g.CurrentState,
			g.SpriteMap,
			g.MapDrawData,
			g.RoomTransition,
			g.Player,
			g.TileSize,
			g.WindowConfig,
			g.ActiveSpaceRectangle,
		)
	}

	return err
}
