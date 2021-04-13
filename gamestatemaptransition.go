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
) {
	inputSystem.DisablePlayer()
	if roomTransition.Style == TransitionSlide && roomTransition.Timer > 0 {
		roomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()

		connectedRooms := roomsMap[*currentRoomID].ConnectedRooms()

		transitionRoomResp := calculateTransitionSlide(
			roomTransition,
			*connectedRooms,
		)

		*nextRoomID = transitionRoomResp.nextRoomID

		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			mapDrawData,
			roomsMap[*currentRoomID].MapName(),
			transitionRoomResp.modX,
			transitionRoomResp.modY,
		)
		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			mapDrawData,
			roomsMap[*nextRoomID].MapName(),
			transitionRoomResp.modXNext,
			transitionRoomResp.modYNext,
		)
		DrawMask(ui.Window)

		// Move player with map transition
		player.ComponentSpatial.Rect = pixel.R(
			player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX,
			player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY,
			player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX+TileSize,
			player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY+TileSize,
		)

		systemsManager.Update()
	} else if roomTransition.Style == TransitionWarp && roomTransition.Timer > 0 {
		roomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

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
) transitionRoomResponse {

	var nextRoomID RoomID
	inc := (roomTransition.Start - float64(roomTransition.Timer))
	incY := inc * (MapH / TileSize)
	incX := inc * (MapW / TileSize)
	modY := 0.0
	modYNext := 0.0
	modX := 0.0
	modXNext := 0.0
	playerModX := 0.0
	playerModY := 0.0
	playerIncY := ((MapH / TileSize) - 1) + 7
	playerIncX := ((MapW / TileSize) - 1) + 7
	if roomTransition.Side == BoundBottom && connectedRooms.Bottom != 0 {
		modY = incY
		modYNext = incY - MapH
		nextRoomID = connectedRooms.Bottom
		playerModY += playerIncY
	} else if roomTransition.Side == BoundTop && connectedRooms.Top != 0 {
		modY = -incY
		modYNext = -incY + MapH
		nextRoomID = connectedRooms.Top
		playerModY -= playerIncY
	} else if roomTransition.Side == BoundLeft && connectedRooms.Left != 0 {
		modX = incX
		modXNext = incX - MapW
		nextRoomID = connectedRooms.Left
		playerModX += playerIncX
	} else if roomTransition.Side == BoundRight && connectedRooms.Right != 0 {
		modX = -incX
		modXNext = -incX + MapW
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
