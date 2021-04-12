package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(
	ui UI,
	gameModel *GameModel,
	roomsMap Rooms,
	systemsManager *SystemsManager,
	inputSystem *InputSystem,
	shouldAddEntities *bool,
	currentRoomID *RoomID,
	currentState *State,
) {
	inputSystem.EnablePlayer()

	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

	DrawMapBackgroundImage(
		ui.Window,
		gameModel.Spritesheet,
		gameModel.AllMapDrawData,
		roomsMap[*currentRoomID].MapName(),
		0, 0)

	if *shouldAddEntities {
		*shouldAddEntities = false
		AddUIHearts(systemsManager, gameModel.Entities.Hearts, gameModel.Entities.Player.ComponentHealth.Total)

		AddUICoin(systemsManager)

		// Draw obstacles on appropriate map tiles
		obstacles := DrawObstaclesPerMapTiles(systemsManager, roomsMap, gameModel.AllMapDrawData, currentRoomID, 0, 0)
		systemsManager.AddEntities(obstacles...)

		gameModel.RoomWarps = map[EntityID]EntityConfig{}

		// Iterate through all entity configurations and build entities and add to systems
		for _, c := range roomsMap[*currentRoomID].(*Room).EntityConfigs {
			entity := BuildEntityFromConfig(c, systemsManager.NewEntityID())
			gameModel.EntitiesMap[entity.ID()] = entity
			systemsManager.AddEntity(entity)

			switch c.Category {
			case CategoryWarp:
				gameModel.RoomWarps[entity.ID()] = c
			}
		}
	}

	DrawMask(ui.Window)

	systemsManager.Update()

	if ui.Window.JustPressed(pixelgl.KeyP) {
		*currentState = StatePause
	}

	if ui.Window.JustPressed(pixelgl.KeyX) {
		*currentState = StateOver
	}
}
