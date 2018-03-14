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
	"github.com/miketmoore/zelduh/collision"
	"github.com/miketmoore/zelduh/enemy"
	"github.com/miketmoore/zelduh/entity"
	"github.com/miketmoore/zelduh/equipment"
	"github.com/miketmoore/zelduh/gamestate"
	"github.com/miketmoore/zelduh/mvmt"
	"github.com/miketmoore/zelduh/palette"
	"github.com/miketmoore/zelduh/player"
	"github.com/nicksnyder/go-i18n/i18n"
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

func g(n float64) float64 {
	return n * spriteSize
}

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

	win.SetSmooth(true)

	// Load sprite sheet graphic
	pic, err := loadPicture(spritePlayerPath)
	if err != nil {
		panic(err)
	}

	sprites := map[string]*pixel.Sprite{
		"playerDownA":  newSpriteIndexed(pic, 109),
		"playerDownB":  newSpriteIndexed(pic, 118),
		"playerUpA":    newSpriteIndexed(pic, 110),
		"playerUpB":    newSpriteIndexed(pic, 119),
		"playerRightA": newSpriteIndexed(pic, 111),
		"playerRightB": newSpriteIndexed(pic, 120),
		"playerLeftA":  newSpriteIndexed(pic, 112),
		"playerLeftB":  newSpriteIndexed(pic, 121),

		"turtleNoShellDownA":  newSpriteIndexed(pic, 1),
		"turtleNoShellDownB":  newSpriteIndexed(pic, 10),
		"turtleNoShellUpA":    newSpriteIndexed(pic, 2),
		"turtleNoShellUpB":    newSpriteIndexed(pic, 11),
		"turtleNoShellRightA": newSpriteIndexed(pic, 3),
		"turtleNoShellRightB": newSpriteIndexed(pic, 12),
		"turtleNoShellLeftA":  newSpriteIndexed(pic, 4),
		"turtleNoShellLeftB":  newSpriteIndexed(pic, 13),

		"sword": newSpriteIndexed(pic, 84),

		"ground": newSpriteIndexed(pic, 8),

		"coinA": newSpriteIndexed(pic, 113),
		"coinB": newSpriteIndexed(pic, 122),
		"coinC": newSpriteIndexed(pic, 131),
	}

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
		"downA": sprites["playerDownA"],
		"downB": sprites["playerDownB"],

		"upA": sprites["playerUpA"],
		"upB": sprites["playerUpB"],

		"rightA": sprites["playerRightA"],
		"rightB": sprites["playerRightB"],

		"leftA": sprites["playerLeftA"],
		"leftB": sprites["playerLeftB"],
	}, pixel.V(mapOrigin.X+(mapW/2), mapOrigin.Y+(mapH/2)))

	// Create enemies
	enemies := []enemy.Enemy{}
	enemySprites := map[string]*pixel.Sprite{
		"downA": sprites["turtleNoShellDownA"],
		"downB": sprites["turtleNoShellDownB"],
		"upA":   sprites["turtleNoShellUpA"],
		"upB":   sprites["turtleNoShellUpB"],

		"rightA": sprites["turtleNoShellRightA"],
		"rightB": sprites["turtleNoShellRightB"],
		"leftA":  sprites["turtleNoShellLeftA"],
		"leftB":  sprites["turtleNoShellLeftB"],
	}
	for i := 0; i < 5; i++ {
		x := float64(r.Intn(int(mapW-spriteSize))) + mapOrigin.X
		y := float64(r.Intn(int(mapH-spriteSize))) + mapOrigin.Y
		var enemy = enemy.NewBlob(win, spriteSize, float64(x), float64(y), 1, 1, 1, enemySprites)
		enemies = append(enemies, enemy)
	}

	currentState := gamestate.Start

	sword := equipment.NewSword(win, spriteSize, sprites["sword"])

	batchDryGround := pixel.NewBatch(&pixel.TrianglesData{}, pic)

	var spriteFrames []pixel.Rect
	fmt.Println(pic.Bounds())
	for x := pic.Bounds().Min.X; x < pic.Bounds().Max.X; x += spriteSize {
		for y := pic.Bounds().Min.Y; y < pic.Bounds().Max.Y; y += spriteSize {
			spriteFrames = append(spriteFrames, pixel.R(x, y, x+spriteSize, y+spriteSize))
		}
	}

	spriteDryGround := pixel.NewSprite(pic, spriteFrames[35])
	// mouse := cam.Unproject(win.MousePosition())
	// tree.Draw(batch, pixel.IM.Scaled(pixel.ZV, 4).Moved(mouse))
	for i := 0.0; i < mapW/spriteSize; i++ {
		for j := 0.0; j < mapH/spriteSize; j++ {
			spriteDryGround.Draw(batchDryGround,
				pixel.IM.Moved(pixel.V(mapOrigin.X+(spriteSize/2)+float64(spriteSize*i),
					mapOrigin.Y+(spriteSize/2)+float64(spriteSize*j))),
			)
		}
	}

	mapOrigin := pixel.V(mapOrigin.X, mapOrigin.Y)

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
			txt.Clear()
			txt.Color = palette.Map[palette.Darkest]

			batchDryGround.Draw(win)

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

func newSprite(pic pixel.Picture, xa, ya, xb, yb float64) *pixel.Sprite {
	return pixel.NewSprite(pic, pixel.Rect{
		Min: pixel.V(g(xa), g(ya)),
		Max: pixel.V(g(xb), g(yb)),
	})
}

func newSpriteIndexed(pic pixel.Picture, index float64) *pixel.Sprite {
	fmt.Printf("newSpriteIndexed index: %f\n", index)
	// iterate over width every spriteSize
	totalRows := pic.Bounds().H() / spriteSize
	totalCols := pic.Bounds().W() / spriteSize
	fmt.Printf("total rows: %f\n", totalRows)
	fmt.Printf("total cols: %f\n", totalCols)

	find := func() (float64, float64) {
		i := 0.0
		var row = 0.0
		var col = 0.0
		for ; row < totalRows; row++ {
			fmt.Printf("i:%f\n", i)
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

	fmt.Printf("found row, col: %f, %f\n", row, col)
	fmt.Println(row*spriteSize, col*spriteSize)

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

	// bottom left point
	s.Push(origin)

	// top right point
	s.Push(pixel.V(origin.X+w, origin.Y+h))

	s.Color = color
	s.Rectangle(0)
	s.Draw(win)
}

// func drawMap(pic pixel.Picture, win *pixelgl.Window, name string) {
// 	if name == "drycracked" {
// 		abc := mapW / spriteSize
// 		def := mapH / spriteSize
// 		var x, y float64
// 		for ; x < abc; x++ {
// 			for y = 0; y < def; y++ {
// 				tileASprite := newSprite(pic, 6, 3, 7, 4)
// 				tileA := imdraw.New(nil)
// 				tileA.Push(pixel.V(0, 0))
// 				tileA.Push(pixel.V(spriteSize, spriteSize))
// 				tileA.Rectangle(0)
// 				a := mapOrigin.X + spriteSize/2
// 				b := mapOrigin.Y + spriteSize/2
// 				tileASprite.Draw(win, pixel.IM.Moved(pixel.V(a+(spriteSize*x), b+(spriteSize*y))))
// 			}
// 		}
// 	}

// }
