package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh/bounds"
	"github.com/miketmoore/zelduh/config"
	"github.com/miketmoore/zelduh/gamemap"
	"github.com/miketmoore/zelduh/rooms"
	"github.com/miketmoore/zelduh/sprites"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/entities"
	"github.com/miketmoore/zelduh/gamestate"
	"github.com/miketmoore/zelduh/systems"
	"github.com/miketmoore/zelduh/tmx"
	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/image/colornames"
)

var (
	win       *pixelgl.Window
	txt       *text.Text
	t         i18n.TranslateFunc
	gameWorld terraform2d.World
)

// GameModel contains data used throughout the game
type GameModel struct {
	AddEntities                           bool
	CurrentRoomID, NextRoomID             rooms.RoomID
	RoomTransition                        rooms.RoomTransition
	CurrentState                          gamestate.Name
	Rand                                  *rand.Rand
	EntitiesMap                           map[terraform2d.EntityID]entities.Entity
	Spritesheet                           map[int]*pixel.Sprite
	Arrow, Bomb, Explosion, Player, Sword entities.Entity
	Hearts                                []entities.Entity
	RoomWarps                             map[terraform2d.EntityID]rooms.EntityConfig
	AllMapDrawData                        map[string]tmx.MapData
	HealthSystem                          *systems.Health
	InputSystem                           *systems.Input
	SpatialSystem                         *systems.Spatial
}

// EntityRemover is a concrete terraform2d.EntityRemover
type EntityRemover struct{}

// Remove removes an entity by category and id
func (r EntityRemover) Remove(w *terraform2d.World, category terraform2d.EntityCategory, id terraform2d.EntityID) {

}

// RemoveAllEntities removes all entities
func (r EntityRemover) RemoveAllEntities(w *terraform2d.World) {

}

func run() {

	entityRemover := EntityRemover{}

	gameWorld = terraform2d.NewWorld(entityRemover)

	gamemap.ProcessMapLayout(roomsMap)

	// Initializations
	t = initI18n()
	txt = initText(20, 50, colornames.Black)
	win = initWindow(t("title"))

	gameModel := GameModel{
		Rand:          rand.New(rand.NewSource(time.Now().UnixNano())),
		EntitiesMap:   map[terraform2d.EntityID]entities.Entity{},
		CurrentState:  gamestate.Start,
		AddEntities:   true,
		CurrentRoomID: 1,
		RoomTransition: rooms.RoomTransition{
			Start: float64(config.TileSize),
		},
		Spritesheet: sprites.LoadAndBuildSpritesheet(),

		// Build entities
		Player:    entities.BuildEntityFromConfig(entities.GetPreset("player")(6, 6), gameWorld.NewEntityID()),
		Bomb:      entities.BuildEntityFromConfig(entities.GetPreset("bomb")(0, 0), gameWorld.NewEntityID()),
		Explosion: entities.BuildEntityFromConfig(entities.GetPreset("explosion")(0, 0), gameWorld.NewEntityID()),
		Sword:     entities.BuildEntityFromConfig(entities.GetPreset("sword")(0, 0), gameWorld.NewEntityID()),
		Arrow:     entities.BuildEntityFromConfig(entities.GetPreset("arrow")(0, 0), gameWorld.NewEntityID()),

		RoomWarps:      map[terraform2d.EntityID]rooms.EntityConfig{},
		AllMapDrawData: tmx.BuildMapDrawData(),

		InputSystem:  &systems.Input{Win: win},
		HealthSystem: &systems.Health{},

		Hearts: entities.BuildEntitiesFromConfigs(
			gameWorld.NewEntityID,
			entities.GetPreset("heart")(1.5, 14),
			entities.GetPreset("heart")(2.15, 14),
			entities.GetPreset("heart")(2.80, 14),
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
			config.MapX,
			config.MapY,
			config.MapX+config.MapW,
			config.MapY+config.MapH,
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

	gameWorld.AddEntitiesToSystem(
		gameModel.Player,
		gameModel.Sword,
		gameModel.Arrow,
		gameModel.Bomb,
	)

	for !win.Closed() {

		allowQuit()

		switch gameModel.CurrentState {
		case gamestate.Start:
			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)
			drawCenterText(t("title"), colornames.Black)

			if win.JustPressed(pixelgl.KeyEnter) {
				gameModel.CurrentState = gamestate.Game
			}
		case gamestate.Game:
			gameModel.InputSystem.EnablePlayer()

			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

			drawMapBGImage(gameModel.Spritesheet, gameModel.AllMapDrawData, roomsMap[gameModel.CurrentRoomID].MapName, 0, 0)

			addHearts(gameModel.Hearts, gameModel.Player.Health.Total)

			if gameModel.AddEntities {
				gameModel.AddEntities = false

				addUICoin()

				// Draw obstacles on appropriate map tiles
				obstacles := drawObstaclesPerMapTiles(gameModel.AllMapDrawData, gameModel.CurrentRoomID, 0, 0)
				gameWorld.AddEntitiesToSystem(obstacles...)

				gameModel.RoomWarps = map[terraform2d.EntityID]rooms.EntityConfig{}

				// Iterate through all entity configurations and build entities and add to systems
				for _, c := range roomsMap[gameModel.CurrentRoomID].EntityConfigs {
					entity := entities.BuildEntityFromConfig(c, gameWorld.NewEntityID())
					gameModel.EntitiesMap[entity.ID] = entity
					gameWorld.AddEntityToSystem(entity)

					switch c.Category {
					case categories.Warp:
						gameModel.RoomWarps[entity.ID] = c
					}
				}
			}

			drawMask()

			gameWorld.Update()

			if win.JustPressed(pixelgl.KeyP) {
				gameModel.CurrentState = gamestate.Pause
			}

			if win.JustPressed(pixelgl.KeyX) {
				gameModel.CurrentState = gamestate.Over
			}

		case gamestate.Pause:
			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)
			drawCenterText(t("paused"), colornames.Black)

			if win.JustPressed(pixelgl.KeyP) {
				gameModel.CurrentState = gamestate.Game
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				gameModel.CurrentState = gamestate.Start
			}
		case gamestate.Over:
			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.Black)
			drawCenterText(t("gameOver"), colornames.White)

			if win.JustPressed(pixelgl.KeyEnter) {
				gameModel.CurrentState = gamestate.Start
			}
		case gamestate.MapTransition:
			gameModel.InputSystem.DisablePlayer()
			if gameModel.RoomTransition.Style == rooms.TransitionSlide && gameModel.RoomTransition.Timer > 0 {
				gameModel.RoomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				gameWorld.RemoveAllEnemies()
				gameWorld.RemoveAllCollisionSwitches()
				gameWorld.RemoveAllMoveableObstacles()
				gameWorld.RemoveAllEntities()

				transitionRoomResp := calculateTransitionSlide(
					&gameModel.RoomTransition,
					roomsMap[gameModel.CurrentRoomID].ConnectedRooms,
					gameModel.CurrentRoomID,
				)

				gameModel.NextRoomID = transitionRoomResp.nextRoomID

				drawMapBGImage(
					gameModel.Spritesheet,
					gameModel.AllMapDrawData,
					roomsMap[gameModel.CurrentRoomID].MapName,
					transitionRoomResp.modX,
					transitionRoomResp.modY,
				)
				drawMapBGImage(
					gameModel.Spritesheet,
					gameModel.AllMapDrawData,
					roomsMap[gameModel.NextRoomID].MapName,
					transitionRoomResp.modXNext,
					transitionRoomResp.modYNext,
				)
				drawMask()

				// Move player with map transition
				gameModel.Player.Spatial.Rect = pixel.R(
					gameModel.Player.Spatial.Rect.Min.X+transitionRoomResp.playerModX,
					gameModel.Player.Spatial.Rect.Min.Y+transitionRoomResp.playerModY,
					gameModel.Player.Spatial.Rect.Min.X+transitionRoomResp.playerModX+config.TileSize,
					gameModel.Player.Spatial.Rect.Min.Y+transitionRoomResp.playerModY+config.TileSize,
				)

				gameWorld.Update()
			} else if gameModel.RoomTransition.Style == rooms.TransitionWarp && gameModel.RoomTransition.Timer > 0 {
				gameModel.RoomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				gameWorld.RemoveAllEnemies()
				gameWorld.RemoveAllCollisionSwitches()
				gameWorld.RemoveAllMoveableObstacles()
				gameWorld.RemoveAllEntities()
			} else {
				gameModel.CurrentState = gamestate.Game
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

func initI18n() i18n.TranslateFunc {
	i18n.LoadTranslationFile(config.TranslationFile)
	T, err := i18n.Tfunc(config.Lang)
	if err != nil {
		fmt.Println("Initializing i18n failed:")
		fmt.Println(err)
		os.Exit(1)
	}
	return T
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
		Bounds: pixel.R(config.WinX, config.WinY, config.WinW, config.WinH),
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
	allMapDrawData map[string]tmx.MapData,
	name string,
	modX, modY float64) {

	d := allMapDrawData[name]
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			sprite := spritesheet[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+config.MapX+modX+config.TileSize/2,
				vec.Y+config.MapY+modY+config.TileSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func drawObstaclesPerMapTiles(allMapDrawData map[string]tmx.MapData, roomID rooms.RoomID, modX, modY float64) []entities.Entity {
	d := allMapDrawData[roomsMap[roomID].MapName]
	obstacles := []entities.Entity{}
	mod := 0.5
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+config.MapX+modX+config.TileSize/2,
				vec.Y+config.MapY+modY+config.TileSize/2,
			)

			if _, ok := config.NonObstacleSprites[spriteData.SpriteID]; !ok {
				x := movedVec.X/config.TileSize - mod
				y := movedVec.Y/config.TileSize - mod
				id := gameWorld.NewEntityID()
				obstacle := entities.BuildEntityFromConfig(entities.GetPreset("obstacle")(x, y), id)
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
	s.Push(pixel.V(0, config.MapY+config.MapH))
	s.Push(pixel.V(config.WinW, config.MapY+config.MapH+(config.WinH-(config.MapY+config.MapH))))
	s.Rectangle(0)
	s.Draw(win)

	// bottom
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(config.WinW, (config.WinH - (config.MapY + config.MapH))))
	s.Rectangle(0)
	s.Draw(win)

	// left
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(0+config.MapX, config.WinH))
	s.Rectangle(0)
	s.Draw(win)

	// right
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(config.MapX+config.MapW, config.MapY))
	s.Push(pixel.V(config.WinW, config.WinH))
	s.Rectangle(0)
	s.Draw(win)
}

var roomsMap = rooms.Rooms{
	1: rooms.NewRoom("overworldFourWallsDoorBottomRight",
		entities.GetPreset("puzzleBox")(5, 5),
		entities.GetPreset("floorSwitch")(5, 6),
		entities.GetPreset("toggleObstacle")(10, 7),
	),
	2: rooms.NewRoom("overworldFourWallsDoorTopBottom",
		entities.GetPreset("skull")(5, 5),
		entities.GetPreset("skeleton")(11, 9),
		entities.GetPreset("spinner")(7, 9),
		entities.GetPreset("eyeburrower")(8, 9),
	),
	3: rooms.NewRoom("overworldFourWallsDoorRightTopBottom",
		entities.WarpStone(3, 7, 6, 5),
	),
	5: rooms.NewRoom("rockWithCaveEntrance",
		rooms.EntityConfig{
			Category:     categories.Warp,
			WarpToRoomID: 11,
			W:            config.TileSize,
			H:            config.TileSize,
			X:            (config.TileSize * 7) + config.TileSize/2,
			Y:            (config.TileSize * 9) + config.TileSize/2,
			Hitbox: &rooms.HitboxConfig{
				Radius: 30,
			},
		},
		rooms.EntityConfig{
			Category:     categories.Warp,
			WarpToRoomID: 11,
			W:            config.TileSize,
			H:            config.TileSize,
			X:            (config.TileSize * 8) + config.TileSize/2,
			Y:            (config.TileSize * 9) + config.TileSize/2,
			Hitbox: &rooms.HitboxConfig{
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
		rooms.EntityConfig{
			Category:     categories.Warp,
			WarpToRoomID: 5,
			W:            config.TileSize,
			H:            config.TileSize,
			X:            (config.TileSize * 6) + config.TileSize + (config.TileSize / 2.5),
			Y:            (config.TileSize * 1) + config.TileSize + (config.TileSize / 2.5),
			Hitbox: &rooms.HitboxConfig{
				Radius: 15,
			},
		},
		rooms.EntityConfig{
			Category:     categories.Warp,
			WarpToRoomID: 5,
			W:            config.TileSize,
			H:            config.TileSize,
			X:            (config.TileSize * 7) + config.TileSize + (config.TileSize / 2.5),
			Y:            (config.TileSize * 1) + config.TileSize + (config.TileSize / 2.5),
			Hitbox: &rooms.HitboxConfig{
				Radius: 15,
			},
		},
	),
}

func addUICoin() {
	coin := entities.BuildEntityFromConfig(entities.GetPreset("uiCoin")(4, 14), gameWorld.NewEntityID())
	gameWorld.AddEntityToSystem(coin)
}

// make sure only correct number of hearts exists in systems
// so, if health is reduced, need to remove a heart entity from the systems,
// the correct one... last one
func addHearts(hearts []entities.Entity, health int) {
	for i, entity := range hearts {
		if i < health {
			gameWorld.AddEntityToSystem(entity)
		}
	}
}

func dropCoin(v pixel.Vec) {
	coin := entities.BuildEntityFromConfig(entities.GetPreset("coin")(v.X/config.TileSize, v.Y/config.TileSize), gameWorld.NewEntityID())
	gameWorld.AddEntityToSystem(coin)
}

// CollisionHandler contains collision handlers
type CollisionHandler struct {
	GameModel *GameModel
}

// OnPlayerCollisionWithBounds handles collisions between player and bounds
func (ch *CollisionHandler) OnPlayerCollisionWithBounds(side bounds.Bound) {
	if !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Side = side
		ch.GameModel.RoomTransition.Style = rooms.TransitionSlide
		ch.GameModel.RoomTransition.Timer = int(ch.GameModel.RoomTransition.Start)
		ch.GameModel.CurrentState = gamestate.MapTransition
		ch.GameModel.AddEntities = true
	}
}

// OnPlayerCollisionWithCoin handles collision between player and coin
func (ch *CollisionHandler) OnPlayerCollisionWithCoin(coinID terraform2d.EntityID) {
	ch.GameModel.Player.Coins.Coins++
	gameWorld.Remove(categories.Coin, coinID)
}

// OnPlayerCollisionWithEnemy handles collision between player and enemy
func (ch *CollisionHandler) OnPlayerCollisionWithEnemy(enemyID terraform2d.EntityID) {
	// TODO repeat what I did with the enemies
	ch.GameModel.SpatialSystem.MovePlayerBack()
	ch.GameModel.Player.Health.Total--

	// remove heart entity
	heartIndex := len(ch.GameModel.Hearts) - 1
	gameWorld.Remove(categories.Heart, ch.GameModel.Hearts[heartIndex].ID)
	ch.GameModel.Hearts = append(ch.GameModel.Hearts[:heartIndex], ch.GameModel.Hearts[heartIndex+1:]...)

	// TODO redraw hearts
	if ch.GameModel.Player.Health.Total == 0 {
		ch.GameModel.CurrentState = gamestate.Over
	}
}

// OnSwordCollisionWithEnemy handles collision between sword and enemy
func (ch *CollisionHandler) OnSwordCollisionWithEnemy(enemyID terraform2d.EntityID) {
	fmt.Printf("SwordCollisionWithEnemy %d\n", enemyID)
	if !ch.GameModel.Sword.Ignore.Value {
		dead := false
		if !ch.GameModel.SpatialSystem.EnemyMovingFromHit(enemyID) {
			dead = ch.GameModel.HealthSystem.Hit(enemyID, 1)
			if dead {
				enemySpatial, _ := ch.GameModel.SpatialSystem.GetEnemySpatial(enemyID)
				ch.GameModel.Explosion.Temporary.Expiration = len(ch.GameModel.Explosion.Animation.Map["default"].Frames)
				ch.GameModel.Explosion.Spatial = &components.Spatial{
					Width:  config.TileSize,
					Height: config.TileSize,
					Rect:   enemySpatial.Rect,
				}
				ch.GameModel.Explosion.Temporary.OnExpiration = func() {
					dropCoin(ch.GameModel.Explosion.Spatial.Rect.Min)
				}
				gameWorld.AddEntityToSystem(ch.GameModel.Explosion)
				gameWorld.RemoveEnemy(enemyID)
			} else {
				ch.GameModel.SpatialSystem.MoveEnemyBack(enemyID, ch.GameModel.Player.Movement.Direction)
			}
		}

	}
}

// OnArrowCollisionWithEnemy handles collision between arrow and enemy
func (ch *CollisionHandler) OnArrowCollisionWithEnemy(enemyID terraform2d.EntityID) {
	if !ch.GameModel.Arrow.Ignore.Value {
		dead := ch.GameModel.HealthSystem.Hit(enemyID, 1)
		ch.GameModel.Arrow.Ignore.Value = true
		if dead {
			fmt.Printf("You killed an enemy with an arrow\n")
			enemySpatial, _ := ch.GameModel.SpatialSystem.GetEnemySpatial(enemyID)
			ch.GameModel.Explosion.Temporary.Expiration = len(ch.GameModel.Explosion.Animation.Map["default"].Frames)
			ch.GameModel.Explosion.Spatial = &components.Spatial{
				Width:  config.TileSize,
				Height: config.TileSize,
				Rect:   enemySpatial.Rect,
			}
			ch.GameModel.Explosion.Temporary.OnExpiration = func() {
				dropCoin(ch.GameModel.Explosion.Spatial.Rect.Min)
			}
			gameWorld.AddEntityToSystem(ch.GameModel.Explosion)
			gameWorld.RemoveEnemy(enemyID)
		} else {
			ch.GameModel.SpatialSystem.MoveEnemyBack(enemyID, ch.GameModel.Player.Movement.Direction)
		}
	}
}

// OnArrowCollisionWithObstacle handles collision between arrow and obstacle
func (ch *CollisionHandler) OnArrowCollisionWithObstacle() {
	ch.GameModel.Arrow.Movement.RemainingMoves = 0
}

// OnPlayerCollisionWithObstacle handles collision between player and obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithObstacle(obstacleID terraform2d.EntityID) {
	// "Block" by undoing rect
	ch.GameModel.Player.Spatial.Rect = ch.GameModel.Player.Spatial.PrevRect
	ch.GameModel.Sword.Spatial.Rect = ch.GameModel.Sword.Spatial.PrevRect
}

// OnPlayerCollisionWithMoveableObstacle handles collision between player and moveable obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithMoveableObstacle(obstacleID terraform2d.EntityID) {
	moved := ch.GameModel.SpatialSystem.MoveMoveableObstacle(obstacleID, ch.GameModel.Player.Movement.Direction)
	if !moved {
		ch.GameModel.Player.Spatial.Rect = ch.GameModel.Player.Spatial.PrevRect
	}
}

// OnMoveableObstacleCollisionWithSwitch handles collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleCollisionWithSwitch(collisionSwitchID terraform2d.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && !entity.Toggler.Enabled() {
			entity.Toggler.Toggle()
		}
	}
}

// OnMoveableObstacleNoCollisionWithSwitch handles *no* collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleNoCollisionWithSwitch(collisionSwitchID terraform2d.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && entity.Toggler.Enabled() {
			entity.Toggler.Toggle()
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
		if id == collisionSwitchID && !entity.Toggler.Enabled() {
			entity.Toggler.Toggle()
		}
	}
}

// OnPlayerNoCollisionWithSwitch handles *no* collision between player and switch
func (ch *CollisionHandler) OnPlayerNoCollisionWithSwitch(collisionSwitchID terraform2d.EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && entity.Toggler.Enabled() {
			entity.Toggler.Toggle()
		}
	}
}

// OnPlayerCollisionWithWarp handles collision between player and warp
func (ch *CollisionHandler) OnPlayerCollisionWithWarp(warpID terraform2d.EntityID) {
	entityConfig, ok := ch.GameModel.RoomWarps[warpID]
	if ok && !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Style = rooms.TransitionWarp
		ch.GameModel.RoomTransition.Timer = 1
		ch.GameModel.CurrentState = gamestate.MapTransition
		ch.GameModel.AddEntities = true
		ch.GameModel.NextRoomID = entityConfig.WarpToRoomID
	}
}

// TransitionRoomResponse contains layout data
type TransitionRoomResponse struct {
	nextRoomID                                             rooms.RoomID
	modX, modY, modXNext, modYNext, playerModX, playerModY float64
}

func calculateTransitionSlide(
	roomTransition *rooms.RoomTransition,
	connectedRooms rooms.ConnectedRooms,
	currentRoomID rooms.RoomID) TransitionRoomResponse {

	var nextRoomID rooms.RoomID
	inc := (roomTransition.Start - float64(roomTransition.Timer))
	incY := inc * (config.MapH / config.TileSize)
	incX := inc * (config.MapW / config.TileSize)
	modY := 0.0
	modYNext := 0.0
	modX := 0.0
	modXNext := 0.0
	playerModX := 0.0
	playerModY := 0.0
	playerIncY := ((config.MapH / config.TileSize) - 1) + 7
	playerIncX := ((config.MapW / config.TileSize) - 1) + 7
	if roomTransition.Side == bounds.Bottom && connectedRooms.Bottom != 0 {
		modY = incY
		modYNext = incY - config.MapH
		nextRoomID = connectedRooms.Bottom
		playerModY += playerIncY
	} else if roomTransition.Side == bounds.Top && connectedRooms.Top != 0 {
		modY = -incY
		modYNext = -incY + config.MapH
		nextRoomID = connectedRooms.Top
		playerModY -= playerIncY
	} else if roomTransition.Side == bounds.Left && connectedRooms.Left != 0 {
		modX = incX
		modXNext = incX - config.MapW
		nextRoomID = connectedRooms.Left
		playerModX += playerIncX
	} else if roomTransition.Side == bounds.Right && connectedRooms.Right != 0 {
		modX = -incX
		modXNext = -incX + config.MapW
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
