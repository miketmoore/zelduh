package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

func GameStateGame(
	ui UISystem,
	levelManager *LevelManager,
	systemsManager *SystemsManager,
	inputSystem *InputSystem,
	shouldAddEntities *bool,
	currentRoomID *RoomID,
	currentState *State,
	entitiesMap EntitiesMap,
	roomWarps RoomWarps,
	tileSize float64,
	frameRate int,
	activeSpaceRectangle ActiveSpaceRectangle,
	entityCreator *EntityCreator,
) error {
	inputSystem.Enable()

	roomName := levelManager.CurrentLevel.RoomByIDMap[*currentRoomID].Name

	ui.DrawLevelBackground(roomName)

	if *shouldAddEntities {
		*shouldAddEntities = false

		entityCreator.CreateUICoin()

		// Draw obstacles on appropriate map tiles
		obstacles := ui.DrawObstaclesPerMapTiles(
			currentRoomID,
			0, 0,
		)
		systemsManager.AddEntities(obstacles...)

		for k := range roomWarps {
			delete(roomWarps, k)
		}

		// Iterate through all entity configurations and build entities and add to systems
		currentRoom := levelManager.CurrentLevel.RoomByIDMap[*currentRoomID]
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

	ui.DrawMask()

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
