package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

type StateGame struct {
	context           *StateContext
	uiSystem          *UISystem
	inputSystem       *InputSystem
	currentRoomID     *RoomID
	shouldAddEntities *bool
	levelManager      *LevelManager
	entityCreator     *EntityCreator
	systemsManager    *SystemsManager
	roomWarps         RoomWarps
	entitiesMap       EntitiesMap
	frameRate         int
}

func NewStateGame(
	context *StateContext,
	uiSystem *UISystem,
	inputSystem *InputSystem,
	currentRoomID *RoomID,
	shouldAddEntities *bool,
	levelManager *LevelManager,
	entityCreator *EntityCreator,
	systemsManager *SystemsManager,
	roomWarps RoomWarps,
	entitiesMap EntitiesMap,
	frameRate int,
) State {
	return StateGame{
		context:           context,
		uiSystem:          uiSystem,
		inputSystem:       inputSystem,
		currentRoomID:     currentRoomID,
		shouldAddEntities: shouldAddEntities,
		levelManager:      levelManager,
		entityCreator:     entityCreator,
		systemsManager:    systemsManager,
		roomWarps:         roomWarps,
		entitiesMap:       entitiesMap,
		frameRate:         frameRate,
	}
}

func (g StateGame) Update() error {
	g.inputSystem.Enable()

	roomName := g.levelManager.CurrentLevel.RoomByIDMap[*g.currentRoomID].Name

	g.uiSystem.DrawLevelBackground(roomName)

	if *g.shouldAddEntities {
		*g.shouldAddEntities = false

		g.entityCreator.CreateUICoin()

		// Draw obstacles on appropriate map tiles
		obstacles := g.uiSystem.DrawObstaclesPerMapTiles(
			g.currentRoomID,
			0, 0,
		)
		g.systemsManager.AddEntities(obstacles...)

		for k := range g.roomWarps {
			delete(g.roomWarps, k)
		}

		// Iterate through all entity configurations and build entities and add to systems
		currentRoom := g.levelManager.CurrentLevel.RoomByIDMap[*g.currentRoomID]
		for _, c := range currentRoom.EntityConfigs {
			entity := g.entityCreator.entityFactory.NewEntity2(c, g.frameRate)
			g.entitiesMap[entity.ID()] = entity
			g.systemsManager.AddEntity(entity)

			switch c.Category {
			case CategoryWarp:
				g.roomWarps[entity.ID()] = c
			}
		}
	}

	g.uiSystem.DrawMask()

	err := g.systemsManager.Update()
	if err != nil {
		return err
	}

	if g.uiSystem.Window.JustPressed(pixelgl.KeyP) {
		err := g.context.SetState(StateNamePause)
		if err != nil {
			return err
		}
	}

	if g.uiSystem.Window.JustPressed(pixelgl.KeyX) {
		err := g.context.SetState(StateNameGameOver)
		if err != nil {
			return err
		}
	}

	return nil
}
