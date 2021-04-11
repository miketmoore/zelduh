package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(ui UI, spritesheet map[int]*pixel.Sprite, roomData *RoomData, entities Entities, roomWarps map[EntityID]Config, entitiesMap EntityByEntityID, allMapDrawData map[string]MapData, inputSystem *SystemInput, roomsMap Rooms, systemsManager *SystemsManager, gameStateManager *GameStateManager) {
	inputSystem.EnablePlayer()

	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)

	DrawMapBackgroundImage(
		ui.Window,
		spritesheet,
		allMapDrawData,
		roomsMap[roomData.CurrentRoomID].MapName(),
		0, 0)

	if systemsManager.GetShouldAddEntities() {
		systemsManager.SetShouldAddEntities(false)
		AddUIHearts(systemsManager, entities.Hearts, entities.Player.ComponentHealth.Total)

		AddUICoin(systemsManager)

		// Draw obstacles on appropriate map tiles
		obstacles := DrawObstaclesPerMapTiles(systemsManager, roomsMap, allMapDrawData, roomData.CurrentRoomID, 0, 0)
		systemsManager.AddEntities(obstacles...)

		// Iterate through all entity configurations and build entities and add to systems
		for _, c := range roomsMap[roomData.CurrentRoomID].(*Room).EntityConfigs {
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
		gameStateManager.CurrentState = StatePause
	}

	if ui.Window.JustPressed(pixelgl.KeyX) {
		gameStateManager.CurrentState = StateOver
	}
}
