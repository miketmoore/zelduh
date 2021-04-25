package main

import (
	"fmt"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/miketmoore/zelduh"
	"golang.org/x/image/colornames"

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

	draw := true

	for !ui.Window.Closed() {

		// Quit application when user input matches
		if ui.Window.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		gameStateManager.Update()

		if draw {
			draw = false

			drawDialog(
				systemsManager,
				entityConfigPresetFnManager,
				entityFactory,
				frameRate,
				tileSize,
			)

			// systemsManager.AddEntities(entityFactory.NewEntity(
			// 	PresetNameDialogSide,
			// 	zelduh.Coordinates{X: 4, Y: 11},
			// 	frameRate,
			// ))
			// systemsManager.AddEntities(entityFactory.NewEntity(
			// 	PresetNameDialogSide,
			// 	zelduh.Coordinates{X: 5, Y: 11},
			// 	frameRate,
			// ))
			// systemsManager.AddEntities(entityFactory.NewEntity(
			// 	PresetNameDialogSide,
			// 	zelduh.Coordinates{X: 6, Y: 11},
			// 	frameRate,
			// ))

			// entityConfigPresetFn := entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)
			// entityConfig := entityConfigPresetFn(zelduh.Coordinates{X: 7, Y: 11})
			// entityConfig.Transform = &zelduh.Transform{
			// 	Rotation: 180,
			// }

			// systemsManager.AddEntities(entityFactory.NewEntity2(
			// 	entityConfig,
			// 	frameRate,
			// ))
		}

		ui.Window.Update()

	}
}

func main() {
	pixelgl.Run(run)
}

func buildRotatedEntityConfig(
	presetName zelduh.PresetName,
	entityConfigPresetFnManager zelduh.EntityConfigPresetFnManager,
	x, y, degrees float64,
) zelduh.EntityConfig {
	entityConfigPresetFn := entityConfigPresetFnManager.GetPreset(presetName)
	entityConfig := entityConfigPresetFn(zelduh.Coordinates{X: x, Y: y})
	entityConfig.Transform = &zelduh.Transform{
		Rotation: degrees,
	}
	return entityConfig
}

func drawDialog(
	systemsManager zelduh.SystemsManager,
	entityConfigPresetFnManager zelduh.EntityConfigPresetFnManager,
	entityFactory zelduh.EntityFactory,
	frameRate int,
	tileSize float64,
) {

	entityConfigs := []zelduh.EntityConfig{
		// Top left corner
		entityConfigPresetFnManager.GetPreset(PresetNameDialogCorner)(zelduh.Coordinates{X: 3, Y: 11}),
		// Top side
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(zelduh.Coordinates{X: 4, Y: 11}),
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(zelduh.Coordinates{X: 5, Y: 11}),
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(zelduh.Coordinates{X: 6, Y: 11}),
		// Top right corner
		buildRotatedEntityConfig(PresetNameDialogCorner, entityConfigPresetFnManager, 7, 11, -90),
		// Left Side
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 3, 10, 90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 3, 9, 90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 3, 8, 90),
		// Right Side
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 7, 10, -90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 7, 9, -90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 7, 8, -90),
		// Bottom left corner
		buildRotatedEntityConfig(PresetNameDialogCorner, entityConfigPresetFnManager, 3, 7, 90),
		// Bottom right corner
		buildRotatedEntityConfig(PresetNameDialogCorner, entityConfigPresetFnManager, 7, 7, 180),
		// Bottom side
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 4, 7, 180),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 5, 7, 180),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 6, 7, 180),

		// Center fill
		{
			Category: zelduh.CategoryRectangle,
			Coordinates: zelduh.Coordinates{
				X: 3,
				Y: 10,
			},
			Dimensions: zelduh.Dimensions{
				Width:  tileSize * 3,
				Height: tileSize * 3,
			},
			Color: colornames.Blue,
		},
	}

	// center fill
	// circle := imdraw.New(nil)
	// circle.Color = colornames.Red
	// circle.Push(0)
	// circle.Circle(64, 0)

	// rect := imdraw.New(nil)
	// rect.Color = colornames.White

	for _, entityConfig := range entityConfigs {
		systemsManager.AddEntity(entityFactory.NewEntity2(entityConfig, frameRate))
	}

}
