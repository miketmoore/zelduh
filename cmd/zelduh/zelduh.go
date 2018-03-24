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
	"github.com/miketmoore/go-pixel-game-template/state"
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
	"github.com/miketmoore/zelduh/playerinput"
	"github.com/miketmoore/zelduh/render"
	"github.com/miketmoore/zelduh/spatial"
	"github.com/miketmoore/zelduh/systems"
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
	mapW float64 = 640
	mapH float64 = 576
	mapX         = (winW - mapW) / 2
	mapY         = (winH - mapH) / 2
)

var (
	win       *pixelgl.Window
	txt       *text.Text
	t         i18n.TranslateFunc
	currState state.State
	pic       pixel.Picture
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
	txt = initText(20, 50, palette.Map[palette.Darkest])
	win = initWindow(t("title"), winX, winY, winW, winH)
	pic = loadPicture(spritePlayerPath)
	sprites = buildSpriteMap(pic, spriteMap)
	messageManager := message.Manager{}
	ecsWorld := ecs.World{}
	ecsWorld.AddSystem(&systems.SpatialSystem{})
	ecsWorld.AddSystem(&systems.RenderSystem{Win: win})
	ecsWorld.AddSystem(&systems.PlayerInputSystem{Win: win})
	// ecsWorld.AddSystem(&systems.CollisionSystem{
	// 	Mailbox: &messageManager,
	// })
	ecsWorld.AddSystem(&systems.CoinsSystem{
		Mailbox: &messageManager,
	})

	customWorld := world.New()

	customWorld.AddSystem(&playerinput.System{Win: win})
	customWorld.AddSystem(&spatial.System{})
	customWorld.AddSystem(&render.System{Win: win})
	customWorld.AddSystem(&collision.System{})

	// Old "entities"... phasing out
	coins := buildCoins()
	player := buildPlayer()
	enemies := buildEnemies()
	sword := buildSword()

	// New entities
	playerEntity := buildPlayerEntity()
	coinEntities := buildCoinEntities(customWorld)

	currentState := gamestate.Start

	// Add entity components to custom ECS systems
	for _, system := range customWorld.Systems() {
		switch sys := system.(type) {
		case *playerinput.System:
			sys.AddPlayer(playerEntity.MovementComponent)
		case *spatial.System:
			sys.AddPlayer(playerEntity.SpatialComponent, playerEntity.MovementComponent)
		case *collision.System:
			sys.AddPlayer(playerEntity.SpatialComponent)
			for _, coin := range coinEntities {
				sys.AddCoin(coin.ID, coin.SpatialComponent)
			}
		case *render.System:
			sys.AddPlayer(playerEntity.AppearanceComponent, playerEntity.SpatialComponent)
		}
	}

	// Add entities and components to systems
	// for _, system := range ecsWorld.Systems() {
	// 	switch sys := system.(type) {
	// 	// case *systems.PlayerInputSystem:
	// 	// sys.Add(&playerEntity.BasicEntity, playerEntity.MovementComponent)
	// 	// case *systems.SpatialSystem:
	// 	// 	sys.Add(&playerEntity.BasicEntity, playerEntity.SpatialComponent, playerEntity.MovementComponent)
	// 	case *systems.RenderSystem:
	// 		// sys.Add(&playerEntity.BasicEntity, playerEntity.SpatialComponent, playerEntity.AppearanceComponent)
	// 		for _, coin := range coinEntities {
	// 			sys.Add(&coin.BasicEntity, coin.SpatialComponent, coin.AppearanceComponent)
	// 		}
	// 		// case *systems.CollisionSystem:
	// 		// 	sys.Add(&playerEntity.BasicEntity, playerEntity.SpatialComponent, playerEntity.EntityTypeComponent)
	// 		// 	for _, coin := range coinEntities {
	// 		// 		sys.Add(&coin.BasicEntity, coin.SpatialComponent, coin.EntityTypeComponent)
	// 		// 	}
	// 		// case *systems.CoinsSystem:
	// 		// 	sys.Add(&playerEntity.BasicEntity, playerEntity.CoinsComponent)
	// 	}
	// }

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

	// messageManager.Listen("DestroyCoinMessage", func(msg message.Message) {
	// 	fmt.Printf("DESTROY COIN NOW!\n")
	// })

	for !win.Closed() {

		allowQuit()

		switch currentState {
		case gamestate.Start:
			win.Clear(palette.Map[palette.Dark])
			txt.Clear()
			drawMapBG(mapX, mapY, mapW, mapH, palette.Map[palette.Lightest])
			txt.Color = palette.Map[palette.Darkest]
			fmt.Fprintln(txt, t("title"))
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
			drawMapBG(mapX, mapY, mapW, mapH, palette.Map[palette.Lightest])
			txt.Clear()
			txt.Color = palette.Map[palette.Darkest]

			ecsWorld.Update(0.125)
			customWorld.Update()

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
					enemies[i].Draw(mapX, mapY, mapX+mapW, mapY+mapH)
				}

			}

			if win.Pressed(pixelgl.KeyUp) {
				player.Move(mvmt.DirectionYPos, mapY+mapH, mapY, mapX+mapW, mapX)
			} else if win.Pressed(pixelgl.KeyRight) {
				player.Move(mvmt.DirectionXPos, mapY+mapH, mapY, mapX+mapW, mapX)
			} else if win.Pressed(pixelgl.KeyDown) {
				player.Move(mvmt.DirectionYNeg, mapY+mapH, mapY, mapX+mapW, mapX)
			} else if win.Pressed(pixelgl.KeyLeft) {
				player.Move(mvmt.DirectionXNeg, mapY+mapH, mapY, mapX+mapW, mapX)
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
			fmt.Fprintln(txt, t("paused"))
			drawMapBG(mapX, mapY, mapW, mapH, palette.Map[palette.Lightest])
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
			drawMapBG(mapX, mapY, mapW, mapH, palette.Map[palette.Darkest])
			txt.Color = palette.Map[palette.Darkest]
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

func buildCoins() []entity.Entity {
	coins := []entity.Entity{}

	coinX := mapX
	coinY := mapY
	for i := 0; i < 12; i++ {
		coin := entity.New(win, spriteSize, pixel.V(coinX, coinY), []*pixel.Sprite{
			sprites["coinA"],
			sprites["coinB"],
			sprites["coinC"],
		}, 7)
		coins = append(coins, coin)
		coinX = mapX + float64(r.Intn(12)*48)
		coinY += 48
	}

	return coins
}

func buildPlayer() player.Player {
	return player.New(win, spriteSize, 4, 3, 3, 1, map[string]*pixel.Sprite{
		"downA":  sprites["playerDownA"],
		"downB":  sprites["playerDownB"],
		"upA":    sprites["playerUpA"],
		"upB":    sprites["playerUpB"],
		"rightA": sprites["playerRightA"],
		"rightB": sprites["playerRightB"],
		"leftA":  sprites["playerLeftA"],
		"leftB":  sprites["playerLeftB"],
	}, pixel.V(mapX+(mapW/2), mapY+(mapH/2)))
}

func buildPlayerEntity() entities.Player {
	return entities.Player{
		BasicEntity: ecs.NewBasic(),
		EntityTypeComponent: &components.EntityTypeComponent{
			Type: "player",
		},
		AppearanceComponent: &systems.AppearanceComponent{
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
			BoundsRect: pixel.R(
				mapX,
				mapY,
				mapX+mapW,
				mapY+mapH,
			),
			Shape: imdraw.New(nil),
		},
		MovementComponent: &components.MovementComponent{
			Moving:    false,
			Direction: direction.Down,
			Speed:     4.0,
		},
		CoinsComponent: &systems.CoinsComponent{
			Coins: 0,
		},
	}
}

func buildCoinEntities(world world.World) []entities.Coin {
	return []entities.Coin{
		entities.Coin{
			ID:          world.NewEntityID(),
			BasicEntity: ecs.NewBasic(),
			EntityTypeComponent: &components.EntityTypeComponent{
				Type: "coin",
			},
			AppearanceComponent: &systems.AppearanceComponent{
				Color: colornames.Yellow,
			},
			SpatialComponent: &components.SpatialComponent{
				Width:  spriteSize,
				Height: spriteSize,
				Rect: pixel.R(
					mapX,
					mapY,
					mapX+spriteSize,
					mapY+spriteSize,
				),
				BoundsRect: pixel.R(
					mapX,
					mapY,
					mapX+mapW,
					mapY+mapH,
				),
				Shape: imdraw.New(nil),
			},
		},
	}
}

func buildEnemies() []enemy.Enemy {
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
		x := float64(r.Intn(int(mapW-spriteSize))) + mapX
		y := float64(r.Intn(int(mapH-spriteSize))) + mapY
		var enemy = enemy.New(win, spriteSize, float64(x), float64(y), 1, 1, 1, enemySprites)
		enemies = append(enemies, enemy)
	}
	return enemies
}

func buildSword() equipment.Sword {
	return equipment.NewSword(win, spriteSize, sprites["sword"])
}
