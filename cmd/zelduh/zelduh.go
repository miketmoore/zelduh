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

	// TileSize defines the width and height of a tile
	const tileSize float64 = 48

	// FrameRate is used to determine which sprite to use for animations
	const frameRate int = 5

	entityConfigPresetFnsMap := zelduh.BuildEntityConfigPresetFnsMap(tileSize)

	entityConfigPresetFnManager := zelduh.NewEntityConfigPresetFnManager(entityConfigPresetFnsMap)

	rooms := zelduh.BuildRooms(&entityConfigPresetFnManager, tileSize)

	zelduh.BuildMapRoomIDToRoom(zelduh.Overworld, rooms)

	systemsManager := zelduh.NewSystemsManager()

	entityFactory := zelduh.NewEntityFactory(&systemsManager, &entityConfigPresetFnManager)

	spatialSystem := zelduh.SpatialSystem{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	healthSystem := &zelduh.HealthSystem{}

	entitiesMap := zelduh.EntitiesMap{}

	roomTransition := zelduh.RoomTransition{
		Start: float64(tileSize),
	}

	roomWarps := zelduh.RoomWarps{}

	shouldAddEntities := true
	var currentRoomID zelduh.RoomID = 1
	var nextRoomID zelduh.RoomID
	currentState := zelduh.StateStart
	spritesheet := zelduh.LoadAndBuildSpritesheet(zelduh.SpritesheetPath, tileSize)

	player := entityFactory.NewEntity("player", 6, 6, frameRate)
	bomb := entityFactory.NewEntity("bomb", 0, 0, frameRate)
	explosion := entityFactory.NewEntity("explosion", 0, 0, frameRate)
	sword := entityFactory.NewEntity("sword", 0, 0, frameRate)
	arrow := entityFactory.NewEntity("arrow", 0, 0, frameRate)
	hearts := []zelduh.Entity{
		entityFactory.NewEntity("heart", 1.5, 14, frameRate),
		entityFactory.NewEntity("heart", 2.15, 14, frameRate),
		entityFactory.NewEntity("heart", 2.80, 14, frameRate),
	}

	collisionSystem := &zelduh.CollisionSystem{
		MapBounds: pixel.R(
			zelduh.MapX,
			zelduh.MapY,
			zelduh.MapX+zelduh.MapW,
			zelduh.MapY+zelduh.MapH,
		),
		CollisionHandler: zelduh.NewCollisionHandler(
			&systemsManager,
			&spatialSystem,
			healthSystem,
			&shouldAddEntities,
			&nextRoomID,
			&currentState,
			&roomTransition,
			entitiesMap,
			&player,
			&sword,
			&explosion,
			&arrow,
			hearts,
			roomWarps,
			&entityConfigPresetFnManager,
			tileSize,
			frameRate,
		),
	}

	ui := zelduh.NewUI(currLocaleMsgs)

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
		player,
		sword,
		arrow,
		bomb,
	)

	mapDrawData := zelduh.BuildMapDrawData(zelduh.TilemapDir, zelduh.TilemapFiles, tileSize)

	gameStateManager := zelduh.NewGameStateManager(
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
		entitiesMap,
		&player,
		hearts,
		roomWarps,
		rooms,
		&entityConfigPresetFnManager,
		tileSize,
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
