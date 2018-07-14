package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miketmoore/zelduh/bounds"
	"github.com/miketmoore/zelduh/csv"
	"github.com/miketmoore/zelduh/gamemap"
	"github.com/miketmoore/zelduh/rooms"
	"github.com/miketmoore/zelduh/sprites"

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
	"github.com/miketmoore/zelduh/tmx"
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
var spriteMap map[string]*pixel.Sprite

type mapDrawData struct {
	Rect     pixel.Rect
	SpriteID int
}

var allMapDrawData map[string]MapData

const frameRate int = 5

type enemyPresetFn = func(xTiles, yTiles float64) rooms.EntityConfig

var spriteSets = map[string][]int{
	"eyeburrower": []int{50, 50, 50, 91, 91, 91, 92, 92, 92, 93, 93, 93, 92, 92, 92},
	"explosion": []int{
		122, 122, 122,
		123, 123, 123,
		124, 124, 124,
		125, 125, 125,
	},
	"uiCoin":           []int{20},
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
	"toggleObstacle":   []int{144, 114},
	"swordUp":          []int{70},
	"swordRight":       []int{67},
	"swordDown":        []int{68},
	"swordLeft":        []int{69},
	"arrowUp":          []int{101},
	"arrowRight":       []int{100},
	"arrowDown":        []int{103},
	"arrowLeft":        []int{102},
	"bomb":             []int{138, 139, 140, 141},
	"coin":             []int{5, 5, 6, 6, 21, 21},
	"heart":            []int{106},
}

var entityPresets = map[string]enemyPresetFn{
	"arrow": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Arrow,
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				Speed:     0.0,
			},
			W: s,
			H: s,
			X: s * xTiles,
			Y: s * yTiles,
			Animation: rooms.AnimationConfig{
				"up":    spriteSets["arrowUp"],
				"right": spriteSets["arrowRight"],
				"down":  spriteSets["arrowDown"],
				"left":  spriteSets["arrowLeft"],
			},
			Hitbox: &rooms.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"bomb": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Bomb,
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				Speed:     0.0,
			},
			W: s,
			H: s,
			X: s * xTiles,
			Y: s * yTiles,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["bomb"],
			},
			Hitbox: &rooms.HitboxConfig{
				Radius: 5,
			},
			Ignore: true,
		}
	},
	"coin": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Coin,
			W:        s,
			H:        s,
			X:        s * xTiles,
			Y:        s * yTiles,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["coin"],
			},
		}
	},
	"explosion": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category:   categories.Explosion,
			Expiration: 12,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["explosion"],
			},
		}
	},
	"obstacle": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Obstacle,
			W:        s,
			H:        s,
			X:        s * xTiles,
			Y:        s * yTiles,
		}
	},
	"player": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Player,
			Health:   3,
			W:        s,
			H:        s,
			X:        s * xTiles,
			Y:        s * yTiles,
			Hitbox: &rooms.HitboxConfig{
				Box:                  imdraw.New(nil),
				Radius:               15,
				CollisionWithRectMod: 5,
			},
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				MaxSpeed:  7.0,
				Speed:     0.0,
			},
			Coins: true,
			Dash: &rooms.DashConfig{
				Charge:    0,
				MaxCharge: 50,
				SpeedMod:  7,
			},
			Animation: rooms.AnimationConfig{
				"up":               spriteSets["playerUp"],
				"right":            spriteSets["playerRight"],
				"down":             spriteSets["playerDown"],
				"left":             spriteSets["playerLeft"],
				"swordAttackUp":    spriteSets["playerSwordUp"],
				"swordAttackRight": spriteSets["playerSwordRight"],
				"swordAttackLeft":  spriteSets["playerSwordLeft"],
				"swordAttackDown":  spriteSets["playerSwordDown"],
			},
		}
	},
	"sword": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Sword,
			Movement: &rooms.MovementConfig{
				Direction: direction.Down,
				Speed:     0.0,
			},
			W: s,
			H: s,
			X: s * xTiles,
			Y: s * yTiles,
			Animation: rooms.AnimationConfig{
				"up":    spriteSets["swordUp"],
				"right": spriteSets["swordRight"],
				"down":  spriteSets["swordDown"],
				"left":  spriteSets["swordLeft"],
			},
			Hitbox: &rooms.HitboxConfig{
				Radius: 20,
			},
			Ignore: true,
		}
	},
	"eyeburrower": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["eyeburrower"],
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Down,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "random",
			},
		}
	},
	"heart": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Heart,
			W:        s,
			H:        s,
			X:        s * xTiles,
			Y:        s * yTiles,
			Hitbox: &rooms.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: rooms.AnimationConfig{
				"default": spriteSets["heart"],
			},
		}

	},
	"skeleton": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["skeleton"],
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Down,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "random",
			},
		}
	},
	"skull": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["skull"],
			},
			Health: 2,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Down,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "random",
			},
		}
	},
	"spinner": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Enemy,
			W:        s, H: s, X: s * xTiles, Y: s * yTiles,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["spinner"],
			},
			Invincible: true,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Movement: &rooms.MovementConfig{
				Direction:    direction.Right,
				Speed:        1.0,
				MaxSpeed:     1.0,
				HitSpeed:     10.0,
				HitBackMoves: 10,
				MaxMoves:     100,
				PatternName:  "left-right",
			},
		}
	},
	"uiCoin": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Heart,
			W:        s,
			H:        s,
			X:        s * xTiles,
			Y:        s * yTiles,
			Hitbox: &rooms.HitboxConfig{
				Box: imdraw.New(nil),
			},
			Animation: rooms.AnimationConfig{
				"default": spriteSets["uiCoin"],
			},
		}
	},
	"warpStone": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.Warp,
			X:        s * xTiles,
			Y:        s * yTiles,
			W:        s,
			H:        s,
			Hitbox: &rooms.HitboxConfig{
				Box:    imdraw.New(nil),
				Radius: 20,
			},
			Animation: rooms.AnimationConfig{
				"default": spriteSets["warpStone"],
			},
		}
	},
	"puzzleBox": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.MovableObstacle,
			X:        s * xTiles,
			Y:        s * yTiles,
			W:        s,
			H:        s,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["puzzleBox"],
			},
			Movement: &rooms.MovementConfig{
				Speed:    1.0,
				MaxMoves: int(s) / 2,
				MaxSpeed: 2.0,
			},
		}
	},
	"floorSwitch": func(xTiles, yTiles float64) rooms.EntityConfig {
		return rooms.EntityConfig{
			Category: categories.CollisionSwitch,
			X:        s * xTiles,
			Y:        s * yTiles,
			W:        s,
			H:        s,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["floorSwitch"],
			},
			Toggleable: true,
		}
	},
	// this is an impassable obstacle that can be toggled "remotely"
	// it has two visual states that coincide with each toggle state
	"toggleObstacle": func(xTiles, yTiles float64) rooms.EntityConfig {
		// TODO get this working again
		return rooms.EntityConfig{
			X: s * xTiles,
			Y: s * yTiles,
			W: s,
			H: s,
			Animation: rooms.AnimationConfig{
				"default": spriteSets["toggleObstacle"],
			},
			// Impassable: true,
			Toggleable: true,
		}
	},
}

func presetWarpStone(X, Y, WarpToRoomID, HitBoxRadius float64) rooms.EntityConfig {
	fmt.Printf("presetWarpStone\n")
	e := entityPresets["warpStone"](X, Y)
	e.WarpToRoomID = 6
	e.Hitbox.Radius = 5
	return e
}

var roomsMap rooms.Rooms
var entitiesMap map[entities.EntityID]entities.Entity

func run() {

	entitiesMap = map[entities.EntityID]entities.Entity{}
	gameWorld = world.New()

	roomsMap = rooms.Rooms{
		1: rooms.NewRoom("overworldFourWallsDoorBottomRight",
			entityPresets["puzzleBox"](5, 5),
			(func() rooms.EntityConfig {
				e := entityPresets["floorSwitch"](5, 6)
				return e
			})(),
			entityPresets["toggleObstacle"](10, 7),
		),
		2: rooms.NewRoom("overworldFourWallsDoorTopBottom",
			entityPresets["skull"](5, 5),
			entityPresets["skeleton"](11, 9),
			entityPresets["spinner"](7, 9),
			entityPresets["eyeburrower"](8, 9),
		),
		3: rooms.NewRoom("overworldFourWallsDoorRightTopBottom",
			presetWarpStone(3, 7, 6, 5),
		),
		5: rooms.NewRoom("rockWithCaveEntrance",
			rooms.EntityConfig{
				Category:     categories.Warp,
				WarpToRoomID: 11,
				W:            s,
				H:            s,
				X:            (s * 7) + s/2,
				Y:            (s * 9) + s/2,
				Hitbox: &rooms.HitboxConfig{
					Radius: 30,
				},
			},
			rooms.EntityConfig{
				Category:     categories.Warp,
				WarpToRoomID: 11,
				W:            s,
				H:            s,
				X:            (s * 8) + s/2,
				Y:            (s * 9) + s/2,
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
				W:            s,
				H:            s,
				X:            (s * 6) + s + (s / 2.5),
				Y:            (s * 1) + s + (s / 2.5),
				Hitbox: &rooms.HitboxConfig{
					Radius: 15,
				},
			},
			rooms.EntityConfig{
				Category:     categories.Warp,
				WarpToRoomID: 5,
				W:            s,
				H:            s,
				X:            (s * 7) + s + (s / 2.5),
				Y:            (s * 1) + s + (s / 2.5),
				Hitbox: &rooms.HitboxConfig{
					Radius: 15,
				},
			},
		),
	}

	gamemap.ProcessMapLayout(roomsMap, overworld)

	// Initializations
	t = initI18n()
	txt = initText(20, 50, colornames.Black)
	win = initWindow(t("title"), winX, winY, winW, winH)

	// load the spritesheet image
	pic = loadPicture(spritesheetPath)
	// build spritesheet
	// this is a map of TMX IDs to sprite instances
	spritesheet = sprites.BuildSpritesheet(pic, s)

	// load all TMX file data for each map
	tmxMapData = tmx.Load(tilemapFiles, tilemapDir)
	allMapDrawData = buildMapDrawData()

	// Build entities
	player := entities.BuildEntityFromConfig(frameRate, entityPresets["player"](6, 6), gameWorld.NewEntityID())
	bomb := entities.BuildEntityFromConfig(frameRate, entityPresets["bomb"](0, 0), gameWorld.NewEntityID())
	explosion := entities.BuildEntityFromConfig(frameRate, entityPresets["explosion"](0, 0), gameWorld.NewEntityID())
	sword := entities.BuildEntityFromConfig(frameRate, entityPresets["sword"](0, 0), gameWorld.NewEntityID())
	arrow := entities.BuildEntityFromConfig(frameRate, entityPresets["arrow"](0, 0), gameWorld.NewEntityID())

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
	gameWorld.AddSystem(inputSystem)
	healthSystem := &systems.Health{}
	gameWorld.AddSystem(healthSystem)
	spatialSystem := &systems.Spatial{
		Rand: r,
	}
	dropCoin := func(v pixel.Vec) {
		coin := entities.BuildEntityFromConfig(frameRate, entityPresets["coin"](v.X/s, v.Y/s), gameWorld.NewEntityID())
		addEntityToSystem(coin)
	}
	gameWorld.AddSystem(spatialSystem)

	hearts := []entities.Entity{
		entities.BuildEntityFromConfig(frameRate, entityPresets["heart"](1.5, 14), gameWorld.NewEntityID()),
		entities.BuildEntityFromConfig(frameRate, entityPresets["heart"](2.15, 14), gameWorld.NewEntityID()),
		entities.BuildEntityFromConfig(frameRate, entityPresets["heart"](2.80, 14), gameWorld.NewEntityID()),
	}

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
			gameWorld.Remove(categories.Coin, coinID)
		},
		PlayerCollisionWithEnemy: func(enemyID entities.EntityID) {
			// TODO repeat what I did with the enemies
			spatialSystem.MovePlayerBack()
			player.Health.Total--

			// remove heart entity
			heartIndex := len(hearts) - 1
			gameWorld.Remove(categories.Heart, hearts[heartIndex].ID)
			hearts = append(hearts[:heartIndex], hearts[heartIndex+1:]...)

			// TODO redraw hearts
			if player.Health.Total == 0 {
				currentState = gamestate.Over
			}
		},
		SwordCollisionWithEnemy: func(enemyID entities.EntityID) {
			fmt.Printf("SwordCollisionWithEnemy %d\n", enemyID)
			if !sword.Ignore.Value {
				dead := false
				if !spatialSystem.EnemyMovingFromHit(enemyID) {
					dead = healthSystem.Hit(enemyID, 1)
					if dead {
						enemySpatial, _ := spatialSystem.GetEnemySpatial(enemyID)
						explosion.Temporary.Expiration = len(explosion.Animation.Map["default"].Frames)
						explosion.Spatial = &components.Spatial{
							Width:  s,
							Height: s,
							Rect:   enemySpatial.Rect,
						}
						explosion.Temporary.OnExpiration = func() {
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
					explosion.Temporary.Expiration = len(explosion.Animation.Map["default"].Frames)
					explosion.Spatial = &components.Spatial{
						Width:  s,
						Height: s,
						Rect:   enemySpatial.Rect,
					}
					explosion.Temporary.OnExpiration = func() {
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
		MoveableObstacleCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range entitiesMap {
				if id == collisionSwitchID && !entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		MoveableObstacleNoCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range entitiesMap {
				if id == collisionSwitchID && entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		EnemyCollisionWithObstacle: func(enemyID, obstacleID entities.EntityID) {
			// Block enemy within the spatial system by reseting current rect to previous rect
			spatialSystem.UndoEnemyRect(enemyID)
		},
		PlayerCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range entitiesMap {
				if id == collisionSwitchID && !entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		PlayerNoCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range entitiesMap {
				if id == collisionSwitchID && entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		PlayerCollisionWithWarp: func(warpID entities.EntityID) {
			entityConfig, ok := roomWarps[warpID]
			if ok && !roomTransition.Active {
				roomTransition.Active = true
				roomTransition.Style = rooms.TransitionWarp
				roomTransition.Timer = 1
				currentState = gamestate.MapTransition
				addEntities = true
				nextRoomID = entityConfig.WarpToRoomID
			}
		},
	}
	gameWorld.AddSystem(collisionSystem)
	gameWorld.AddSystem(&systems.Render{
		Win:         win,
		Spritesheet: spritesheet,
	})

	// make sure only correct number of hearts exists in systems
	// so, if health is reduced, need to remove a heart entity from the systems,
	// the correct one... last one
	addHearts := func(health int) {
		for i, entity := range hearts {
			if i < health {
				addEntityToSystem(entity)
			}
		}
	}

	addUICoin := func() {
		coin := entities.BuildEntityFromConfig(frameRate, entityPresets["uiCoin"](4, 14), gameWorld.NewEntityID())
		addEntityToSystem(coin)
	}

	addEntityToSystem(player)
	addEntityToSystem(sword)
	addEntityToSystem(arrow)
	addEntityToSystem(bomb)

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

			addHearts(player.Health.Total)

			if addEntities {
				addEntities = false

				addUICoin()

				// Draw obstacles on appropriate map tiles
				obstacles := drawObstaclesPerMapTiles(currentRoomID, 0, 0)
				for _, entity := range obstacles {
					addEntityToSystem(entity)
				}

				roomWarps = map[entities.EntityID]rooms.EntityConfig{}

				// Iterate through all entity configurations and build entities and add to systems
				for _, c := range roomsMap[currentRoomID].EntityConfigs {
					entity := entities.BuildEntityFromConfig(frameRate, c, gameWorld.NewEntityID())
					entitiesMap[entity.ID] = entity
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
				removeAllEntitiesFromSystems()

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
				removeAllEntitiesFromSystems()
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

func addEntityToSystem(entity entities.Entity) {
	for _, system := range gameWorld.Systems() {
		system.AddEntity(entity)
	}
}

func removeAllEntitiesFromSystems() {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Render:
			sys.RemoveAllEntities()
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

func removeAllCollisionSwitchesFromSystems() {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
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

			records := csv.Parse(strings.TrimSpace(layer.Data.Value) + ",")
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
	mod := 0.5
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+mapX+modX+s/2,
				vec.Y+mapY+modY+s/2,
			)

			if _, ok := nonObstacleSprites[spriteData.SpriteID]; !ok {
				x := movedVec.X/s - mod
				y := movedVec.Y/s - mod
				id := gameWorld.NewEntityID()
				obstacle := entities.BuildEntityFromConfig(frameRate, entityPresets["obstacle"](x, y), id)
				obstacles = append(obstacles, obstacle)
			}
		}
	}
	return obstacles
}
