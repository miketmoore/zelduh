package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/zelduh/equipment"
	"github.com/miketmoore/zelduh/gamestate"
	"github.com/miketmoore/zelduh/mvmt"
	"github.com/miketmoore/zelduh/npc"
	"github.com/miketmoore/zelduh/palette"
	"github.com/miketmoore/zelduh/pc"
	"github.com/nicksnyder/go-i18n/i18n"
)

const winW float64 = 800
const winH float64 = 800

const mapW float64 = 320
const mapH float64 = 288

var mapOrigin = pixel.V((winW-mapW)/2, (winH-mapH)/2)

const characterSize float64 = 16

const translationFile = "i18n/zelduh/en-US.all.json"
const lang = "en-US"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var spritePlayerPath = "assets/bink-spritesheet-01.png"

func g(n float64) float64 {
	return n * 16
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

	// Load sprite sheet graphic
	pic, err := loadPicture(spritePlayerPath)
	if err != nil {
		panic(err)
	}

	// Init player character
	player := pc.New(win, characterSize, 2, 3, 3, map[string]*pixel.Sprite{
		"downA": newSprite(pic, 0, g(7), g(1), g(8)),
		"downB": newSprite(pic, g(8), g(7), g(9), g(8)),

		"upA": newSprite(pic, g(1), g(7), g(2), g(8)),
		"upB": newSprite(pic, g(9), g(7), g(10), g(8)),

		"rightA": newSprite(pic, g(2), g(7), g(3), g(8)),
		"rightB": newSprite(pic, g(10), g(7), g(11), g(8)),

		"leftA": newSprite(pic, g(3), g(7), g(4), g(8)),
		"leftB": newSprite(pic, g(11), g(7), g(12), g(8)),
	}, pixel.V(mapOrigin.X+(mapW/2), mapOrigin.Y+(mapH/2)))

	// Create enemies
	enemies := []npc.Blob{}
	enemySprites := map[string]*pixel.Sprite{
		"downA": newSprite(pic, 0, g(6), g(1), g(7)),
		"downB": newSprite(pic, g(8), g(6), g(9), g(7)),
		"upA":   newSprite(pic, 0, g(6), g(1), g(7)),
		"upB":   newSprite(pic, g(8), g(6), g(9), g(7)),

		"rightA": newSprite(pic, g(1), g(6), g(2), g(7)),
		"rightB": newSprite(pic, g(9), g(6), g(10), g(7)),
		"leftA":  newSprite(pic, g(1), g(6), g(2), g(7)),
		"leftB":  newSprite(pic, g(9), g(6), g(10), g(7)),
	}
	for i := 0; i < 5; i++ {
		x := float64(r.Intn(int(mapW-characterSize))) + mapOrigin.X
		y := float64(r.Intn(int(mapH-characterSize))) + mapOrigin.Y
		var enemy = npc.NewBlob(win, characterSize, float64(x), float64(y), 1, 1, enemySprites)
		enemies = append(enemies, enemy)
	}

	currentState := gamestate.Start

	playerSword := equipment.NewSword()

	for !win.Closed() {

		// For every state, allow quiting by pressing <q>
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		// // Get mouse position and log to screen
		// mpos := win.MousePosition()
		// coordDebugTxt.Clear()
		// fmt.Fprintln(coordDebugTxt, fmt.Sprintf("%d, %d", int(math.Ceil(mpos.X)), int(math.Ceil(mpos.Y))))
		// coordDebugTxt.Draw(win, pixel.IM.Moved(coordDebugTxtOrig))

		switch currentState {
		case gamestate.Start:
			win.Clear(palette.Map[palette.Dark])
			txt.Clear()
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

			mapOrigin := pixel.V(mapOrigin.X, mapOrigin.Y)
			drawMapBG(win, mapOrigin, mapW, mapH)

			// draw tiles

			abc := mapW / 16
			var i float64
			// fmt.Println(abc)
			for ; i < abc; i++ {
				tileASprite := newSprite(pic, g(6), g(3), g(7), g(4))
				tileA := imdraw.New(nil)
				tileA.Push(pixel.V(0, 0))
				tileA.Push(pixel.V(16, 16))
				tileA.Rectangle(0)
				a := mapOrigin.X + 16/2
				b := mapOrigin.Y + 16/2
				// fmt.Printf("%f %f\n", a, b)
				tileASprite.Draw(win, pixel.IM.Moved(pixel.V(a+(16*i), b)))
			}

			// matrix := pixel.IM.Moved(pixel.V(player.Last.X+player.Size/2, player.Last.Y+player.Size/2))
			// sprite.Draw(player.Win, matrix)

			player.Draw()
			for i := 0; i < len(enemies); i++ {

				// collision detection
				collision := player.Last.Y > (enemies[i].Last.Y+enemies[i].Size) ||
					(player.Last.Y+player.Size) < enemies[i].Last.Y ||
					player.Last.X > (enemies[i].Last.X+enemies[i].Size) ||
					(player.Last.X+player.Size) < enemies[i].Last.X

				if !collision {
					fmt.Printf("Collision\n")
					// TODO move character back x pixels, in opposite direction that enemy
					// is facing.
					player.Hit(enemies[i].AttackPower)
					if player.IsDead() {
						currentState = gamestate.Over
					}
				} else {
					// fmt.Printf("No collision\n")
					enemies[i].Draw(mapOrigin.X, mapOrigin.Y, mapOrigin.X+mapW, mapOrigin.Y+mapH)
				}
			}

			if win.Pressed(pixelgl.KeyUp) && player.Last.Y+player.Stride <= ((mapOrigin.Y+mapH)-player.Size) {
				player.Last.Y += player.Stride
				player.LastDir = mvmt.DirectionYPos
			} else if win.Pressed(pixelgl.KeyDown) && player.Last.Y-player.Stride >= mapOrigin.Y {
				player.Last.Y -= player.Stride
				player.LastDir = mvmt.DirectionYNeg
			} else if win.Pressed(pixelgl.KeyRight) && player.Last.X+player.Stride <= ((mapOrigin.X+mapW)-player.Size) {
				player.Last.X += player.Stride
				player.LastDir = mvmt.DirectionXPos
			} else if win.Pressed(pixelgl.KeyLeft) && player.Last.X-player.Stride >= mapOrigin.X {
				player.Last.X -= player.Stride
				player.LastDir = mvmt.DirectionXNeg
			}

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Pause
			}

			if win.JustPressed(pixelgl.KeyX) {
				currentState = gamestate.Over
			}

			if win.JustPressed(pixelgl.KeySpace) {
				// Attack with sword
				// TODO sword should appear on screen longer
				// I think a state machine for the sword makes sense
				// [sheathed]
				// [attacking] - this would go away after x ticks
				fmt.Printf("Sword attack direction: %s\n", player.LastDir)

				playerSword.Shape.Clear()
				playerSword.Shape.Color = palette.Map[palette.Lightest]

				// Attack in direction player last moved
				switch player.LastDir {
				case mvmt.DirectionXPos:
					playerSword.Shape.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y))
					playerSword.Shape.Push(pixel.V(player.Last.X+(player.SwordSize*2), player.Last.Y+player.SwordSize))
				case mvmt.DirectionXNeg:
					playerSword.Shape.Push(pixel.V(player.Last.X-player.SwordSize, player.Last.Y))
					playerSword.Shape.Push(pixel.V(player.Last.X, player.Last.Y+player.SwordSize))
				case mvmt.DirectionYPos:
					playerSword.Shape.Push(pixel.V(player.Last.X, player.Last.Y+player.SwordSize))
					playerSword.Shape.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+(player.SwordSize*2)))
				case mvmt.DirectionYNeg:
					playerSword.Shape.Push(pixel.V(player.Last.X, player.Last.Y-player.SwordSize))
					playerSword.Shape.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y))
				}

				playerSword.Shape.Rectangle(0)
				playerSword.Shape.Draw(win)
			}
		case gamestate.Pause:
			win.Clear(palette.Map[palette.Dark])
			txt.Clear()
			fmt.Fprintln(txt, T("paused"))
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Game
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				currentState = gamestate.Start
			}
		case gamestate.Over:
			win.Clear(palette.Map[palette.Darkest])
			txt.Clear()
			txt.Color = palette.Map[palette.Lightest]
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
		Min: pixel.V(xa, ya),
		Max: pixel.V(xb, yb),
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

func drawMapBG(win *pixelgl.Window, origin pixel.Vec, w, h float64) {
	s := imdraw.New(nil)

	// bottom left point
	s.Push(origin)

	// top right point
	s.Push(pixel.V(origin.X+w, origin.Y+h))

	s.Color = palette.Map[palette.Lightest]
	s.Rectangle(0)
	s.Draw(win)
}
