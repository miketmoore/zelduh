package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

var (
	win       *pixelgl.Window
	txt       *text.Text
	gameWorld zelduh.World
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

	gameWorld = zelduh.New()

	zelduh.BuildMapRoomIDToRoom(zelduh.Overworld, roomsMap)

	txt = initText(20, 50, colornames.Black)
	win = initWindow(currLocaleMsgs["gameTitle"])

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

	collisionHandler := CollisionHandler{
		GameModel: &gameModel,
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

		allowQuit()

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

func initText(x, y float64, color color.RGBA) *text.Text {
	orig := pixel.V(x, y)
	txt := text.New(orig, text.Atlas7x13)
	txt.Color = color
	return txt
}

func initWindow(title string) *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(zelduh.WinX, zelduh.WinY, zelduh.WinW, zelduh.WinH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		fmt.Println("Initializing GUI window failed:")
		fmt.Println(err)
		os.Exit(1)
	}
	return win
}

func allowQuit() {
	if win.JustPressed(pixelgl.KeyQ) {
		os.Exit(1)
	}
}

func dropCoin(v pixel.Vec) {
	coin := zelduh.BuildEntityFromConfig(zelduh.GetPreset("coin")(v.X/zelduh.TileSize, v.Y/zelduh.TileSize), gameWorld.NewEntityID())
	gameWorld.AddEntity(coin)
}

// CollisionHandler contains collision handlers
type CollisionHandler struct {
	GameModel *zelduh.GameModel
}

// OnPlayerCollisionWithBounds handles collisions between player and bounds
func (ch *CollisionHandler) OnPlayerCollisionWithBounds(side zelduh.Bound) {
	if !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Side = side
		ch.GameModel.RoomTransition.Style = zelduh.TransitionSlide
		ch.GameModel.RoomTransition.Timer = int(ch.GameModel.RoomTransition.Start)
		ch.GameModel.CurrentState = zelduh.StateMapTransition
		ch.GameModel.AddEntities = true
	}
}

// OnPlayerCollisionWithCoin handles collision between player and coin
func (ch *CollisionHandler) OnPlayerCollisionWithCoin(coinID zelduh.EntityID) {
	ch.GameModel.Player.ComponentCoins.Coins++
	gameWorld.Remove(zelduh.CategoryCoin, coinID)
}

// OnPlayerCollisionWithEnemy handles collision between player and enemy
func (ch *CollisionHandler) OnPlayerCollisionWithEnemy(enemyID zelduh.EntityID) {
	// TODO repeat what I did with the enemies
	ch.GameModel.SpatialSystem.MovePlayerBack()
	ch.GameModel.Player.ComponentHealth.Total--

	// remove heart entity
	heartIndex := len(ch.GameModel.Hearts) - 1
	gameWorld.Remove(zelduh.CategoryHeart, ch.GameModel.Hearts[heartIndex].ID())
	ch.GameModel.Hearts = append(ch.GameModel.Hearts[:heartIndex], ch.GameModel.Hearts[heartIndex+1:]...)

	if ch.GameModel.Player.ComponentHealth.Total == 0 {
		ch.GameModel.CurrentState = zelduh.StateOver
	}
}

// OnSwordCollisionWithEnemy handles collision between sword and enemy
func (ch *CollisionHandler) OnSwordCollisionWithEnemy(enemyID zelduh.EntityID) {
	if !ch.GameModel.Sword.ComponentIgnore.Value {
		dead := false
		if !ch.GameModel.SpatialSystem.EnemyMovingFromHit(enemyID) {
			dead = ch.GameModel.HealthSystem.Hit(enemyID, 1)
			if dead {
				enemySpatial, _ := ch.GameModel.SpatialSystem.GetEnemySpatial(enemyID)
				ch.GameModel.Explosion.ComponentTemporary.Expiration = len(ch.GameModel.Explosion.ComponentAnimation.Map["default"].Frames)
				ch.GameModel.Explosion.ComponentSpatial = &zelduh.ComponentSpatial{
					Width:  zelduh.TileSize,
					Height: zelduh.TileSize,
					Rect:   enemySpatial.Rect,
				}
				ch.GameModel.Explosion.ComponentTemporary.OnExpiration = func() {
					dropCoin(ch.GameModel.Explosion.ComponentSpatial.Rect.Min)
				}
				gameWorld.AddEntity(ch.GameModel.Explosion)
				gameWorld.RemoveEnemy(enemyID)
			} else {
				ch.GameModel.SpatialSystem.MoveEnemyBack(enemyID, ch.GameModel.Player.ComponentMovement.Direction)
			}
		}

	}
}

// OnArrowCollisionWithEnemy handles collision between arrow and enemy
func (ch *CollisionHandler) OnArrowCollisionWithEnemy(enemyID zelduh.EntityID) {
	if !ch.GameModel.Arrow.ComponentIgnore.Value {
		dead := ch.GameModel.HealthSystem.Hit(enemyID, 1)
		ch.GameModel.Arrow.ComponentIgnore.Value = true
		if dead {
			enemySpatial, _ := ch.GameModel.SpatialSystem.GetEnemySpatial(enemyID)
			ch.GameModel.Explosion.ComponentTemporary.Expiration = len(ch.GameModel.Explosion.ComponentAnimation.Map["default"].Frames)
			ch.GameModel.Explosion.ComponentSpatial = &zelduh.ComponentSpatial{
				Width:  zelduh.TileSize,
				Height: zelduh.TileSize,
				Rect:   enemySpatial.Rect,
			}
			ch.GameModel.Explosion.ComponentTemporary.OnExpiration = func() {
				dropCoin(ch.GameModel.Explosion.ComponentSpatial.Rect.Min)
			}
			gameWorld.AddEntity(ch.GameModel.Explosion)
			gameWorld.RemoveEnemy(enemyID)
		} else {
			ch.GameModel.SpatialSystem.MoveEnemyBack(enemyID, ch.GameModel.Player.ComponentMovement.Direction)
		}
	}
}

// OnArrowCollisionWithObstacle handles collision between arrow and obstacle
func (ch *CollisionHandler) OnArrowCollisionWithObstacle() {
	ch.GameModel.Arrow.ComponentMovement.RemainingMoves = 0
}

// OnPlayerCollisionWithObstacle handles collision between player and obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithObstacle(obstacleID zelduh.EntityID) {
	// "Block" by undoing rect
	ch.GameModel.Player.ComponentSpatial.Rect = ch.GameModel.Player.ComponentSpatial.PrevRect
	ch.GameModel.Sword.ComponentSpatial.Rect = ch.GameModel.Sword.ComponentSpatial.PrevRect
}

// OnPlayerCollisionWithMoveableObstacle handles collision between player and moveable obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithMoveableObstacle(obstacleID zelduh.EntityID) {
	moved := ch.GameModel.SpatialSystem.MoveMoveableObstacle(obstacleID, ch.GameModel.Player.ComponentMovement.Direction)
	if !moved {
		ch.GameModel.Player.ComponentSpatial.Rect = ch.GameModel.Player.ComponentSpatial.PrevRect
	}
}

// OnMoveableObstacleCollisionWithSwitch handles collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleCollisionWithSwitch(collisionSwitchID zelduh.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnMoveableObstacleNoCollisionWithSwitch handles *no* collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleNoCollisionWithSwitch(collisionSwitchID zelduh.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnEnemyCollisionWithObstacle handles collision between enemy and obstacle
func (ch *CollisionHandler) OnEnemyCollisionWithObstacle(enemyID, obstacleID zelduh.EntityID) {
	// Block enemy within the spatial system by reseting current rect to previous rect
	ch.GameModel.SpatialSystem.UndoEnemyRect(enemyID)
}

// OnPlayerCollisionWithSwitch handles collision between player and switch
func (ch *CollisionHandler) OnPlayerCollisionWithSwitch(collisionSwitchID zelduh.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerNoCollisionWithSwitch handles *no* collision between player and switch
func (ch *CollisionHandler) OnPlayerNoCollisionWithSwitch(collisionSwitchID zelduh.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerCollisionWithWarp handles collision between player and warp
func (ch *CollisionHandler) OnPlayerCollisionWithWarp(warpID zelduh.EntityID) {
	entityConfig, ok := ch.GameModel.RoomWarps[warpID]
	if ok && !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Style = zelduh.TransitionWarp
		ch.GameModel.RoomTransition.Timer = 1
		ch.GameModel.CurrentState = zelduh.StateMapTransition
		ch.GameModel.AddEntities = true
		ch.GameModel.NextRoomID = entityConfig.WarpToRoomID
	}
}
