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

// TODO move this to a higher level configuration location
func BuildEntityConfigPresetFnsMap(tileSize float64) map[string]zelduh.EntityConfigPresetFn {
	return map[string]zelduh.EntityConfigPresetFn{
		"arrow": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryArrow,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				W: tileSize,
				H: tileSize,
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
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
		"bomb": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryBomb,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				W: tileSize,
				H: tileSize,
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("bomb"),
				},
				Hitbox: &zelduh.HitboxConfig{
					Radius: 5,
				},
				Ignore: true,
			}
		},
		"coin": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryCoin,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("coin"),
				},
			}
		},
		"explosion": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category:   zelduh.CategoryExplosion,
				Expiration: 12,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("explosion"),
				},
			}
		},
		"obstacle": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryObstacle,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
			}
		},
		"player": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryPlayer,
				Health:   3,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
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
		"sword": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategorySword,
				Movement: &zelduh.MovementConfig{
					Direction: zelduh.DirectionDown,
					Speed:     0.0,
				},
				W: tileSize,
				H: tileSize,
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
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
		"eyeburrower": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
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
		"heart": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryHeart,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				Hitbox: &zelduh.HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("heart"),
				},
			}

		},
		"skeleton": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
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
		"skull": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
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
		"spinner": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryEnemy,
				W:        tileSize, H: tileSize, X: tileSize * xTiles, Y: tileSize * yTiles,
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
		"uiCoin": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryHeart,
				W:        tileSize,
				H:        tileSize,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				Hitbox: &zelduh.HitboxConfig{
					Box: imdraw.New(nil),
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("uiCoin"),
				},
			}
		},
		"warpStone": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryWarp,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				W:        tileSize,
				H:        tileSize,
				Hitbox: &zelduh.HitboxConfig{
					Box:    imdraw.New(nil),
					Radius: 20,
				},
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("warpStone"),
				},
			}
		},
		"puzzleBox": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryMovableObstacle,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				W:        tileSize,
				H:        tileSize,
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
		"floorSwitch": func(xTiles, yTiles float64) zelduh.EntityConfig {
			return zelduh.EntityConfig{
				Category: zelduh.CategoryCollisionSwitch,
				X:        tileSize * xTiles,
				Y:        tileSize * yTiles,
				W:        tileSize,
				H:        tileSize,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("floorSwitch"),
				},
				Toggleable: true,
			}
		},
		// this is an impassable obstacle that can be toggled "remotely"
		// it has two visual states that coincide with each toggle state
		"toggleObstacle": func(xTiles, yTiles float64) zelduh.EntityConfig {
			// TODO get this working again
			return zelduh.EntityConfig{
				X: tileSize * xTiles,
				Y: tileSize * yTiles,
				W: tileSize,
				H: tileSize,
				Animation: zelduh.AnimationConfig{
					"default": zelduh.GetSpriteSet("toggleObstacle"),
				},
				// Impassable: true,
				Toggleable: true,
			}
		},
	}
}
