package zelduh

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/core/direction"
	"github.com/miketmoore/zelduh/core/entity"
	"golang.org/x/image/colornames"
)

type Main struct {
	debugMode      bool
	currLocaleMsgs LocaleMessagesMap
	tileSize       float64
	frameRate      int
}

func NewMain(
	debugMode bool,
	tileSize float64,
	frameRate int,
) (*Main, error) {

	currLocaleMsgs, err := GetLocaleMessageMapByLanguage("en")
	if err != nil {
		return nil, err
	}

	return &Main{
		debugMode:      debugMode,
		currLocaleMsgs: currLocaleMsgs,
		tileSize:       tileSize,
		frameRate:      frameRate,
	}, nil
}

func (m *Main) Run() error {

	entityConfigPresetFnsMap := BuildEntityConfigPresetFnsMap(m.tileSize)

	entityConfigPresetFnManager := NewEntityConfigPresetFnManager(entityConfigPresetFnsMap)

	roomFactory := NewRoomFactory()

	testLevel := buildTestLevel(
		roomFactory,
		&entityConfigPresetFnManager,
		m.tileSize,
	)

	levelManager := NewLevelManager(&testLevel)

	systemsManager := NewSystemsManager()

	movementSystem := NewMovementSystem(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		m.tileSize,
	)

	healthSystem := NewHealthSystem()

	entitiesMap := NewEntitiesMap()

	roomTransitionManager := NewRoomTransitionManager(m.tileSize)

	roomWarps := NewRoomWarps()

	shouldAddEntities := true

	roomManager := NewRoomManager(1)

	// currentState := StateStart
	spriteMap := LoadAndBuildSpritesheet("assets/spritesheet.png", m.tileSize)

	temporarySystem := NewTemporarySystem()

	entityFactory := NewEntityFactory(
		&systemsManager,
		&entityConfigPresetFnManager,
		&temporarySystem,
		&movementSystem,
		m.tileSize,
		m.frameRate,
	)

	player := entityFactory.NewEntityFromPresetName("player", NewCoordinates(6, 6), m.frameRate)
	playerID := player.ID()

	bomb := entityFactory.NewEntityFromPresetName("bomb", NewCoordinates(0, 0), m.frameRate)
	// explosion := entityFactory.NewEntityFromPresetName("explosion", NewCoordinates(0, 0), m.frameRate)

	sword := entityFactory.NewEntityFromPresetName("sword", NewCoordinates(0, 0), m.frameRate)
	swordID := sword.ID()

	arrow := entityFactory.NewEntityFromPresetName("arrow", NewCoordinates(0, 0), m.frameRate)
	arrowID := arrow.ID()

	windowConfig := NewWindowConfig(0, 0, 800, 800)

	activeSpaceRectangle := NewActiveSpaceRectangle(0, 0, m.tileSize*14, m.tileSize*12)

	activeSpaceRectangle.X = (windowConfig.Width - activeSpaceRectangle.Width) / 2
	activeSpaceRectangle.Y = (windowConfig.Height - activeSpaceRectangle.Height) / 2

	mapDrawData, mapDrawDataErr := BuildMapDrawData(
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
			"overworldFourWallsDoorBottomLeft",
			"rockPatternTest",
			"rockPathOpenLeft",
			"rockWithCaveEntrance",
			"rockPathLeftRightEntrance",
			"test",
			"dungeonFourDoors",
		},
		m.tileSize,
	)

	if mapDrawDataErr != nil {
		fmt.Println(mapDrawData)
		os.Exit(0)
	}

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

	ui := NewUISystem(
		m.currLocaleMsgs,
		windowConfig,
		activeSpaceRectangle,
		spriteMap,
		mapDrawData,
		m.tileSize,
		m.frameRate,
		&entityConfigPresetFnManager,
		&levelManager,
		nonObstacleSprites,
		&entityFactory,
	)

	renderSystem := NewRenderSystem(
		ui.Window,
		spriteMap,
		activeSpaceRectangle,
		m.tileSize,
		&temporarySystem,
	)

	mapBounds := pixel.R(
		activeSpaceRectangle.X,
		activeSpaceRectangle.Y,
		activeSpaceRectangle.X+activeSpaceRectangle.Width,
		activeSpaceRectangle.Y+activeSpaceRectangle.Height,
	)

	var stateContext *StateContext

	ignoreSystem := NewIgnoreSystem()
	coinsSystem := NewCoinsSystem()
	toggleSystem := NewToggleSystem()

	collisionSystem := NewCollisionSystem(
		mapBounds,
		activeSpaceRectangle,
		ui.Window,
		&OnCollisionHandlers{
			PlayerWithEnemy: func(entityID entity.EntityID) {
				// TODO repeat what I did with the enemies
				movementSystem.MovePlayerBack()
				healthSystem.Hit(playerID, 1)

				if healthSystem.Health(playerID) == 0 {
					// currentState = StateOver
					err := stateContext.SetState("gameOver")
					if err != nil {
						fmt.Println("Error: ", err)
						os.Exit(0)
					}
				}
			},
			PlayerWithCoin: func(coinEntityID entity.EntityID) {
				coinsSystem.AddCoins(player.ID(), 1)
				systemsManager.Remove(CategoryCoin, coinEntityID)
			},
			SwordWithEnemy: func(enemyEntityID entity.EntityID) {
				if !ignoreSystem.IsCurrentlyIgnored(sword.ID()) {
					// if !sword.componentIgnore.Value {
					dead := false
					if !movementSystem.EnemyMovingFromHit(enemyEntityID) {
						dead = healthSystem.Hit(enemyEntityID, 1)
						if dead {
							entityFactory.CreateExplosion(enemyEntityID)
							systemsManager.RemoveEnemy(enemyEntityID)
						} else {
							playerDirection, err := movementSystem.Direction(player.ID())
							if err != nil {
								fmt.Println("Error: ", err)
								os.Exit(0)
							}
							movementSystem.MoveEnemyBack(enemyEntityID, playerDirection)
						}
					}

				}
			},
			ArrowWithEnemy: func(enemyEntityID entity.EntityID) {
				if !ignoreSystem.IsCurrentlyIgnored(arrow.ID()) {
					dead := healthSystem.Hit(enemyEntityID, 1)
					// arrow.componentIgnore.Value = true
					ignoreSystem.Ignore(arrow.ID())
					if dead {
						entityFactory.CreateExplosion(enemyEntityID)
						systemsManager.RemoveEnemy(enemyEntityID)
					} else {
						playerDirection, err := movementSystem.Direction(player.ID())
						if err != nil {
							fmt.Println("Error: ", err)
							os.Exit(0)
						}
						movementSystem.MoveEnemyBack(enemyEntityID, playerDirection)
					}
				}
			},
			ArrowWithObstacle: func(arrowID entity.EntityID) {
				movementSystem.SetRemainingMoves(arrowID, 0)
			},
			PlayerWithObstacle: func(obstacleID entity.EntityID) {
				// "Block" by undoing rect
				movementSystem.UsePreviousRectangle(player.ID())
				movementSystem.SetZeroSpeed(player.ID())

				movementSystem.UsePreviousRectangle(sword.ID())
			},
			PlayerWithMoveableObstacle: func(moveableObstacleID entity.EntityID) {
				playerDirection, err := movementSystem.Direction(player.ID())
				if err != nil {
					fmt.Println("Error: ", err)
					os.Exit(0)
				}
				moved := movementSystem.MoveMoveableObstacle(moveableObstacleID, playerDirection)
				if !moved {
					movementSystem.UsePreviousRectangle(player.ID())
				}
			},
			MoveableObstacleWithSwitch: func(collisionSwitchID entity.EntityID) {
				if !toggleSystem.Enabled(collisionSwitchID) {
					toggleSystem.Toggle(collisionSwitchID)
				}
			},
			MoveableObstacleWithSwitchNoCollision: func(collisionSwitchID entity.EntityID) {
				if toggleSystem.Enabled(collisionSwitchID) {
					toggleSystem.Toggle(collisionSwitchID)
				}
			},
			PlayerWithSwitch: func(collisionSwitchID entity.EntityID) {
				if !toggleSystem.Enabled(collisionSwitchID) {
					toggleSystem.Toggle(collisionSwitchID)
				}
			},
			PlayerWithSwitchNoCollision: func(collisionSwitchID entity.EntityID) {
				if toggleSystem.Enabled(collisionSwitchID) {
					toggleSystem.Toggle(collisionSwitchID)
				}
			},
			EnemyWithObstacle: func(enemyEntityID entity.EntityID) {
				// Block enemy within the spatial system by reseting current rect to previous rect
				movementSystem.UndoEnemyRect(enemyEntityID)
			},
			PlayerWithWarp: func(warpID entity.EntityID) {
				entityConfig, ok := roomWarps[warpID]
				if ok && !roomTransitionManager.Active() {
					roomTransitionManager.Enable()
					roomTransitionManager.SetWarp()
					roomTransitionManager.SetTimer(1)
					// currentState = StateMapTransition
					err := stateContext.SetState("transition")
					if err != nil {
						fmt.Println("Error: ", err)
						os.Exit(0)
					}
					shouldAddEntities = true
					roomManager.SetNext(entityConfig.WarpToRoomID)
				}
			},
		},
	)

	input := InputImpl{window: ui.Window}

	inputHandlers := InputHandlers{
		OnUp: func() {
			movementSystem.SetMaxSpeed(playerID)
			movementSystem.ChangeDirection(playerID, direction.DirectionUp)
		},
		OnRight: func() {
			movementSystem.SetMaxSpeed(playerID)
			movementSystem.ChangeDirection(playerID, direction.DirectionRight)
		},
		OnDown: func() {
			movementSystem.SetMaxSpeed(playerID)
			movementSystem.ChangeDirection(playerID, direction.DirectionDown)
		},
		OnLeft: func() {
			movementSystem.SetMaxSpeed(playerID)
			movementSystem.ChangeDirection(playerID, direction.DirectionLeft)
		},
		OnNoDirection: func() {
			movementSystem.SetZeroSpeed(playerID)
		},
		OnPrimaryAttack: func() {
			movementSystem.MatchDirectionToPlayer(swordID)
			movementSystem.ChangeSpeed(swordID, 1.0)
			ignoreSystem.DoNotIgnore(swordID)
		},
		OnNoPrimaryAttack: func() {
			movementSystem.MatchDirectionToPlayer(swordID)
			movementSystem.SetZeroSpeed(swordID)
			ignoreSystem.Ignore(swordID)
		},
		OnSecondaryAttack: func() {
			if movementSystem.RemainingMoves(arrow.ID()) == 0 {
				movementSystem.MatchDirectionToPlayer(arrowID)
				movementSystem.ChangeSpeed(arrowID, 7.0)
				movementSystem.SetRemainingMoves(arrowID, 100)
				ignoreSystem.DoNotIgnore(arrowID)
			} else {
				movementSystem.DecrementRemainingMoves(arrowID)
			}
		},
		OnNoSecondaryAttack: func() {
			if movementSystem.RemainingMoves(arrow.ID()) == 0 {
				movementSystem.MatchDirectionToPlayer(arrowID)
				movementSystem.SetZeroSpeed(arrowID)
				movementSystem.SetRemainingMoves(arrowID, 0)
				ignoreSystem.Ignore(arrowID)
			} else {
				movementSystem.DecrementRemainingMoves(arrowID)
			}
		},
	}

	inputSystem := NewInputSystem(
		input,
		inputHandlers,
	)

	uiHeartSystem := NewUIHeartSystem(
		&systemsManager,
		entityFactory,
		m.frameRate,
	)

	stateContext = NewStateContext(
		&ui,
		&inputSystem,
		&roomTransitionManager,
		&collisionSystem,
		&systemsManager,
		&levelManager,
		&entityFactory,
		roomWarps,
		entitiesMap,
		roomManager,
		&shouldAddEntities,
		m.tileSize,
		m.frameRate,
		activeSpaceRectangle,
		player,
	)

	boundsCollisionSystem := NewBoundsCollisionSystem(
		ui.Window,
		mapBounds,
		func(side Bound) {
			// TODO prevent room transition if no room exists on this side
			if !roomTransitionManager.Active() {
				fmt.Printf("room transition start %d %d\n", roomManager.Current(), roomManager.Next())
				// roomTransitionManager.Start(side, RoomTransitionSlide)
				roomTransitionManager.Enable()
				roomTransitionManager.SetSide(side)
				roomTransitionManager.SetSlide()
				roomTransitionManager.ResetTimer()
				// currentState = StateMapTransition
				err := stateContext.SetState("transition")
				if err != nil {
					fmt.Println("Error: ", err)
					os.Exit(0)
				}
				shouldAddEntities = true
			}
		},
	)

	systemsManager.AddSystems(
		&inputSystem,
		&healthSystem,
		&movementSystem,
		&collisionSystem,
		&renderSystem,
		&movementSystem,
		&ignoreSystem,
		&temporarySystem,
		&boundsCollisionSystem,
		&coinsSystem,
		&uiHeartSystem,
	)

	systemsManager.AddEntities(
		player,
		sword,
		arrow,
		bomb,
	)

	debug := NewDebug(
		activeSpaceRectangle,
		m.tileSize,
		ui.Window,
		m.debugMode,
	)

	for !ui.Window.Closed() {

		// Quit application when user input matches
		if ui.Window.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		// err := gameStateManager.Update()
		err := stateContext.Update()
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(0)
		}

		movementSystem.UpdateLastDirection(playerID)

		// draw grid after everything else is drawn?
		debug.DrawGrid()

		// drawDialog(
		// 	systemsManager,
		// 	entityConfigPresetFnManager,
		// 	entityFactory,
		// 	m.frameRate,
		// 	m.tileSize,
		// )

		ui.Window.Update()

	}

	return nil
}

func buildRotatedEntityConfig(
	presetName PresetName,
	entityConfigPresetFnManager EntityConfigPresetFnManager,
	x, y, degrees float64,
) EntityConfig {
	entityConfigPresetFn := entityConfigPresetFnManager.GetPreset(presetName)
	entityConfig := entityConfigPresetFn(Coordinates{X: x, Y: y})
	entityConfig.Transform = &Transform{
		Rotation: degrees,
	}
	return entityConfig
}

func drawDialog(
	systemsManager SystemsManager,
	entityConfigPresetFnManager EntityConfigPresetFnManager,
	entityFactory EntityFactory,
	frameRate int,
	tileSize float64,
) {

	entityConfigs := []EntityConfig{
		// Top left corner
		entityConfigPresetFnManager.GetPreset(PresetNameDialogCorner)(NewCoordinates(3, 9)),
		// Top side
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(NewCoordinates(4, 9)),
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(NewCoordinates(5, 9)),
		entityConfigPresetFnManager.GetPreset(PresetNameDialogSide)(NewCoordinates(6, 9)),
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
			Category: CategoryRectangle,
			Coordinates: Coordinates{
				X: 4,
				Y: 6,
			},
			Dimensions: Dimensions{
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
		systemsManager.AddEntity(entityFactory.NewEntityFromConfig(entityConfig, frameRate))
	}

}
