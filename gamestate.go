package zelduh

import "fmt"

type GameStateName string

const (
	GameStateNameStart      GameStateName = "start"
	GameStateNameGame       GameStateName = "game"
	GameStateNamePause      GameStateName = "pause"
	GameStateNameGameOver   GameStateName = "gameOver"
	GameStateNameTransition GameStateName = "transition"
)

type GameStateMap map[GameStateName]GameState

type GameStateContext struct {
	current GameState

	stateStart, stateGameOver, statePause, stateTransition, stateGame GameState
}

func (g *GameStateContext) Update() error {
	return g.current.Update()
}

func NewGameStateContext(
	uiSystem *UISystem,
	inputSystem *InputSystem,
	roomTransitionManager *RoomTransitionManager,
	collisionSystem *CollisionSystem,
	systemsManager *SystemsManager,
	levelManager *LevelManager,
	entityCreator *EntityCreator,
	roomWarps RoomWarps,
	entitiesMap EntitiesMap,
	currentRoomID *RoomID,
	nextRoomID *RoomID,
	shouldAddEntities *bool,
	tileSize float64,
	frameRate int,
	activeSpaceRectangle ActiveSpaceRectangle,
	player Entity,
) *GameStateContext {

	context := &GameStateContext{}

	context.stateStart = NewGameStateStart(context, uiSystem)
	context.current = context.stateStart

	context.stateGameOver = NewGameStateGameOver(context, uiSystem)
	context.statePause = NewGameStatePause(context, uiSystem)
	context.stateTransition = NewGameStateTransition(
		context,
		uiSystem,
		inputSystem,
		roomTransitionManager,
		collisionSystem,
		systemsManager,
		levelManager,
		currentRoomID,
		nextRoomID,
		tileSize,
		activeSpaceRectangle,
		player,
	)
	context.stateGame = NewGameStateGame(
		context,
		uiSystem,
		inputSystem,
		currentRoomID,
		shouldAddEntities,
		levelManager,
		entityCreator,
		systemsManager,
		roomWarps,
		entitiesMap,
		frameRate,
	)

	return context

}

func (g *GameStateContext) SetState(name GameStateName) {
	// fmt.Println("SetState ", name)
	// state, ok := g.stateMap[name]
	// if !ok {
	// 	panic(fmt.Sprintf("state not found, value=%s", name))
	// }
	// g.current = state

	var state GameState

	switch name {
	case GameStateNameStart:
		state = g.stateStart
	case GameStateNameGame:
		state = g.stateGame
	case GameStateNameGameOver:
		state = g.stateGameOver
	case GameStateNameTransition:
		state = g.stateTransition
	case GameStateNamePause:
		state = g.statePause
	default:
		panic(fmt.Sprintf("state not found, value=%s", name))
	}

	g.current = state
}

type GameState interface {
	Update() error
}
