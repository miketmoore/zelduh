package zelduh

import (
	"fmt"

	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

func (s *GameStateManager) stateMapTransition() error {

	ui := s.UI
	systemsManager := s.SystemsManager

	s.InputSystem.Disable()

	if s.RoomTransition.Style == TransitionSlide && s.RoomTransition.Timer > 0 {
		s.RoomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		ui.DrawMapBackground(colornames.White)

		s.CollisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()

		connectedRooms := s.LevelManager.CurrentLevel.RoomByIDMap[*s.CurrentRoomID].ConnectedRooms()

		transitionRoomResp := calculateTransitionSlide(
			s.RoomTransition,
			*connectedRooms,
			s.TileSize,
			s.ActiveSpaceRectangle,
		)

		*s.NextRoomID = transitionRoomResp.nextRoomID

		if s.CurrentRoomID == nil {
			return fmt.Errorf("current room ID is nil")
		}
		currentRoom, currentRoomOk := s.LevelManager.CurrentLevel.RoomByIDMap[*s.CurrentRoomID]
		if !currentRoomOk {
			return fmt.Errorf("current room not found by ID=%d", *s.CurrentRoomID)
		}
		ui.DrawMapBackgroundImage(
			currentRoom.Name,
			transitionRoomResp.modX,
			transitionRoomResp.modY,
		)

		if s.NextRoomID == nil {
			return fmt.Errorf("next room ID is nil")
		}
		nextRoom, nextRoomOk := s.LevelManager.CurrentLevel.RoomByIDMap[*s.NextRoomID]
		if !nextRoomOk {
			return fmt.Errorf("next room not found by ID=%d", *s.NextRoomID)
		}
		ui.DrawMapBackgroundImage(
			nextRoom.Name,
			transitionRoomResp.modXNext,
			transitionRoomResp.modYNext,
		)
		ui.DrawMask()

		// Move player with map transition
		s.Player.componentRectangle.Rect = pixel.R(
			s.Player.componentRectangle.Rect.Min.X+transitionRoomResp.playerModX,
			s.Player.componentRectangle.Rect.Min.Y+transitionRoomResp.playerModY,
			s.Player.componentRectangle.Rect.Min.X+transitionRoomResp.playerModX+s.TileSize,
			s.Player.componentRectangle.Rect.Min.Y+transitionRoomResp.playerModY+s.TileSize,
		)

		err := systemsManager.Update()
		if err != nil {
			return err
		}
	} else if s.RoomTransition.Style == TransitionWarp && s.RoomTransition.Timer > 0 {
		s.RoomTransition.Timer--
		ui.Window.Clear(colornames.Darkgray)
		ui.DrawMapBackground(colornames.White)

		s.CollisionSystem.RemoveAll(CategoryObstacle)
		systemsManager.RemoveAllEnemies()
		systemsManager.RemoveAllCollisionSwitches()
		systemsManager.RemoveAllMoveableObstacles()
		systemsManager.RemoveAllEntities()
	} else {
		*s.CurrentState = StateGame
		if *s.NextRoomID != 0 {
			*s.CurrentRoomID = *s.NextRoomID
		}
		s.RoomTransition.Active = false
	}

	return nil
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
