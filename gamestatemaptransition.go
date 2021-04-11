package zelduh

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func GameStateMapTransition(
	ui UI,
	spritesheet map[int]*pixel.Sprite,
	entities Entities,
	allMapDrawData map[string]MapData,
	inputSystem *SystemInput,
	systemsManager *SystemsManager,
	roomsMap Rooms,
	collisionSystem *SystemCollision,
	gameModel *GameModel,
	gameStateManager *GameStateManager,
	roomData *RoomData,
) {
	inputSystem.DisablePlayer()
	if gameModel.RoomTransition.Style == TransitionSlide && gameModel.RoomTransition.Timer > 0 {
		gameModel.RoomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()

		currentRoomID := roomData.CurrentRoomID

		connectedRooms := roomsMap[currentRoomID].ConnectedRooms()

		transitionRoomResp := calculateTransitionSlide(
			gameModel.RoomTransition,
			*connectedRooms,
			roomData.CurrentRoomID,
		)

		roomData.NextRoomID = transitionRoomResp.nextRoomID

		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			allMapDrawData,
			roomsMap[roomData.CurrentRoomID].MapName(),
			transitionRoomResp.modX,
			transitionRoomResp.modY,
		)
		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			allMapDrawData,
			roomsMap[roomData.NextRoomID].MapName(),
			transitionRoomResp.modXNext,
			transitionRoomResp.modYNext,
		)
		DrawMask(ui.Window)

		// Move player with map transition
		entities.Player.ComponentSpatial.Rect = pixel.R(
			entities.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX,
			entities.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY,
			entities.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX+TileSize,
			entities.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY+TileSize,
		)

		systemsManager.Update()
	} else if gameModel.RoomTransition.Style == TransitionWarp && gameModel.RoomTransition.Timer > 0 {
		gameModel.RoomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()
	} else {
		gameStateManager.CurrentState = StateGame
		if roomData.NextRoomID != 0 {
			roomData.CurrentRoomID = roomData.NextRoomID
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
