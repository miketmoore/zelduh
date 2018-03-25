package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

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
	"github.com/miketmoore/zelduh/physics"
	"github.com/miketmoore/zelduh/playerinput"
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
	spritePlayerPath string  = "assets/bink-spritesheet-01.png"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	coinEntities := buildCoinEntities()
	enemyEntities := buildEnemyEntities()

	obstacles := buildLevelObstacles("fourWalls")

	findEnemy := func(id int) (entities.Enemy, bool) {
		for _, e := range enemyEntities {
			if e.ID == id {
				return e, true
			}
		}
		return entities.Enemy{}, false
	}

	// Create systems and add to game world
	gameWorld.AddSystem(&playerinput.System{Win: win})
	gameWorld.AddSystem(&spatial.System{
		Rand: r,
	})
	gameWorld.AddSystem(&physics.System{})
	gameWorld.AddSystem(&collision.System{
		PlayerCollisionWithCoin: func(coinID int) {
			// Player gets a coin!
			playerEntity.CoinsComponent.Coins++
			fmt.Printf("Player coins: %d\n", playerEntity.CoinsComponent.Coins)
			gameWorld.RemoveCoin(coinID)
		},
		PlayerCollisionWithEnemy: func(enemyID int) {
			// fmt.Printf("Player collided with enemy ID:%d\n", enemyID)
			playerEntity.SpatialComponent.Rect = playerEntity.SpatialComponent.PrevRect
			enemy, ok := findEnemy(enemyID)
			if ok {
				enemy.SpatialComponent.Rect = enemy.SpatialComponent.PrevRect
			}
		},
		PlayerCollisionWithObstacle: func(obstacleID int) {
			// "Block" by undoing rect
			playerEntity.SpatialComponent.Rect = playerEntity.SpatialComponent.PrevRect
		},
		EnemyCollisionWithObstacle: func(enemyID, obstacleID int) {
			// "Block" by undoing rect
			enemy, ok := findEnemy(enemyID)
			if ok {
				enemy.SpatialComponent.Rect = enemy.SpatialComponent.PrevRect
			}
		},
	})
	gameWorld.AddSystem(&render.System{Win: win})

	currentState := gamestate.Start

	// Add entity components to custom ECS systems
	for _, system := range gameWorld.Systems() {
		switch sys := system.(type) {
		case *playerinput.System:
			sys.AddPlayer(playerEntity.PhysicsComponent)
		case *spatial.System:
			sys.AddPlayer(playerEntity.SpatialComponent, playerEntity.MovementComponent)
			for _, enemy := range enemyEntities {
				sys.AddEnemy(enemy.ID, enemy.SpatialComponent, enemy.MovementComponent)
			}
		case *collision.System:
			sys.AddPlayer(playerEntity.SpatialComponent)
			for _, coin := range coinEntities {
				sys.AddCoin(coin.ID, coin.SpatialComponent)
			}
			for _, enemy := range enemyEntities {
				sys.AddEnemy(enemy.ID, enemy.SpatialComponent)
			}
			for _, obstacle := range obstacles {
				sys.AddObstacle(obstacle.ID, obstacle.SpatialComponent)
			}
		case *physics.System:
			sys.AddPlayer(playerEntity.PhysicsComponent, playerEntity.MovementComponent)
		case *render.System:
			sys.AddPlayer(playerEntity.AppearanceComponent, playerEntity.SpatialComponent)
			for _, coin := range coinEntities {
				sys.AddCoin(coin.ID, coin.AppearanceComponent, coin.SpatialComponent)
			}
			for _, enemy := range enemyEntities {
				sys.AddEnemy(enemy.ID, enemy.AppearanceComponent, enemy.SpatialComponent)
			}
			for _, obstacle := range obstacles {
				sys.AddObstacle(obstacle.ID, obstacle.AppearanceComponent, obstacle.SpatialComponent)
			}
		}
	}

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
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Green,
		},
		PhysicsComponent: &components.PhysicsComponent{
			ForceDown:  0,
			ForceLeft:  0,
			ForceRight: 0,
			ForceUp:    0,
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
	}
}

func buildCoinEntities() []entities.Coin {
	coins := []entities.Coin{}
	yInc := spriteSize
	x := mapX
	y := mapY
	totalCoins := 12
	for i := 0; i < totalCoins; i++ {
		coins = append(coins, entities.Coin{
			ID: gameWorld.NewEntityID(),
			AppearanceComponent: &components.AppearanceComponent{
				Color: colornames.Yellow,
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
		})
		x = mapX + float64(r.Intn(totalCoins))*spriteSize
		y += yInc
	}
	return coins
}

func buildEnemyEntities() []entities.Enemy {
	enemies := []entities.Enemy{}

	x := float64(r.Intn(int(mapW-spriteSize))) + mapX
	y := mapY

	for i := 0; i < 5; i++ {
		yInc := float64(i) * spriteSize
		enemies = append(enemies, entities.Enemy{
			ID: gameWorld.NewEntityID(),
			AppearanceComponent: &components.AppearanceComponent{
				Color: colornames.Red,
			},
			PhysicsComponent: &components.PhysicsComponent{
				ForceDown:  0,
				ForceLeft:  0,
				ForceRight: 0,
				ForceUp:    0,
			},
			SpatialComponent: &components.SpatialComponent{
				Width:  spriteSize,
				Height: spriteSize,
				Rect: pixel.R(
					x,
					y+yInc,
					x+spriteSize,
					y+yInc+spriteSize,
				),
				Shape: imdraw.New(nil),
			},
			MovementComponent: &components.MovementComponent{
				Direction: direction.Down,
				Speed:     1.0,
			},
		})
	}

	return enemies
}

func buildLevelObstacles(level string) []entities.Obstacle {
	obstacles := []entities.Obstacle{}

	w := spriteSize
	h := spriteSize

	switch level {
	case "fourWalls":
		for i := 0.0; i < (mapW/w)-2; i++ {
			// Build top wall
			obstacles = append(obstacles, buildObstacle(mapX+w+(w*i), mapY))
			// Build bottom wall
			obstacles = append(obstacles, buildObstacle(mapX+w+(w*i), mapY+mapH-h))
		}

		for i := 0.0; i < (mapH/h)-2; i++ {
			// Build left wall
			obstacles = append(obstacles, buildObstacle(mapX, (mapY+h)+(h*i)))
			// Build right wall
			obstacles = append(obstacles, buildObstacle(mapX+mapW-w, (mapY+h)+(h*i)))

		}
	}

	return obstacles
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
