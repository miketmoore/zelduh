package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/deanobob/tmxreader"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/go-pixel-game-template/state"
	"github.com/miketmoore/zelduh/collision"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"github.com/miketmoore/zelduh/entities"
	"github.com/miketmoore/zelduh/gamestate"
	"github.com/miketmoore/zelduh/input"
	"github.com/miketmoore/zelduh/render"
	"github.com/miketmoore/zelduh/spatial"
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
	spriteSize       float64 = 48
	spritePlayerPath string  = "assets/spritesheet.png"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var tilemapDir = "assets/tilemaps/"
var tilemapFiles = []string{
	"overworldOpen",
	"overworldOpenCircleOfTrees",
	"overworldFourWallsDoorBottom",
	"overworldFourWallsDoorTop",
}

var tmxMapData map[string]tmxreader.TmxMap

var spriteMap = map[string]float64{
	"playerDownA":         361,
	"playerDownB":         376,
	"playerUpA":           362,
	"playerUpB":           377,
	"playerRightA":        363,
	"playerRightB":        378,
	"playerLeftA":         364,
	"playerLeftB":         379,
	"turtleNoShellDownA":  316,
	"turtleNoShellDownB":  331,
	"turtleNoShellUpA":    316,
	"turtleNoShellUpB":    331,
	"turtleNoShellRightA": 317,
	"turtleNoShellRightB": 332,
	"turtleNoShellLeftA":  317,
	"turtleNoShellLeftB":  332,
	"sword":               84,
	"ground":              8,
	"coinA":               365,
	"coinB":               380,
	"coinC":               395,
}

var sprites map[string]*pixel.Sprite

func run() {
	// Initializations
	t = initI18n()
	txt = initText(20, 50, colornames.Black)
	win = initWindow(t("title"), winX, winY, winW, winH)
	pic = loadPicture(spritePlayerPath)
	sprites = buildSpriteMap(pic, spriteMap)

	gameWorld = world.New()

	// New entities
	playerEntity := buildPlayerEntity()
	sword := entities.Sword{
		Ignore: &components.Ignore{
			Value: true,
		},
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Deeppink,
		},
		SpatialComponent: &components.SpatialComponent{
			Width:  spriteSize,
			Height: spriteSize,
			Rect:   pixel.R(0, 0, 0, 0),
			Shape:  imdraw.New(nil),
		},
		MovementComponent: &components.MovementComponent{
			Direction: playerEntity.MovementComponent.Direction,
			Speed:     0.0,
		},
	}
	arrow := entities.Arrow{
		Ignore: &components.Ignore{
			Value: true,
		},
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Deeppink,
		},
		SpatialComponent: &components.SpatialComponent{
			Width:  spriteSize,
			Height: spriteSize,
			Rect:   pixel.R(0, 0, 0, 0),
			Shape:  imdraw.New(nil),
		},
		MovementComponent: &components.MovementComponent{
			Direction: direction.Down,
			Speed:     0.0,
		},
	}
	enemyEntities := buildEnemyEntities()

	loadTmxData()
	obstacles := buildLevelObstacles("fourWalls")

	moveableObstacles := []entities.MoveableObstacle{
		buildMoveableObstacle(mapX+(spriteSize*5), mapY+(spriteSize*5)),
	}

	collisionSwitches := buildCollisionSwitches()

	findEnemy := func(id int) (entities.Enemy, bool) {
		for _, e := range enemyEntities {
			if e.ID == id {
				return e, true
			}
		}
		return entities.Enemy{}, false
	}

	findMoveableObstacle := func(id int) (entities.MoveableObstacle, bool) {
		for _, e := range moveableObstacles {
			if e.ID == id {
				return e, true
			}
		}
		return entities.MoveableObstacle{}, false
	}

	findCollisionSwitchIndex := func(id int) int {
		for i, e := range collisionSwitches {
			if e.ID == id {
				return i
			}
		}
		return -1
	}

	currentState := gamestate.Start

	// Create systems and add to game world
	gameWorld.AddSystem(&input.System{Win: win})
	spatialSystem := &spatial.System{
		Rand: r,
	}
	gameWorld.AddSystem(spatialSystem)
	gameWorld.AddSystem(&collision.System{
		PlayerCollisionWithCoin: func(coinID int) {
			playerEntity.CoinsComponent.Coins++
			fmt.Printf("Player coins: %d\n", playerEntity.CoinsComponent.Coins)
			gameWorld.RemoveCoin(coinID)
		},
		PlayerCollisionWithEnemy: func(enemyID int) {
			spatialSystem.MovePlayerBack()
			playerEntity.Health.Total--
			if playerEntity.Health.Total == 0 {
				currentState = gamestate.Over
			}
		},
		SwordCollisionWithEnemy: func(enemyID int) {
			if !sword.Ignore.Value {
				enemy, ok := findEnemy(enemyID)
				if ok {
					enemy.Health.Total--
					if enemy.Health.Total == 0 {
						gameWorld.RemoveEnemy(enemy.ID)
						r := enemy.SpatialComponent.Rect
						coin := buildCoin(r.Min.X, r.Min.Y)
						addCoinToSystem(coin)
					} else {
						spatialSystem.MoveEnemyBack(enemyID, playerEntity.MovementComponent.Direction)
					}
				}
			}
		},
		ArrowCollisionWithEnemy: func(enemyID int) {
			if !arrow.Ignore.Value {
				enemy, ok := findEnemy(enemyID)
				if ok {
					enemy.Health.Total--
					if enemy.Health.Total == 0 {
						gameWorld.RemoveEnemy(enemy.ID)
						r := enemy.SpatialComponent.Rect
						coin := buildCoin(r.Min.X, r.Min.Y)
						addCoinToSystem(coin)
					} else {
						spatialSystem.MoveEnemyBack(enemyID, playerEntity.MovementComponent.Direction)
					}
				}
			}
		},
		ArrowCollisionWithObstacle: func() {
			arrow.MovementComponent.MoveCount = 0
		},
		PlayerCollisionWithObstacle: func(obstacleID int) {
			// "Block" by undoing rect
			playerEntity.SpatialComponent.Rect = playerEntity.SpatialComponent.PrevRect
			sword.SpatialComponent.Rect = sword.SpatialComponent.PrevRect
		},
		PlayerCollisionWithMoveableObstacle: func(obstacleID int) {
			spatialSystem.MoveMoveableObstacle(obstacleID, playerEntity.MovementComponent.Direction)
		},
		EnemyCollisionWithObstacle: func(enemyID, obstacleID int) {
			// "Block" by undoing rect
			enemy, ok := findEnemy(enemyID)
			if ok {
				enemy.SpatialComponent.Rect = enemy.SpatialComponent.PrevRect
			}
		},
		EnemyCollisionWithMoveableObstacle: func(enemyID int) {
			// "Block" by undoing rect
			enemy, ok := findEnemy(enemyID)
			if ok {
				enemy.SpatialComponent.Rect = enemy.SpatialComponent.PrevRect
			}
		},
		MoveableObstacleCollisionWithObstacle: func(moveableObstacleID int) {
			obstacle, ok := findMoveableObstacle(moveableObstacleID)
			if ok {
				obstacle.SpatialComponent.Rect = obstacle.SpatialComponent.PrevRect
			}
		},
		PlayerCollisionWithSwitch: func(switchID int) {
			i := findCollisionSwitchIndex(switchID)
			if i != -1 {
				s := &collisionSwitches[i]
				if !s.Enabled {
					fmt.Printf("Enabled switch %d!\n", switchID)
					s.Enabled = true
				}
			}
		},
		PlayerNoCollisionWithSwitch: func(switchID int) {
			i := findCollisionSwitchIndex(switchID)
			if i != -1 {
				s := &collisionSwitches[i]
				if s.Enabled {
					fmt.Printf("Disabled switch %d!\n", switchID)
					s.Enabled = false
				}
			}
		},
	})
	gameWorld.AddSystem(&render.System{Win: win})

	addPlayerToSystems(playerEntity)
	addSwordToSystems(sword)
	addArrowToSystems(arrow)
	addEnemiesToSystem(enemyEntities)
	addObstaclesToSystem(obstacles)
	addMoveableObstaclesToSystem(moveableObstacles)
	addCollisionSwitchesToSystem(collisionSwitches)

	for !win.Closed() {

		allowQuit()

		switch currentState {
		case gamestate.Start:
			win.Clear(colornames.Darkgray)
			txt.Clear()
			drawMapBG(mapX, mapY, mapW, mapH, colornames.White)
			txt.Color = colornames.Black
			fmt.Fprintln(txt, t("title"))
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = gamestate.Game
			}
		case gamestate.Game:

			win.Clear(colornames.Darkgray)
			drawMapBG(mapX, mapY, mapW, mapH, colornames.White)

			gameWorld.Update()

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Pause
			}

			if win.JustPressed(pixelgl.KeyX) {
				currentState = gamestate.Over
			}

		case gamestate.Pause:
			win.Clear(colornames.Darkgray)
			txt.Clear()
			fmt.Fprintln(txt, t("paused"))
			drawMapBG(mapX, mapY, mapW, mapH, colornames.White)
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Game
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				currentState = gamestate.Start
			}
		case gamestate.Over:
			win.Clear(colornames.Darkgray)
			txt.Clear()
			drawMapBG(mapX, mapY, mapW, mapH, colornames.Black)
			txt.Color = colornames.White
			fmt.Fprintln(txt, t("gameOver"))
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = gamestate.Start
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

func buildSpriteMap(pic pixel.Picture, config map[string]float64) map[string]*pixel.Sprite {
	spriteMap := map[string]*pixel.Sprite{}
	for k, v := range config {
		spriteMap[k] = newSpriteIndexed(pic, v)
	}
	return spriteMap
}

func newSpriteIndexed(pic pixel.Picture, index float64) *pixel.Sprite {
	totalRows := pic.Bounds().H() / spriteSize
	totalCols := pic.Bounds().W() / spriteSize

	find := func() (float64, float64) {
		i := 0.0
		var row = 0.0
		var col = 0.0
		for ; row < totalRows; row++ {
			if i == index {
				return row, col
			}
			for col = 0.0; col < totalCols; col++ {
				i++
				if i == index {
					return row, col
				}
			}
		}
		return row, col
	}

	row, col := find()

	return pixel.NewSprite(pic, pixel.Rect{
		Min: pixel.V(col*spriteSize, row*spriteSize),
		Max: pixel.V(col*spriteSize+spriteSize, row*spriteSize+spriteSize),
	})
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

func buildPlayerEntity() entities.Player {
	return entities.Player{
		Health: &components.Health{
			Total: 3,
		},
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Green,
		},
		SpatialComponent: &components.SpatialComponent{
			Width:  spriteSize,
			Height: spriteSize,
			Rect: pixel.R(
				mapX+(mapW/2),
				mapY+(mapH/2),
				mapX+(mapW/2)+spriteSize,
				mapY+(mapH/2)+spriteSize,
			),
			Shape: imdraw.New(nil),
		},
		MovementComponent: &components.MovementComponent{
			Direction: direction.Down,
			Speed:     4.0,
		},
		CoinsComponent: &components.CoinsComponent{
			Coins: 0,
		},
		Dash: &components.Dash{
			Charge:    0,
			MaxCharge: 50,
			SpeedMod:  7,
		},
	}
}

func buildCoinEntities() []entities.Coin {
	w := spriteSize
	h := spriteSize
	return []entities.Coin{
		buildCoin(mapX+w, mapY+h),
		buildCoin(mapX+w*10, mapY+h*7),
	}
}

func buildCoin(x, y float64) entities.Coin {
	w := spriteSize
	h := spriteSize
	return entities.Coin{
		ID: gameWorld.NewEntityID(),
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Purple,
		},
		SpatialComponent: &components.SpatialComponent{
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

func buildEnemyEntities() []entities.Enemy {
	w := spriteSize
	h := spriteSize
	return []entities.Enemy{
		buildEnemy(mapX+(w*2), mapY+(h*3)),
		buildEnemy(mapX+(w*10), mapY+(h*3)),
	}
}

func buildEnemy(x, y float64) entities.Enemy {
	return entities.Enemy{
		ID:     gameWorld.NewEntityID(),
		Health: &components.Health{Total: 1},
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Red,
		},
		SpatialComponent: &components.SpatialComponent{
			Width:  spriteSize,
			Height: spriteSize,
			Rect: pixel.R(
				x,
				y,
				x+spriteSize,
				y+spriteSize,
			),
			Shape: imdraw.New(nil),
		},
		MovementComponent: &components.MovementComponent{
			Direction: direction.Down,
			Speed:     1.0,
		},
	}
}

func buildObstacle(x, y float64) entities.Obstacle {
	return entities.Obstacle{
		ID: gameWorld.NewEntityID(),
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Black,
		},
		SpatialComponent: &components.SpatialComponent{
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
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Purple,
		},
		SpatialComponent: &components.SpatialComponent{
			Width:  spriteSize,
			Height: spriteSize,
			Rect:   pixel.R(x, y, x+spriteSize, y+spriteSize),
			Shape:  imdraw.New(nil),
		},
		MovementComponent: &components.MovementComponent{
			Direction: direction.Down,
			Speed:     1.0,
		},
	}
}

func buildCollisionSwitches() []entities.CollisionSwitch {
	switches := []entities.CollisionSwitch{}
	x := mapX + (spriteSize * 3)
	y := mapY + (spriteSize * 3)
	switches = append(switches, buildCollisionSwitch(x, y))
	return switches
}

func buildCollisionSwitch(x, y float64) entities.CollisionSwitch {
	return entities.CollisionSwitch{
		ID:      gameWorld.NewEntityID(),
		Enabled: false,
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Sandybrown,
		},
		SpatialComponent: &components.SpatialComponent{
			Width:  spriteSize,
			Height: spriteSize,
			Rect:   pixel.R(x, y, x+spriteSize, y+spriteSize),
			Shape:  imdraw.New(nil),
		},
	}
}

func addCoinToSystem(coin entities.Coin) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *collision.System:
			sys.AddCoin(coin.ID, coin.SpatialComponent)
		case *render.System:
			sys.AddCoin(coin.ID, coin.AppearanceComponent, coin.SpatialComponent)
		}
	}
}

func addPlayerToSystems(playerEntity entities.Player) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *input.System:
			sys.AddPlayer(playerEntity.MovementComponent, playerEntity.Dash)
		case *spatial.System:
			sys.AddPlayer(playerEntity.SpatialComponent, playerEntity.MovementComponent, playerEntity.Dash)
		case *collision.System:
			sys.AddPlayer(playerEntity.SpatialComponent)
		case *render.System:
			sys.AddPlayer(playerEntity.AppearanceComponent, playerEntity.SpatialComponent)
		}
	}
}

func addSwordToSystems(sword entities.Sword) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *input.System:
			sys.AddSword(sword.MovementComponent, sword.Ignore)
		case *spatial.System:
			sys.AddSword(sword.SpatialComponent, sword.MovementComponent)
		case *collision.System:
			sys.AddSword(sword.SpatialComponent)
		case *render.System:
			sys.AddSword(sword.AppearanceComponent, sword.SpatialComponent)
		}
	}
}

func addArrowToSystems(arrow entities.Arrow) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *input.System:
			sys.AddArrow(arrow.MovementComponent, arrow.Ignore)
		case *spatial.System:
			sys.AddArrow(arrow.SpatialComponent, arrow.MovementComponent)
		case *collision.System:
			sys.AddArrow(arrow.SpatialComponent)
		case *render.System:
			sys.AddArrow(arrow.AppearanceComponent, arrow.SpatialComponent)
		}
	}
}

func addEnemiesToSystem(enemyEntities []entities.Enemy) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *spatial.System:
			for _, enemy := range enemyEntities {
				sys.AddEnemy(enemy.ID, enemy.SpatialComponent, enemy.MovementComponent)
			}
		case *collision.System:
			for _, enemy := range enemyEntities {
				sys.AddEnemy(enemy.ID, enemy.SpatialComponent)
			}
		case *render.System:
			for _, enemy := range enemyEntities {
				sys.AddEnemy(enemy.ID, enemy.AppearanceComponent, enemy.SpatialComponent)
			}
		}
	}
}

func addObstaclesToSystem(obstacles []entities.Obstacle) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *collision.System:
			for _, obstacle := range obstacles {
				sys.AddObstacle(obstacle.ID, obstacle.SpatialComponent)
			}
		case *render.System:
			for _, obstacle := range obstacles {
				sys.AddObstacle(obstacle.ID, obstacle.AppearanceComponent, obstacle.SpatialComponent)
			}
		}
	}
}

func addMoveableObstaclesToSystem(moveableObstacles []entities.MoveableObstacle) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *spatial.System:
			for _, moveable := range moveableObstacles {
				sys.AddMoveableObstacle(moveable.ID, moveable.SpatialComponent, moveable.MovementComponent)
			}
		case *collision.System:
			for _, moveable := range moveableObstacles {
				sys.AddMoveableObstacle(moveable.ID, moveable.SpatialComponent)
			}
		case *render.System:
			for _, moveable := range moveableObstacles {
				sys.AddMoveableObstacle(moveable.ID, moveable.AppearanceComponent, moveable.SpatialComponent)
			}
		}
	}
}

func addCollisionSwitchesToSystem(collisionSwitches []entities.CollisionSwitch) {
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *collision.System:
			for _, collisionSwitch := range collisionSwitches {
				sys.AddCollisionSwitch(collisionSwitch.ID, collisionSwitch.SpatialComponent)
			}
		case *render.System:
			for _, collisionSwitch := range collisionSwitches {
				sys.AddCollisionSwitch(collisionSwitch.AppearanceComponent, collisionSwitch.SpatialComponent)
			}
		}
	}
}

func loadTmxData() {
	tmxMapData = map[string]tmxreader.TmxMap{}
	for i, name := range tilemapFiles {
		path := fmt.Sprintf("%s%s.tmx", tilemapDir, name)
		fmt.Printf("Loading TMX tile map %d/%d %s: %s\n", i+1, len(tilemapFiles), name, path)
		tmxMapData[name] = parseTmxFile(path)
	}
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

func buildLevelObstacles(level string) []entities.Obstacle {

	var layout [][]int

	switch level {
	case "fourWalls":
		layout = [][]int{
			[]int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		}
	case "fourWallsDoorBottom":
		layout = [][]int{
			[]int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0},
		}
	case "fourWallsDoorTop":
		layout = [][]int{
			[]int{0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			[]int{0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0},
		}
	}

	allCoords := buildCoordsSliceFromLayout(layout)
	obstacles := []entities.Obstacle{}
	for _, coords := range allCoords {
		obstacles = append(obstacles, buildObstacle(coords[0], coords[1]))
	}
	return obstacles
}

func buildMapFromLayout(layout [][]int) []entities.Obstacle {
	obstacles := []entities.Obstacle{}
	y := mapY
	s := ""
	for r := len(layout) - 1; r >= 0; r-- {

		row := layout[r]

		for i, c := range row {
			x := mapX + (float64(i) * spriteSize)
			o := buildObstacle(x, y)
			s = fmt.Sprintf("%s %d", s, c)
			if c == 1 {
				obstacles = append(obstacles, o)
			}

		}

		s = fmt.Sprintf("%s\n", s)
		y += spriteSize
	}
	fmt.Println(s)
	return obstacles
}

func buildCoordsSliceFromLayout(layout [][]int) [][]float64 {
	coords := [][]float64{}
	y := mapY
	s := ""
	for r := len(layout) - 1; r >= 0; r-- {

		row := layout[r]

		for i, c := range row {
			x := mapX + (float64(i) * spriteSize)
			s = fmt.Sprintf("%s %d", s, c)
			if c == 1 {
				coords = append(coords, []float64{x, y})
			}

		}

		s = fmt.Sprintf("%s\n", s)
		y += spriteSize
	}
	fmt.Println(s)
	return coords
}
