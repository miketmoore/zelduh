package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(ui UI, gameModel *GameModel, roomsMap Rooms, gameWorld *World) {
	gameModel.InputSystem.EnablePlayer()

	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

	DrawMapBackgroundImage(
		ui.Window,
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

	DrawMask(ui.Window)

	gameWorld.Update()

	if ui.Window.JustPressed(pixelgl.KeyP) {
		gameModel.CurrentState = StatePause
	}

	if ui.Window.JustPressed(pixelgl.KeyX) {
		gameModel.CurrentState = StateOver
	}
}
