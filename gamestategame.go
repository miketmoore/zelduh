package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

func (g *GameStateManager) stateGame() error {
	g.InputSystem.Enable()

	roomName := g.LevelManager.CurrentLevel.RoomByIDMap[*g.CurrentRoomID].Name

	g.UI.DrawLevelBackground(roomName)

	if *g.ShouldAddEntities {
		*g.ShouldAddEntities = false

		g.entityCreator.CreateUICoin()

		// Draw obstacles on appropriate map tiles
		obstacles := g.UI.DrawObstaclesPerMapTiles(
			g.CurrentRoomID,
			0, 0,
		)
		g.SystemsManager.AddEntities(obstacles...)

		for k := range g.RoomWarps {
			delete(g.RoomWarps, k)
		}

		// Iterate through all entity configurations and build entities and add to systems
		currentRoom := g.LevelManager.CurrentLevel.RoomByIDMap[*g.CurrentRoomID]
		for _, c := range currentRoom.EntityConfigs {
			entity := g.entityCreator.entityFactory.NewEntity2(c, g.FrameRate)
			g.EntitiesMap[entity.ID()] = entity
			g.SystemsManager.AddEntity(entity)

			switch c.Category {
			case CategoryWarp:
				g.RoomWarps[entity.ID()] = c
			}
		}
	}

	g.UI.DrawMask()

	err := g.SystemsManager.Update()
	if err != nil {
		return err
	}

	if g.UI.Window.JustPressed(pixelgl.KeyP) {
		g.setCurrentState(StatePause)
	}

	if g.UI.Window.JustPressed(pixelgl.KeyX) {
		g.setCurrentState(StateOver)
	}

	return nil
}
