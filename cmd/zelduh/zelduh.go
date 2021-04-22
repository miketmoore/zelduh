package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {

	currLocaleMsgs, err := zelduh.GetLocaleMessageMapByLanguage("en")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// TileSize defines the width and height of a tile
	const tileSize float64 = 48

	// FrameRate is used to determine which sprite to use for animations
	const frameRate int = 5

	entityConfigPresetFnsMap := BuildEntityConfigPresetFnsMap(tileSize)

	entityConfigPresetFnManager := zelduh.NewEntityConfigPresetFnManager(entityConfigPresetFnsMap)

	testLevel := buildTestLevel(
		&entityConfigPresetFnManager,
		tileSize,
	)

	levelManager := zelduh.NewLevelManager(&testLevel)

	systemsManager := zelduh.NewSystemsManager()

	entityFactory := zelduh.NewEntityFactory(&systemsManager, &entityConfigPresetFnManager)

	spatialSystem := zelduh.SpatialSystem{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	healthSystem := zelduh.NewHealthSystem()

	entitiesMap := zelduh.NewEntitiesMap()

	roomTransition := zelduh.NewRoomTransition(float64(tileSize))

	roomWarps := zelduh.NewRoomWarps()

	shouldAddEntities := true
	var currentRoomID zelduh.RoomID = 1
	var nextRoomID zelduh.RoomID
	currentState := zelduh.StateStart
	spritesheet := zelduh.LoadAndBuildSpritesheet("assets/spritesheet.png", tileSize)

	player := entityFactory.NewEntity("player", zelduh.Coordinates{X: 6, Y: 6}, frameRate)
	bomb := entityFactory.NewEntity("bomb", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	explosion := entityFactory.NewEntity("explosion", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	sword := entityFactory.NewEntity("sword", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	arrow := entityFactory.NewEntity("arrow", zelduh.Coordinates{X: 0, Y: 0}, frameRate)
	hearts := []zelduh.Entity{
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 1.5, Y: 14}, frameRate),
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 2.15, Y: 14}, frameRate),
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 2.80, Y: 14}, frameRate),
	}

	windowConfig := zelduh.NewWindowConfig(0, 0, 800, 800)

	activeSpaceRectangle := zelduh.NewActiveSpaceRectangle(0, 0, tileSize*14, tileSize*12)

	activeSpaceRectangle.X = (windowConfig.Width - activeSpaceRectangle.Width) / 2
	activeSpaceRectangle.Y = (windowConfig.Height - activeSpaceRectangle.Height) / 2

	collisionHandler := zelduh.NewCollisionHandler(
		&systemsManager,
		&spatialSystem,
		&healthSystem,
		&shouldAddEntities,
		&nextRoomID,
		&currentState,
		&roomTransition,
		entitiesMap,
		&player,
		&sword,
		&explosion,
		&arrow,
		hearts,
		roomWarps,
		&entityConfigPresetFnManager,
		tileSize,
		frameRate,
	)

	collisionSystem := zelduh.NewCollisionSystem(
		pixel.R(
			activeSpaceRectangle.X,
			activeSpaceRectangle.Y,
			activeSpaceRectangle.X+activeSpaceRectangle.Width,
			activeSpaceRectangle.Y+activeSpaceRectangle.Height,
		),
		&collisionHandler,
	)

	ui := zelduh.NewUI(currLocaleMsgs, windowConfig)

	input := Input{window: ui.Window}

	inputSystem := zelduh.NewInputSystem(input)

	systemsManager.AddSystems(
		&inputSystem,
		&healthSystem,
		&spatialSystem,
		&collisionSystem,
		&zelduh.RenderSystem{
			Win:         ui.Window,
			Spritesheet: spritesheet,
		},
	)

	systemsManager.AddEntities(
		player,
		sword,
		arrow,
		bomb,
	)

	mapDrawData := zelduh.BuildMapDrawData(
		"assets/tilemaps/",
		[]string{
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
		},
		tileSize,
	)

	// NonObstacleSprites defines which sprites are not obstacles
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

	gameStateManager := zelduh.NewGameStateManager(
		&systemsManager,
		ui,
		currLocaleMsgs,
		&collisionSystem,
		&inputSystem,
		&shouldAddEntities,
		&currentRoomID,
		&nextRoomID,
		&currentState,
		spritesheet,
		mapDrawData,
		&roomTransition,
		entitiesMap,
		&player,
		hearts,
		roomWarps,
		&levelManager,
		&entityConfigPresetFnManager,
		tileSize,
		frameRate,
		nonObstacleSprites,
		windowConfig,
		activeSpaceRectangle,
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
