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
	spriteSize      float64 = 48
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

var roomID RoomID

// RoomID is a room ID
type RoomID int

// ConnectedRooms is used to configure adjacent rooms
type ConnectedRooms struct {
	Top    RoomID
	Right  RoomID
	Bottom RoomID
	Left   RoomID
}

// EnemyConfig is used to build enemies
type EnemyConfig struct {
	W, H, X, Y, HitBoxRadius float64
}

// WarpConfig is used to build warps
type WarpConfig struct {
	W, H, X, Y, HitBoxRadius float64
	WarpToRoomID             RoomID
	IsAnimated               bool
}

// MoveableObstacleConfig is used to build moveable obstacles
type MoveableObstacleConfig struct {
	W, H, X, Y float64
	IsAnimated bool
}

// Room represents one map section
type Room struct {
	MapName                 string
	ConnectedRooms          ConnectedRooms
	EnemyConfigs            []EnemyConfig
	WarpConfigs             []WarpConfig
	MoveableObstacleConfigs []MoveableObstacleConfig
}

var rooms map[RoomID]Room

// Multi-dimensional array representing the overworld
// Each room ID should be unique
var overworld = [][]RoomID{
	[]RoomID{1, 10},
	[]RoomID{2, 0, 0, 8},
	[]RoomID{3, 5, 6, 7},
	[]RoomID{9},
	[]RoomID{11},
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

type transitionStyle string

const (
	transitionSlide transitionStyle = "slide"
	transitionWarp  transitionStyle = "warp"
)

func run() {

	gameWorld = world.New()

	fmt.Printf("build room configurations...\n")
	rooms = map[RoomID]Room{
		1: Room{
			MapName: "overworldFourWallsDoorBottomRight",
			EnemyConfigs: []EnemyConfig{
				EnemyConfig{spriteSize, spriteSize, spriteSize * 5, spriteSize * 5, 20},
				EnemyConfig{spriteSize, spriteSize, spriteSize * 11, spriteSize * 9, 20},
			},
			WarpConfigs: []WarpConfig{
				WarpConfig{
					WarpToRoomID: 6,
					W:            spriteSize,
					H:            spriteSize,
					X:            spriteSize * 3,
					Y:            spriteSize * 7,
					IsAnimated:   true,
					HitBoxRadius: 5,
				},
			},
		},
		2: Room{
			MapName: "overworldFourWallsDoorTopBottom",
			EnemyConfigs: []EnemyConfig{
				EnemyConfig{spriteSize, spriteSize, spriteSize * 5, spriteSize * 5, 20},
				EnemyConfig{spriteSize, spriteSize, spriteSize * 11, spriteSize * 9, 20},
			},
		},
		3: Room{
			MapName: "overworldFourWallsDoorRightTopBottom",
			EnemyConfigs: []EnemyConfig{
				EnemyConfig{spriteSize, spriteSize, spriteSize * 5, spriteSize * 5, 20},
				EnemyConfig{spriteSize, spriteSize, spriteSize * 11, spriteSize * 9, 20},
			},
		},
		5: Room{
			MapName: "rockWithCaveEntrance",
			WarpConfigs: []WarpConfig{
				WarpConfig{
					WarpToRoomID: 11,
					W:            spriteSize,
					H:            spriteSize,
					X:            (spriteSize * 7) + spriteSize/2,
					Y:            (spriteSize * 9) + spriteSize/2,
					HitBoxRadius: 30,
				},
				WarpConfig{
					WarpToRoomID: 11,
					W:            spriteSize,
					H:            spriteSize,
					X:            (spriteSize * 8) + spriteSize/2,
					Y:            (spriteSize * 9) + spriteSize/2,
					HitBoxRadius: 30,
				},
			},
		},
		6: Room{MapName: "rockPathLeftRightEntrance"},
		7: Room{
			MapName: "overworldFourWallsDoorLeftTop",
			MoveableObstacleConfigs: []MoveableObstacleConfig{
				MoveableObstacleConfig{
					W: spriteSize,
					H: spriteSize,
					X: (spriteSize * 8) + spriteSize/2,
					Y: (spriteSize * 9) + spriteSize/2,
				},
			},
		},
		8: Room{MapName: "overworldFourWallsDoorBottom"},
		9: Room{MapName: "overworldFourWallsDoorTop"},
		10: Room{
			MapName: "overworldFourWallsDoorLeft",
			EnemyConfigs: []EnemyConfig{
				EnemyConfig{spriteSize, spriteSize, spriteSize * 5, spriteSize * 9, 20},
			},
		},
		11: Room{
			MapName: "dungeonFourDoors",
			WarpConfigs: []WarpConfig{
				// South door of cave - warp to cave entrance
				WarpConfig{
					WarpToRoomID: 5,
					W:            spriteSize,
					H:            spriteSize,
					X:            (spriteSize * 6) + spriteSize + (spriteSize / 2.5),
					Y:            (spriteSize * 1) + spriteSize + (spriteSize / 2.5),
					HitBoxRadius: 15,
				},
				WarpConfig{
					WarpToRoomID: 5,
					W:            spriteSize,
					H:            spriteSize,
					X:            (spriteSize * 7) + spriteSize + (spriteSize / 2.5),
					Y:            (spriteSize * 1) + spriteSize + (spriteSize / 2.5),
					HitBoxRadius: 15,
				},
			},
		},
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
	frameRate := 5
	player := entities.BuildPlayer(spriteSize, spriteSize, mapW/2, mapH/2)
	player.Animation = &components.Animation{
		Up: &components.AnimationData{
			Frames:    []int{4, 195},
			FrameRate: frameRate,
		},
		Right: &components.AnimationData{
			Frames:    []int{3, 194},
			FrameRate: frameRate,
		},
		Down: &components.AnimationData{
			Frames:    []int{1, 192},
			FrameRate: frameRate,
		},
		Left: &components.AnimationData{
			Frames:    []int{2, 193},
			FrameRate: frameRate,
		},
		SwordAttackUp: &components.AnimationData{
			Frames: []int{165},
		},
		SwordAttackRight: &components.AnimationData{
			Frames: []int{164},
		},
		SwordAttackLeft: &components.AnimationData{
			Frames: []int{179},
		},
		SwordAttackDown: &components.AnimationData{
			Frames: []int{180},
		},
	}

	explosion := entities.Generic{
		ID: gameWorld.NewEntityID(),
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

	sword := entities.BuildSword(spriteSize, spriteSize, player.Movement.Direction)
	sword.Animation = &components.Animation{
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
	}

	arrow := entities.BuildArrow(spriteSize, spriteSize, 0.0, 0.0, player.Movement.Direction)
	arrow.Animation = &components.Animation{
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
	}

	isTransitioning := false
	var transitionSide bounds.Bound
	var transitionTimerStart = float64(spriteSize)
	var transitionTimer int
	var transitionStyle transitionStyle
	currentState := gamestate.Start
	addEntities := true
	var currentRoomID RoomID = 1
	var nextRoomID RoomID

	var roomWarps map[entities.EntityID]WarpConfig

	heartSprite := spritesheet[106]

	// heart := entities.Generic{
	// 	Category: categories.UI,
	// 	Animation: &components.Animation{
	// 		Default: &components.AnimationData{
	// 			Frames: []int{106},
	// 		},
	// 	},
	// }

	// addHeartToSystems(heart)

	// Create systems and add to game world
	inputSystem := &systems.Input{Win: win}
	gameWorld.AddSystem(inputSystem)
	healthSystem := &systems.Health{}
	gameWorld.AddSystem(healthSystem)
	spatialSystem := &systems.Spatial{
		Rand: r,
	}
	appendCoinAnimation := func(coin *entities.Coin) {
		coin.Animation = &components.Animation{
			Default: &components.AnimationData{
				Frames:    []int{5, 5, 6, 6, 21, 21},
				FrameRate: frameRate,
			},
		}
	}
	dropCoin := func(v pixel.Vec) {
		fmt.Printf("Drop coin\n")
		coin := buildCoin(v.X, v.Y)
		appendCoinAnimation(&coin)
		addCoinToSystem(coin)
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
			if !isTransitioning {
				isTransitioning = true
				transitionSide = side
				transitionStyle = transitionSlide
				transitionTimer = int(transitionTimerStart)
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
			spatialSystem.MovePlayerBack()
			player.Health.Total--
			if player.Health.Total == 0 {
				currentState = gamestate.Over
			}
		},
		SwordCollisionWithEnemy: func(enemyID entities.EntityID) {
			if !sword.Ignore.Value {
				dead := healthSystem.Hit(enemyID, 1)
				if dead {
					fmt.Printf("You killed an enemy with a sword\n")
					enemySpatial, _ := spatialSystem.GetEnemySpatial(enemyID)
					explosion.Animation.Expiration = len(explosion.Animation.Default.Frames)
					explosion.Spatial = &components.Spatial{
						Width:  spriteSize,
						Height: spriteSize,
						Rect:   enemySpatial.Rect,
					}
					explosion.OnExpiration = func() {
						dropCoin(explosion.Spatial.Rect.Min)
					}
					addGenericToSystems(categories.Explosion, explosion, enemySpatial.Rect.Min)
					gameWorld.RemoveEnemy(enemyID)
				} else {
					// TODO - instead of just moving back in the spatial system,
					// what about setting a velocity (direction * speed)
					// decrease speed over time
					// collision system will handle ... collisions!
					// Currently, just moving from A to B without any in-between movement allows
					// the entity to pass through obstacles.

					spatialSystem.MoveEnemyBack(enemyID, player.Movement.Direction, spriteSize)
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
						Width:  spriteSize,
						Height: spriteSize,
						Rect:   enemySpatial.Rect,
					}
					explosion.OnExpiration = func() {
						dropCoin(explosion.Spatial.Rect.Min)
					}
					addGenericToSystems(categories.Explosion, explosion, enemySpatial.Rect.Min)
					gameWorld.RemoveEnemy(enemyID)
				} else {
					spatialSystem.MoveEnemyBack(enemyID, player.Movement.Direction, spriteSize*3)
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
			spatialSystem.MoveMoveableObstacle(obstacleID, player.Movement.Direction)
		},
		EnemyCollisionWithObstacle: func(enemyID, obstacleID entities.EntityID) {
			// Block enemy within the spatial system by reseting current rect to previous rect
			spatialSystem.UndoEnemyRect(enemyID)
		},
		PlayerCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			fmt.Printf("PlayerCollisionWithSwitch %d\n", collisionSwitchID)
			warpConfig, ok := roomWarps[collisionSwitchID]
			if ok {
				fmt.Printf("Warp Config: %v\n", warpConfig)
				fmt.Printf("Warp to room ID %v\n", warpConfig.WarpToRoomID)
				if !isTransitioning {
					isTransitioning = true
					transitionStyle = transitionWarp
					transitionTimer = 1
					currentState = gamestate.MapTransition
					addEntities = true
					nextRoomID = warpConfig.WarpToRoomID
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

	addPlayerToSystems(player)
	addSwordToSystems(sword)
	addArrowToSystems(arrow)

	// TODO move
	drawHeart := func(offsetX, offsetY float64) {
		v := pixel.V(
			mapX+offsetX,
			mapY+mapH+offsetY,
		)
		matrix := pixel.IM.Moved(v)
		heartSprite.Draw(win, matrix)
	}
	drawHearts := func(health int) {
		switch health {
		case 3:
			drawHeart(96, 16)
			fallthrough
		case 2:
			drawHeart(64, 16)
			fallthrough
		case 1:
			drawHeart(32, 16)
		}
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

			drawMapBGImage(rooms[currentRoomID].MapName, 0, 0)

			if addEntities {
				addEntities = false
				obstacles := drawObstaclesPerMapTiles(currentRoomID, 0, 0)
				addObstaclesToSystem(obstacles)

				for _, c := range rooms[currentRoomID].MoveableObstacleConfigs {
					entity := entities.BuildMoveableObstacle(gameWorld.NewEntityID(), c.W, c.H, c.X, c.Y)
					entity.Animation = &components.Animation{
						Default: &components.AnimationData{
							Frames:    []int{63},
							FrameRate: frameRate,
						},
					}
					addMoveableObstaclesToSystem([]entities.MoveableObstacle{entity})
				}

				for _, c := range rooms[currentRoomID].EnemyConfigs {
					enemy := entities.BuildEnemy(gameWorld.NewEntityID(), c.W, c.H, c.X, c.Y, c.HitBoxRadius)
					enemy.Animation = &components.Animation{
						Default: &components.AnimationData{
							Frames:    []int{36, 37, 38, 39},
							FrameRate: frameRate,
						},
					}
					addEnemiesToSystem([]entities.Enemy{enemy})
				}

				roomWarps = map[entities.EntityID]WarpConfig{}
				for _, c := range rooms[currentRoomID].WarpConfigs {
					warp := entities.BuildCollisionSwitch(gameWorld.NewEntityID(), c.W, c.H, c.X, c.Y)
					if c.IsAnimated {
						warp.Animation = &components.Animation{
							Default: &components.AnimationData{
								Frames:    []int{61},
								FrameRate: frameRate,
							},
						}
					}
					warp.Spatial.HitBoxRadius = c.HitBoxRadius
					addCollisionSwitchesToSystem([]entities.CollisionSwitch{warp})
					roomWarps[warp.ID] = c
				}
			}

			drawMask()
			drawHearts(player.Health.Total)

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
			if transitionStyle == transitionSlide && transitionTimer > 0 {
				transitionTimer--
				win.Clear(colornames.Darkgray)
				drawMapBG(mapX, mapY, mapW, mapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				removeAllEnemiesFromSystems()
				removeAllCollisionSwitchesFromSystems()
				removeAllMoveableObstaclesFromSystems()

				inc := (transitionTimerStart - float64(transitionTimer))
				incY := inc * (mapH / spriteSize)
				incX := inc * (mapW / spriteSize)
				modY := 0.0
				modYNext := 0.0
				modX := 0.0
				modXNext := 0.0
				playerModX := 0.0
				playerModY := 0.0
				playerIncY := ((mapH / spriteSize) - 1) + 7
				playerIncX := ((mapW / spriteSize) - 1) + 7
				if transitionSide == "bottom" && rooms[currentRoomID].ConnectedRooms.Bottom != 0 {
					modY = incY
					modYNext = incY - mapH
					nextRoomID = rooms[currentRoomID].ConnectedRooms.Bottom

					playerModY += playerIncY
				} else if transitionSide == "top" && rooms[currentRoomID].ConnectedRooms.Top != 0 {
					modY = -incY
					modYNext = -incY + mapH
					nextRoomID = rooms[currentRoomID].ConnectedRooms.Top
					playerModY -= playerIncY
				} else if transitionSide == "left" && rooms[currentRoomID].ConnectedRooms.Left != 0 {
					modX = incX
					modXNext = incX - mapW
					nextRoomID = rooms[currentRoomID].ConnectedRooms.Left
					playerModX += playerIncX
				} else if transitionSide == "right" && rooms[currentRoomID].ConnectedRooms.Right != 0 {
					modX = -incX
					modXNext = -incX + mapW
					nextRoomID = rooms[currentRoomID].ConnectedRooms.Right
					playerModX -= playerIncX
				} else {
					nextRoomID = 0
				}

				drawMapBGImage(rooms[currentRoomID].MapName, modX, modY)
				drawMapBGImage(rooms[nextRoomID].MapName, modXNext, modYNext)
				drawMask()
				drawHearts(player.Health.Total)

				// Move player with map transition
				player.Spatial.Rect = pixel.R(
					player.Spatial.Rect.Min.X+playerModX,
					player.Spatial.Rect.Min.Y+playerModY,
					player.Spatial.Rect.Min.X+playerModX+spriteSize,
					player.Spatial.Rect.Min.Y+playerModY+spriteSize,
				)

				gameWorld.Update()
			} else if transitionStyle == transitionWarp && transitionTimer > 0 {
				transitionTimer--
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
				isTransitioning = false
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

func buildCoin(x, y float64) entities.Coin {
	w := spriteSize
	h := spriteSize
	return entities.Coin{
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
	}
}

func buildObstacle(x, y float64) entities.Obstacle {
	return entities.Obstacle{
		ID:       gameWorld.NewEntityID(),
		Category: categories.Obstacle,
		Spatial: &components.Spatial{
			Width:  spriteSize,
			Height: spriteSize,
			Rect:   pixel.R(x, y, x+spriteSize, y+spriteSize),
			Shape:  imdraw.New(nil),
		},
	}
}

func buildMoveableObstacle(x, y float64) entities.MoveableObstacle {
	return entities.MoveableObstacle{
		ID: gameWorld.NewEntityID(),
		Spatial: &components.Spatial{
			Width:  spriteSize,
			Height: spriteSize,
			Rect:   pixel.R(x, y, x+spriteSize, y+spriteSize),
			Shape:  imdraw.New(nil),
		},
		Appearance: &components.Appearance{
			Color: colornames.Blueviolet,
		},
	}
}

func addCoinToSystem(coin entities.Coin) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
			sys.Add(coin.Category, coin.ID, coin.Spatial)
		case *systems.Render:
			sys.Add(coin.Category, coin.ID, coin.Appearance, coin.Spatial, coin.Animation, nil, nil)
		}
	}
}

func addPlayerToSystems(player entities.Player) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Input:
			sys.AddPlayer(player.Movement, player.Dash)
		case *systems.Spatial:
			sys.Add(categories.Player, 0, player.Spatial, player.Movement, player.Dash)
		case *systems.Collision:
			sys.Add(player.Category, 0, player.Spatial)
		case *systems.Render:
			sys.Add(player.Category, 0, player.Appearance, player.Spatial, player.Animation, player.Movement, nil)
		}
	}
}

func addSwordToSystems(sword entities.Sword) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Input:
			sys.AddSword(sword.Movement, sword.Ignore)
		case *systems.Spatial:
			sys.Add(categories.Sword, 0, sword.Spatial, sword.Movement, nil)
		case *systems.Collision:
			sys.Add(sword.Category, 0, sword.Spatial)
		case *systems.Render:
			sys.Add(sword.Category, 0, sword.Appearance, sword.Spatial, sword.Animation, nil, sword.Ignore)
		}
	}
}

func addArrowToSystems(arrow entities.Arrow) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Input:
			sys.AddArrow(arrow.Movement, arrow.Ignore)
		case *systems.Spatial:
			sys.Add(categories.Arrow, 0, arrow.Spatial, arrow.Movement, nil)
		case *systems.Collision:
			sys.Add(arrow.Category, 0, arrow.Spatial)
		case *systems.Render:
			sys.Add(arrow.Category, 0, arrow.Appearance, arrow.Spatial, arrow.Animation, nil, arrow.Ignore)
		}
	}
}

func addEnemiesToSystem(enemies []entities.Enemy) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Spatial:
			for _, enemy := range enemies {
				sys.Add(categories.Enemy, enemy.ID, enemy.Spatial, enemy.Movement, nil)
			}
		case *systems.Collision:
			for _, enemy := range enemies {
				sys.Add(enemy.Category, enemy.ID, enemy.Spatial)
			}
		case *systems.Health:
			for _, enemy := range enemies {
				sys.AddEntity(enemy.ID, enemy.Health)
			}
		case *systems.Render:
			for _, enemy := range enemies {
				sys.Add(enemy.Category, enemy.ID, enemy.Appearance, enemy.Spatial, enemy.Animation, enemy.Movement, nil)
			}
		}
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

func addObstaclesToSystem(obstacles []entities.Obstacle) {
	for i, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
			fmt.Printf("addObstaclesToSystem %d\n", i)
			for _, obstacle := range obstacles {
				sys.Add(obstacle.Category, obstacle.ID, obstacle.Spatial)
			}
		}
	}
}

func addMoveableObstaclesToSystem(moveableObstacles []entities.MoveableObstacle) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Spatial:
			for _, moveable := range moveableObstacles {
				sys.Add(categories.MovableObstacle, moveable.ID, moveable.Spatial, moveable.Movement, nil)
			}
		case *systems.Collision:
			for _, moveable := range moveableObstacles {
				sys.Add(moveable.Category, moveable.ID, moveable.Spatial)
			}
		case *systems.Render:
			for _, moveable := range moveableObstacles {
				sys.Add(moveable.Category, moveable.ID, moveable.Appearance, moveable.Spatial, moveable.Animation, nil, nil)
			}
		}
	}
}

func addCollisionSwitchesToSystem(collisionSwitches []entities.CollisionSwitch) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
			for _, collisionSwitch := range collisionSwitches {
				sys.Add(collisionSwitch.Category, collisionSwitch.ID, collisionSwitch.Spatial)
			}
		case *systems.Render:
			for _, collisionSwitch := range collisionSwitches {
				sys.Add(collisionSwitch.Category, 0, collisionSwitch.Appearance, collisionSwitch.Spatial, collisionSwitch.Animation, nil, nil)
			}
		}
	}
}

func addGenericToSystems(category categories.Category, generic entities.Generic, v pixel.Vec) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Render:
			sys.Add(category, 0, nil, generic.Spatial, generic.Animation, nil, nil)
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
	cols := pic.Bounds().W() / spriteSize
	rows := pic.Bounds().H() / spriteSize

	maxIndex := (rows * cols) - 1.0

	index := maxIndex
	id := maxIndex + 1
	spritesheet := map[int]*pixel.Sprite{}
	for row := (rows - 1); row >= 0; row-- {
		for col := (cols - 1); col >= 0; col-- {
			x := col
			y := math.Abs(rows-row) - 1
			spritesheet[int(id)] = pixel.NewSprite(pic, pixel.R(
				x*spriteSize,
				y*spriteSize,
				x*spriteSize+spriteSize,
				y*spriteSize+spriteSize,
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
						y := float64(11-row) * spriteSize
						x := float64(col) * spriteSize

						record := records[row][col]
						spriteID, err := strconv.Atoi(record)
						if err != nil {
							panic(err)
						}
						mrd := mapDrawData{
							Rect:     pixel.R(x, y, x+spriteSize, y+spriteSize),
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
				vec.X+mapX+modX+spriteSize/2,
				vec.Y+mapY+modY+spriteSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func drawObstaclesPerMapTiles(roomID RoomID, modX, modY float64) []entities.Obstacle {
	d := allMapDrawData[rooms[roomID].MapName]
	obstacles := []entities.Obstacle{}
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+mapX+modX+spriteSize/2,
				vec.Y+mapY+modY+spriteSize/2,
			)
			if _, ok := nonObstacleSprites[spriteData.SpriteID]; !ok {
				obstacle := buildObstacle(movedVec.X-spriteSize/2, movedVec.Y-spriteSize/2)
				obstacles = append(obstacles, obstacle)
			}
		}
	}
	return obstacles
}

func drawMoveableObstaclesPerMapTiles(roomID RoomID, modX, modY float64) []entities.MoveableObstacle {
	d := allMapDrawData[rooms[roomID].MapName]
	entities := []entities.MoveableObstacle{}
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+mapX+modX+spriteSize/2,
				vec.Y+mapY+modY+spriteSize/2,
			)
			if _, ok := nonObstacleSprites[spriteData.SpriteID]; !ok {
				entity := buildMoveableObstacle(movedVec.X-spriteSize/2, movedVec.Y-spriteSize/2)
				entities = append(entities, entity)
			}
		}
	}
	return entities
}

func indexRoom(a, b RoomID, dir direction.Name) {
	// fmt.Printf("indexRoom a:%d b:%d dir:%s\n", a, b, dir)
	roomA, okA := rooms[a]
	roomB, okB := rooms[b]
	if okA && okB {
		switch dir {
		case direction.Up:
			// b is above a
			roomA.ConnectedRooms.Top = b
			rooms[a] = roomA
			roomB.ConnectedRooms.Bottom = a
			rooms[b] = roomB
		case direction.Right:
			// b is right of a
			roomA, ok := rooms[a]
			if ok {
				roomA.ConnectedRooms.Right = b
				rooms[a] = roomA
				roomB.ConnectedRooms.Left = a
				rooms[b] = roomB
			}
		case direction.Down:
			// b is below a
			roomA, ok := rooms[a]
			if ok {
				roomA.ConnectedRooms.Bottom = b
				rooms[a] = roomA
				roomB.ConnectedRooms.Top = a
				rooms[b] = roomB
			}
		case direction.Left:
			// b is left of a
			roomA, ok := rooms[a]
			if ok {
				roomA.ConnectedRooms.Left = b
				rooms[a] = roomA
				roomB.ConnectedRooms.Right = a
				rooms[b] = roomB
			}
		}
	}

}

func processMapLayout(layout [][]RoomID) {
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
