package main

import (
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {

	currLocaleMsgs := zelduh.LocaleMessages["en"]

	systemsManager := zelduh.NewSystemsManager()

	zelduh.BuildMapRoomIDToRoom(zelduh.Overworld, zelduh.RoomsMap)

	ui := zelduh.NewUI(currLocaleMsgs)

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
		Player:    zelduh.BuildEntityFromConfig(zelduh.GetPreset("player")(6, 6), systemsManager.NewEntityID()),
		Bomb:      zelduh.BuildEntityFromConfig(zelduh.GetPreset("bomb")(0, 0), systemsManager.NewEntityID()),
		Explosion: zelduh.BuildEntityFromConfig(zelduh.GetPreset("explosion")(0, 0), systemsManager.NewEntityID()),
		Sword:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("sword")(0, 0), systemsManager.NewEntityID()),
		Arrow:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("arrow")(0, 0), systemsManager.NewEntityID()),

		RoomWarps:      map[zelduh.EntityID]zelduh.Config{},
		AllMapDrawData: zelduh.BuildMapDrawData(zelduh.TilemapDir, zelduh.TilemapFiles, zelduh.TileSize),

		InputSystem:  &zelduh.SystemInput{Win: ui.Window},
		HealthSystem: &zelduh.SystemHealth{},

		Hearts: zelduh.BuildEntitiesFromConfigs(
			systemsManager.NewEntityID,
			zelduh.GetPreset("heart")(1.5, 14),
			zelduh.GetPreset("heart")(2.15, 14),
			zelduh.GetPreset("heart")(2.80, 14),
		),
	}

	gameModel.SpatialSystem = &zelduh.SystemSpatial{
		Rand: gameModel.Rand,
	}

	collisionHandler := zelduh.CollisionHandler{
		GameModel:      &gameModel,
		SystemsManager: &systemsManager,
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

	systemsManager.AddSystems(
		gameModel.InputSystem,
		gameModel.HealthSystem,
		gameModel.SpatialSystem,
		collisionSystem,
		&zelduh.SystemRender{
			Win:         ui.Window,
			Spritesheet: gameModel.Spritesheet,
		},
	)

	systemsManager.AddEntities(
		gameModel.Player,
		gameModel.Sword,
		gameModel.Arrow,
		gameModel.Bomb,
	)

	gameStateManager := zelduh.NewGameStateManager(
		&gameModel,
		&systemsManager,
		ui,
		currLocaleMsgs,
		collisionSystem,
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
