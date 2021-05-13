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

	// entityConfigPresetFnsMap := BuildEntityConfigPresetFnsMap(m.tileSize)

	roomFactory := NewRoomFactory()

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
	spritesheetPicture, spriteMap := LoadAndBuildSpritesheet("assets/spritesheet.png", m.tileSize)

	temporarySystem := NewTemporarySystem()

	// "nes_": "overworldFourWallsDoorLeft",
	// "ne_w": "overworldFourWallsDoorBottom",
	// "n_sw": "overworldFourWallsDoorRight",
	// "_esw": "overworldFourWallsDoorTop",
	// "ne__": "overworldFourWallsDoorBottomLeft",
	// "n__w": "overworldFourWallsDoorBottomRight",
	// "_es_": "overworldFourWallsDoorLeftTop",
	// "__sw": "overworldFourWallsDoorRightTop",
	// "_e__": "overworldFourWallsDoorTopBottomLeft",
	// "n_s_": "overworldFourWallsDoorRightLeft",
	// "n___": "overworldFourWallsDoorRightBottomLeft",
	// "_e_w": "overworldFourWallsDoorTopBottom",
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
			"overworldFourWallsDoorTopBottomLeft",
			"overworldFourWallsDoorRightLeft",
			"overworldFourWallsDoorRightBottomLeft",
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
		// for debug
		// tree
		// 77: true,
	}

	activeSpaceRectangle := NewActiveSpaceRectangle(0, 0, m.tileSize*14, m.tileSize*12)

	levelManager := &LevelManager{}

	entityFactory := NewEntityFactory(
		&systemsManager,
		&temporarySystem,
		&movementSystem,
		m.tileSize,
		m.frameRate,
		levelManager,
		mapDrawData,
		activeSpaceRectangle,
		nonObstacleSprites,
	)

	// testLevel := buildTestLevel(
	// 	roomFactory,
	// 	&entityFactory,
	// 	m.tileSize,
	// )
	testLevel := buildLevelMaze(roomFactory, &entityFactory, m.tileSize)
	levelManager.SetCurrentLevel(&testLevel)
	// levelManager = NewLevelManager(&testLevel)

	player := entityFactory.NewEntityFromConfig(entityFactory.PresetPlayer()(NewCoordinates(6, 6)), m.frameRate)
	playerID := player.ID()

	bomb := entityFactory.NewEntityFromConfig(entityFactory.PresetBomb()(NewCoordinates(0, 0)), m.frameRate)
	// explosion := entityFactory.NewEntityFromPresetName("explosion", NewCoordinates(0, 0), m.frameRate)

	sword := entityFactory.NewEntityFromConfig(entityFactory.PresetSword()(NewCoordinates(0, 0)), m.frameRate)
	swordID := sword.ID()

	arrow := entityFactory.NewEntityFromConfig(entityFactory.PresetArrow()(NewCoordinates(0, 0)), m.frameRate)
	arrowID := arrow.ID()

	windowConfig := NewWindowConfig(0, 0, 800, 800)

	activeSpaceRectangle.X = (windowConfig.Width - activeSpaceRectangle.Width) / 2
	activeSpaceRectangle.Y = (windowConfig.Height - activeSpaceRectangle.Height) / 2

	ui := NewUISystem(
		m.currLocaleMsgs,
		windowConfig,
		activeSpaceRectangle,
		spriteMap,
		mapDrawData,
		m.tileSize,
		m.frameRate,
		levelManager,
		nonObstacleSprites,
		&entityFactory,
	)

	renderSystem := NewRenderSystem(
		ui.Window,
		spriteMap,
		activeSpaceRectangle,
		m.tileSize,
		&temporarySystem,
		spritesheetPicture,
		levelManager,
		roomManager,
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
				// fmt.Printf("PlayerWithObstacle collision handler obstacleID=%d\n", obstacleID)
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
		levelManager,
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
			fmt.Printf("collision with bound=%s currentRoomID=%d nextRoomID=%d\n", side, roomManager.Current(), roomManager.Next())
			// TODO prevent room transition if no room exists on this side

			// // fmt.Printf("bound collision currentRoomID=%d nextRoomID=%d\n", roomManager.Current(), roomManager.Next())
			// if roomManager.Next() == 0 {
			// 	// player has collided with a wall where no adjacent room exists
			// 	// want to prevent the room transition from starting
			// 	// and want to prevent player from moving out of the active space
			// 	movementSystem.MovePlayerBack()
			// 	return
			// }

			if !roomTransitionManager.Active() {
				// fmt.Printf("room transition start %d %d\n", roomManager.Current(), roomManager.Next())
				// // roomTransitionManager.Start(side, RoomTransitionSlide)
				// roomTransitionManager.Enable()
				roomTransitionManager.SetSide(side)
				// roomTransitionManager.SetSlide()
				// roomTransitionManager.ResetTimer()
				// // currentState = StateMapTransition

				// shouldAddEntities = true

				err := stateContext.SetState(StateNamePrepareMapTransition)
				if err != nil {
					fmt.Println(err)
					os.Exit(0)
				}

				// err := stateContext.SetState("transition")
				// if err != nil {
				// 	fmt.Println("Error: ", err)
				// 	os.Exit(0)
				// }
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

func drawDialog(
	systemsManager SystemsManager,
	entityFactory *EntityFactory,
	frameRate int,
	tileSize float64,
) {

	entityConfigs := []EntityConfig{
		// Top left corner
		entityFactory.PresetDialogCorner(0)(NewCoordinates(3, 9)),
		// Top side
		entityFactory.PresetDialogSide(0)(NewCoordinates(4, 9)),
		entityFactory.PresetDialogSide(0)(NewCoordinates(5, 9)),
		entityFactory.PresetDialogSide(0)(NewCoordinates(6, 9)),
		// Top right corner
		entityFactory.PresetDialogCorner(-90)(NewCoordinates(7, 9)),
		// Left Side
		entityFactory.PresetDialogSide(90)(NewCoordinates(3, 8)),
		entityFactory.PresetDialogSide(90)(NewCoordinates(3, 7)),
		entityFactory.PresetDialogSide(90)(NewCoordinates(3, 6)),
		// Right Side
		entityFactory.PresetDialogSide(-90)(NewCoordinates(7, 8)),
		entityFactory.PresetDialogSide(-90)(NewCoordinates(7, 7)),
		entityFactory.PresetDialogSide(-90)(NewCoordinates(7, 6)),
		// Bottom left corner
		entityFactory.PresetDialogCorner(90)(NewCoordinates(3, 5)),
		// Bottom right corner
		entityFactory.PresetDialogCorner(-180)(NewCoordinates(7, 5)),
		// Bottom side
		entityFactory.PresetDialogSide(180)(NewCoordinates(4, 5)),
		entityFactory.PresetDialogSide(180)(NewCoordinates(5, 5)),
		entityFactory.PresetDialogSide(180)(NewCoordinates(6, 5)),

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
