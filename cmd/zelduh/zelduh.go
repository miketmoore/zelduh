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
	"github.com/miketmoore/zelduh/config"
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

var (
	win       *pixelgl.Window
	txt       *text.Text
	t         i18n.TranslateFunc
	currState state.State
	pic       pixel.Picture
	gameWorld world.World
)

const (
	spritesheetPath string = "assets/spritesheet.png"
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

var entitiesMap map[entities.EntityID]entities.Entity

func run() {

	entitiesMap = map[entities.EntityID]entities.Entity{}
	gameWorld = world.New()

	gamemap.ProcessMapLayout(roomsMap, overworld)

	// Initializations
	t = initI18n()
	txt = initText(20, 50, colornames.Black)
	win = initWindow(t("title"))

	// load the spritesheet image
	pic = loadPicture(spritesheetPath)
	// build spritesheet
	// this is a map of TMX IDs to sprite instances
	spritesheet = sprites.BuildSpritesheet(pic, config.TileSize)

	// load all TMX file data for each map
	tmxMapData = tmx.Load(tilemapFiles, tilemapDir)
	allMapDrawData = buildMapDrawData()

	// Build entities
	player := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("player")(6, 6), gameWorld.NewEntityID())
	bomb := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("bomb")(0, 0), gameWorld.NewEntityID())
	explosion := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("explosion")(0, 0), gameWorld.NewEntityID())
	sword := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("sword")(0, 0), gameWorld.NewEntityID())
	arrow := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("arrow")(0, 0), gameWorld.NewEntityID())

	roomTransition := rooms.RoomTransition{
		Start: float64(config.TileSize),
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
		coin := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("coin")(v.X/config.TileSize, v.Y/config.TileSize), gameWorld.NewEntityID())
		gameWorld.AddEntityToSystem(coin)
	}
	gameWorld.AddSystem(spatialSystem)

	hearts := []entities.Entity{
		entities.BuildEntityFromConfig(frameRate, entities.GetPreset("heart")(1.5, 14), gameWorld.NewEntityID()),
		entities.BuildEntityFromConfig(frameRate, entities.GetPreset("heart")(2.15, 14), gameWorld.NewEntityID()),
		entities.BuildEntityFromConfig(frameRate, entities.GetPreset("heart")(2.80, 14), gameWorld.NewEntityID()),
	}

	collisionSystem := &systems.Collision{
		MapBounds: pixel.R(
			config.MapX,
			config.MapY,
			config.MapX+config.MapW,
			config.MapY+config.MapH,
		),
		OnPlayerCollisionWithBounds: func(side bounds.Bound) {
			if !roomTransition.Active {
				roomTransition.Active = true
				roomTransition.Side = side
				roomTransition.Style = rooms.TransitionSlide
				roomTransition.Timer = int(roomTransition.Start)
				currentState = gamestate.MapTransition
				addEntities = true
			}
		},
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
				currentState = gamestate.Over
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
			for id, entity := range entitiesMap {
				if id == collisionSwitchID && !entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		OnMoveableObstacleNoCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range entitiesMap {
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
			for id, entity := range entitiesMap {
				if id == collisionSwitchID && !entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		OnPlayerNoCollisionWithSwitch: func(collisionSwitchID entities.EntityID) {
			for id, entity := range entitiesMap {
				if id == collisionSwitchID && entity.Toggler.Enabled() {
					entity.Toggler.Toggle()
				}
			}
		},
		OnPlayerCollisionWithWarp: func(warpID entities.EntityID) {
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
				gameWorld.AddEntityToSystem(entity)
			}
		}
	}

	addUICoin := func() {
		coin := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("uiCoin")(4, 14), gameWorld.NewEntityID())
		gameWorld.AddEntityToSystem(coin)
	}

	gameWorld.AddEntityToSystem(player)
	gameWorld.AddEntityToSystem(sword)
	gameWorld.AddEntityToSystem(arrow)
	gameWorld.AddEntityToSystem(bomb)

	for !win.Closed() {

		allowQuit()

		switch currentState {
		case gamestate.Start:
			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)
			drawCenterText(t("title"), colornames.Black)

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = gamestate.Game
			}
		case gamestate.Game:
			inputSystem.EnablePlayer()

			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

			drawMapBGImage(roomsMap[currentRoomID].MapName, 0, 0)

			addHearts(player.Health.Total)

			if addEntities {
				addEntities = false

				addUICoin()

				// Draw obstacles on appropriate map tiles
				obstacles := drawObstaclesPerMapTiles(currentRoomID, 0, 0)
				for _, entity := range obstacles {
					gameWorld.AddEntityToSystem(entity)
				}

				roomWarps = map[entities.EntityID]rooms.EntityConfig{}

				// Iterate through all entity configurations and build entities and add to systems
				for _, c := range roomsMap[currentRoomID].EntityConfigs {
					entity := entities.BuildEntityFromConfig(frameRate, c, gameWorld.NewEntityID())
					entitiesMap[entity.ID] = entity
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
				currentState = gamestate.Pause
			}

			if win.JustPressed(pixelgl.KeyX) {
				currentState = gamestate.Over
			}

		case gamestate.Pause:
			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)
			drawCenterText(t("paused"), colornames.Black)

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Game
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				currentState = gamestate.Start
			}
		case gamestate.Over:
			win.Clear(colornames.Darkgray)
			drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.Black)
			drawCenterText(t("gameOver"), colornames.White)

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = gamestate.Start
			}
		case gamestate.MapTransition:
			inputSystem.DisablePlayer()
			if roomTransition.Style == rooms.TransitionSlide && roomTransition.Timer > 0 {
				roomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				gameWorld.RemoveAllEnemies()
				gameWorld.RemoveAllCollisionSwitches()
				gameWorld.RemoveAllMoveableObstacles()
				gameWorld.RemoveAllEntities()

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
				if roomTransition.Side == bounds.Bottom && roomsMap[currentRoomID].ConnectedRooms.Bottom != 0 {
					modY = incY
					modYNext = incY - config.MapH
					nextRoomID = roomsMap[currentRoomID].ConnectedRooms.Bottom

					playerModY += playerIncY
				} else if roomTransition.Side == bounds.Top && roomsMap[currentRoomID].ConnectedRooms.Top != 0 {
					modY = -incY
					modYNext = -incY + config.MapH
					nextRoomID = roomsMap[currentRoomID].ConnectedRooms.Top
					playerModY -= playerIncY
				} else if roomTransition.Side == bounds.Left && roomsMap[currentRoomID].ConnectedRooms.Left != 0 {
					modX = incX
					modXNext = incX - config.MapW
					nextRoomID = roomsMap[currentRoomID].ConnectedRooms.Left
					playerModX += playerIncX
				} else if roomTransition.Side == bounds.Right && roomsMap[currentRoomID].ConnectedRooms.Right != 0 {
					modX = -incX
					modXNext = -incX + config.MapW
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
					player.Spatial.Rect.Min.X+playerModX+config.TileSize,
					player.Spatial.Rect.Min.Y+playerModY+config.TileSize,
				)

				gameWorld.Update()
			} else if roomTransition.Style == rooms.TransitionWarp && roomTransition.Timer > 0 {
				roomTransition.Timer--
				win.Clear(colornames.Darkgray)
				drawMapBG(config.MapX, config.MapY, config.MapW, config.MapH, colornames.White)

				collisionSystem.RemoveAll(categories.Obstacle)
				gameWorld.RemoveAllEnemies()
				gameWorld.RemoveAllCollisionSwitches()
				gameWorld.RemoveAllMoveableObstacles()
				gameWorld.RemoveAllEntities()
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

// MapData represents data for one map
type MapData struct {
	Name string
	Data []mapDrawData
}

func buildMapDrawData() map[string]MapData {
	all := map[string]MapData{}

	for mapName, mapData := range tmxMapData {
		// fmt.Printf("Building map draw data for map %config.TileSize\n", mapName)

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
						y := float64(11-row) * config.TileSize
						x := float64(col) * config.TileSize

						record := records[row][col]
						spriteID, err := strconv.Atoi(record)
						if err != nil {
							panic(err)
						}
						mrd := mapDrawData{
							Rect:     pixel.R(x, y, x+config.TileSize, y+config.TileSize),
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
				vec.X+config.MapX+modX+config.TileSize/2,
				vec.Y+config.MapY+modY+config.TileSize/2,
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
				vec.X+config.MapX+modX+config.TileSize/2,
				vec.Y+config.MapY+modY+config.TileSize/2,
			)

			if _, ok := nonObstacleSprites[spriteData.SpriteID]; !ok {
				x := movedVec.X/config.TileSize - mod
				y := movedVec.Y/config.TileSize - mod
				id := gameWorld.NewEntityID()
				obstacle := entities.BuildEntityFromConfig(frameRate, entities.GetPreset("obstacle")(x, y), id)
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
