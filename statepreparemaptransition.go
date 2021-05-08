package zelduh

import (
	"fmt"

	"golang.org/x/image/colornames"
)

type StatePrepareMapTransition struct {
	context               *StateContext
	systemsManager        *SystemsManager
	inputSystem           *InputSystem
	uiSystem              *UISystem
	collisionSystem       *CollisionSystem
	roomManager           *RoomManager
	roomTransitionManager *RoomTransitionManager
	levelManager          *LevelManager
	activeSpaceRectangle  ActiveSpaceRectangle
	tileSize              float64
}

func NewStatePrepareMapTransition(
	context *StateContext,
	systemsManager *SystemsManager,
	inputSystem *InputSystem,
	uiSystem *UISystem,
	collisionSystem *CollisionSystem,
	roomManager *RoomManager,
	roomTransitionManager *RoomTransitionManager,
	levelManager *LevelManager,
	activeSpaceRectangle ActiveSpaceRectangle,
	tileSize float64,
) *StatePrepareMapTransition {
	return &StatePrepareMapTransition{
		context:               context,
		systemsManager:        systemsManager,
		inputSystem:           inputSystem,
		uiSystem:              uiSystem,
		collisionSystem:       collisionSystem,
		roomManager:           roomManager,
		roomTransitionManager: roomTransitionManager,
		levelManager:          levelManager,
		activeSpaceRectangle:  activeSpaceRectangle,
		tileSize:              tileSize,
	}
}

func (s StatePrepareMapTransition) Update() error {

	s.inputSystem.Disable()

	s.roomTransitionManager.Enable()
	// s.roomTransitionManager.SetSide(side)
	s.roomTransitionManager.SetSlide()
	s.roomTransitionManager.ResetTimer()

	s.roomTransitionManager.DecrementTimer()
	s.uiSystem.Window.Clear(colornames.Darkgray)
	s.uiSystem.DrawMapBackground(colornames.White)

	// s.collisionSystem.RemoveAll(CategoryObstacle)
	// s.systemsManager.RemoveAllEnemies()
	// s.systemsManager.RemoveAllCollisionSwitches()
	// s.systemsManager.RemoveAllMoveableObstacles()
	// s.systemsManager.RemoveAllEntities()
	s.systemsManager.RemoveAllByCategory(CategoryObstacle)

	fmt.Printf("StatePrepareMapTransition: currentRoomID=%d nextRoomID=%d\n", s.roomManager.Current(), s.roomManager.Next())

	currentRoomID := s.roomManager.Current()

	connectedRooms := s.levelManager.CurrentLevel.RoomByIDMap[currentRoomID].ConnectedRooms()

	fmt.Printf("%s\n", connectedRooms.String())

	transitionRoomResp := calculateTransitionSlide(
		s.roomTransitionManager,
		connectedRooms,
		s.tileSize,
		s.activeSpaceRectangle,
	)

	currentRoom, currentRoomOk := s.levelManager.CurrentLevel.RoomByIDMap[currentRoomID]
	if !currentRoomOk {
		return fmt.Errorf("current room not found by ID=%d", currentRoomID)
	}
	s.uiSystem.DrawMapBackgroundImage(
		currentRoom.TMXFileName,
		transitionRoomResp.modX,
		transitionRoomResp.modY,
	)

	s.uiSystem.DrawMask()

	fmt.Println(transitionRoomResp)

	fmt.Printf("StatePrepareMapTransition: nextRoomID before=%d after=%d\n", s.roomManager.Next(), transitionRoomResp.nextRoomID)

	// panic("blop")

	err := s.systemsManager.Update()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error occured when calling updating on system manager")
	}

	err = s.context.SetState(StateNameTransition)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error occured when setting state from StatePrepareMapTransition to StateNameTransition")
	}

	return nil
}
