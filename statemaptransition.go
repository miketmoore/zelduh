package zelduh

import (
	"fmt"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

type transitionRoomResponse struct {
	nextRoomID                                             RoomID
	modX, modY, modXNext, modYNext, playerModX, playerModY float64
}

func calculateTransitionSlide(
	roomTransitionManager *RoomTransitionManager,
	connectedRooms *ConnectedRooms,
	tileSize float64,
	activeSpaceRectangle ActiveSpaceRectangle,
) transitionRoomResponse {

	var nextRoomID RoomID
	inc := (roomTransitionManager.Start() - float64(roomTransitionManager.Timer()))
	incY := inc * (activeSpaceRectangle.Height / tileSize)
	incX := inc * (activeSpaceRectangle.Width / tileSize)
	modY := 0.0
	modYNext := 0.0
	modX := 0.0
	modXNext := 0.0
	playerModX := 0.0
	playerModY := 0.0
	playerIncY := ((activeSpaceRectangle.Height / tileSize) - 1) + 7
	playerIncX := ((activeSpaceRectangle.Width / tileSize) - 1) + 7

	// fmt.Printf("calculateTransitionSlide [%s] %s\n", roomTransitionManager.Side(), connectedRooms.String())

	if roomTransitionManager.Side() == BoundBottom && connectedRooms.Bottom != 0 {
		// fmt.Printf("calculateTransitionSlide bottom\n")
		modY = incY
		modYNext = incY - activeSpaceRectangle.Height
		nextRoomID = connectedRooms.Bottom
		playerModY += playerIncY
	} else if roomTransitionManager.Side() == BoundTop && connectedRooms.Top != 0 {
		// fmt.Printf("calculateTransitionSlide top\n")
		modY = -incY
		modYNext = -incY + activeSpaceRectangle.Height
		nextRoomID = connectedRooms.Top
		playerModY -= playerIncY
	} else if roomTransitionManager.Side() == BoundLeft && connectedRooms.Left != 0 {
		// fmt.Printf("calculateTransitionSlide left\n")
		modX = incX
		modXNext = incX - activeSpaceRectangle.Width
		nextRoomID = connectedRooms.Left
		playerModX += playerIncX
	} else if roomTransitionManager.Side() == BoundRight && connectedRooms.Right != 0 {
		// fmt.Printf("calculateTransitionSlide right\n")
		modX = -incX
		modXNext = -incX + activeSpaceRectangle.Width
		nextRoomID = connectedRooms.Right
		playerModX -= playerIncX
	} else {
		fmt.Printf("calculateTransitionSlide default\n")
		nextRoomID = 0
	}

	return transitionRoomResponse{
		nextRoomID,
		modX, modY, modXNext, modYNext, playerModX, playerModY,
	}
}

type StateTransition struct {
	context               *StateContext
	uiSystem              *UISystem
	inputSystem           *InputSystem
	roomTransitionManager *RoomTransitionManager
	collisionSystem       *CollisionSystem
	systemsManager        *SystemsManager
	levelManager          *LevelManager
	roomManager           *RoomManager
	tileSize              float64
	activeSpaceRectangle  ActiveSpaceRectangle
	player                Entity
	shouldAddEntities     *bool
}

func NewStateTransition(
	context *StateContext,
	uiSystem *UISystem,
	inputSystem *InputSystem,
	roomTransitionManager *RoomTransitionManager,
	collisionSystem *CollisionSystem,
	systemsManager *SystemsManager,
	levelManager *LevelManager,
	roomManager *RoomManager,
	tileSize float64,
	activeSpaceRectangle ActiveSpaceRectangle,
	player Entity,
	shouldAddEntities *bool,
) State {
	return StateTransition{
		context:               context,
		uiSystem:              uiSystem,
		inputSystem:           inputSystem,
		roomTransitionManager: roomTransitionManager,
		collisionSystem:       collisionSystem,
		systemsManager:        systemsManager,
		levelManager:          levelManager,
		roomManager:           roomManager,
		tileSize:              tileSize,
		activeSpaceRectangle:  activeSpaceRectangle,
		player:                player,
		shouldAddEntities:     shouldAddEntities,
	}
}

func (g StateTransition) Update() error {

	g.inputSystem.Disable()

	currentRoomID := g.roomManager.Current()

	if g.roomTransitionManager.Style() == RoomTransitionSlide && g.roomTransitionManager.Timer() > 0 {
		g.common()
		err := g.slideTransition(currentRoomID)
		if err != nil {
			fmt.Println(err)
			// return fmt.Errorf("error during slide transition currentRoomID=%d", currentRoomID)

		}
		return nil
	} else if g.roomTransitionManager.Style() == RoomTransitionWarp && g.roomTransitionManager.Timer() > 0 {
		g.common()
		return nil
	} else {
		fmt.Println("transition done")
		fmt.Printf("g.roomManager.Next()=%d\n", g.roomManager.Next())
		nextRoom := g.levelManager.CurrentLevel.RoomByIDMap[g.roomManager.Next()]
		fmt.Println(nextRoom.TMXFileName)

		err := g.context.SetState(StateNameGame)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("error occured when changing to game state after map transitiong was finished")
		}
		if g.roomManager.Next() != 0 {
			g.roomManager.MoveToNext()
		}
		g.roomTransitionManager.Disable()
		return nil
	}

}

func (g StateTransition) slideTransition(currentRoomID RoomID) error {

	connectedRooms := g.levelManager.CurrentLevel.RoomByIDMap[currentRoomID].ConnectedRooms()

	transitionRoomResp := calculateTransitionSlide(
		g.roomTransitionManager,
		connectedRooms,
		g.tileSize,
		g.activeSpaceRectangle,
	)

	// nextRoomID = transitionRoomResp.nextRoomID
	// fmt.Printf("slideTransition calling SetNext; currentRoomID=%d nextRoomID=%d\n", currentRoomID, transitionRoomResp.nextRoomID)
	g.roomManager.SetNext(transitionRoomResp.nextRoomID)

	currentRoom, currentRoomOk := g.levelManager.CurrentLevel.RoomByIDMap[currentRoomID]
	if !currentRoomOk {
		return fmt.Errorf("current room not found by ID=%d", currentRoomID)
	}
	g.uiSystem.DrawMapBackgroundImage(
		currentRoom.TMXFileName,
		transitionRoomResp.modX,
		transitionRoomResp.modY,
	)

	nextRoom, nextRoomOk := g.levelManager.CurrentLevel.RoomByIDMap[g.roomManager.Next()]
	if !nextRoomOk {
		return fmt.Errorf("next room not found by ID=%d", g.roomManager.Next())
	}
	if g.roomTransitionManager.Timer() == 0 {
		fmt.Printf("drawing nextRoom.TMXFileName=%s\n", nextRoom.TMXFileName)
	}
	g.uiSystem.DrawMapBackgroundImage(
		nextRoom.TMXFileName,
		transitionRoomResp.modXNext,
		transitionRoomResp.modYNext,
	)
	g.uiSystem.DrawMask()

	// Move player with map transition
	g.player.componentRectangle.Rect = pixel.R(
		g.player.componentRectangle.Rect.Min.X+transitionRoomResp.playerModX,
		g.player.componentRectangle.Rect.Min.Y+transitionRoomResp.playerModY,
		g.player.componentRectangle.Rect.Min.X+transitionRoomResp.playerModX+g.tileSize,
		g.player.componentRectangle.Rect.Min.Y+transitionRoomResp.playerModY+g.tileSize,
	)

	err := g.systemsManager.Update()
	if err != nil {
		return err
	}

	*g.shouldAddEntities = true

	return nil
}

func (g StateTransition) common() {
	g.roomTransitionManager.DecrementTimer()
	g.uiSystem.Window.Clear(colornames.Darkgray)
	g.uiSystem.DrawMapBackground(colornames.White)

	// g.collisionSystem.RemoveAll(CategoryObstacle)
	g.systemsManager.RemoveAllByCategory(CategoryObstacle)
	// g.systemsManager.RemoveAllEnemies()
	// g.systemsManager.RemoveAllCollisionSwitches()
	// g.systemsManager.RemoveAllMoveableObstacles()
	// g.systemsManager.RemoveAllEntities()
}
