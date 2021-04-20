package zelduh

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func GameStateMapTransition(
	ui UI,
	systemsManager *SystemsManager,
	roomsMap Rooms,
	collisionSystem *CollisionSystem,
	inputSystem *InputSystem,
	currentRoomID *RoomID,
	nextRoomID *RoomID,
	currentState *State,
	spritesheet Spritesheet,
	mapDrawData MapDrawData,
	roomTransition *RoomTransition,
	player *Entity,
	tileSize float64,
	windowConfig WindowConfig,
	activeSpaceRectangle ActiveSpaceRectangle,
) {
	inputSystem.DisablePlayer()
	if roomTransition.Style == TransitionSlide && roomTransition.Timer > 0 {
		roomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, activeSpaceRectangle, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()

		connectedRooms := roomsMap[*currentRoomID].ConnectedRooms()

		transitionRoomResp := calculateTransitionSlide(
			roomTransition,
			*connectedRooms,
			tileSize,
			activeSpaceRectangle,
		)

		*nextRoomID = transitionRoomResp.nextRoomID

		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			mapDrawData,
			roomsMap[*currentRoomID].RoomName(),
			transitionRoomResp.modX,
			transitionRoomResp.modY,
			tileSize,
			activeSpaceRectangle,
		)
		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			mapDrawData,
			roomsMap[*nextRoomID].RoomName(),
			transitionRoomResp.modXNext,
			transitionRoomResp.modYNext,
			tileSize,
			activeSpaceRectangle,
		)
		DrawMask(ui.Window, windowConfig, activeSpaceRectangle)

		// Move player with map transition
		player.ComponentSpatial.Rect = pixel.R(
			player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX,
			player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY,
			player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX+tileSize,
			player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY+tileSize,
		)

		systemsManager.Update()
	} else if roomTransition.Style == TransitionWarp && roomTransition.Timer > 0 {
		roomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, activeSpaceRectangle, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()
	} else {
		*currentState = StateGame
		if *nextRoomID != 0 {
			*currentRoomID = *nextRoomID
		}
		roomTransition.Active = false
	}
}

type transitionRoomResponse struct {
	nextRoomID                                             RoomID
	modX, modY, modXNext, modYNext, playerModX, playerModY float64
}

func calculateTransitionSlide(
	roomTransition *RoomTransition,
	connectedRooms ConnectedRooms,
	tileSize float64,
	activeSpaceRectangle ActiveSpaceRectangle,
) transitionRoomResponse {

	var nextRoomID RoomID
	inc := (roomTransition.Start - float64(roomTransition.Timer))
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
	if roomTransition.Side == BoundBottom && connectedRooms.Bottom != 0 {
		modY = incY
		modYNext = incY - activeSpaceRectangle.Height
		nextRoomID = connectedRooms.Bottom
		playerModY += playerIncY
	} else if roomTransition.Side == BoundTop && connectedRooms.Top != 0 {
		modY = -incY
		modYNext = -incY + activeSpaceRectangle.Height
		nextRoomID = connectedRooms.Top
		playerModY -= playerIncY
	} else if roomTransition.Side == BoundLeft && connectedRooms.Left != 0 {
		modX = incX
		modXNext = incX - activeSpaceRectangle.Width
		nextRoomID = connectedRooms.Left
		playerModX += playerIncX
	} else if roomTransition.Side == BoundRight && connectedRooms.Right != 0 {
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
