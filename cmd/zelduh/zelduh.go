package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"
	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// Map of RoomID to a Room configuration
var roomsMap = zelduh.Rooms{
	1: zelduh.NewRoom("overworldFourWallsDoorBottomRight",
		zelduh.GetPreset("puzzleBox")(5, 5),
		zelduh.GetPreset("floorSwitch")(5, 6),
		zelduh.GetPreset("toggleObstacle")(10, 7),
	),
	2: zelduh.NewRoom("overworldFourWallsDoorTopBottom",
		zelduh.GetPreset("skull")(5, 5),
		zelduh.GetPreset("skeleton")(11, 9),
		zelduh.GetPreset("spinner")(7, 9),
		zelduh.GetPreset("eyeburrower")(8, 9),
	),
	3: zelduh.NewRoom("overworldFourWallsDoorRightTopBottom",
		zelduh.WarpStone(3, 7, 6, 5),
	),
	5: zelduh.NewRoom("rockWithCaveEntrance",
		zelduh.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 11,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 7) + zelduh.TileSize/2,
			Y:            (zelduh.TileSize * 9) + zelduh.TileSize/2,
			Hitbox: &zelduh.HitboxConfig{
				Radius: 30,
			},
		},
		zelduh.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 11,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 8) + zelduh.TileSize/2,
			Y:            (zelduh.TileSize * 9) + zelduh.TileSize/2,
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
		zelduh.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 5,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 6) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Y:            (zelduh.TileSize * 1) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Hitbox: &zelduh.HitboxConfig{
				Radius: 15,
			},
		},
		zelduh.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 5,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 7) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Y:            (zelduh.TileSize * 1) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Hitbox: &zelduh.HitboxConfig{
				Radius: 15,
			},
		},
	),
}

func run() {

	// Just a stub for now since English is the only language supported at this time
	localeMsgs := map[string]map[string]string{
		"en": {
			"gameTitle":             "Zelduh",
			"pauseScreenMessage":    "Paused",
			"gameOverScreenMessage": "Game Over",
		},
		"es": {
			"gameTitle":             "Zelduh",
			"pauseScreenMessage":    "Paused",
			"gameOverScreenMessage": "Game Over",
		},
	}

	currLocaleMsgs := localeMsgs["en"]

	gameWorld := zelduh.New()

	zelduh.BuildMapRoomIDToRoom(zelduh.Overworld, roomsMap)

	// Initialize text
	orig := pixel.V(20, 50)
	txt := text.New(orig, text.Atlas7x13)
	txt.Color = colornames.Black

	// Initialize window
	cfg := pixelgl.WindowConfig{
		Title:  currLocaleMsgs["gameTitle"],
		Bounds: pixel.R(zelduh.WinX, zelduh.WinY, zelduh.WinW, zelduh.WinH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		fmt.Println("Initializing GUI window failed:")
		fmt.Println(err)
		os.Exit(1)
	}

	gameModel := zelduh.GameModel{
		Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		EntitiesMap:   map[zelduh.EntityID]zelduh.Entity{},
		CurrentState:  zelduh.StateStart,
		AddEntities:   true,
		CurrentRoomID: 1,
		RoomTransition: &zelduh.RoomTransition{
			Start: float64(zelduh.TileSize),
		},
		Spritesheet: zelduh.LoadAndBuildSpritesheet(zelduh.SpritesheetPath, zelduh.TileSize),

		// Build entities
		Player:    zelduh.BuildEntityFromConfig(zelduh.GetPreset("player")(6, 6), gameWorld.NewEntityID()),
		Bomb:      zelduh.BuildEntityFromConfig(zelduh.GetPreset("bomb")(0, 0), gameWorld.NewEntityID()),
		Explosion: zelduh.BuildEntityFromConfig(zelduh.GetPreset("explosion")(0, 0), gameWorld.NewEntityID()),
		Sword:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("sword")(0, 0), gameWorld.NewEntityID()),
		Arrow:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("arrow")(0, 0), gameWorld.NewEntityID()),

		RoomWarps:      map[zelduh.EntityID]zelduh.Config{},
		AllMapDrawData: zelduh.BuildMapDrawData(zelduh.TilemapDir, zelduh.TilemapFiles, zelduh.TileSize),

		InputSystem:  &zelduh.SystemInput{Win: win},
		HealthSystem: &zelduh.SystemHealth{},

		Hearts: zelduh.BuildEntitiesFromConfigs(
			gameWorld.NewEntityID,
			zelduh.GetPreset("heart")(1.5, 14),
			zelduh.GetPreset("heart")(2.15, 14),
			zelduh.GetPreset("heart")(2.80, 14),
		),
	}

	gameModel.SpatialSystem = &zelduh.SystemSpatial{
		Rand: gameModel.Rand,
	}

	collisionHandler := zelduh.CollisionHandler{
		GameModel: &gameModel,
		GameWorld: &gameWorld,
	}

	collisionSystem := &zelduh.SystemCollision{
		MapBounds: pixel.R(
			zelduh.MapX,
			zelduh.MapY,
			zelduh.MapX+zelduh.MapW,
			zelduh.MapY+zelduh.MapH,
		),
		OnPlayerCollisionWithBounds:             collisionHandler.OnPlayerCollisionWithBounds,
		OnPlayerCollisionWithCoin:               collisionHandler.OnPlayerCollisionWithCoin,
		OnPlayerCollisionWithEnemy:              collisionHandler.OnPlayerCollisionWithEnemy,
		OnSwordCollisionWithEnemy:               collisionHandler.OnSwordCollisionWithEnemy,
		OnArrowCollisionWithEnemy:               collisionHandler.OnArrowCollisionWithEnemy,
		OnArrowCollisionWithObstacle:            collisionHandler.OnArrowCollisionWithObstacle,
		OnPlayerCollisionWithObstacle:           collisionHandler.OnPlayerCollisionWithObstacle,
		OnPlayerCollisionWithMoveableObstacle:   collisionHandler.OnPlayerCollisionWithMoveableObstacle,
		OnMoveableObstacleCollisionWithSwitch:   collisionHandler.OnMoveableObstacleCollisionWithSwitch,
		OnMoveableObstacleNoCollisionWithSwitch: collisionHandler.OnMoveableObstacleNoCollisionWithSwitch,
		OnEnemyCollisionWithObstacle:            collisionHandler.OnEnemyCollisionWithObstacle,
		OnPlayerCollisionWithSwitch:             collisionHandler.OnPlayerCollisionWithSwitch,
		OnPlayerNoCollisionWithSwitch:           collisionHandler.OnPlayerNoCollisionWithSwitch,
		OnPlayerCollisionWithWarp:               collisionHandler.OnPlayerCollisionWithWarp,
	}

	gameWorld.AddSystems(
		gameModel.InputSystem,
		gameModel.HealthSystem,
		gameModel.SpatialSystem,
		collisionSystem,
		&zelduh.SystemRender{
			Win:         win,
			Spritesheet: gameModel.Spritesheet,
		},
	)

	gameWorld.AddEntities(
		gameModel.Player,
		gameModel.Sword,
		gameModel.Arrow,
		gameModel.Bomb,
	)

	for !win.Closed() {

		// Quit application when user input matches
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		switch gameModel.CurrentState {
		case zelduh.StateStart:
			zelduh.GameStateStart(win, txt, currLocaleMsgs, &gameModel)
		case zelduh.StateGame:
			zelduh.GameStateGame(win, &gameModel, roomsMap, &gameWorld)
		case zelduh.StatePause:
			zelduh.GameStatePause(win, txt, currLocaleMsgs, &gameModel)
		case zelduh.StateOver:
			zelduh.GameStateOver(win, txt, currLocaleMsgs, &gameModel)
		case zelduh.StateMapTransition:
			zelduh.GameStateMapTransition(win, &gameWorld, roomsMap, collisionSystem, &gameModel)
		}

		win.Update()

	}
}

func main() {
	pixelgl.Run(run)
}
