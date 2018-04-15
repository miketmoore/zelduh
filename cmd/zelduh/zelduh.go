package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miketmoore/zelduh/bounds"
	"github.com/miketmoore/zelduh/rooms"

	"github.com/deanobob/tmxreader"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/go-pixel-game-template/state"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"github.com/miketmoore/zelduh/entities"
	"github.com/miketmoore/zelduh/gamestate"
	"github.com/miketmoore/zelduh/systems"
	"github.com/miketmoore/zelduh/world"
	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/image/colornames"
)

const (
	translationFile = "i18n/zelduh/en-US.all.json"
	lang            = "en-US"
)

const (
	winX float64 = 0
	winY float64 = 0
	winW float64 = 800
	winH float64 = 800
)

const (
	mapW float64 = 672 // 48 * 14
	mapH float64 = 576 // 48 * 12
	mapX         = (winW - mapW) / 2
	mapY         = (winH - mapH) / 2
)

var (
	win       *pixelgl.Window
	txt       *text.Text
	t         i18n.TranslateFunc
	currState state.State
	pic       pixel.Picture
	gameWorld world.World
)

const (
	s               float64 = 48
	spritesheetPath string  = "assets/spritesheet.png"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var tilemapDir = "assets/tilemaps/"
var tilemapFiles = []string{
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
}

var roomID rooms.RoomID

// Multi-dimensional array representing the overworld
// Each room ID should be unique
var overworld = [][]rooms.RoomID{
	[]rooms.RoomID{1, 10},
	[]rooms.RoomID{2, 0, 0, 8},
	[]rooms.RoomID{3, 5, 6, 7},
	[]rooms.RoomID{9},
	[]rooms.RoomID{11},
}

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

var spritesheet map[int]*pixel.Sprite
var tmxMapData map[string]tmxreader.TmxMap
var sprites map[string]*pixel.Sprite

type mapDrawData struct {
	Rect     pixel.Rect
	SpriteID int
}

var allMapDrawData map[string]MapData

const frameRate int = 5

type enemyPresetFn = func(xTiles, yTiles float64) rooms.EntityConfig

var spriteSets = map[string][]int{
	"eyeburrower":      []int{20, 20, 20, 91, 91, 91, 92, 92, 92, 93, 93, 93, 92, 92, 92},
	"skeleton":         []int{31, 32},
	"skull":            []int{36, 37, 38, 39},
	"spinner":          []int{51, 52},
	"puzzleBox":        []int{63},
	"warpStone":        []int{61},
	"playerUp":         []int{4, 195},
	"playerRight":      []int{3, 194},
	"playerDown":       []int{1, 192},
	"playerLeft":       []int{2, 193},
	"playerSwordUp":    []int{165},
	"playerSwordRight": []int{164},
	"playerSwordLeft":  []int{179},
	"playerSwordDown":  []int{180},
	"floorSwitch":      []int{112, 127},
}

var entityPresets = map[string]enemyPresetFn{
	"eyeburrower": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			HitBoxRadius: 20,
			SpriteFrames: spriteSets["eyeburrower"],
			Invincible:   false,
			PatternName:  "random",
			Direction:    direction.Down,
		}
	},
	"skeleton": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			HitBoxRadius: 20,
			SpriteFrames: spriteSets["skeleton"],
			Invincible:   false,
			PatternName:  "random",
			Direction:    direction.Down,
		}
	},
	"skull": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			HitBoxRadius: 20,
			SpriteFrames: spriteSets["skull"],
			Invincible:   false,
			PatternName:  "random",
			Direction:    direction.Down,
		}
	},
	"spinner": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			HitBoxRadius: 20,
			SpriteFrames: spriteSets["spinner"],
			Invincible:   true,
			PatternName:  "left-right",
			Direction:    direction.Right,
		}
	},
	"warpstone": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category:     categories.Warp,
			X:            s * xTiles,
			Y:            s * yTiles,
			W:            s,
			H:            s,
			SpriteFrames: spriteSets["warpStone"],
		}
	},
	"puzzleBox": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category:     categories.MovableObstacle,
			X:            s * xTiles,
			Y:            s * yTiles,
			W:            s,
			H:            s,
			SpriteFrames: spriteSets["puzzleBox"],
		}
	},
	"floorSwitch": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category:     categories.CollisionSwitch,
			X:            s * xTiles,
			Y:            s * yTiles,
			W:            s,
			H:            s,
			SpriteFrames: spriteSets["floorSwitch"],
		}
	},
}

func presetWarpStone(X, Y, WarpToRoomID, HitBoxRadius float64) rooms.EntityConfig {
	e := entityPresets["warpstone"](X, Y)
	e.WarpToRoomID = 6
	e.HitBoxRadius = 5
	return e
}

var roomsMap rooms.Rooms

func run() {

	gameWorld = world.New()

	fmt.Printf("build room configurations...\n")
	roomsMap = rooms.Rooms{
		1: rooms.NewRoom("overworldFourWallsDoorBottomRight",
			presetWarpStone(3, 7, 6, 5),
			entityPresets["skull"](5, 5),
			entityPresets["skeleton"](11, 9),
			entityPresets["spinner"](7, 9),
			entityPresets["eyeburrower"](8, 9),
		),
		2: rooms.NewRoom("overworldFourWallsDoorTopBottom",
			entityPresets["puzzleBox"](5, 5),
			entityPresets["floorSwitch"](10, 10),
		),
		3: rooms.NewRoom("overworldFourWallsDoorRightTopBottom"),
		5: rooms.NewRoom("rockWithCaveEntrance",
			rooms.EntityConfig{
				Category:     categories.Warp,
				WarpToRoomID: 11,
				W:            s,
				H:            s,
				X:            (s * 7) + s/2,
				Y:            (s * 9) + s/2,
				HitBoxRadius: 30,
			},
			rooms.EntityConfig{
				Category:     categories.Warp,
				WarpToRoomID: 11,
				W:            s,
				H:            s,
				X:            (s * 8) + s/2,
				Y:            (s * 9) + s/2,
				HitBoxRadius: 30,
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
				W:            s,
				H:            s,
				X:            (s * 6) + s + (s / 2.5),
				Y:            (s * 1) + s + (s / 2.5),
				HitBoxRadius: 15,
			},
			rooms.EntityConfig{
				Category:     categories.Warp,
				WarpToRoomID: 5,
				W:            s,
				H:            s,
				X:            (s * 7) + s + (s / 2.5),
				Y:            (s * 1) + s + (s / 2.5),
				HitBoxRadius: 15,
			},
		),
	}

	processMapLayout(overworld)

	// Initializations
	t = initI18n()
	txt = initText(20, 50, colornames.Black)
	win = initWindow(t("title"), winX, winY, winW, winH)

	// load the spritesheet image
	pic = loadPicture(spritesheetPath)
	// build spritesheet
	// this is a map of TMX IDs to sprite instances
	spritesheet = buildSpritesheet()

	// load all TMX file data for each map
	tmxMapData = loadTmxData()
	allMapDrawData = buildMapDrawData()

	// Build entities
	player := NewPlayer(gameWorld.NewEntityID(), s, s, mapW/2, mapH/2)
	sword := NewSword(gameWorld.NewEntityID(), s, s, player.Movement.Direction)
	arrow := NewArrow(gameWorld.NewEntityID(), 0.0, 0.0, s, s, player.Movement.Direction)
	bomb := NewBomb(gameWorld.NewEntityID(), 0.0, 0.0, s, s)
	explosion := NewExplosion(gameWorld.NewEntityID())

	hearts := []entities.Entity{
		NewHeart(gameWorld.NewEntityID(), mapX+16, mapY+mapH, s, s),
		NewHeart(gameWorld.NewEntityID(), mapX+48, mapY+mapH, s, s),
		NewHeart(gameWorld.NewEntityID(), mapX+80, mapY+mapH, s, s),
	}

	roomTransition := rooms.RoomTransition{
		Start: float64(s),
	}

	currentState := gamestate.Start
	addEntities := true
	var currentRoomID rooms.RoomID = 1
	var nextRoomID rooms.RoomID

	var roomWarps map[entities.EntityID]rooms.EntityConfig

	// Create systems and add to game world
	inputSystem := &systems.Input{Win: win}
	// gameWorld.SystemsMap["input"] = inputSystem
	gameWorld.AddSystem(inputSystem)
	healthSystem := &systems.Health{}
	// gameWorld.SystemsMap["health"] = healthSystem
	gameWorld.AddSystem(healthSystem)
	spatialSystem := &systems.Spatial{
		Rand: r,
	}
	dropCoin := func(v pixel.Vec) {
		fmt.Printf("Drop coin\n")
		coin := NewCoin(gameWorld.NewEntityID(), v.X, v.Y, s, s)
		addEntityToSystem(coin)
	}
	gameWorld.AddSystem(spatialSystem)
	collisionSystem := &systems.Collision{
		MapBounds: pixel.R(
			mapX,
			mapY,
			mapX+mapW,
			mapY+mapH,
		),
		PlayerCollisionWithBounds: func(side bounds.Bound) {
			if !roomTransition.Active {
				roomTransition.Active = true
				roomTransition.Side = side
				roomTransition.Style = rooms.TransitionSlide
				roomTransition.Timer = int(roomTransition.Start)
				currentState = gamestate.MapTransition
				addEntities = true
			}
		},
		PlayerCollisionWithCoin: func(coinID entities.EntityID) {
			player.Coins.Coins++
			fmt.Printf("Player coins: %d\n", player.Coins.Coins)
			gameWorld.Remove(categories.Coin, coinID)
		},
		PlayerCollisionWithEnemy: func(enemyID entities.EntityID) {
			// TODO repeat what I did with the enemies
			spatialSystem.MovePlayerBack()
			player.Health.Total--
			if player.Health.Total == 0 {
				currentState = gamestate.Over
			}
		},
		SwordCollisionWithEnemy: func(enemyID entities.EntityID) {
			if !sword.Ignore.Value {
				dead := false
				if !spatialSystem.EnemyMovingFromHit(enemyID) {
					fmt.Printf("HIT!\n")
					dead = healthSystem.Hit(enemyID, 1)
					if dead {
						fmt.Printf("You killed an enemy with a sword\n")
						enemySpatial, _ := spatialSystem.GetEnemySpatial(enemyID)
						explosion.Animation.Expiration = len(explosion.Animation.Default.Frames)
						explosion.Spatial = &components.Spatial{
							Width:  s,
							Height: s,
							Rect:   enemySpatial.Rect,
						}
						explosion.OnExpiration = func() {
							dropCoin(explosion.Spatial.Rect.Min)
						}
						addEntityToSystem(explosion)
						gameWorld.RemoveEnemy(enemyID)
					} else {
						spatialSystem.MoveEnemyBack(enemyID, player.Movement.Direction)
					}
				}

			}
		},
		ArrowCollisionWithEnemy: func(enemyID entities.EntityID) {
			if !arrow.Ignore.Value {
				dead := healthSystem.Hit(enemyID, 1)
				arrow.Ignore.Value = true
				if dead {
					fmt.Printf("You killed an enemy with an arrow\n")
					enemySpatial, _ := spatialSystem.GetEnemySpatial(enemyID)
					explosion.Animation.Expiration = len(explosion.Animation.Default.Frames)
					explosion.Spatial = &components.Spatial{
						Width:  s,
						Height: s,
						Rect:   enemySpatial.Rect,
					}
					explosion.OnExpiration = func() {
						dropCoin(explosion.Spatial.Rect.Min)
					}
					addEntityToSystem(explosion)
					gameWorld.RemoveEnemy(enemyID)
				} else {
					spatialSystem.MoveEnemyBack(enemyID, player.Movement.Direction)
				}
			}
		},
		ArrowCollisionWithObstacle: func() {
			arrow.Movement.RemainingMoves = 0
		},
		PlayerCollisionWithObstacle: func(obstacleID entities.EntityID) {
			// "Block" by undoing rect
			player.Spatial.Rect = player.Spatial.PrevRect
			sword.Spatial.Rect = sword.Spatial.PrevRect
		},
		PlayerCollisionWithMoveableObstacle: func(obstacleID entities.EntityID) {
			moved := spatialSystem.MoveMoveableObstacle(obstacleID, player.Movement.Direction)
			if !moved {
				player.Spatial.Rect = player.Spatial.PrevRect
			}
		},
		EnemyCollisionWithObstacle: func(enemyID, obstacleID entities.EntityID) {
			// Block enemy within the spatial system by reseting current rect to previous rect
			spatialSystem.UndoEnemyRect(enemyID)
		},
		PlayerCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			fmt.Printf("PlayerCollisionWithSwitch %d\n", collisionSwitchID)
			entityConfig, ok := roomWarps[collisionSwitchID]
			if ok {
				fmt.Printf("Warp Config: %v\n", entityConfig)
				fmt.Printf("Warp to room ID %v\n", entityConfig.WarpToRoomID)
				if !roomTransition.Active {
					roomTransition.Active = true
					roomTransition.Style = rooms.TransitionWarp
					roomTransition.Timer = 1
					currentState = gamestate.MapTransition
					addEntities = true
					nextRoomID = entityConfig.WarpToRoomID
				}
			}
		},
		PlayerNoCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			// fmt.Printf("PlayerNoCollisionWithSwitch\n")
		},
	}
	gameWorld.AddSystem(collisionSystem)
	gameWorld.AddSystem(&systems.Render{
		Win:         win,
		Spritesheet: spritesheet,
	})

	addEntityToSystem(player)
	addEntityToSystem(sword)
	addEntityToSystem(arrow)
	addEntityToSystem(bomb)
	for _, entity := range hearts {
		addEntityToSystem(entity)
	}

	drawMask := func() {
		// top
		s := imdraw.New(nil)
		s.Color = colornames.White
		s.Push(pixel.V(0, mapY+mapH))
		s.Push(pixel.V(winW, mapY+mapH+(winH-(mapY+mapH))))
		s.Rectangle(0)
		s.Draw(win)

		// bottom
		s = imdraw.New(nil)
		s.Color = colornames.White
		s.Push(pixel.V(0, 0))
		s.Push(pixel.V(winW, (winH - (mapY + mapH))))
		s.Rectangle(0)
		s.Draw(win)

		// left
		s = imdraw.New(nil)
		s.Color = colornames.White
		s.Push(pixel.V(0, 0))
		s.Push(pixel.V(0+mapX, winH))
		s.Rectangle(0)
		s.Draw(win)

		// right
		s = imdraw.New(nil)
		s.Color = colornames.White
		s.Push(pixel.V(mapX+mapW, mapY))
		s.Push(pixel.V(winW, winH))
		s.Rectangle(0)
		s.Draw(win)
	}

	for !win.Closed() {

		allowQuit()

		switch currentState {
		case gamestate.Start:
			win.Clear(colornames.Darkgray)
			drawMapBG(mapX, mapY, mapW, mapH, colornames.White)
			drawCenterText(t("title"), colornames.Black)

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = gamestate.Game
			}
		case gamestate.Game:
			inputSystem.EnablePlayer()

			win.Clear(colornames.Darkgray)
			drawMapBG(mapX, mapY, mapW, mapH, colornames.White)

			drawMapBGImage(roomsMap[currentRoomID].MapName, 0, 0)

			if addEntities {
				addEntities = false

				// Draw obstacles on appropriate map tiles
				obstacles := drawObstaclesPerMapTiles(currentRoomID, 0, 0)
				for _, entity := range obstacles {
					addEntityToSystem(entity)
				}

				roomWarps = map[entities.EntityID]rooms.EntityConfig{}

				// Iterate through all entity configurations and build entities and add to systems
				for _, c := range roomsMap[currentRoomID].EntityConfigs {
					var entity entities.Entity
					switch c.Category {
					case categories.CollisionSwitch:
						entity = NewCollisionSwitch(gameWorld.NewEntityID(), c.W, c.H, c.X, c.Y)
					case categories.MovableObstacle:
						entity = NewMoveableObstacle(gameWorld.NewEntityID(), c.W, c.H, c.X, c.Y)
					case categories.Warp:
						entity = NewCollisionSwitch(gameWorld.NewEntityID(), c.W, c.H, c.X, c.Y)
						entity.Spatial.HitBoxRadius = c.HitBoxRadius
					case categories.Enemy:
						fmt.Printf("Add enemy to room...\n")
						entity = NewEnemy(gameWorld.NewEntityID(), c.W, c.H, c.X, c.Y, c.HitBoxRadius, c.SpriteFrames, c.Invincible, c.PatternName, c.Direction)
					}

					if len(c.SpriteFrames) > 0 {
						entity.Animation = &components.Animation{
							Default: &components.AnimationData{
								Frames:    c.SpriteFrames,
								FrameRate: frameRate,
							},
						}
					}

					addEntityToSystem(entity)

					switch c.Category {
					case categories.Warp:
						roomWarps[entity.ID] = c
					}
				}
			}

			drawMask()

			gameWorld.Update()

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Pause
			}

			if win.JustPressed(pixelgl.KeyX) {
				currentState = gamestate.Over
			}

		case gamestate.Pause:
			win.Clear(colornames.Darkgray)
			drawMapBG(mapX, mapY, mapW, mapH, colornames.White)
			drawCenterText(t("paused"), colornames.Black)

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Game
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				currentState = gamestate.Start
			}
		case gamestate.Over:
			win.Clear(colornames.Darkgray)
			drawMapBG(mapX, mapY, mapW, mapH, colornames.Black)
			drawCenterText(t("gameOver"), colornames.White)

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = gamestate.Start
			}
		case gamestate.MapTransition:
			inputSystem.DisablePlayer()
			if roomTransition.Style == rooms.TransitionSlide && roomTransition.Timer > 0 {
				roomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(mapX, mapY, mapW, mapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				removeAllEnemiesFromSystems()
				removeAllCollisionSwitchesFromSystems()
				removeAllMoveableObstaclesFromSystems()

				inc := (roomTransition.Start - float64(roomTransition.Timer))
				incY := inc * (mapH / s)
				incX := inc * (mapW / s)
				modY := 0.0
				modYNext := 0.0
				modX := 0.0
				modXNext := 0.0
				playerModX := 0.0
				playerModY := 0.0
				playerIncY := ((mapH / s) - 1) + 7
				playerIncX := ((mapW / s) - 1) + 7
				if roomTransition.Side == bounds.Bottom && roomsMap[currentRoomID].ConnectedRooms.Bottom != 0 {
					modY = incY
					modYNext = incY - mapH
					nextRoomID = roomsMap[currentRoomID].ConnectedRooms.Bottom

					playerModY += playerIncY
				} else if roomTransition.Side == bounds.Top && roomsMap[currentRoomID].ConnectedRooms.Top != 0 {
					modY = -incY
					modYNext = -incY + mapH
					nextRoomID = roomsMap[currentRoomID].ConnectedRooms.Top
					playerModY -= playerIncY
				} else if roomTransition.Side == bounds.Left && roomsMap[currentRoomID].ConnectedRooms.Left != 0 {
					modX = incX
					modXNext = incX - mapW
					nextRoomID = roomsMap[currentRoomID].ConnectedRooms.Left
					playerModX += playerIncX
				} else if roomTransition.Side == bounds.Right && roomsMap[currentRoomID].ConnectedRooms.Right != 0 {
					modX = -incX
					modXNext = -incX + mapW
					nextRoomID = roomsMap[currentRoomID].ConnectedRooms.Right
					playerModX -= playerIncX
				} else {
					nextRoomID = 0
				}

				drawMapBGImage(roomsMap[currentRoomID].MapName, modX, modY)
				drawMapBGImage(roomsMap[nextRoomID].MapName, modXNext, modYNext)
				drawMask()

				// Move player with map transition
				player.Spatial.Rect = pixel.R(
					player.Spatial.Rect.Min.X+playerModX,
					player.Spatial.Rect.Min.Y+playerModY,
					player.Spatial.Rect.Min.X+playerModX+s,
					player.Spatial.Rect.Min.Y+playerModY+s,
				)

				gameWorld.Update()
			} else if roomTransition.Style == rooms.TransitionWarp && roomTransition.Timer > 0 {
				roomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(mapX, mapY, mapW, mapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				removeAllEnemiesFromSystems()
				removeAllCollisionSwitchesFromSystems()
			} else {
				currentState = gamestate.Game
				if nextRoomID != 0 {
					currentRoomID = nextRoomID
				}
				roomTransition.Active = false
			}

		}

		win.Update()

	}
}

func main() {
	pixelgl.Run(run)
}

func initI18n() i18n.TranslateFunc {
	i18n.LoadTranslationFile(translationFile)
	T, err := i18n.Tfunc(lang)
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

func initWindow(title string, x, y, w, h float64) *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(x, y, w, h),
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

// NewObstacle builds a new obstacle
func NewObstacle(id entities.EntityID, x, y float64) entities.Entity {
	return entities.Entity{
		ID:       gameWorld.NewEntityID(),
		Category: categories.Obstacle,
		Spatial: &components.Spatial{
			Width:  s,
			Height: s,
			Rect:   pixel.R(x, y, x+s, y+s),
			Shape:  imdraw.New(nil),
		},
	}
}

// NewPlayer builds a new Entity as a Player
func NewPlayer(id entities.EntityID, w, h, x, y float64) entities.Entity {
	return entities.Entity{
		ID:       id,
		Category: categories.Player,
		Health: &components.Health{
			Total: 3,
		},
		Appearance: &components.Appearance{
			Color: colornames.Green,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect: pixel.R(
				x,
				y,
				x+w,
				y+h,
			),
			Shape:                imdraw.New(nil),
			HitBox:               imdraw.New(nil),
			HitBoxRadius:         15,
			CollisionWithRectMod: 5,
		},
		Movement: &components.Movement{
			Direction: direction.Down,
			MaxSpeed:  7.0,
			Speed:     0.0,
		},
		Coins: &components.Coins{
			Coins: 0,
		},
		Dash: &components.Dash{
			Charge:    0,
			MaxCharge: 50,
			SpeedMod:  7,
		},
		Animation: &components.Animation{
			Up: &components.AnimationData{
				Frames:    spriteSets["playerUp"],
				FrameRate: frameRate,
			},
			Right: &components.AnimationData{
				Frames:    spriteSets["playerRight"],
				FrameRate: frameRate,
			},
			Down: &components.AnimationData{
				Frames:    spriteSets["playerDown"],
				FrameRate: frameRate,
			},
			Left: &components.AnimationData{
				Frames:    spriteSets["playerLeft"],
				FrameRate: frameRate,
			},
			SwordAttackUp: &components.AnimationData{
				Frames: spriteSets["playerSwordUp"],
			},
			SwordAttackRight: &components.AnimationData{
				Frames: spriteSets["playerSwordRight"],
			},
			SwordAttackLeft: &components.AnimationData{
				Frames: spriteSets["playerSwordLeft"],
			},
			SwordAttackDown: &components.AnimationData{
				Frames: spriteSets["playerSwordDown"],
			},
		},
	}
}

// NewSword builds a sword entity
func NewSword(id entities.EntityID, w, h float64, dir direction.Name) entities.Entity {
	return entities.Entity{
		ID:       id,
		Category: categories.Sword,
		Ignore: &components.Ignore{
			Value: true,
		},
		Appearance: &components.Appearance{
			Color: colornames.Deeppink,
		},
		Spatial: &components.Spatial{
			Width:        w,
			Height:       h,
			Rect:         pixel.R(0, 0, 0, 0),
			Shape:        imdraw.New(nil),
			HitBox:       imdraw.New(nil),
			HitBoxRadius: 20,
		},
		Movement: &components.Movement{
			Direction: dir,
			Speed:     0.0,
		},
		Animation: &components.Animation{
			Up: &components.AnimationData{
				Frames:    []int{70},
				FrameRate: frameRate,
			},
			Right: &components.AnimationData{
				Frames:    []int{67},
				FrameRate: frameRate,
			},
			Down: &components.AnimationData{
				Frames:    []int{68},
				FrameRate: frameRate,
			},
			Left: &components.AnimationData{
				Frames:    []int{69},
				FrameRate: frameRate,
			},
		},
	}
}

// NewArrow builds an arrow entity
func NewArrow(id entities.EntityID, x, y, w, h float64, dir direction.Name) entities.Entity {
	return entities.Entity{
		ID:       id,
		Category: categories.Arrow,
		Ignore: &components.Ignore{
			Value: true,
		},
		Appearance: &components.Appearance{
			Color: colornames.Deeppink,
		},
		Spatial: &components.Spatial{
			Width:        w,
			Height:       h,
			Rect:         pixel.R(x, y, x+w, y+h),
			Shape:        imdraw.New(nil),
			HitBox:       imdraw.New(nil),
			HitBoxRadius: 5,
		},
		Movement: &components.Movement{
			Direction: dir,
			Speed:     0.0,
		},
		Animation: &components.Animation{
			Up: &components.AnimationData{
				Frames:    []int{101},
				FrameRate: frameRate,
			},
			Right: &components.AnimationData{
				Frames:    []int{100},
				FrameRate: frameRate,
			},
			Down: &components.AnimationData{
				Frames:    []int{103},
				FrameRate: frameRate,
			},
			Left: &components.AnimationData{
				Frames:    []int{102},
				FrameRate: frameRate,
			},
		},
	}
}

// NewBomb builds a bomb entity
func NewBomb(id entities.EntityID, x, y, w, h float64) entities.Entity {
	return entities.Entity{
		ID:       id,
		Category: categories.Bomb,
		Ignore: &components.Ignore{
			Value: true,
		},
		Appearance: &components.Appearance{
			Color: colornames.Deeppink,
		},
		Spatial: &components.Spatial{
			Width:        w,
			Height:       h,
			Rect:         pixel.R(x, y, x+w, y+h),
			Shape:        imdraw.New(nil),
			HitBox:       imdraw.New(nil),
			HitBoxRadius: 5,
		},
		Animation: &components.Animation{
			Default: &components.AnimationData{
				Frames:    []int{138, 139, 140, 141},
				FrameRate: frameRate,
			},
		},
	}
}

// NewHeart builds a heart entity
func NewHeart(id entities.EntityID, x, y, w, h float64) entities.Entity {
	fmt.Printf("NewHeart %d\n", id)
	return entities.Entity{
		ID:       gameWorld.NewEntityID(),
		Category: categories.Heart,
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect:   pixel.R(x, y, x+w, y+h),
			Shape:  imdraw.New(nil),
			HitBox: imdraw.New(nil),
		},
		Animation: &components.Animation{
			Default: &components.AnimationData{
				Frames: []int{106},
			},
		},
	}
}

// NewExplosion builds an explosion entity
func NewExplosion(id entities.EntityID) entities.Entity {
	return entities.Entity{
		ID:       gameWorld.NewEntityID(),
		Category: categories.Explosion,
		Animation: &components.Animation{
			Expiration: 12,
			Default: &components.AnimationData{
				Frames: []int{
					122, 122, 122,
					123, 123, 123,
					124, 124, 124,
					125, 125, 125,
				},
			},
		},
	}
}

// NewMoveableObstacle builds a moveable obstacle
func NewMoveableObstacle(id entities.EntityID, w, h, x, y float64) entities.Entity {
	return entities.Entity{
		ID:       id,
		Category: categories.MovableObstacle,
		Appearance: &components.Appearance{
			Color: colornames.Purple,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect:   pixel.R(x, y, x+w, y+h),
			Shape:  imdraw.New(nil),
		},
		Movement: &components.Movement{
			Direction: direction.Down,
			Speed:     1.0,
			MaxMoves:  int(s) / 2,
			MaxSpeed:  2.0,
		},
	}
}

// NewCollisionSwitch builds a new collision switch
func NewCollisionSwitch(id entities.EntityID, w, h, x, y float64) entities.Entity {
	return entities.Entity{
		ID:       id,
		Category: categories.CollisionSwitch,
		Enabled: &components.Enabled{
			Value: false,
		},
		Appearance: &components.Appearance{
			Color: colornames.Sandybrown,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect:   pixel.R(x, y, x+w, y+h),
			Shape:  imdraw.New(nil),
		},
	}
}

// NewEnemy builds a new enemy
func NewEnemy(id entities.EntityID, w, h, x, y, hitRadius float64, frames []int, invincible bool, patternName string, dir direction.Name) entities.Entity {
	return entities.Entity{
		ID:       id,
		Category: categories.Enemy,
		Health:   &components.Health{Total: 2},
		Invincible: &components.Invincible{
			Enabled: invincible,
		},
		Appearance: &components.Appearance{
			Color: colornames.Red,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect: pixel.R(
				x,
				y,
				x+w,
				y+h,
			),
			Shape:        imdraw.New(nil),
			HitBox:       imdraw.New(nil),
			HitBoxRadius: hitRadius,
		},
		Movement: &components.Movement{
			Direction:    dir,
			Speed:        1.0,
			MaxSpeed:     1.0,
			HitSpeed:     10.0,
			HitBackMoves: 10,
			MaxMoves:     100,
			PatternName:  patternName,
		},
		Animation: &components.Animation{
			Default: &components.AnimationData{
				Frames:    frames,
				FrameRate: frameRate,
			},
		},
	}
}

// NewCoin builds a coin entity
func NewCoin(id entities.EntityID, x, y, w, h float64) entities.Entity {
	return entities.Entity{
		ID:       gameWorld.NewEntityID(),
		Category: categories.Coin,
		Appearance: &components.Appearance{
			Color: colornames.Purple,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect: pixel.R(
				x,
				y,
				x+w,
				y+h,
			),
			Shape: imdraw.New(nil),
		},
		Animation: &components.Animation{
			Default: &components.AnimationData{
				Frames:    []int{5, 5, 6, 6, 21, 21},
				FrameRate: frameRate,
			},
		},
	}
}

func addEntityToSystem(entity entities.Entity) {
	for _, system := range gameWorld.Systems() {
		system.AddEntity(entity)
	}
}

func removeAllEnemiesFromSystems() {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Spatial:
			sys.RemoveAll(categories.Enemy)
		case *systems.Collision:
			sys.RemoveAll(categories.Enemy)
		case *systems.Render:
			sys.RemoveAll(categories.Enemy)
		}
	}
}

func removeAllCollisionSwitchesFromSystems() {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
			sys.RemoveAll(categories.CollisionSwitch)
		case *systems.Render:
			sys.RemoveAll(categories.CollisionSwitch)
		}
	}
}

func removeAllMoveableObstaclesFromSystems() {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
			sys.RemoveAll(categories.MovableObstacle)
		case *systems.Render:
			sys.RemoveAll(categories.MovableObstacle)
		}
	}
}

func loadTmxData() map[string]tmxreader.TmxMap {
	tmxMapData := map[string]tmxreader.TmxMap{}
	for _, name := range tilemapFiles {
		path := fmt.Sprintf("%s%s.tmx", tilemapDir, name)
		tmxMapData[name] = parseTmxFile(path)
	}
	return tmxMapData
}

func parseTmxFile(filename string) tmxreader.TmxMap {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	tmxMap, err := tmxreader.Parse(raw)
	if err != nil {
		panic(err)
	}

	return tmxMap
}

// this is a map of pixel engine sprites
func buildSpritesheet() map[int]*pixel.Sprite {
	cols := pic.Bounds().W() / s
	rows := pic.Bounds().H() / s

	maxIndex := (rows * cols) - 1.0

	index := maxIndex
	id := maxIndex + 1
	spritesheet := map[int]*pixel.Sprite{}
	for row := (rows - 1); row >= 0; row-- {
		for col := (cols - 1); col >= 0; col-- {
			x := col
			y := math.Abs(rows-row) - 1
			spritesheet[int(id)] = pixel.NewSprite(pic, pixel.R(
				x*s,
				y*s,
				x*s+s,
				y*s+s,
			))
			index--
			id--
		}
	}
	return spritesheet
}

// MapData represents data for one map
type MapData struct {
	Name string
	Data []mapDrawData
}

func buildMapDrawData() map[string]MapData {
	all := map[string]MapData{}

	for mapName, mapData := range tmxMapData {
		// fmt.Printf("Building map draw data for map %s\n", mapName)

		md := MapData{
			Name: mapName,
			Data: []mapDrawData{},
		}

		layers := mapData.Layers
		for _, layer := range layers {

			records := parseCsv(strings.TrimSpace(layer.Data.Value) + ",")
			for row := 0; row <= len(records); row++ {
				if len(records) > row {
					for col := 0; col < len(records[row])-1; col++ {
						y := float64(11-row) * s
						x := float64(col) * s

						record := records[row][col]
						spriteID, err := strconv.Atoi(record)
						if err != nil {
							panic(err)
						}
						mrd := mapDrawData{
							Rect:     pixel.R(x, y, x+s, y+s),
							SpriteID: spriteID,
						}
						md.Data = append(md.Data, mrd)
					}
				}

			}
			all[mapName] = md
		}
	}

	return all
}

func parseCsv(in string) [][]string {
	r := csv.NewReader(strings.NewReader(in))

	records := [][]string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, record)
	}

	return records
}

func drawMapBGImage(name string, modX, modY float64) {
	d := allMapDrawData[name]
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			sprite := spritesheet[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+mapX+modX+s/2,
				vec.Y+mapY+modY+s/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func drawObstaclesPerMapTiles(roomID rooms.RoomID, modX, modY float64) []entities.Entity {
	d := allMapDrawData[roomsMap[roomID].MapName]
	obstacles := []entities.Entity{}
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+mapX+modX+s/2,
				vec.Y+mapY+modY+s/2,
			)
			if _, ok := nonObstacleSprites[spriteData.SpriteID]; !ok {
				obstacle := NewObstacle(gameWorld.NewEntityID(), movedVec.X-s/2, movedVec.Y-s/2)
				obstacles = append(obstacles, obstacle)
			}
		}
	}
	return obstacles
}

func indexRoom(a, b rooms.RoomID, dir direction.Name) {
	// fmt.Printf("indexRoom a:%d b:%d dir:%s\n", a, b, dir)
	roomA, okA := roomsMap[a]
	roomB, okB := roomsMap[b]
	if okA && okB {
		switch dir {
		case direction.Up:
			// b is above a
			roomA.ConnectedRooms.Top = b
			roomsMap[a] = roomA
			roomB.ConnectedRooms.Bottom = a
			roomsMap[b] = roomB
		case direction.Right:
			// b is right of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.ConnectedRooms.Right = b
				roomsMap[a] = roomA
				roomB.ConnectedRooms.Left = a
				roomsMap[b] = roomB
			}
		case direction.Down:
			// b is below a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.ConnectedRooms.Bottom = b
				roomsMap[a] = roomA
				roomB.ConnectedRooms.Top = a
				roomsMap[b] = roomB
			}
		case direction.Left:
			// b is left of a
			roomA, ok := roomsMap[a]
			if ok {
				roomA.ConnectedRooms.Left = b
				roomsMap[a] = roomA
				roomB.ConnectedRooms.Right = a
				roomsMap[b] = roomB
			}
		}
	}

}

func processMapLayout(layout [][]rooms.RoomID) {
	// transform multi-dimensional array into map of Room structs, indexed by ID
	for row := 0; row < len(layout); row++ {
		for col := 0; col < len(layout[row]); col++ {
			roomID := layout[row][col]
			// fmt.Printf("Room ID: %d\n", roomID)
			// Top
			if row > 0 {
				if len(layout[row-1]) > col {
					n := layout[row-1][col]
					if n > 0 {
						// fmt.Printf("\t%d is below %d\n", roomID, n)
						indexRoom(roomID, n, direction.Up)
					}
				}
			}
			// Right
			if len(layout[row]) > col+1 {
				n := layout[row][col+1]
				if n > 0 {
					// fmt.Printf("\t%d is left of %d\n", roomID, n)
					indexRoom(roomID, n, direction.Right)
				}
			}
			// Bottom
			if len(layout) > row+1 {
				if len(layout[row+1]) > col {
					n := layout[row+1][col]
					if n > 0 {
						// fmt.Printf("\t%d is above %d\n", roomID, n)
						indexRoom(roomID, n, direction.Down)
					}
				}
			}
			// Left
			if col > 0 {
				n := layout[row][col-1]
				if n > 0 {
					// fmt.Printf("\t%d is right of %d\n", roomID, n)
					indexRoom(roomID, n, direction.Left)
				}
			}
		}
	}
}
