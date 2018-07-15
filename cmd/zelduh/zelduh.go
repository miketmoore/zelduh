package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

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
	"github.com/miketmoore/zelduh/world"
	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/image/colornames"
)

var (
	win       *pixelgl.Window
	txt       *text.Text
	t         i18n.TranslateFunc
	gameWorld world.World
)

// GameModel contains data used throughout the game
type GameModel struct {
	AddEntities               bool
	CurrentRoomID, NextRoomID rooms.RoomID
	RoomTransition            rooms.RoomTransition
	CurrentState              gamestate.Name
	Rand                      *rand.Rand
	EntitiesMap               map[entities.EntityID]entities.Entity
	Spritesheet               map[int]*pixel.Sprite
}

func run() {

	gameModel := GameModel{
		CurrentState:  gamestate.Start,
		AddEntities:   true,
		CurrentRoomID: 1,
		RoomTransition: rooms.RoomTransition{
			Start: float64(config.TileSize),
		},
	}

	gameModel.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	gameModel.EntitiesMap = map[entities.EntityID]entities.Entity{}
	gameWorld = world.New()

	gamemap.ProcessMapLayout(roomsMap)

	// Initializations
	t = initI18n()
	txt = initText(20, 50, colornames.Black)
	win = initWindow(t("title"))

	// load the spritesheet image
	pic := loadPicture(config.SpritesheetPath)
	// build spritesheet
	// this is a map of TMX IDs to sprite instances
	gameModel.Spritesheet = sprites.BuildSpritesheet(pic, config.TileSize)

	allMapDrawData := tmx.BuildMapDrawData()

	// Build entities
	player := entities.BuildEntityFromConfig(entities.GetPreset("player")(6, 6), gameWorld.NewEntityID())
	bomb := entities.BuildEntityFromConfig(entities.GetPreset("bomb")(0, 0), gameWorld.NewEntityID())
	explosion := entities.BuildEntityFromConfig(entities.GetPreset("explosion")(0, 0), gameWorld.NewEntityID())
	sword := entities.BuildEntityFromConfig(entities.GetPreset("sword")(0, 0), gameWorld.NewEntityID())
	arrow := entities.BuildEntityFromConfig(entities.GetPreset("arrow")(0, 0), gameWorld.NewEntityID())

	var roomWarps map[entities.EntityID]rooms.EntityConfig

	// Create systems and add to game world
	inputSystem := &systems.Input{Win: win}
	gameWorld.AddSystem(inputSystem)
	healthSystem := &systems.Health{}
	gameWorld.AddSystem(healthSystem)
	spatialSystem := &systems.Spatial{
		Rand: gameModel.Rand,
	}
	dropCoin := func(v pixel.Vec) {
		coin := entities.BuildEntityFromConfig(entities.GetPreset("coin")(v.X/config.TileSize, v.Y/config.TileSize), gameWorld.NewEntityID())
		gameWorld.AddEntityToSystem(coin)
	}
	gameWorld.AddSystem(spatialSystem)

	hearts := []entities.Entity{
		entities.BuildEntityFromConfig(entities.GetPreset("heart")(1.5, 14), gameWorld.NewEntityID()),
		entities.BuildEntityFromConfig(entities.GetPreset("heart")(2.15, 14), gameWorld.NewEntityID()),
		entities.BuildEntityFromConfig(entities.GetPreset("heart")(2.80, 14), gameWorld.NewEntityID()),
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
		OnPlayerCollisionWithBounds: collisionHandler.OnPlayerCollisionWithBounds,
		OnPlayerCollisionWithCoin: func(coinID entities.EntityID) {
			player.Coins.Coins++
			gameWorld.Remove(categories.Coin, coinID)
		},
		OnPlayerCollisionWithEnemy: func(enemyID entities.EntityID) {
			// TODO repeat what I did with the enemies
			spatialSystem.MovePlayerBack()
			player.Health.Total--

			// remove heart entity
			heartIndex := len(hearts) - 1
			gameWorld.Remove(categories.Heart, hearts[heartIndex].ID)
			hearts = append(hearts[:heartIndex], hearts[heartIndex+1:]...)

			// TODO redraw hearts
			if player.Health.Total == 0 {
				gameModel.CurrentState = gamestate.Over
			}
		},
		OnSwordCollisionWithEnemy: func(enemyID entities.EntityID) {
			fmt.Printf("SwordCollisionWithEnemy %d\n", enemyID)
			if !sword.Ignore.Value {
				dead := false
				if !spatialSystem.EnemyMovingFromHit(enemyID) {
					dead = healthSystem.Hit(enemyID, 1)
					if dead {
						enemySpatial, _ := spatialSystem.GetEnemySpatial(enemyID)
						explosion.Temporary.Expiration = len(explosion.Animation.Map["default"].Frames)
						explosion.Spatial = &components.Spatial{
							Width:  config.TileSize,
							Height: config.TileSize,
							Rect:   enemySpatial.Rect,
						}
						explosion.Temporary.OnExpiration = func() {
							dropCoin(explosion.Spatial.Rect.Min)
						}
						gameWorld.AddEntityToSystem(explosion)
						gameWorld.RemoveEnemy(enemyID)
					} else {
						spatialSystem.MoveEnemyBack(enemyID, player.Movement.Direction)
					}
				}

			}
		},
		OnArrowCollisionWithEnemy: func(enemyID entities.EntityID) {
			if !arrow.Ignore.Value {
				dead := healthSystem.Hit(enemyID, 1)
				arrow.Ignore.Value = true
				if dead {
					fmt.Printf("You killed an enemy with an arrow\n")
					enemySpatial, _ := spatialSystem.GetEnemySpatial(enemyID)
					explosion.Temporary.Expiration = len(explosion.Animation.Map["default"].Frames)
					explosion.Spatial = &components.Spatial{
						Width:  config.TileSize,
						Height: config.TileSize,
						Rect:   enemySpatial.Rect,
					}
					explosion.Temporary.OnExpiration = func() {
						dropCoin(explosion.Spatial.Rect.Min)
					}
					gameWorld.AddEntityToSystem(explosion)
					gameWorld.RemoveEnemy(enemyID)
				} else {
					spatialSystem.MoveEnemyBack(enemyID, player.Movement.Direction)
				}
			}
		},
		OnArrowCollisionWithObstacle: func() {
			arrow.Movement.RemainingMoves = 0
		},
		OnPlayerCollisionWithObstacle: func(obstacleID entities.EntityID) {
			// "Block" by undoing rect
			player.Spatial.Rect = player.Spatial.PrevRect
			sword.Spatial.Rect = sword.Spatial.PrevRect
		},
		OnPlayerCollisionWithMoveableObstacle: func(obstacleID entities.EntityID) {
			moved := spatialSystem.MoveMoveableObstacle(obstacleID, player.Movement.Direction)
			if !moved {
				player.Spatial.Rect = player.Spatial.PrevRect
			}
		},
		OnMoveableObstacleCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range gameModel.EntitiesMap {
				if id == collisionSwitchID && !entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		OnMoveableObstacleNoCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range gameModel.EntitiesMap {
				if id == collisionSwitchID && entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		OnEnemyCollisionWithObstacle: func(enemyID, obstacleID entities.EntityID) {
			// Block enemy within the spatial system by reseting current rect to previous rect
			spatialSystem.UndoEnemyRect(enemyID)
		},
		OnPlayerCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range gameModel.EntitiesMap {
				if id == collisionSwitchID && !entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		OnPlayerNoCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range gameModel.EntitiesMap {
				if id == collisionSwitchID && entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		OnPlayerCollisionWithWarp: func(warpID entities.EntityID) {
			entityConfig, ok := roomWarps[warpID]
			if ok && !gameModel.RoomTransition.Active {
				gameModel.RoomTransition.Active = true
				gameModel.RoomTransition.Style = rooms.TransitionWarp
				gameModel.RoomTransition.Timer = 1
				gameModel.CurrentState = gamestate.MapTransition
				gameModel.AddEntities = true
				gameModel.NextRoomID = entityConfig.WarpToRoomID
			}
		},
	}
	gameWorld.AddSystem(collisionSystem)
	gameWorld.AddSystem(&systems.Render{
		Win:         win,
		Spritesheet: gameModel.Spritesheet,
	})

	gameWorld.AddEntitiesToSystem([]entities.Entity{player, sword, arrow, bomb})

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
			inputSystem.EnablePlayer()

			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

			drawMapBGImage(gameModel.Spritesheet, allMapDrawData, roomsMap[gameModel.CurrentRoomID].MapName, 0, 0)

			addHearts(hearts, player.Health.Total)

			if gameModel.AddEntities {
				gameModel.AddEntities = false

				addUICoin()

				// Draw obstacles on appropriate map tiles
				obstacles := drawObstaclesPerMapTiles(allMapDrawData, gameModel.CurrentRoomID, 0, 0)
				gameWorld.AddEntitiesToSystem(obstacles)

				roomWarps = map[entities.EntityID]rooms.EntityConfig{}

				// Iterate through all entity configurations and build entities and add to systems
				for _, c := range roomsMap[gameModel.CurrentRoomID].EntityConfigs {
					entity := entities.BuildEntityFromConfig(c, gameWorld.NewEntityID())
					gameModel.EntitiesMap[entity.ID] = entity
					gameWorld.AddEntityToSystem(entity)

					switch c.Category {
					case categories.Warp:
						roomWarps[entity.ID] = c
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
			inputSystem.DisablePlayer()
			if gameModel.RoomTransition.Style == rooms.TransitionSlide && gameModel.RoomTransition.Timer > 0 {
				gameModel.RoomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				gameWorld.RemoveAllEnemies()
				gameWorld.RemoveAllCollisionSwitches()
				gameWorld.RemoveAllMoveableObstacles()
				gameWorld.RemoveAllEntities()

				inc := (gameModel.RoomTransition.Start - float64(gameModel.RoomTransition.Timer))
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
				if gameModel.RoomTransition.Side == bounds.Bottom && roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Bottom != 0 {
					modY = incY
					modYNext = incY - config.MapH
					gameModel.NextRoomID = roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Bottom

					playerModY += playerIncY
				} else if gameModel.RoomTransition.Side == bounds.Top && roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Top != 0 {
					modY = -incY
					modYNext = -incY + config.MapH
					gameModel.NextRoomID = roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Top
					playerModY -= playerIncY
				} else if gameModel.RoomTransition.Side == bounds.Left && roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Left != 0 {
					modX = incX
					modXNext = incX - config.MapW
					gameModel.NextRoomID = roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Left
					playerModX += playerIncX
				} else if gameModel.RoomTransition.Side == bounds.Right && roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Right != 0 {
					modX = -incX
					modXNext = -incX + config.MapW
					gameModel.NextRoomID = roomsMap[gameModel.CurrentRoomID].ConnectedRooms.Right
					playerModX -= playerIncX
				} else {
					gameModel.NextRoomID = 0
				}

				drawMapBGImage(gameModel.Spritesheet, allMapDrawData, roomsMap[gameModel.CurrentRoomID].MapName, modX, modY)
				drawMapBGImage(gameModel.Spritesheet, allMapDrawData, roomsMap[gameModel.NextRoomID].MapName, modXNext, modYNext)
				drawMask()

				// Move player with map transition
				player.Spatial.Rect = pixel.R(
					player.Spatial.Rect.Min.X+playerModX,
					player.Spatial.Rect.Min.Y+playerModY,
					player.Spatial.Rect.Min.X+playerModX+config.TileSize,
					player.Spatial.Rect.Min.Y+playerModY+config.TileSize,
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

func loadPicture(path string) pixel.Picture {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Could not open the picture:")
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Could not decode the picture:")
		fmt.Println(err)
		os.Exit(1)
	}
	return pixel.PictureDataFromImage(img)
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
