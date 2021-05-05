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
	UI                          UISystem
	CollisionSystem             *CollisionSystem
	InputSystem                 *InputSystem
	ShouldAddEntities           *bool
	CurrentRoomID, NextRoomID   *RoomID
	CurrentState                *State
	SpriteMap                   SpriteMap
	MapDrawData                 MapDrawData
	EntitiesMap                 EntitiesMap
	Player                      *Entity
	RoomWarps                   RoomWarps
	LevelManager                *LevelManager
	EntityConfigPresetFnManager *EntityConfigPresetFnManager
	TileSize                    float64
	FrameRate                   int
	NonObstacleSprites          map[int]bool
	ActiveSpaceRectangle        ActiveSpaceRectangle
	entityCreator               *EntityCreator
	roomTransitionManager       *RoomTransitionManager
}

func NewGameStateManager(
	systemsManager *SystemsManager,
	ui UISystem,
	collisionSystem *CollisionSystem,
	inputSystem *InputSystem,
	shouldAddEntities *bool,
	currentRoomID *RoomID,
	nextRoomID *RoomID,
	currentState *State,
	spriteMap SpriteMap,
	mapDrawData MapDrawData,
	entitiesMap EntitiesMap,
	player *Entity,
	roomWarps RoomWarps,
	levelManager *LevelManager,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
	frameRate int,
	nonObstacleSprites map[int]bool,
	activeSpaceRectangle ActiveSpaceRectangle,
	entityCreator *EntityCreator,
	roomTransitionManager *RoomTransitionManager,
) GameStateManager {
	return GameStateManager{
		SystemsManager:              systemsManager,
		UI:                          ui,
		CollisionSystem:             collisionSystem,
		InputSystem:                 inputSystem,
		ShouldAddEntities:           shouldAddEntities,
		CurrentRoomID:               currentRoomID,
		NextRoomID:                  nextRoomID,
		CurrentState:                currentState,
		SpriteMap:                   spriteMap,
		MapDrawData:                 mapDrawData,
		EntitiesMap:                 entitiesMap,
		Player:                      player,
		RoomWarps:                   roomWarps,
		LevelManager:                levelManager,
		EntityConfigPresetFnManager: entityConfigPresetFnManager,
		TileSize:                    tileSize,
		FrameRate:                   frameRate,
		NonObstacleSprites:          nonObstacleSprites,
		ActiveSpaceRectangle:        activeSpaceRectangle,
		entityCreator:               entityCreator,
		roomTransitionManager:       roomTransitionManager,
	}
}

func (g *GameStateManager) Update() error {
	var err error
	switch *g.CurrentState {
	case StateStart:
		err = g.stateStart()
	case StateGame:
		err = g.stateGame()
	case StatePause:
		err = g.statePause()
	case StateOver:
		err = g.stateOver()
	case StateMapTransition:
		err = g.stateMapTransition()
	}

	return err
}

func (g *GameStateManager) setCurrentState(state State) {
	*g.CurrentState = state
}
