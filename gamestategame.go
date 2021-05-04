package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateGame(
	ui UISystem,
	roomByIDMap RoomByIDMap,
	systemsManager *SystemsManager,
	inputSystem *InputSystem,
	shouldAddEntities *bool,
	currentRoomID *RoomID,
	currentState *State,
	spriteMap SpriteMap,
	mapDrawData MapDrawData,
	entitiesMap EntitiesMap,
	player *Entity,
	roomWarps RoomWarps,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
	frameRate int,
	nonObstacleSprites map[int]bool,
	windowConfig WindowConfig,
	activeSpaceRectangle ActiveSpaceRectangle,
) error {
	inputSystem.Enable()

	ui.Window.Clear(colornames.Darkgray)
	ui.DrawMapBackground(colornames.White)

	ui.DrawMapBackgroundImage(
		mapDrawData,
		roomByIDMap[*currentRoomID].Name,
		0, 0,
		tileSize,
		activeSpaceRectangle,
	)

	if *shouldAddEntities {
		*shouldAddEntities = false

		AddUICoin(systemsManager, entityConfigPresetFnManager, frameRate)

		// Draw obstacles on appropriate map tiles
		obstacles := DrawObstaclesPerMapTiles(
			systemsManager,
			entityConfigPresetFnManager,
			roomByIDMap,
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
		currentRoom := roomByIDMap[*currentRoomID]
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

	err := systemsManager.Update()
	if err != nil {
		return err
	}

	if ui.Window.JustPressed(pixelgl.KeyP) {
		*currentState = StatePause
	}

	if ui.Window.JustPressed(pixelgl.KeyX) {
		*currentState = StateOver
	}

	return nil
}
