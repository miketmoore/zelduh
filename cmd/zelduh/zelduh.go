package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"engo.io/ecs"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/zelduh/collision"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"github.com/miketmoore/zelduh/enemy"
	"github.com/miketmoore/zelduh/entities"
	"github.com/miketmoore/zelduh/entity"
	"github.com/miketmoore/zelduh/equipment"
	"github.com/miketmoore/zelduh/gamestate"
	"github.com/miketmoore/zelduh/message"
	"github.com/miketmoore/zelduh/mvmt"
	"github.com/miketmoore/zelduh/palette"
	"github.com/miketmoore/zelduh/player"
	"github.com/miketmoore/zelduh/systems"
	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/image/colornames"
)

const winW float64 = 800
const winH float64 = 800

const mapW float64 = 640
const mapH float64 = 576

var mapOrigin = pixel.V((winW-mapW)/2, (winH-mapH)/2)

const spriteSize float64 = 48

const translationFile = "i18n/zelduh/en-US.all.json"
const lang = "en-US"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var spritePlayerPath = "assets/bink-spritesheet-01.png"

func run() {
	// i18n
	i18n.MustLoadTranslationFile(translationFile)
	T, err := i18n.Tfunc(lang)
	if err != nil {
		panic(err)
	}

	// Setup Text
	orig := pixel.V(20, 50)
	txt := text.New(orig, text.Atlas7x13)
	txt.Color = palette.Map[palette.Darkest]

	coordDebugTxtOrig := pixel.V(5, 5)
	coordDebugTxt := text.New(coordDebugTxtOrig, text.Atlas7x13)
	coordDebugTxt.Color = palette.Map[palette.Darkest]

	// Setup GUI window
	cfg := pixelgl.WindowConfig{
		Title:  T("title"),
		Bounds: pixel.R(0, 0, winW, winH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Load sprite sheet graphic
	pic, err := loadPicture(spritePlayerPath)
	if err != nil {
		panic(err)
	}

	sprites := buildSpriteMap(pic, map[string]float64{
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
	})

	messageManager := message.Manager{}

	world := ecs.World{}
	world.AddSystem(&systems.SpatialSystem{})
	world.AddSystem(&systems.RenderSystem{Win: win})
	world.AddSystem(&systems.PlayerInputSystem{Win: win})
	world.AddSystem(&systems.CollisionSystem{
		Mailbox: &messageManager,
	})

	coins := []entity.Entity{}

	coinX := mapOrigin.X
	coinY := mapOrigin.Y
	for i := 0; i < 12; i++ {
		coin := entity.New(win, spriteSize, pixel.V(coinX, coinY), []*pixel.Sprite{
			sprites["coinA"],
			sprites["coinB"],
			sprites["coinC"],
		}, 7)
		coins = append(coins, coin)
		coinX = mapOrigin.X + float64(r.Intn(12)*48)
		coinY += 48
	}

	// Init player character
	player := player.New(win, spriteSize, 4, 3, 3, 1, map[string]*pixel.Sprite{
		"downA":  sprites["playerDownA"],
		"downB":  sprites["playerDownB"],
		"upA":    sprites["playerUpA"],
		"upB":    sprites["playerUpB"],
		"rightA": sprites["playerRightA"],
		"rightB": sprites["playerRightB"],
		"leftA":  sprites["playerLeftA"],
		"leftB":  sprites["playerLeftB"],
	}, pixel.V(mapOrigin.X+(mapW/2), mapOrigin.Y+(mapH/2)))

	playerEntity := entities.Player{
		BasicEntity: ecs.NewBasic(),
		AppearanceComponent: &components.AppearanceComponent{
			Color: colornames.Green,
		},
		SpatialComponent: &components.SpatialComponent{
			Width:  spriteSize,
			Height: spriteSize,
			Rect: pixel.R(
				mapOrigin.X+(mapW/2),
				mapOrigin.Y+(mapH/2),
				mapOrigin.X+(mapW/2)+spriteSize,
				mapOrigin.Y+(mapH/2)+spriteSize,
			),
			BoundsRect: pixel.R(
				mapOrigin.X,
				mapOrigin.Y,
				mapOrigin.X+mapW,
				mapOrigin.Y+mapH,
			),
			Shape: imdraw.New(nil),
		},
		MovementComponent: &components.MovementComponent{
			Moving:    false,
			Direction: direction.Down,
			Speed:     4.0,
		},
	}

	coinEntities := []entities.Coin{
		entities.Coin{
			BasicEntity: ecs.NewBasic(),
			AppearanceComponent: &components.AppearanceComponent{
				Color: colornames.Yellow,
			},
			SpatialComponent: &components.SpatialComponent{
				Width:  spriteSize,
				Height: spriteSize,
				Rect: pixel.R(
					mapOrigin.X,
					mapOrigin.Y,
					mapOrigin.X+spriteSize,
					mapOrigin.Y+spriteSize,
				),
				BoundsRect: pixel.R(
					mapOrigin.X,
					mapOrigin.Y,
					mapOrigin.X+mapW,
					mapOrigin.Y+mapH,
				),
				Shape: imdraw.New(nil),
			},
		},
	}

	// Create enemies
	enemies := []enemy.Enemy{}
	enemySprites := map[string]*pixel.Sprite{
		"downA":  sprites["turtleNoShellDownA"],
		"downB":  sprites["turtleNoShellDownB"],
		"upA":    sprites["turtleNoShellUpA"],
		"upB":    sprites["turtleNoShellUpB"],
		"rightA": sprites["turtleNoShellRightA"],
		"rightB": sprites["turtleNoShellRightB"],
		"leftA":  sprites["turtleNoShellLeftA"],
		"leftB":  sprites["turtleNoShellLeftB"],
	}
	for i := 0; i < 5; i++ {
		x := float64(r.Intn(int(mapW-spriteSize))) + mapOrigin.X
		y := float64(r.Intn(int(mapH-spriteSize))) + mapOrigin.Y
		var enemy = enemy.New(win, spriteSize, float64(x), float64(y), 1, 1, 1, enemySprites)
		enemies = append(enemies, enemy)
	}

	currentState := gamestate.Start

	sword := equipment.NewSword(win, spriteSize, sprites["sword"])

	mapOrigin := pixel.V(mapOrigin.X, mapOrigin.Y)

	// Add entities and components to systems
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *systems.PlayerInputSystem:
			sys.Add(&playerEntity.BasicEntity, playerEntity.MovementComponent)
		case *systems.SpatialSystem:
			sys.Add(&playerEntity.BasicEntity, playerEntity.SpatialComponent, playerEntity.MovementComponent)
		case *systems.RenderSystem:
			sys.Add(&playerEntity.BasicEntity, playerEntity.SpatialComponent, playerEntity.AppearanceComponent)
			for _, coin := range coinEntities {
				sys.Add(&coin.BasicEntity, coin.SpatialComponent, coin.AppearanceComponent)
			}
		case *systems.CollisionSystem:
			sys.Add(&playerEntity.BasicEntity, playerEntity.SpatialComponent)
			for _, coin := range coinEntities {
				sys.Add(&coin.BasicEntity, coin.SpatialComponent)
			}
		}
	}

	// TODO
	// https://github.com/EngoEngine/engo/blob/1b75afe871eca0c876f5884d0e24b86084b968f0/demos/pong/pong.go
	// Listen inside systems
	// Create "New" methods for systems and start listening inside this method
	// When message is received, loop through this system's entities and look for an ID match
	// messageManager.Listen("CollisionMessage", func(msg message.Message) {
	// fmt.Printf("Inbox alert: %v", msg.Type())

	// collision, isCollision := msg.(systems.CollisionMessage)
	// fmt.Printf("%v %v\n", x, y)
	// })

	for !win.Closed() {

		// For every state, allow quiting by pressing <q>
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		switch currentState {
		case gamestate.Start:
			win.Clear(palette.Map[palette.Dark])
			txt.Clear()
			drawMapBG(win, mapOrigin, mapW, mapH, palette.Map[palette.Lightest])
			txt.Color = palette.Map[palette.Darkest]
			fmt.Fprintln(txt, T("title"))
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			// Reset characters to starting positions
			player.Reset()
			for i := 0; i < len(enemies); i++ {
				enemies[i].Reset()
			}

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = gamestate.Game
			}
		case gamestate.Game:

			win.Clear(palette.Map[palette.Dark])
			drawMapBG(win, mapOrigin, mapW, mapH, palette.Map[palette.Lightest])
			txt.Clear()
			txt.Color = palette.Map[palette.Darkest]

			world.Update(0.125)

			player.Draw()

			for i := 0; i < len(coins); i++ {
				coins[i].Draw()
			}

			// check if player picked up a coin
			for i := len(coins) - 1; i > 0; i-- {
				collision := collision.IsColliding(player.Last, coins[i].Last, spriteSize)
				if collision {
					player.Deposit(1)
					// destroy coin
					coins = append(coins[:i], coins[i+1:]...)
					fmt.Printf("Coins remaining: %d\n", len(coins))
				}
			}

			for i := 0; i < len(enemies); i++ {
				if !enemies[i].IsDead() {
					// Check for collisions with enemy
					if sword.IsAttacking() && collision.IsColliding(sword.Last, enemies[i].Last, spriteSize) {
						// Sword hit enemy
						enemies[i].Hit(player.AttackPower)
					} else if collision.IsColliding(player.Last, enemies[i].Last, spriteSize) {
						// Enemy hit player
						player.Hit(enemies[i].AttackPower)
						if player.IsDead() {
							currentState = gamestate.Over
						}
					}
				}

				// Draw enemy if not dead
				if !enemies[i].IsDead() {
					enemies[i].Draw(mapOrigin.X, mapOrigin.Y, mapOrigin.X+mapW, mapOrigin.Y+mapH)
				}

			}

			if win.Pressed(pixelgl.KeyUp) {
				player.Move(mvmt.DirectionYPos, mapOrigin.Y+mapH, mapOrigin.Y, mapOrigin.X+mapW, mapOrigin.X)
			} else if win.Pressed(pixelgl.KeyRight) {
				player.Move(mvmt.DirectionXPos, mapOrigin.Y+mapH, mapOrigin.Y, mapOrigin.X+mapW, mapOrigin.X)
			} else if win.Pressed(pixelgl.KeyDown) {
				player.Move(mvmt.DirectionYNeg, mapOrigin.Y+mapH, mapOrigin.Y, mapOrigin.X+mapW, mapOrigin.X)
			} else if win.Pressed(pixelgl.KeyLeft) {
				player.Move(mvmt.DirectionXNeg, mapOrigin.Y+mapH, mapOrigin.Y, mapOrigin.X+mapW, mapOrigin.X)
			}

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Pause
			}

			if win.JustPressed(pixelgl.KeyX) {
				currentState = gamestate.Over
			}
			if win.JustPressed(pixelgl.KeySpace) {
				sword.Attack()
			}

			// Attack in direction player last moved
			switch player.LastDir {
			case mvmt.DirectionXPos:
				sword.Last = pixel.V(player.Last.X+player.SwordSize, player.Last.Y)
			case mvmt.DirectionXNeg:
				sword.Last = pixel.V(player.Last.X-player.SwordSize, player.Last.Y)
			case mvmt.DirectionYPos:
				sword.Last = pixel.V(player.Last.X, player.Last.Y+player.SwordSize)
			case mvmt.DirectionYNeg:
				sword.Last = pixel.V(player.Last.X, player.Last.Y-player.SwordSize)
			}

			sword.LastDir = player.LastDir

			sword.Draw()
		case gamestate.Pause:
			win.Clear(palette.Map[palette.Dark])
			txt.Clear()
			fmt.Fprintln(txt, T("paused"))
			drawMapBG(win, mapOrigin, mapW, mapH, palette.Map[palette.Lightest])
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Game
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				currentState = gamestate.Start
			}
		case gamestate.Over:
			win.Clear(palette.Map[palette.Dark])
			txt.Clear()
			drawMapBG(win, mapOrigin, mapW, mapH, palette.Map[palette.Darkest])
			txt.Color = palette.Map[palette.Darkest]
			fmt.Fprintln(txt, T("gameOver"))
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

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func drawMapBG(win *pixelgl.Window, origin pixel.Vec, w, h float64, color color.RGBA) {
	s := imdraw.New(nil)
	s.Push(origin)
	s.Push(pixel.V(origin.X+w, origin.Y+h))
	s.Color = color
	s.Rectangle(0)
	s.Draw(win)
}
