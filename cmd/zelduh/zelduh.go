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

	rooms := BuildRooms(&entityConfigPresetFnManager, tileSize)

	zelduh.BuildMapRoomIDToRoom(
		// Overworld is a multi-dimensional array representing the overworld
		// Each room ID should be unique
		[][]zelduh.RoomID{
			{1, 10},
			{2, 0, 0, 8},
			{3, 5, 6, 7},
			{9},
			{11},
		},
		rooms,
	)

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
	spritesheet := zelduh.LoadAndBuildSpritesheet("assets/spritesheet.png", tileSize)

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

	windowConfig := zelduh.WindowConfig{
		X:      0,
		Y:      0,
		Width:  800,
		Height: 800,
	}

	mapConfig := zelduh.MapConfig{
		Width:  tileSize * 14,
		Height: tileSize * 12,
	}

	mapConfig.X = (windowConfig.Width - mapConfig.Width) / 2
	mapConfig.Y = (windowConfig.Height - mapConfig.Height) / 2

	collisionSystem := &zelduh.CollisionSystem{
		MapBounds: pixel.R(
			mapConfig.X,
			mapConfig.Y,
			mapConfig.X+mapConfig.Width,
			mapConfig.Y+mapConfig.Height,
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

	ui := zelduh.NewUI(currLocaleMsgs, windowConfig)

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

	mapDrawData := zelduh.BuildMapDrawData(
		"assets/tilemaps/",
		[]string{
			"overworldOpen",
			"overworldOpenCircleOfTrees",
			"overworldFourWallsDoorBottom",
			"overworldFourWallsDoorLeftTop",
			"overworldFourWallsDoorRightTop",
			"overworldFourWallsDoorTopBottom",
			"overworldFourWallsDoorRightTopBottom",
			"overworldFourWallsDoorBottomRight",
			"overworldFourWallsDoorTop",
			"overworldFourWallsDoorRight",
			"overworldFourWallsDoorLeft",
			"overworldTreeClusterTopRight",
			"overworldFourWallsClusterTrees",
			"overworldFourWallsDoorsAllSides",
			"rockPatternTest",
			"rockPathOpenLeft",
			"rockWithCaveEntrance",
			"rockPathLeftRightEntrance",
			"test",
			"dungeonFourDoors",
		},
		tileSize,
	)

	// NonObstacleSprites defines which sprites are not obstacles
	var nonObstacleSprites = map[int]bool{
		8:   true,
		9:   true,
		24:  true,
		37:  true,
		38:  true,
		52:  true,
		53:  true,
		66:  true,
		86:  true,
		136: true,
		137: true,
	}

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
		nonObstacleSprites,
		windowConfig,
		mapConfig,
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

// TODO move to zelduh cmd file since it is configuration
// Map of RoomID to a Room configuration
func BuildRooms(entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager, tileSize float64) zelduh.Rooms {
	return zelduh.Rooms{
		1: zelduh.NewRoom("overworldFourWallsDoorBottomRight",
			entityConfigPresetFnManager.GetPreset("puzzleBox")(5, 5),
			entityConfigPresetFnManager.GetPreset("floorSwitch")(5, 6),
			entityConfigPresetFnManager.GetPreset("toggleObstacle")(10, 7),
		),
		2: zelduh.NewRoom("overworldFourWallsDoorTopBottom",
			entityConfigPresetFnManager.GetPreset("skull")(5, 5),
			entityConfigPresetFnManager.GetPreset("skeleton")(11, 9),
			entityConfigPresetFnManager.GetPreset("spinner")(7, 9),
			entityConfigPresetFnManager.GetPreset("eyeburrower")(8, 9),
		),
		3: zelduh.NewRoom("overworldFourWallsDoorRightTopBottom",
			WarpStone(entityConfigPresetFnManager, 3, 7, 6, 5),
		),
		5: zelduh.NewRoom("rockWithCaveEntrance",
			zelduh.EntityConfig{
				Category:     zelduh.CategoryWarp,
				WarpToRoomID: 11,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 7) + tileSize/2,
				Y:            (tileSize * 9) + tileSize/2,
				Hitbox: &zelduh.HitboxConfig{
					Radius: 30,
				},
			},
			zelduh.EntityConfig{
				Category:     zelduh.CategoryWarp,
				WarpToRoomID: 11,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 8) + tileSize/2,
				Y:            (tileSize * 9) + tileSize/2,
				Hitbox: &zelduh.HitboxConfig{
					Radius: 30,
				},
			},
		),
		6:  zelduh.NewRoom("rockPathLeftRightEntrance"),
		7:  zelduh.NewRoom("overworldFourWallsDoorLeftTop"),
		8:  zelduh.NewRoom("overworldFourWallsDoorBottom"),
		9:  zelduh.NewRoom("overworldFourWallsDoorTop"),
		10: zelduh.NewRoom("overworldFourWallsDoorLeft"),
		11: zelduh.NewRoom("dungeonFourDoors",
			// South door of cave - warp to cave entrance
			zelduh.EntityConfig{
				Category:     zelduh.CategoryWarp,
				WarpToRoomID: 5,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 6) + tileSize + (tileSize / 2.5),
				Y:            (tileSize * 1) + tileSize + (tileSize / 2.5),
				Hitbox: &zelduh.HitboxConfig{
					Radius: 15,
				},
			},
			zelduh.EntityConfig{
				Category:     zelduh.CategoryWarp,
				WarpToRoomID: 5,
				W:            tileSize,
				H:            tileSize,
				X:            (tileSize * 7) + tileSize + (tileSize / 2.5),
				Y:            (tileSize * 1) + tileSize + (tileSize / 2.5),
				Hitbox: &zelduh.HitboxConfig{
					Radius: 15,
				},
			},
		),
	}
}

// WarpStone returns an entity config for a warp stone
func WarpStone(entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager, X, Y, WarpToRoomID, HitBoxRadius float64) zelduh.EntityConfig {
	presetFn := entityConfigPresetFnManager.GetPreset("warpStone")
	e := presetFn(X, Y)
	e.WarpToRoomID = 6
	e.Hitbox.Radius = 5
	return e
}
