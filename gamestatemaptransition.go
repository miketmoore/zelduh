package zelduh

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func GameStateMapTransition(ui UI, gameWorld *World, roomsMap Rooms, collisionSystem *SystemCollision, gameModel *GameModel) {
	gameModel.InputSystem.DisablePlayer()
	if gameModel.RoomTransition.Style == TransitionSlide && gameModel.RoomTransition.Timer > 0 {
		gameModel.RoomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		gameWorld.RemoveAllEnemies()
		gameWorld.RemoveAllCollisionSwitches()
		gameWorld.RemoveAllMoveableObstacles()
		gameWorld.RemoveAllEntities()

		currentRoomID := gameModel.CurrentRoomID

		connectedRooms := roomsMap[currentRoomID].ConnectedRooms()

		transitionRoomResp := calculateTransitionSlide(
			gameModel.RoomTransition,
			*connectedRooms,
			gameModel.CurrentRoomID,
		)

		gameModel.NextRoomID = transitionRoomResp.nextRoomID

		DrawMapBackgroundImage(
			ui.Window,
			gameModel.Spritesheet,
			gameModel.AllMapDrawData,
			roomsMap[gameModel.CurrentRoomID].MapName(),
			transitionRoomResp.modX,
			transitionRoomResp.modY,
		)
		DrawMapBackgroundImage(
			ui.Window,
			gameModel.Spritesheet,
			gameModel.AllMapDrawData,
			roomsMap[gameModel.NextRoomID].MapName(),
			transitionRoomResp.modXNext,
			transitionRoomResp.modYNext,
		)
		DrawMask(ui.Window)

		// Move player with map transition
		gameModel.Player.ComponentSpatial.Rect = pixel.R(
			gameModel.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX,
			gameModel.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY,
			gameModel.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX+TileSize,
			gameModel.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY+TileSize,
		)

		gameWorld.Update()
	} else if gameModel.RoomTransition.Style == TransitionWarp && gameModel.RoomTransition.Timer > 0 {
		gameModel.RoomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		gameWorld.RemoveAllEnemies()
		gameWorld.RemoveAllCollisionSwitches()
		gameWorld.RemoveAllMoveableObstacles()
		gameWorld.RemoveAllEntities()
	} else {
		gameModel.CurrentState = StateGame
		if gameModel.NextRoomID != 0 {
			gameModel.CurrentRoomID = gameModel.NextRoomID
		}
		gameModel.RoomTransition.Active = false
	}
}

type transitionRoomResponse struct {
	nextRoomID                                             RoomID
	modX, modY, modXNext, modYNext, playerModX, playerModY float64
}

func calculateTransitionSlide(
	roomTransition *RoomTransition,
	connectedRooms ConnectedRooms,
	currentRoomID RoomID) transitionRoomResponse {

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
