package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(
	ui UI,
	roomsMap Rooms,
	systemsManager *SystemsManager,
	inputSystem *InputSystem,
	shouldAddEntities *bool,
	currentRoomID *RoomID,
	currentState *State,
	spritesheet Spritesheet,
	mapDrawData MapDrawData,
	entitiesMap EntitiesMap,
	player *Entity,
	hearts []Entity,
	roomWarps RoomWarps,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
	frameRate int,
	nonObstacleSprites map[int]bool,
	windowConfig WindowConfig,
	activeSpaceRectangle ActiveSpaceRectangle,
) {
	inputSystem.EnablePlayer()

	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, activeSpaceRectangle, colornames.White)

	DrawMapBackgroundImage(
		ui.Window,
		spritesheet,
		mapDrawData,
		roomsMap[*currentRoomID].RoomName(),
		0, 0,
		tileSize,
		activeSpaceRectangle,
	)

	if *shouldAddEntities {
		*shouldAddEntities = false
		AddUIHearts(systemsManager, hearts, player.ComponentHealth.Total)

		AddUICoin(systemsManager, entityConfigPresetFnManager, frameRate)

		// Draw obstacles on appropriate map tiles
		obstacles := DrawObstaclesPerMapTiles(
			systemsManager,
			entityConfigPresetFnManager,
			roomsMap,
			mapDrawData,
			currentRoomID,
			0, 0,
			tileSize,
			frameRate,
			nonObstacleSprites,
			activeSpaceRectangle,
		)
		systemsManager.AddEntities(obstacles...)

		for k := range roomWarps {
			delete(roomWarps, k)
		}

		// Iterate through all entity configurations and build entities and add to systems
		currentRoom := roomsMap[*currentRoomID]
		for _, c := range currentRoom.EntityConfigs {
			entity := BuildEntityFromConfig(c, systemsManager.NewEntityID(), frameRate)
			entitiesMap[entity.ID()] = entity
			systemsManager.AddEntity(entity)

			switch c.Category {
			case CategoryWarp:
				roomWarps[entity.ID()] = c
			}
		}
	}

	DrawMask(ui.Window, windowConfig, activeSpaceRectangle)

	systemsManager.Update()

	if ui.Window.JustPressed(pixelgl.KeyP) {
		*currentState = StatePause
	}

	if ui.Window.JustPressed(pixelgl.KeyX) {
		*currentState = StateOver
	}
}
