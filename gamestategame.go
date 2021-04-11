package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(
	ui UI,
	mapConfig MapConfig,
	windowConfig WindowConfig,
	spritesheet map[int]*pixel.Sprite,
	roomData *RoomData,
	entities Entities,
	roomWarps map[EntityID]EntityConfig,
	entitiesMap EntityByEntityID,
	allMapDrawData map[string]MapData,
	inputSystem *SystemInput,
	roomsMap Rooms,
	systemsManager *SystemsManager,
	gameStateManager *GameStateManager,
	frameRate int,
) {
	inputSystem.EnablePlayer()

	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, mapConfig, colornames.White)

	DrawMapBackgroundImage(
		ui.Window,
		spritesheet,
		allMapDrawData,
		roomsMap[roomData.CurrentRoomID].MapName(),
		0, 0,
		mapConfig,
	)

	if systemsManager.GetShouldAddEntities() {
		systemsManager.SetShouldAddEntities(false)
		AddUIHearts(systemsManager, entities.Hearts, entities.Player.ComponentHealth.Total)

		AddUICoin(systemsManager, frameRate)

		// Draw obstacles on appropriate map tiles
		obstacles := DrawObstaclesPerMapTiles(systemsManager, roomsMap, allMapDrawData, roomData.CurrentRoomID, 0, 0, mapConfig, frameRate)
		systemsManager.AddEntities(obstacles...)

		// Iterate through all entity configurations and build entities and add to systems
		for _, c := range roomsMap[roomData.CurrentRoomID].(*Room).EntityConfigs {
			entity := BuildEntityFromConfig(c, systemsManager.NewEntityID(), frameRate)
			entitiesMap[entity.ID()] = entity
			systemsManager.AddEntity(entity)

			switch c.Category {
			case CategoryWarp:
				roomWarps[entity.ID()] = c
			}
		}
	}

	DrawMask(ui.Window, windowConfig, mapConfig)

	systemsManager.Update()

	if ui.Window.JustPressed(pixelgl.KeyP) {
		gameStateManager.CurrentState = StatePause
	}

	if ui.Window.JustPressed(pixelgl.KeyX) {
		gameStateManager.CurrentState = StateOver
	}
}
