package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(ui UI, spritesheet map[int]*pixel.Sprite, gameModel *GameModel, entities Entities, roomWarps map[EntityID]Config, entitiesMap EntityByEntityID, allMapDrawData map[string]MapData, inputSystem *SystemInput, roomsMap Rooms, systemsManager *SystemsManager) {
	inputSystem.EnablePlayer()

	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

	DrawMapBackgroundImage(
		ui.Window,
		spritesheet,
		allMapDrawData,
		roomsMap[gameModel.CurrentRoomID].MapName(),
		0, 0)

	if gameModel.AddEntities {
		gameModel.AddEntities = false
		AddUIHearts(systemsManager, entities.Hearts, entities.Player.ComponentHealth.Total)

		AddUICoin(systemsManager)

		// Draw obstacles on appropriate map tiles
		obstacles := DrawObstaclesPerMapTiles(systemsManager, roomsMap, allMapDrawData, gameModel.CurrentRoomID, 0, 0)
		systemsManager.AddEntities(obstacles...)

		// Iterate through all entity configurations and build entities and add to systems
		for _, c := range roomsMap[gameModel.CurrentRoomID].(*Room).EntityConfigs {
			entity := BuildEntityFromConfig(c, systemsManager.NewEntityID())
			entitiesMap[entity.ID()] = entity
			systemsManager.AddEntity(entity)

			switch c.Category {
			case CategoryWarp:
				roomWarps[entity.ID()] = c
			}
		}
	}

	DrawMask(ui.Window)

	systemsManager.Update()

	if ui.Window.JustPressed(pixelgl.KeyP) {
		gameModel.CurrentState = StatePause
	}

	if ui.Window.JustPressed(pixelgl.KeyX) {
		gameModel.CurrentState = StateOver
	}
}
