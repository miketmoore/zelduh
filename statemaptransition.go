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
	connectedRooms ConnectedRooms,
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
	if roomTransitionManager.Side() == BoundBottom && connectedRooms.Bottom != 0 {
		modY = incY
		modYNext = incY - activeSpaceRectangle.Height
		nextRoomID = connectedRooms.Bottom
		playerModY += playerIncY
	} else if roomTransitionManager.Side() == BoundTop && connectedRooms.Top != 0 {
		modY = -incY
		modYNext = -incY + activeSpaceRectangle.Height
		nextRoomID = connectedRooms.Top
		playerModY -= playerIncY
	} else if roomTransitionManager.Side() == BoundLeft && connectedRooms.Left != 0 {
		modX = incX
		modXNext = incX - activeSpaceRectangle.Width
		nextRoomID = connectedRooms.Left
		playerModX += playerIncX
	} else if roomTransitionManager.Side() == BoundRight && connectedRooms.Right != 0 {
		modX = -incX
		modXNext = -incX + activeSpaceRectangle.Width
		nextRoomID = connectedRooms.Right
		playerModX -= playerIncX
	} else {
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
	currentRoomID         *RoomID
	nextRoomID            *RoomID
	tileSize              float64
	activeSpaceRectangle  ActiveSpaceRectangle
	player                Entity
}

func NewStateTransition(
	context *StateContext,
	uiSystem *UISystem,
	inputSystem *InputSystem,
	roomTransitionManager *RoomTransitionManager,
	collisionSystem *CollisionSystem,
	systemsManager *SystemsManager,
	levelManager *LevelManager,
	currentRoomID *RoomID,
	nextRoomID *RoomID,
	tileSize float64,
	activeSpaceRectangle ActiveSpaceRectangle,
	player Entity,
) State {
	return StateTransition{
		context:               context,
		uiSystem:              uiSystem,
		inputSystem:           inputSystem,
		roomTransitionManager: roomTransitionManager,
		collisionSystem:       collisionSystem,
		systemsManager:        systemsManager,
		levelManager:          levelManager,
		currentRoomID:         currentRoomID,
		nextRoomID:            nextRoomID,
		tileSize:              tileSize,
		activeSpaceRectangle:  activeSpaceRectangle,
		player:                player,
	}
}

func (g StateTransition) Update() error {

	g.inputSystem.Disable()

	if g.roomTransitionManager.Style() == RoomTransitionSlide && g.roomTransitionManager.Timer() > 0 {
		g.roomTransitionManager.DecrementTimer()
		g.uiSystem.Window.Clear(colornames.Darkgray)
		g.uiSystem.DrawMapBackground(colornames.White)

		g.collisionSystem.RemoveAll(CategoryObstacle)
		g.systemsManager.RemoveAllEnemies()
		g.systemsManager.RemoveAllCollisionSwitches()
		g.systemsManager.RemoveAllMoveableObstacles()
		g.systemsManager.RemoveAllEntities()

		connectedRooms := g.levelManager.CurrentLevel.RoomByIDMap[*g.currentRoomID].ConnectedRooms()

		transitionRoomResp := calculateTransitionSlide(
			g.roomTransitionManager,
			*connectedRooms,
			g.tileSize,
			g.activeSpaceRectangle,
		)

		*g.nextRoomID = transitionRoomResp.nextRoomID

		if g.currentRoomID == nil {
			return fmt.Errorf("current room ID is nil")
		}
		currentRoom, currentRoomOk := g.levelManager.CurrentLevel.RoomByIDMap[*g.currentRoomID]
		if !currentRoomOk {
			return fmt.Errorf("current room not found by ID=%d", *g.currentRoomID)
		}
		g.uiSystem.DrawMapBackgroundImage(
			currentRoom.Name,
			transitionRoomResp.modX,
			transitionRoomResp.modY,
		)

		if g.nextRoomID == nil {
			return fmt.Errorf("next room ID is nil")
		}
		nextRoom, nextRoomOk := g.levelManager.CurrentLevel.RoomByIDMap[*g.nextRoomID]
		if !nextRoomOk {
			return fmt.Errorf("next room not found by ID=%d", *g.nextRoomID)
		}
		g.uiSystem.DrawMapBackgroundImage(
			nextRoom.Name,
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
	} else if g.roomTransitionManager.Style() == RoomTransitionWarp && g.roomTransitionManager.Timer() > 0 {
		g.roomTransitionManager.DecrementTimer()
		g.uiSystem.Window.Clear(colornames.Darkgray)
		g.uiSystem.DrawMapBackground(colornames.White)

		g.collisionSystem.RemoveAll(CategoryObstacle)
		g.systemsManager.RemoveAllEnemies()
		g.systemsManager.RemoveAllCollisionSwitches()
		g.systemsManager.RemoveAllMoveableObstacles()
		g.systemsManager.RemoveAllEntities()
	} else {
		fmt.Println("state: transition => game")
		err := g.context.SetState(StateNameGame)
		if err != nil {
			return err
		}
		if *g.nextRoomID != 0 {
			*g.currentRoomID = *g.nextRoomID
		}
		g.roomTransitionManager.Disable()
	}

	return nil
}
