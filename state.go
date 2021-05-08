package zelduh

import "fmt"

type StateName string

const (
	StateNameStart                StateName = "start"
	StateNameGame                 StateName = "game"
	StateNamePause                StateName = "pause"
	StateNameGameOver             StateName = "gameOver"
	StateNamePrepareMapTransition StateName = "prepareMapTransition"
	StateNameTransition           StateName = "transition"
)

type StateMap map[StateName]State

type StateContext struct {
	current State

	stateStart, stateGameOver, statePause, statePrepareMapTransition, stateTransition, stateGame State
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
	entityFactory *EntityFactory,
	roomWarps RoomWarps,
	entitiesMap EntitiesMap,
	roomManager *RoomManager,
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
	context.statePrepareMapTransition = NewStatePrepareMapTransition(
		context,
		systemsManager,
		inputSystem,
		uiSystem,
		collisionSystem,
		roomManager,
		roomTransitionManager,
		levelManager,
		activeSpaceRectangle,
		tileSize,
		shouldAddEntities,
	)
	context.stateTransition = NewStateTransition(
		context,
		uiSystem,
		inputSystem,
		roomTransitionManager,
		collisionSystem,
		systemsManager,
		levelManager,
		roomManager,
		tileSize,
		activeSpaceRectangle,
		player,
		shouldAddEntities,
	)
	context.stateGame = NewStateGame(
		context,
		uiSystem,
		inputSystem,
		roomManager,
		shouldAddEntities,
		levelManager,
		entityFactory,
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
	case StateNamePrepareMapTransition:
		state = g.statePrepareMapTransition
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
