package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"

	"github.com/faiface/pixel/pixelgl"
)

const tilemapDir = "assets/tilemaps/"
const spritesheetPath = "assets/spritesheet.png"

func run() {

	// frameRate is used to determine which sprite to use for animations
	const frameRate int = 5

	const tileSize float64 = 48

	var windowConfig zelduh.WindowConfig = zelduh.WindowConfig{
		X:      0,
		Y:      0,
		Width:  800,
		Height: 800,
	}

	var mapConfig zelduh.MapConfig = zelduh.MapConfig{
		Width:  tileSize * 14,
		Height: tileSize * 12,
	}

	mapConfig.X = (windowConfig.Width - mapConfig.Width) / 2
	mapConfig.Y = (windowConfig.Height - mapConfig.Height) / 2

	mapBoundsConfig := zelduh.MapBoundsConfig{
		X:      mapConfig.X,
		Y:      mapConfig.Y,
		Width:  mapConfig.X + mapConfig.Width,
		Height: mapConfig.Y + mapConfig.Height,
	}

	currLocaleMsgs, err := zelduh.GetLocaleMessageMapByLanguage("en")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	systemsManager := zelduh.NewSystemsManager()

	zelduh.BuildMapRoomIDToRoom(zelduh.Overworld, zelduh.RoomsMap)

	ui := zelduh.NewUI(currLocaleMsgs, windowConfig)

	allMapDrawData := zelduh.BuildMapDrawData(tilemapDir, zelduh.TilemapFiles, zelduh.TileSize)

	roomData := zelduh.NewRoomData()

	roomTransitionManager := zelduh.NewRoomTransitionManager()

	entities := zelduh.Entities{
		Player:    zelduh.BuildEntityFromConfig(zelduh.GetPreset("player")(6, 6), systemsManager.NewEntityID(), frameRate),
		Bomb:      zelduh.BuildEntityFromConfig(zelduh.GetPreset("bomb")(0, 0), systemsManager.NewEntityID(), frameRate),
		Explosion: zelduh.BuildEntityFromConfig(zelduh.GetPreset("explosion")(0, 0), systemsManager.NewEntityID(), frameRate),
		Sword:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("sword")(0, 0), systemsManager.NewEntityID(), frameRate),
		Arrow:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("arrow")(0, 0), systemsManager.NewEntityID(), frameRate),
		Hearts: zelduh.BuildEntitiesFromConfigs(
			systemsManager.NewEntityID,
			frameRate,
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

	collisionSystem := zelduh.NewSystemCollision(
		mapBoundsConfig,
		&roomTransitionManager,
		&systemsManager,
		healthSystem,
		spatialSystem,
		entitiesMap,
		roomWarps,
		entities,
	)

	inputSystem := &zelduh.SystemInput{Win: ui.Window}

	spritesheet := zelduh.LoadAndBuildSpritesheet(spritesheetPath, zelduh.TileSize)

	systemsManager.AddSystems(
		inputSystem,
		healthSystem,
		spatialSystem,
		&collisionSystem,
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
		&collisionSystem,
		inputSystem,
		spritesheet,
		entitiesMap,
		allMapDrawData,
		roomWarps,
		entities,
		&roomData,
		mapConfig,
		windowConfig,
		frameRate,
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
