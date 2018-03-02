package main

import (
	"fmt"
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
	"github.com/miketmoore/zelduh/pc"
	"github.com/nicksnyder/go-i18n/i18n"
	"golang.org/x/image/colornames"
)

// These should be multiples of 8 for now
const screenW float64 = 320
const screenH float64 = 288

const characterSize float64 = 16

const translationFile = "i18n/zelduh/en-US.all.json"
const lang = "en-US"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	txt.Color = colornames.White

	coordDebugTxtOrig := pixel.V(5, 5)
	coordDebugTxt := text.New(coordDebugTxtOrig, text.Atlas7x13)
	coordDebugTxt.Color = colornames.White

	// Setup GUI window
	cfg := pixelgl.WindowConfig{
		Title:  T("title"),
		Bounds: pixel.R(0, 0, screenW, screenH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Init player character
	player := pc.Player{
		Win:       win,
		Size:      characterSize,
		Shape:     imdraw.New(nil),
		SwordSize: characterSize,
		Health:    3,
		MaxHealth: 3,
	}
	player.Start = pixel.V((screenW/2)-player.Size, (screenH/2)-player.Size)
	player.Last = player.Start
	player.Stride = player.Size

	// Create enemies
	enemies := []npc.Blob{}
	for i := 0; i < 3; i++ {
		x := r.Intn(int(screenW - characterSize))
		y := r.Intn(int(screenH - characterSize))
		var enemy = npc.NewBlob(win, characterSize, float64(x), float64(y), 1, 1)
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
			win.Clear(colornames.Darkgreen)
			txt.Clear()
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
			win.Clear(colornames.Darkgreen)
			txt.Clear()

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
					enemies[i].Draw(screenW, screenH)
				}
			}

			// Detect edge of window
			if win.JustPressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
				if player.Last.Y+player.Stride < screenH {
					player.Last.Y += player.Stride
					player.LastDir = mvmt.DirectionYPos
				}
			} else if win.JustPressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
				if player.Last.Y-player.Stride >= 0 {
					player.Last.Y -= player.Stride
					player.LastDir = mvmt.DirectionYNeg
				}
			} else if win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
				if player.Last.X+player.Stride < screenW {
					player.Last.X += player.Stride
					player.LastDir = mvmt.DirectionXPos
				}
			} else if win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
				if player.Last.X-player.Stride >= 0 {
					player.Last.X -= player.Stride
					player.LastDir = mvmt.DirectionXNeg
				}
			}

			if win.JustPressed(pixelgl.KeyP) {
				currentState = gamestate.Pause
			}

			if win.JustPressed(pixelgl.KeyX) {
				currentState = gamestate.Over
			}

			if win.JustPressed(pixelgl.KeySpace) {
				// Attack with sword
				fmt.Printf("Sword attack direction: %s\n", player.LastDir)

				playerSword.Shape.Clear()
				playerSword.Shape.Color = colornames.Darkgray

				// Attack in direction player last moved
				switch player.LastDir {
				case mvmt.DirectionXPos:
					playerSword.Shape.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y))
					playerSword.Shape.Push(pixel.V(player.Last.X+(player.SwordSize*2), player.Last.Y+player.SwordSize))
				case mvmt.DirectionXNeg:
					playerSword.Shape.Push(pixel.V(player.Last.X-player.SwordSize, player.Last.Y))
					playerSword.Shape.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+player.SwordSize))
				case mvmt.DirectionYPos:
					playerSword.Shape.Push(pixel.V(player.Last.X, player.Last.Y+player.SwordSize))
					playerSword.Shape.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+(player.SwordSize*2)))
				case mvmt.DirectionYNeg:
					playerSword.Shape.Push(pixel.V(player.Last.X, player.Last.Y-player.SwordSize))
					playerSword.Shape.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+player.SwordSize))
				}

				playerSword.Shape.Rectangle(0)
				playerSword.Shape.Draw(win)
			}
		case gamestate.Pause:
			win.Clear(colornames.Darkblue)
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
			win.Clear(colornames.Black)
			txt.Clear()
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
