package zelduh

import "fmt"

type StateName string

const (
	StateNameStart      StateName = "start"
	StateNameGame       StateName = "game"
	StateNamePause      StateName = "pause"
	StateNameGameOver   StateName = "gameOver"
	StateNameTransition StateName = "transition"
)

type StateMap map[StateName]State

type StateContext struct {
	current State

	stateStart, stateGameOver, statePause, stateTransition, stateGame State
}

func (g *StateContext) Update() error {
	return g.current.Update()
}

func NewStateContext(
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
) *StateContext {

	context := &StateContext{}

	context.stateStart = NewStateStart(context, uiSystem)

	// Set the first value for the current state
	context.current = context.stateStart

	context.stateGameOver = NewStateGameOver(context, uiSystem)
	context.statePause = NewStatePause(context, uiSystem)
	context.stateTransition = NewStateTransition(
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
	context.stateGame = NewStateGame(
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

func (g *StateContext) SetState(name StateName) error {

	var state State

	switch name {
	case StateNameStart:
		state = g.stateStart
	case StateNameGame:
		state = g.stateGame
	case StateNameGameOver:
		state = g.stateGameOver
	case StateNameTransition:
		state = g.stateTransition
	case StateNamePause:
		state = g.statePause
	default:
		return fmt.Errorf("state not found, value=%s", name)
	}

	g.current = state

	return nil
}

type State interface {
	Update() error
}
