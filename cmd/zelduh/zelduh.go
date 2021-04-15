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

	rooms := zelduh.BuildRooms()

	zelduh.BuildMapRoomIDToRoom(zelduh.Overworld, rooms)

	systemsManager := zelduh.NewSystemsManager()

	entityFactory := zelduh.NewEntityFactory(&systemsManager)

	spatialSystem := zelduh.SpatialSystem{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	healthSystem := &zelduh.HealthSystem{}

	entitiesMap := zelduh.EntitiesMap{}

	roomTransition := zelduh.RoomTransition{
		Start: float64(zelduh.TileSize),
	}

	roomWarps := zelduh.RoomWarps{}

	shouldAddEntities := true
	var currentRoomID zelduh.RoomID = 1
	var nextRoomID zelduh.RoomID
	currentState := zelduh.StateStart
	spritesheet := zelduh.LoadAndBuildSpritesheet(zelduh.SpritesheetPath, zelduh.TileSize)

	player := entityFactory.NewEntity("player", 6, 6)
	bomb := entityFactory.NewEntity("bomb", 0, 0)
	explosion := entityFactory.NewEntity("explosion", 0, 0)
	sword := entityFactory.NewEntity("sword", 0, 0)
	arrow := entityFactory.NewEntity("arrow", 0, 0)
	hearts := []zelduh.Entity{
		entityFactory.NewEntity("heart", 1.5, 14),
		entityFactory.NewEntity("heart", 2.15, 14),
		entityFactory.NewEntity("heart", 2.80, 14),
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

	mapDrawData := zelduh.BuildMapDrawData(zelduh.TilemapDir, zelduh.TilemapFiles, zelduh.TileSize)

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
