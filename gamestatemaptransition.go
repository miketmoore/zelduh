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
	gameStateManager *GameStateManager,
	roomData *RoomData,
	roomTransitionManager *RoomTransitionManager,
	windowConfig WindowConfig,
	mapConfig MapConfig,
) {
	inputSystem.DisablePlayer()
	if roomTransitionManager.Style() == TransitionSlide && roomTransitionManager.Timer() > 0 {
		roomTransitionManager.DecrementTimer()
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, mapConfig, colornames.White)

		collisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()

		currentRoomID := roomData.CurrentRoomID

		connectedRooms := roomsMap[currentRoomID].ConnectedRooms()

		transitionRoomResp := calculateTransitionSlide(
			roomTransitionManager,
			*connectedRooms,
			roomData.CurrentRoomID,
			mapConfig,
		)

		roomData.NextRoomID = transitionRoomResp.nextRoomID

		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			allMapDrawData,
			roomsMap[roomData.CurrentRoomID].MapName(),
			transitionRoomResp.modX,
			transitionRoomResp.modY,
			mapConfig,
		)
		DrawMapBackgroundImage(
			ui.Window,
			spritesheet,
			allMapDrawData,
			roomsMap[roomData.NextRoomID].MapName(),
			transitionRoomResp.modXNext,
			transitionRoomResp.modYNext,
			mapConfig,
		)
		DrawMask(ui.Window, windowConfig, mapConfig)

		// Move player with map transition
		entities.Player.ComponentSpatial.Rect = pixel.R(
			entities.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX,
			entities.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY,
			entities.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX+TileSize,
			entities.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY+TileSize,
		)

		systemsManager.Update()
	} else if roomTransitionManager.Style() == TransitionWarp && roomTransitionManager.Timer() > 0 {
		roomTransitionManager.DecrementTimer()
		ui.Window.Clear(colornames.Darkgray)
		DrawMapBackground(ui.Window, mapConfig, colornames.White)

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
		roomTransitionManager.SetActive(false)
	}
}

type transitionRoomResponse struct {
	nextRoomID                                             RoomID
	modX, modY, modXNext, modYNext, playerModX, playerModY float64
}

func calculateTransitionSlide(
	roomTransitionManager *RoomTransitionManager,
	connectedRooms ConnectedRooms,
	currentRoomID RoomID,
	mapConfig MapConfig,
) transitionRoomResponse {

	var nextRoomID RoomID
	inc := (roomTransitionManager.Start() - float64(roomTransitionManager.Timer()))
	incY := inc * (mapConfig.Height / TileSize)
	incX := inc * (mapConfig.Width / TileSize)
	modY := 0.0
	modYNext := 0.0
	modX := 0.0
	modXNext := 0.0
	playerModX := 0.0
	playerModY := 0.0
	playerIncY := ((mapConfig.Height / TileSize) - 1) + 7
	playerIncX := ((mapConfig.Width / TileSize) - 1) + 7

	side := roomTransitionManager.Side()

	if side == BoundBottom && connectedRooms.Bottom != 0 {
		modY = incY
		modYNext = incY - mapConfig.Height
		nextRoomID = connectedRooms.Bottom
		playerModY += playerIncY
	} else if side == BoundTop && connectedRooms.Top != 0 {
		modY = -incY
		modYNext = -incY + mapConfig.Height
		nextRoomID = connectedRooms.Top
		playerModY -= playerIncY
	} else if side == BoundLeft && connectedRooms.Left != 0 {
		modX = incX
		modXNext = incX - mapConfig.Width
		nextRoomID = connectedRooms.Left
		playerModX += playerIncX
	} else if side == BoundRight && connectedRooms.Right != 0 {
		modX = -incX
		modXNext = -incX + mapConfig.Width
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
