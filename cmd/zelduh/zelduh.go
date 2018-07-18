package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/entityconfig"
	"github.com/miketmoore/zelduh/rooms"
	"github.com/miketmoore/zelduh/world"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/zelduh/systems"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/image/colornames"
	"golang.org/x/text/language"
)

var (
	win       *pixelgl.Window
	txt       *text.Text
	gameWorld world.World
)

// GameModel contains data used throughout the game
type GameModel struct {
	AddEntities                           bool
	CurrentRoomID, NextRoomID             terraform2d.RoomID
	RoomTransition                        *terraform2d.RoomTransition
	CurrentState                          terraform2d.State
	Rand                                  *rand.Rand
	EntitiesMap                           map[terraform2d.EntityID]zelduh.Entity
	Spritesheet                           map[int]*pixel.Sprite
	Arrow, Bomb, Explosion, Player, Sword zelduh.Entity
	Hearts                                []zelduh.Entity
	RoomWarps                             map[terraform2d.EntityID]entityconfig.Config
	AllMapDrawData                        map[string]terraform2d.MapData
	HealthSystem                          *systems.Health
	InputSystem                           *systems.Input
	SpatialSystem                         *systems.Spatial
}

func run() {

	bundle := &i18n.Bundle{DefaultLanguage: language.English}

	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile(zelduh.TranslationFile)

	localizer := i18n.NewLocalizer(bundle, "en")

	i18nTitle := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Title",
		},
	})
	i18nPaused := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "Paused",
		},
	})
	i18nGameOver := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: "GameOver",
		},
	})

	gameWorld = world.New()

	rooms.ProcessMapLayout(zelduh.Overworld, roomsMap)

	txt = initText(20, 50, colornames.Black)
	win = initWindow(i18nTitle)

	gameModel := GameModel{
		Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		EntitiesMap:   map[terraform2d.EntityID]zelduh.Entity{},
		CurrentState:  terraform2d.StateStart,
		AddEntities:   true,
		CurrentRoomID: 1,
		RoomTransition: &terraform2d.RoomTransition{
			Start: float64(zelduh.TileSize),
		},
		Spritesheet: terraform2d.LoadAndBuildSpritesheet(zelduh.SpritesheetPath, zelduh.TileSize),

		// Build entities
		Player:    zelduh.BuildEntityFromConfig(zelduh.GetPreset("player")(6, 6), gameWorld.NewEntityID()),
		Bomb:      zelduh.BuildEntityFromConfig(zelduh.GetPreset("bomb")(0, 0), gameWorld.NewEntityID()),
		Explosion: zelduh.BuildEntityFromConfig(zelduh.GetPreset("explosion")(0, 0), gameWorld.NewEntityID()),
		Sword:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("sword")(0, 0), gameWorld.NewEntityID()),
		Arrow:     zelduh.BuildEntityFromConfig(zelduh.GetPreset("arrow")(0, 0), gameWorld.NewEntityID()),

		RoomWarps:      map[terraform2d.EntityID]entityconfig.Config{},
		AllMapDrawData: terraform2d.BuildMapDrawData(zelduh.TilemapDir, zelduh.TilemapFiles, zelduh.TileSize),

		InputSystem:  &systems.Input{Win: win},
		HealthSystem: &systems.Health{},

		Hearts: zelduh.BuildEntitiesFromConfigs(
			gameWorld.NewEntityID,
			zelduh.GetPreset("heart")(1.5, 14),
			zelduh.GetPreset("heart")(2.15, 14),
			zelduh.GetPreset("heart")(2.80, 14),
		),
	}

	gameModel.SpatialSystem = &systems.Spatial{
		Rand: gameModel.Rand,
	}

	collisionHandler := CollisionHandler{
		GameModel: &gameModel,
	}

	collisionSystem := &systems.Collision{
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
		&systems.Render{
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
		case terraform2d.StateStart:
			win.Clear(colornames.Darkgray)
			drawMapBG(zelduh.MapX, zelduh.MapY, zelduh.MapW, zelduh.MapH, colornames.White)
			drawCenterText(i18nTitle, colornames.Black)

			if win.JustPressed(pixelgl.KeyEnter) {
				gameModel.CurrentState = terraform2d.StateGame
			}
		case terraform2d.StateGame:
			gameModel.InputSystem.EnablePlayer()

			win.Clear(colornames.Darkgray)
			drawMapBG(zelduh.MapX, zelduh.MapY, zelduh.MapW, zelduh.MapH, colornames.White)

			drawMapBGImage(
				gameModel.Spritesheet,
				gameModel.AllMapDrawData,
				roomsMap[gameModel.CurrentRoomID].MapName(),
				0, 0)

			if gameModel.AddEntities {
				gameModel.AddEntities = false
				addUIHearts(gameModel.Hearts, gameModel.Player.ComponentHealth.Total)

				addUICoin()

				// Draw obstacles on appropriate map tiles
				obstacles := drawObstaclesPerMapTiles(gameModel.AllMapDrawData, gameModel.CurrentRoomID, 0, 0)
				gameWorld.AddEntities(obstacles...)

				gameModel.RoomWarps = map[terraform2d.EntityID]entityconfig.Config{}

				// Iterate through all entity configurations and build entities and add to systems
				for _, c := range roomsMap[gameModel.CurrentRoomID].(*rooms.Room).EntityConfigs {
					entity := zelduh.BuildEntityFromConfig(c, gameWorld.NewEntityID())
					gameModel.EntitiesMap[entity.ID()] = entity
					gameWorld.AddEntity(entity)

					switch c.Category {
					case zelduh.CategoryWarp:
						gameModel.RoomWarps[entity.ID()] = c
					}
				}
			}

			drawMask()

			gameWorld.Update()

			if win.JustPressed(pixelgl.KeyP) {
				gameModel.CurrentState = terraform2d.StatePause
			}

			if win.JustPressed(pixelgl.KeyX) {
				gameModel.CurrentState = terraform2d.StateOver
			}

		case terraform2d.StatePause:
			win.Clear(colornames.Darkgray)
			drawMapBG(zelduh.MapX, zelduh.MapY, zelduh.MapW, zelduh.MapH, colornames.White)
			drawCenterText(i18nPaused, colornames.Black)

			if win.JustPressed(pixelgl.KeyP) {
				gameModel.CurrentState = terraform2d.StateGame
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				gameModel.CurrentState = terraform2d.StateStart
			}
		case terraform2d.StateOver:
			win.Clear(colornames.Darkgray)
			drawMapBG(zelduh.MapX, zelduh.MapY, zelduh.MapW, zelduh.MapH, colornames.Black)
			drawCenterText(i18nGameOver, colornames.White)

			if win.JustPressed(pixelgl.KeyEnter) {
				gameModel.CurrentState = terraform2d.StateStart
			}
		case terraform2d.StateMapTransition:
			gameModel.InputSystem.DisablePlayer()
			if gameModel.RoomTransition.Style == terraform2d.TransitionSlide && gameModel.RoomTransition.Timer > 0 {
				gameModel.RoomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(zelduh.MapX, zelduh.MapY, zelduh.MapW, zelduh.MapH, colornames.White)

				collisionSystem.RemoveAll(zelduh.CategoryObstacle)
				gameWorld.RemoveAllEnemies()
				gameWorld.RemoveAllCollisionSwitches()
				gameWorld.RemoveAllMoveableObstacles()
				gameWorld.RemoveAllEntities()

				currentRoomID := gameModel.CurrentRoomID

				connectedRooms := roomsMap[currentRoomID].ConnectedRooms()

				transitionRoomResp := calculateTransitionSlide(
					gameModel.RoomTransition,
					*connectedRooms,
					gameModel.CurrentRoomID,
				)

				gameModel.NextRoomID = transitionRoomResp.nextRoomID

				drawMapBGImage(
					gameModel.Spritesheet,
					gameModel.AllMapDrawData,
					roomsMap[gameModel.CurrentRoomID].MapName(),
					transitionRoomResp.modX,
					transitionRoomResp.modY,
				)
				drawMapBGImage(
					gameModel.Spritesheet,
					gameModel.AllMapDrawData,
					roomsMap[gameModel.NextRoomID].MapName(),
					transitionRoomResp.modXNext,
					transitionRoomResp.modYNext,
				)
				drawMask()

				// Move player with map transition
				gameModel.Player.ComponentSpatial.Rect = pixel.R(
					gameModel.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX,
					gameModel.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY,
					gameModel.Player.ComponentSpatial.Rect.Min.X+transitionRoomResp.playerModX+zelduh.TileSize,
					gameModel.Player.ComponentSpatial.Rect.Min.Y+transitionRoomResp.playerModY+zelduh.TileSize,
				)

				gameWorld.Update()
			} else if gameModel.RoomTransition.Style == terraform2d.TransitionWarp && gameModel.RoomTransition.Timer > 0 {
				gameModel.RoomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(zelduh.MapX, zelduh.MapY, zelduh.MapW, zelduh.MapH, colornames.White)

				collisionSystem.RemoveAll(zelduh.CategoryObstacle)
				gameWorld.RemoveAllEnemies()
				gameWorld.RemoveAllCollisionSwitches()
				gameWorld.RemoveAllMoveableObstacles()
				gameWorld.RemoveAllEntities()
			} else {
				gameModel.CurrentState = terraform2d.StateGame
				if gameModel.NextRoomID != 0 {
					gameModel.CurrentRoomID = gameModel.NextRoomID
				}
				gameModel.RoomTransition.Active = false
			}

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

func drawCenterText(s string, c color.RGBA) {
	txt.Clear()
	txt.Color = c
	fmt.Fprintln(txt, s)
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
}

func drawMapBG(x, y, w, h float64, color color.Color) {
	s := imdraw.New(nil)
	s.Color = color
	s.Push(pixel.V(x, y))
	s.Push(pixel.V(x+w, y+h))
	s.Rectangle(0)
	s.Draw(win)
}

func drawMapBGImage(
	spritesheet map[int]*pixel.Sprite,
	allMapDrawData map[string]terraform2d.MapData,
	name string,
	modX, modY float64) {

	d := allMapDrawData[name]
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			sprite := spritesheet[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+zelduh.MapX+modX+zelduh.TileSize/2,
				vec.Y+zelduh.MapY+modY+zelduh.TileSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func drawObstaclesPerMapTiles(allMapDrawData map[string]terraform2d.MapData, roomID terraform2d.RoomID, modX, modY float64) []zelduh.Entity {
	d := allMapDrawData[roomsMap[roomID].MapName()]
	obstacles := []zelduh.Entity{}
	mod := 0.5
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+zelduh.MapX+modX+zelduh.TileSize/2,
				vec.Y+zelduh.MapY+modY+zelduh.TileSize/2,
			)

			if _, ok := zelduh.NonObstacleSprites[spriteData.SpriteID]; !ok {
				x := movedVec.X/zelduh.TileSize - mod
				y := movedVec.Y/zelduh.TileSize - mod
				id := gameWorld.NewEntityID()
				obstacle := zelduh.BuildEntityFromConfig(zelduh.GetPreset("obstacle")(x, y), id)
				obstacles = append(obstacles, obstacle)
			}
		}
	}
	return obstacles
}

func drawMask() {
	// top
	s := imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, zelduh.MapY+zelduh.MapH))
	s.Push(pixel.V(zelduh.WinW, zelduh.MapY+zelduh.MapH+(zelduh.WinH-(zelduh.MapY+zelduh.MapH))))
	s.Rectangle(0)
	s.Draw(win)

	// bottom
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(zelduh.WinW, (zelduh.WinH - (zelduh.MapY + zelduh.MapH))))
	s.Rectangle(0)
	s.Draw(win)

	// left
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(0+zelduh.MapX, zelduh.WinH))
	s.Rectangle(0)
	s.Draw(win)

	// right
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(zelduh.MapX+zelduh.MapW, zelduh.MapY))
	s.Push(pixel.V(zelduh.WinW, zelduh.WinH))
	s.Rectangle(0)
	s.Draw(win)
}

var roomsMap = rooms.Rooms{
	1: rooms.NewRoom("overworldFourWallsDoorBottomRight",
		zelduh.GetPreset("puzzleBox")(5, 5),
		zelduh.GetPreset("floorSwitch")(5, 6),
		zelduh.GetPreset("toggleObstacle")(10, 7),
	),
	2: rooms.NewRoom("overworldFourWallsDoorTopBottom",
		zelduh.GetPreset("skull")(5, 5),
		zelduh.GetPreset("skeleton")(11, 9),
		zelduh.GetPreset("spinner")(7, 9),
		zelduh.GetPreset("eyeburrower")(8, 9),
	),
	3: rooms.NewRoom("overworldFourWallsDoorRightTopBottom",
		zelduh.WarpStone(3, 7, 6, 5),
	),
	5: rooms.NewRoom("rockWithCaveEntrance",
		entityconfig.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 11,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 7) + zelduh.TileSize/2,
			Y:            (zelduh.TileSize * 9) + zelduh.TileSize/2,
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 30,
			},
		},
		entityconfig.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 11,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 8) + zelduh.TileSize/2,
			Y:            (zelduh.TileSize * 9) + zelduh.TileSize/2,
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 30,
			},
		},
	),
	6:  rooms.NewRoom("rockPathLeftRightEntrance"),
	7:  rooms.NewRoom("overworldFourWallsDoorLeftTop"),
	8:  rooms.NewRoom("overworldFourWallsDoorBottom"),
	9:  rooms.NewRoom("overworldFourWallsDoorTop"),
	10: rooms.NewRoom("overworldFourWallsDoorLeft"),
	11: rooms.NewRoom("dungeonFourDoors",
		// South door of cave - warp to cave entrance
		entityconfig.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 5,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 6) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Y:            (zelduh.TileSize * 1) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 15,
			},
		},
		entityconfig.Config{
			Category:     zelduh.CategoryWarp,
			WarpToRoomID: 5,
			W:            zelduh.TileSize,
			H:            zelduh.TileSize,
			X:            (zelduh.TileSize * 7) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Y:            (zelduh.TileSize * 1) + zelduh.TileSize + (zelduh.TileSize / 2.5),
			Hitbox: &entityconfig.HitboxConfig{
				Radius: 15,
			},
		},
	),
}

func addUICoin() {
	coin := zelduh.BuildEntityFromConfig(zelduh.GetPreset("uiCoin")(4, 14), gameWorld.NewEntityID())
	gameWorld.AddEntity(coin)
}

// make sure only correct number of hearts exists in systems
// so, if health is reduced, need to remove a heart entity from the systems,
// the correct one... last one
func addUIHearts(hearts []zelduh.Entity, health int) {
	for i, entity := range hearts {
		if i < health {
			gameWorld.AddEntity(entity)
		}
	}
}

func dropCoin(v pixel.Vec) {
	coin := zelduh.BuildEntityFromConfig(zelduh.GetPreset("coin")(v.X/zelduh.TileSize, v.Y/zelduh.TileSize), gameWorld.NewEntityID())
	gameWorld.AddEntity(coin)
}

// CollisionHandler contains collision handlers
type CollisionHandler struct {
	GameModel *GameModel
}

// OnPlayerCollisionWithBounds handles collisions between player and bounds
func (ch *CollisionHandler) OnPlayerCollisionWithBounds(side terraform2d.Bound) {
	if !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Side = side
		ch.GameModel.RoomTransition.Style = terraform2d.TransitionSlide
		ch.GameModel.RoomTransition.Timer = int(ch.GameModel.RoomTransition.Start)
		ch.GameModel.CurrentState = terraform2d.StateMapTransition
		ch.GameModel.AddEntities = true
	}
}

// OnPlayerCollisionWithCoin handles collision between player and coin
func (ch *CollisionHandler) OnPlayerCollisionWithCoin(coinID terraform2d.EntityID) {
	ch.GameModel.Player.ComponentCoins.Coins++
	gameWorld.Remove(zelduh.CategoryCoin, coinID)
}

// OnPlayerCollisionWithEnemy handles collision between player and enemy
func (ch *CollisionHandler) OnPlayerCollisionWithEnemy(enemyID terraform2d.EntityID) {
	// TODO repeat what I did with the enemies
	ch.GameModel.SpatialSystem.MovePlayerBack()
	ch.GameModel.Player.ComponentHealth.Total--

	// remove heart entity
	heartIndex := len(ch.GameModel.Hearts) - 1
	gameWorld.Remove(zelduh.CategoryHeart, ch.GameModel.Hearts[heartIndex].ID())
	ch.GameModel.Hearts = append(ch.GameModel.Hearts[:heartIndex], ch.GameModel.Hearts[heartIndex+1:]...)

	if ch.GameModel.Player.ComponentHealth.Total == 0 {
		ch.GameModel.CurrentState = terraform2d.StateOver
	}
}

// OnSwordCollisionWithEnemy handles collision between sword and enemy
func (ch *CollisionHandler) OnSwordCollisionWithEnemy(enemyID terraform2d.EntityID) {
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
func (ch *CollisionHandler) OnArrowCollisionWithEnemy(enemyID terraform2d.EntityID) {
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
func (ch *CollisionHandler) OnPlayerCollisionWithObstacle(obstacleID terraform2d.EntityID) {
	// "Block" by undoing rect
	ch.GameModel.Player.ComponentSpatial.Rect = ch.GameModel.Player.ComponentSpatial.PrevRect
	ch.GameModel.Sword.ComponentSpatial.Rect = ch.GameModel.Sword.ComponentSpatial.PrevRect
}

// OnPlayerCollisionWithMoveableObstacle handles collision between player and moveable obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithMoveableObstacle(obstacleID terraform2d.EntityID) {
	moved := ch.GameModel.SpatialSystem.MoveMoveableObstacle(obstacleID, ch.GameModel.Player.ComponentMovement.Direction)
	if !moved {
		ch.GameModel.Player.ComponentSpatial.Rect = ch.GameModel.Player.ComponentSpatial.PrevRect
	}
}

// OnMoveableObstacleCollisionWithSwitch handles collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleCollisionWithSwitch(collisionSwitchID terraform2d.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnMoveableObstacleNoCollisionWithSwitch handles *no* collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleNoCollisionWithSwitch(collisionSwitchID terraform2d.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnEnemyCollisionWithObstacle handles collision between enemy and obstacle
func (ch *CollisionHandler) OnEnemyCollisionWithObstacle(enemyID, obstacleID terraform2d.EntityID) {
	// Block enemy within the spatial system by reseting current rect to previous rect
	ch.GameModel.SpatialSystem.UndoEnemyRect(enemyID)
}

// OnPlayerCollisionWithSwitch handles collision between player and switch
func (ch *CollisionHandler) OnPlayerCollisionWithSwitch(collisionSwitchID terraform2d.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerNoCollisionWithSwitch handles *no* collision between player and switch
func (ch *CollisionHandler) OnPlayerNoCollisionWithSwitch(collisionSwitchID terraform2d.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerCollisionWithWarp handles collision between player and warp
func (ch *CollisionHandler) OnPlayerCollisionWithWarp(warpID terraform2d.EntityID) {
	entityConfig, ok := ch.GameModel.RoomWarps[warpID]
	if ok && !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Style = terraform2d.TransitionWarp
		ch.GameModel.RoomTransition.Timer = 1
		ch.GameModel.CurrentState = terraform2d.StateMapTransition
		ch.GameModel.AddEntities = true
		ch.GameModel.NextRoomID = entityConfig.WarpToRoomID
	}
}

// TransitionRoomResponse contains layout data
type TransitionRoomResponse struct {
	nextRoomID                                             terraform2d.RoomID
	modX, modY, modXNext, modYNext, playerModX, playerModY float64
}

func calculateTransitionSlide(
	roomTransition *terraform2d.RoomTransition,
	connectedRooms terraform2d.ConnectedRooms,
	currentRoomID terraform2d.RoomID) TransitionRoomResponse {

	var nextRoomID terraform2d.RoomID
	inc := (roomTransition.Start - float64(roomTransition.Timer))
	incY := inc * (zelduh.MapH / zelduh.TileSize)
	incX := inc * (zelduh.MapW / zelduh.TileSize)
	modY := 0.0
	modYNext := 0.0
	modX := 0.0
	modXNext := 0.0
	playerModX := 0.0
	playerModY := 0.0
	playerIncY := ((zelduh.MapH / zelduh.TileSize) - 1) + 7
	playerIncX := ((zelduh.MapW / zelduh.TileSize) - 1) + 7
	if roomTransition.Side == terraform2d.BoundBottom && connectedRooms.Bottom != 0 {
		modY = incY
		modYNext = incY - zelduh.MapH
		nextRoomID = connectedRooms.Bottom
		playerModY += playerIncY
	} else if roomTransition.Side == terraform2d.BoundTop && connectedRooms.Top != 0 {
		modY = -incY
		modYNext = -incY + zelduh.MapH
		nextRoomID = connectedRooms.Top
		playerModY -= playerIncY
	} else if roomTransition.Side == terraform2d.BoundLeft && connectedRooms.Left != 0 {
		modX = incX
		modXNext = incX - zelduh.MapW
		nextRoomID = connectedRooms.Left
		playerModX += playerIncX
	} else if roomTransition.Side == terraform2d.BoundRight && connectedRooms.Right != 0 {
		modX = -incX
		modXNext = -incX + zelduh.MapW
		nextRoomID = connectedRooms.Right
		playerModX -= playerIncX
	} else {
		nextRoomID = 0
	}

	return TransitionRoomResponse{
		nextRoomID,
		modX, modY, modXNext, modYNext, playerModX, playerModY,
	}
}
