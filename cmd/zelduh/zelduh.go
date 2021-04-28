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
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
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

	player := entityFactory.NewEntity("player", zelduh.NewCoordinates(6, 6), frameRate)
	bomb := entityFactory.NewEntity("bomb", zelduh.NewCoordinates(0, 0), frameRate)
	explosion := entityFactory.NewEntity("explosion", zelduh.NewCoordinates(0, 0), frameRate)
	sword := entityFactory.NewEntity("sword", zelduh.NewCoordinates(0, 0), frameRate)
	arrow := entityFactory.NewEntity("arrow", zelduh.NewCoordinates(0, 0), frameRate)
	hearts := []zelduh.Entity{
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 1.5, Y: 14}, frameRate),
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 2.15, Y: 14}, frameRate),
		entityFactory.NewEntity("heart", zelduh.Coordinates{X: 2.80, Y: 14}, frameRate),
	}

	windowConfig := zelduh.NewWindowConfig(0, 0, 800, 800)

	activeSpaceRectangle := zelduh.NewActiveSpaceRectangle(0, 0, tileSize*14, tileSize*12)

	activeSpaceRectangle.X = (windowConfig.Width - activeSpaceRectangle.Width) / 2
	activeSpaceRectangle.Y = (windowConfig.Height - activeSpaceRectangle.Height) / 2

	fmt.Println("active space ", activeSpaceRectangle)

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

	ui := zelduh.NewUI(currLocaleMsgs, windowConfig)

	collisionSystem := zelduh.NewCollisionSystem(
		pixel.R(
			activeSpaceRectangle.X,
			activeSpaceRectangle.Y,
			activeSpaceRectangle.X+activeSpaceRectangle.Width,
			activeSpaceRectangle.Y+activeSpaceRectangle.Height,
		),
		&collisionHandler,
		activeSpaceRectangle,
		ui.Window,
	)

	input := Input{window: ui.Window}

	inputSystem := zelduh.NewInputSystem(input)

	systemsManager.AddSystems(
		&inputSystem,
		&healthSystem,
		&spatialSystem,
		&collisionSystem,
		&zelduh.RenderSystem{
			Win:                  ui.Window,
			Spritesheet:          spritesheet,
			ActiveSpaceRectangle: activeSpaceRectangle,
			TileSize:             tileSize,
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

	// totalCells := int(activeSpaceRectangle.Width / tileSize * activeSpaceRectangle.Height / tileSize)
	// fmt.Println(totalCells)
	// debugGridCellCachePopulated := false
	// // debugGridCellCache := []*imdraw.IMDraw{}
	// var debugGridCellCache []*imdraw.IMDraw = make([]*imdraw.IMDraw, totalCells)

	// debugTxtOrigin := pixel.V(20, 50)
	// debugTxt := text.New(debugTxtOrigin, text.Atlas7x13)

	for !ui.Window.Closed() {

		// Quit application when user input matches
		if ui.Window.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		gameStateManager.Update()

		// draw grid after everything else is drawn?
		// drawDebugGrid(
		// 	&debugGridCellCachePopulated,
		// 	debugGridCellCache,
		// 	ui.Window,
		// 	debugTxt,
		// 	activeSpaceRectangle,
		// 	tileSize,
		// )

		// drawDialog(
		// 	systemsManager,
		// 	entityConfigPresetFnManager,
		// 	entityFactory,
		// 	frameRate,
		// 	tileSize,
		// )

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

func buildDebugGridCell(win *pixelgl.Window, rect pixel.Rect, tileSize float64) *imdraw.IMDraw {

	imdraw := imdraw.New(nil)
	imdraw.Color = colornames.Blue

	imdraw.Push(rect.Min)
	imdraw.Push(rect.Max)

	imdraw.Rectangle(1)
	// imdraw.Draw(win)

	return imdraw
}

// Draw and overlay representing the virtual grid
func drawDebugGrid(debugGridCellCachePopulated *bool, cache []*imdraw.IMDraw, win *pixelgl.Window, txt *text.Text, activeSpaceRectangle zelduh.ActiveSpaceRectangle, tileSize float64) {
	// fmt.Println("draw debug grid")
	// win.Clear(colornames.White)

	actualOriginX := activeSpaceRectangle.X
	actualOriginY := activeSpaceRectangle.Y

	totalColumns := activeSpaceRectangle.Width / tileSize
	totalRows := activeSpaceRectangle.Height / tileSize

	cacheIndex := 0

	var x float64 = 0
	var y float64 = 0

	if !(*debugGridCellCachePopulated) {
		fmt.Println("building cache")
		for ; x < totalColumns; x++ {
			cellX := actualOriginX + (x * tileSize)
			cellY := actualOriginY + (y * tileSize)

			rect := pixel.R(cellX, cellY, cellX+tileSize, cellY+tileSize)

			imdraw := buildDebugGridCell(win, rect, tileSize)
			cache[cacheIndex] = imdraw
			cacheIndex++

			if (x == (totalColumns - 1)) && (y < (totalRows - 1)) {
				x = -1
				y++
			}
		}
		*debugGridCellCachePopulated = true
		fmt.Println("cache built")
	} else {
		for _, imdraw := range cache {
			imdraw.Draw(win)
		}
	}

	// fmt.Println("drawing debug grid text, ", totalColumns, totalRows)
	txt.Clear()
	for ; x < totalColumns; x++ {
		// fmt.Println("drawing text...", x)
		message := fmt.Sprintf("%d,%d", int(x), int(y))
		fmt.Fprintln(txt, message)
		matrix := pixel.IM.Moved(
			pixel.V(
				(actualOriginX-18)+(x*tileSize),
				(actualOriginY-48)+(y*tileSize),
			),
		)
		txt.Color = colornames.White
		txt.Draw(win, matrix)
		txt.Clear()

		if (x == (totalColumns - 1)) && (y < (totalRows - 1)) {
			x = -1
			y++
		}
	}

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
		entityConfigPresetFnManager.GetPreset(PresetNameDialogCorner)(zelduh.NewCoordinates(3, 9)),
		// Top side
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(zelduh.NewCoordinates(4, 9)),
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(zelduh.NewCoordinates(5, 9)),
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(zelduh.NewCoordinates(6, 9)),
		// Top right corner
		buildRotatedEntityConfig(PresetNameDialogCorner, entityConfigPresetFnManager, 7, 9, -90),
		// Left Side
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 3, 8, 90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 3, 7, 90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 3, 6, 90),
		// Right Side
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 7, 8, -90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 7, 7, -90),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 7, 6, -90),
		// Bottom left corner
		buildRotatedEntityConfig(PresetNameDialogCorner, entityConfigPresetFnManager, 3, 5, 90),
		// Bottom right corner
		buildRotatedEntityConfig(PresetNameDialogCorner, entityConfigPresetFnManager, 7, 5, 180),
		// Bottom side
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 4, 5, 180),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 5, 5, 180),
		buildRotatedEntityConfig(PresetNameDialogSide, entityConfigPresetFnManager, 6, 5, 180),

		// Center fill
		{
			Category: zelduh.CategoryRectangle,
			Coordinates: zelduh.Coordinates{
				X: 4,
				Y: 6,
			},
			Dimensions: zelduh.Dimensions{
				Width:  3,
				Height: 3,
			},
			Color: colornames.White,
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
