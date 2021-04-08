package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(win *pixelgl.Window, gameModel *GameModel, roomsMap Rooms, gameWorld *World) {
	gameModel.InputSystem.EnablePlayer()

	win.Clear(colornames.Darkgray)
	DrawMapBackground(win, MapX, MapY, MapW, MapH, colornames.White)

	DrawMapBackgroundImage(
		win,
		gameModel.Spritesheet,
		gameModel.AllMapDrawData,
		roomsMap[gameModel.CurrentRoomID].MapName(),
		0, 0)

	if gameModel.AddEntities {
		gameModel.AddEntities = false
		AddUIHearts(gameWorld, gameModel.Hearts, gameModel.Player.ComponentHealth.Total)

		AddUICoin(gameWorld)

		// Draw obstacles on appropriate map tiles
		obstacles := DrawObstaclesPerMapTiles(gameWorld, roomsMap, gameModel.AllMapDrawData, gameModel.CurrentRoomID, 0, 0)
		gameWorld.AddEntities(obstacles...)

		gameModel.RoomWarps = map[EntityID]Config{}

		// Iterate through all entity configurations and build entities and add to systems
		for _, c := range roomsMap[gameModel.CurrentRoomID].(*Room).EntityConfigs {
			entity := BuildEntityFromConfig(c, gameWorld.NewEntityID())
			gameModel.EntitiesMap[entity.ID()] = entity
			gameWorld.AddEntity(entity)

			switch c.Category {
			case CategoryWarp:
				gameModel.RoomWarps[entity.ID()] = c
			}
		}
	}

	DrawMask(win)

	gameWorld.Update()

	if win.JustPressed(pixelgl.KeyP) {
		gameModel.CurrentState = StatePause
	}

	if win.JustPressed(pixelgl.KeyX) {
		gameModel.CurrentState = StateOver
	}
}
