package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {

	currLocaleMsgs, err := zelduh.GetLocaleMessageMapByLanguage("en")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	systemsManager := zelduh.NewSystemsManager()

	zelduh.BuildMapRoomIDToRoom(zelduh.Overworld, zelduh.RoomsMap)

	ui := zelduh.NewUI(currLocaleMsgs)

	shouldAddEntities := true
	var currentRoomID zelduh.RoomID = 1
	var nextRoomID zelduh.RoomID
	currentState := zelduh.StateStart
	spritesheet := zelduh.LoadAndBuildSpritesheet(zelduh.SpritesheetPath, zelduh.TileSize)

	roomTransition := zelduh.RoomTransition{
		Start: float64(zelduh.TileSize),
	}

	gameModel := zelduh.GameModel{
		EntitiesMap: map[zelduh.EntityID]zelduh.Entity{},

		RoomWarps: map[zelduh.EntityID]zelduh.EntityConfig{},

		Entities: zelduh.Entities{
			Player:    zelduh.BuildEntityFromConfig(zelduh.GetPreset("player")(6, 6), systemsManager.NewEntityID()),
			Bomb:      zelduh.BuildEntityFromConfig(zelduh.GetPreset("bomb")(0, 0), systemsManager.NewEntityID()),
			Explosion: zelduh.BuildEntityFromConfig(zelduh.GetPreset("explosion")(0, 0), systemsManager.NewEntityID()),
			Sword:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("sword")(0, 0), systemsManager.NewEntityID()),
			Arrow:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("arrow")(0, 0), systemsManager.NewEntityID()),
			Hearts: zelduh.BuildEntitiesFromConfigs(
				systemsManager.NewEntityID,
				zelduh.GetPreset("heart")(1.5, 14),
				zelduh.GetPreset("heart")(2.15, 14),
				zelduh.GetPreset("heart")(2.80, 14),
			),
		},
	}

	mapDrawData := zelduh.BuildMapDrawData(zelduh.TilemapDir, zelduh.TilemapFiles, zelduh.TileSize)

	spatialSystem := zelduh.SpatialSystem{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	healthSystem := &zelduh.HealthSystem{}

	collisionSystem := &zelduh.CollisionSystem{
		MapBounds: pixel.R(
			zelduh.MapX,
			zelduh.MapY,
			zelduh.MapX+zelduh.MapW,
			zelduh.MapY+zelduh.MapH,
		),
		CollisionHandler: zelduh.CollisionHandler{
			GameModel:         &gameModel,
			SystemsManager:    &systemsManager,
			SpatialSystem:     &spatialSystem,
			HealthSystem:      healthSystem,
			ShouldAddEntities: &shouldAddEntities,
			CurrentState:      &currentState,
			RoomTransition:    &roomTransition,
			NextRoomID:        &nextRoomID,
		},
	}

	inputSystem := &zelduh.InputSystem{Win: ui.Window}

	systemsManager.AddSystems(
		inputSystem,
		healthSystem,
		&spatialSystem,
		collisionSystem,
		&zelduh.RenderSystem{
			Win:         ui.Window,
			Spritesheet: spritesheet,
		},
	)

	systemsManager.AddEntities(
		gameModel.Entities.Player,
		gameModel.Entities.Sword,
		gameModel.Entities.Arrow,
		gameModel.Entities.Bomb,
	)

	gameStateManager := zelduh.NewGameStateManager(
		&gameModel,
		&systemsManager,
		ui,
		currLocaleMsgs,
		collisionSystem,
		inputSystem,
		&shouldAddEntities,
		&currentRoomID,
		&nextRoomID,
		&currentState,
		spritesheet,
		mapDrawData,
		&roomTransition,
	)

	for !ui.Window.Closed() {

		// Quit application when user input matches
		if ui.Window.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		gameStateManager.Update()

		ui.Window.Update()

	}
}

func main() {
	pixelgl.Run(run)
}
