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

	allMapDrawData := zelduh.BuildMapDrawData(zelduh.TilemapDir, zelduh.TilemapFiles, zelduh.TileSize)

	roomData := zelduh.NewRoomData()

	roomTransitionManager := zelduh.NewRoomTransitionManager()

	entities := zelduh.Entities{
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
	}

	healthSystem := &zelduh.SystemHealth{}

	spatialSystem := &zelduh.SystemSpatial{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	entitiesMap := zelduh.NewEntityByEntityIDMap()

	roomWarps := map[zelduh.EntityID]zelduh.EntityConfig{}

	collisionSystem := &zelduh.SystemCollision{
		MapBounds: pixel.R(
			zelduh.MapX,
			zelduh.MapY,
			zelduh.MapX+zelduh.MapW,
			zelduh.MapY+zelduh.MapH,
		),
		CollisionHandler: zelduh.CollisionHandler{
			RoomTransitionManager: &roomTransitionManager,
			SystemsManager:        &systemsManager,
			HealthSystem:          healthSystem,
			SpatialSystem:         spatialSystem,
			EntitiesMap:           entitiesMap,
			RoomWarps:             roomWarps,
			Entities:              entities,
		},
	}

	inputSystem := &zelduh.SystemInput{Win: ui.Window}

	spritesheet := zelduh.LoadAndBuildSpritesheet(zelduh.SpritesheetPath, zelduh.TileSize)

	systemsManager.AddSystems(
		inputSystem,
		healthSystem,
		spatialSystem,
		collisionSystem,
		&zelduh.SystemRender{
			Win:         ui.Window,
			Spritesheet: spritesheet,
		},
	)

	systemsManager.AddEntities(
		entities.Player,
		entities.Sword,
		entities.Arrow,
		entities.Bomb,
	)

	gameStateManager := zelduh.NewGameStateManager(
		&roomTransitionManager,
		&systemsManager,
		ui,
		currLocaleMsgs,
		collisionSystem,
		inputSystem,
		spritesheet,
		entitiesMap,
		allMapDrawData,
		roomWarps,
		entities,
		&roomData,
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
