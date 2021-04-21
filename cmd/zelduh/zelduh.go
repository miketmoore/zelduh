package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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

	entityConfigPresetFnsMap := BuildEntityConfigPresetFnsMap(tileSize)

	entityConfigPresetFnManager := zelduh.NewEntityConfigPresetFnManager(entityConfigPresetFnsMap)

	testLevel := buildTestLevel(
		&entityConfigPresetFnManager,
		tileSize,
	)

	levelManager := zelduh.NewLevelManager(&testLevel)

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

	player := entityFactory.NewEntity("player", zelduh.Coordinates{X: 6, Y: 6}, frameRate)
	bomb := entityFactory.NewEntity("bomb", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	explosion := entityFactory.NewEntity("explosion", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	sword := entityFactory.NewEntity("sword", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	arrow := entityFactory.NewEntity("arrow", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	hearts := []zelduh.Entity{
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 1.5, Y: 14}, frameRate),
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 2.15, Y: 14}, frameRate),
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 2.80, Y: 14}, frameRate),
	}

	windowConfig := zelduh.WindowConfig{
		X:      0,
		Y:      0,
		Width:  800,
		Height: 800,
	}

	activeSpaceRectangle := zelduh.ActiveSpaceRectangle{
		Width:  tileSize * 14,
		Height: tileSize * 12,
	}

	activeSpaceRectangle.X = (windowConfig.Width - activeSpaceRectangle.Width) / 2
	activeSpaceRectangle.Y = (windowConfig.Height - activeSpaceRectangle.Height) / 2

	collisionHandler := zelduh.NewCollisionHandler(
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
	)

	collisionSystem := zelduh.NewCollisionSystem(
		pixel.R(
			activeSpaceRectangle.X,
			activeSpaceRectangle.Y,
			activeSpaceRectangle.X+activeSpaceRectangle.Width,
			activeSpaceRectangle.Y+activeSpaceRectangle.Height,
		),
		&collisionHandler,
	)

	ui := zelduh.NewUI(currLocaleMsgs, windowConfig)

	inputSystem := &zelduh.InputSystem{Win: ui.Window}

	systemsManager.AddSystems(
		inputSystem,
		healthSystem,
		&spatialSystem,
		&collisionSystem,
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
		&collisionSystem,
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
		&levelManager,
		&entityConfigPresetFnManager,
		tileSize,
		frameRate,
		nonObstacleSprites,
		windowConfig,
		activeSpaceRectangle,
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

func buildTestLevel(
	entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager,
	tileSize float64,
) zelduh.Level {
	// Build a map of RoomIDs to Room structs
	roomByIDMap := BuildRooms(entityConfigPresetFnManager, tileSize)

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
		// This is mutated
		roomByIDMap,
	)

	return zelduh.Level{
		RoomByIDMap: roomByIDMap,
	}
}

// TODO move to zelduh cmd file since it is configuration
// Map of RoomID to a Room configuration
func BuildRooms(entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager, tileSize float64) zelduh.RoomByIDMap {

	buildWarpStone := buildWarpStoneFnFactory(entityConfigPresetFnManager)

	buildWarp := buildWarpFnFactory(tileSize, zelduh.Dimensions{
		Width:  tileSize,
		Height: tileSize,
	})

	return zelduh.RoomByIDMap{
		1: zelduh.NewRoom("overworldFourWallsDoorBottomRight",
			entityConfigPresetFnManager.GetPreset(PresetNamePuzzleBox)(zelduh.Coordinates{X: 5, Y: 5}),
			entityConfigPresetFnManager.GetPreset(PresetNameFloorSwitch)(zelduh.Coordinates{X: 5, Y: 6}),
			entityConfigPresetFnManager.GetPreset(PresetNameToggleObstacle)(zelduh.Coordinates{X: 10, Y: 7}),
		),
		2: zelduh.NewRoom("overworldFourWallsDoorTopBottom",
			entityConfigPresetFnManager.GetPreset(PresetNameEnemySkull)(zelduh.Coordinates{X: 5, Y: 5}),
			entityConfigPresetFnManager.GetPreset(PresetNameEnemySkeleton)(zelduh.Coordinates{X: 11, Y: 9}),
			entityConfigPresetFnManager.GetPreset(PresetNameEnemySpinner)(zelduh.Coordinates{X: 7, Y: 9}),
			entityConfigPresetFnManager.GetPreset(PresetNameEnemyEyeBurrower)(zelduh.Coordinates{X: 8, Y: 9}),
		),
		3: zelduh.NewRoom("overworldFourWallsDoorRightTopBottom",
			buildWarpStone(6, zelduh.Coordinates{X: 3, Y: 7}, 5),
		),
		5: zelduh.NewRoom("rockWithCaveEntrance",
			buildWarp(
				11,
				zelduh.Coordinates{
					X: (tileSize * 7) + tileSize/2,
					Y: (tileSize * 9) + tileSize/2,
				},
				30,
			),
			buildWarp(
				11,
				zelduh.Coordinates{
					X: (tileSize * 8) + tileSize/2,
					Y: (tileSize * 9) + tileSize/2,
				},
				30,
			),
		),
		6:  zelduh.NewRoom("rockPathLeftRightEntrance"),
		7:  zelduh.NewRoom("overworldFourWallsDoorLeftTop"),
		8:  zelduh.NewRoom("overworldFourWallsDoorBottom"),
		9:  zelduh.NewRoom("overworldFourWallsDoorTop"),
		10: zelduh.NewRoom("overworldFourWallsDoorLeft"),
		11: zelduh.NewRoom("dungeonFourDoors",
			// South door of cave - warp to cave entrance
			buildWarp(
				5,
				zelduh.Coordinates{
					X: (tileSize * 6) + tileSize + (tileSize / 2.5),
					Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
				},
				15,
			),
			buildWarp(
				5,
				zelduh.Coordinates{
					X: (tileSize * 7) + tileSize + (tileSize / 2.5),
					Y: (tileSize * 1) + tileSize + (tileSize / 2.5),
				},
				15,
			),
		),
	}
}

type BuildWarpFn func(
	warpToRoomID zelduh.RoomID,
	coordinates zelduh.Coordinates,
	hitboxRadius float64,
) zelduh.EntityConfig

func buildWarpFnFactory(
	tileSize float64,
	dimensions zelduh.Dimensions,
) BuildWarpFn {

	return func(
		warpToRoomID zelduh.RoomID,
		coordinates zelduh.Coordinates,
		hitboxRadius float64,
	) zelduh.EntityConfig {
		return zelduh.EntityConfig{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: warpToRoomID,
			Dimensions:   dimensions,
			Coordinates:  coordinates,
			Hitbox: &zelduh.HitboxConfig{
				Radius: hitboxRadius,
			},
		}
	}
}

type BuildWarpStoneFn func(
	WarpToRoomID zelduh.RoomID,
	coordinates zelduh.Coordinates,
	HitBoxRadius float64,
) zelduh.EntityConfig

func buildWarpStoneFnFactory(
	entityConfigPresetFnManager *zelduh.EntityConfigPresetFnManager,
) BuildWarpStoneFn {
	return func(
		WarpToRoomID zelduh.RoomID,
		coordinates zelduh.Coordinates,
		HitBoxRadius float64,
	) zelduh.EntityConfig {
		presetFn := entityConfigPresetFnManager.GetPreset("warpStone")
		e := presetFn(zelduh.Coordinates{X: coordinates.X, Y: coordinates.Y})
		e.WarpToRoomID = 6
		e.Hitbox.Radius = 5
		return e
	}
}

const (
	PresetNameArrow            zelduh.PresetName = "arrow"
	PresetNameBomb             zelduh.PresetName = "bomb"
	PresetNameCoin             zelduh.PresetName = "coin"
	PresetNameExplosion        zelduh.PresetName = "explosion"
	PresetNameObstacle         zelduh.PresetName = "obstacle"
	PresetNamePlayer           zelduh.PresetName = "player"
	PresetNameFloorSwitch      zelduh.PresetName = "floorSwitch"
	PresetNameToggleObstacle   zelduh.PresetName = "toggleObstacle"
	PresetNamePuzzleBox        zelduh.PresetName = "puzzleBox"
	PresetNameWarpStone        zelduh.PresetName = "warpStone"
	PresetNameUICoin           zelduh.PresetName = "uiCoin"
	PresetNameEnemySpinner     zelduh.PresetName = "spinner"
	PresetNameEnemySkull       zelduh.PresetName = "skull"
	PresetNameEnemySkeleton    zelduh.PresetName = "skeleton"
	PresetNameHeart            zelduh.PresetName = "heart"
	PresetNameEnemyEyeBurrower zelduh.PresetName = "eyeBurrower"
	PresetNameSword            zelduh.PresetName = "sword"
)

// TODO move this to a higher level configuration location
func BuildEntityConfigPresetFnsMap(tileSize float64) map[zelduh.PresetName]zelduh.EntityConfigPresetFn {

	dimensions := zelduh.Dimensions{
		Width:  tileSize,
		Height: tileSize,
	}

	buildCoordinates := func(coordinates zelduh.Coordinates) zelduh.Coordinates {
		return zelduh.Coordinates{
			X: tileSize * coordinates.X,
			Y: tileSize * coordinates.Y,
		}
	}

	return map[zelduh.PresetName]zelduh.EntityConfigPresetFn{
		PresetNameArrow: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryArrow,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"up":    zelduh.GetSpriteSet("arrowUp"),
					"right": zelduh.GetSpriteSet("arrowRight"),
					"down":  zelduh.GetSpriteSet("arrowDown"),
					"left":  zelduh.GetSpriteSet("arrowLeft"),
				},
				Hitbox: &zelduh.HitboxConfig{
					Radius: 5,
				},
				Ignore: true,
			}
		},
		PresetNameBomb: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryBomb,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("bomb"),
				},
				Hitbox: &zelduh.HitboxConfig{
					Radius: 5,
				},
				Ignore: true,
			}
		},
		PresetNameCoin: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryCoin,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("coin"),
				},
			}
		},
		PresetNameExplosion: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:   zelduh.CategoryExplosion,
				Expiration: 12,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("explosion"),
				},
			}
		},
		PresetNameObstacle: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryObstacle,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
			}
		},
		PresetNamePlayer: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryPlayer,
				Health:      3,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &zelduh.HitboxConfig{
					Box:                  imdraw.New(nil),
					Radius:               15,
					CollisionWithRectMod: 5,
				},
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					MaxSpeed:  7.0,
					Speed:     0.0,
				},
				Coins: true,
				Dash: &zelduh.DashConfig{
					Charge:    0,
					MaxCharge: 50,
					SpeedMod:  7,
				},
				Animation: zelduh.AnimationConfig{
					"up":               zelduh.GetSpriteSet("playerUp"),
					"right":            zelduh.GetSpriteSet("playerRight"),
					"down":             zelduh.GetSpriteSet("playerDown"),
					"left":             zelduh.GetSpriteSet("playerLeft"),
					"swordAttackUp":    zelduh.GetSpriteSet("playerSwordUp"),
					"swordAttackRight": zelduh.GetSpriteSet("playerSwordRight"),
					"swordAttackLeft":  zelduh.GetSpriteSet("playerSwordLeft"),
					"swordAttackDown":  zelduh.GetSpriteSet("playerSwordDown"),
				},
			}
		},
		PresetNameSword: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategorySword,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"up":    zelduh.GetSpriteSet("swordUp"),
					"right": zelduh.GetSpriteSet("swordRight"),
					"down":  zelduh.GetSpriteSet("swordDown"),
					"left":  zelduh.GetSpriteSet("swordLeft"),
				},
				Hitbox: &zelduh.HitboxConfig{
					Radius: 20,
				},
				Ignore: true,
			}
		},
		PresetNameEnemyEyeBurrower: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("eyeburrower"),
				},
				Health: 2,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		PresetNameHeart: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryHeart,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &zelduh.HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("heart"),
				},
			}

		},
		PresetNameEnemySkeleton: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("skeleton"),
				},
				Health: 2,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		PresetNameEnemySkull: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("skull"),
				},
				Health: 2,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionDown,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "random",
				},
			}
		},
		PresetNameEnemySpinner: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryEnemy,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("spinner"),
				},
				Invincible: true,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Movement: &zelduh.MovementConfig{
					Direction:    zelduh.DirectionRight,
					Speed:        1.0,
					MaxSpeed:     1.0,
					HitSpeed:     10.0,
					HitBackMoves: 10,
					MaxMoves:     100,
					PatternName:  "left-right",
				},
			}
		},
		PresetNameUICoin: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryHeart,
				Dimensions:  dimensions,
				Coordinates: buildCoordinates(coordinates),
				Hitbox: &zelduh.HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("uiCoin"),
				},
			}
		},
		PresetNameWarpStone: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryWarp,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("warpStone"),
				},
			}
		},
		PresetNamePuzzleBox: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryMovableObstacle,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("puzzleBox"),
				},
				Movement: &zelduh.MovementConfig{
					Speed:    1.0,
					MaxMoves: int(tileSize) / 2,
					MaxSpeed: 2.0,
				},
			}
		},
		PresetNameFloorSwitch: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:    zelduh.CategoryCollisionSwitch,
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("floorSwitch"),
				},
				Toggleable: true,
			}
		},
		// this is an impassable obstacle that can be toggled "remotely"
		// it has two visual states that coincide with each toggle state
		PresetNameToggleObstacle: func(coordinates zelduh.Coordinates) zelduh.EntityConfig {
			// TODO get this working again
			return zelduh.EntityConfig{
				Coordinates: buildCoordinates(coordinates),
				Dimensions:  dimensions,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("toggleObstacle"),
				},
				// Impassable: true,
				Toggleable: true,
			}
		},
	}
}
