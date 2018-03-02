package main

import (
	"fmt"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/zelduh/gamestate"
	"github.com/miketmoore/zelduh/mvmt"
	"github.com/miketmoore/zelduh/npc"
	"github.com/miketmoore/zelduh/pc"
	"golang.org/x/image/colornames"
)

const screenW = 160
const screenH = 144

var title = "Zelduh"

func run() {
	// Setup Text
	orig := pixel.V(20, 50)
	txt := text.New(orig, text.Atlas7x13)
	txt.Color = colornames.White

	coordDebugTxtOrig := pixel.V(5, 5)
	coordDebugTxt := text.New(coordDebugTxtOrig, text.Atlas7x13)
	coordDebugTxt.Color = colornames.White

	// Setup GUI window
	cfg := pixelgl.WindowConfig{
		Title:  title,
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
		Size:      8,
		Shape:     imdraw.New(nil),
		SwordSize: 8,
	}
	player.Start = pixel.V((screenW/2)-player.Size, (screenH/2)-player.Size)
	player.Last = player.Start
	player.Stride = player.Size

	// Init non-player character
	var blob = npc.Blob{
		Win:   win,
		Size:  8,
		Start: pixel.V(0, 0),
		Last:  pixel.V(0, 0),
		Shape: imdraw.New(nil),
	}

	currentState := gamestate.Start

	playerSword := imdraw.New(nil)

	for !win.Closed() {

		// For every state, allow quiting by pressing <q>
		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		// Get mouse position and log to screen
		mpos := win.MousePosition()
		coordDebugTxt.Clear()
		fmt.Fprintln(coordDebugTxt, fmt.Sprintf("%d, %d", int(math.Ceil(mpos.X)), int(math.Ceil(mpos.Y))))
		coordDebugTxt.Draw(win, pixel.IM.Moved(coordDebugTxtOrig))

		switch currentState {
		case gamestate.Start:
			win.Clear(colornames.Darkgreen)
			txt.Clear()
			fmt.Fprintln(txt, title)
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			// player.Last.X = playerStartX
			// player.Last.Y = playerStartY
			player.Last = player.Start

			// npcLastX = npcStartX
			// npcLastY = npcStartY
			blob.Last = blob.Start

			if win.JustPressed(pixelgl.KeyEnter) {
				fmt.Println("Transition from state %s to %s\n", currentState, gamestate.Game)
				currentState = gamestate.Game
			}
		case gamestate.Game:
			win.Clear(colornames.Darkgreen)
			txt.Clear()

			blob.Draw(pixel.V(player.Last.X, player.Last.Y))
			player.Draw()

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

				playerSword.Clear()
				playerSword.Color = colornames.Darkgray

				// Attack in direction player last moved
				switch player.LastDir {
				case mvmt.DirectionXPos:
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y))
					playerSword.Push(pixel.V(player.Last.X+(player.SwordSize*2), player.Last.Y+player.SwordSize))
				case mvmt.DirectionXNeg:
					playerSword.Push(pixel.V(player.Last.X-player.SwordSize, player.Last.Y))
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+player.SwordSize))
				case mvmt.DirectionYPos:
					playerSword.Push(pixel.V(player.Last.X, player.Last.Y+player.SwordSize))
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+(player.SwordSize*2)))
				case mvmt.DirectionYNeg:
					playerSword.Push(pixel.V(player.Last.X, player.Last.Y-player.SwordSize))
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+player.SwordSize))
				}

				playerSword.Rectangle(0)
				playerSword.Draw(win)
			}
		case gamestate.Pause:
			win.Clear(colornames.Darkblue)
			txt.Clear()
			fmt.Fprintln(txt, "Pause")
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
			fmt.Fprintln(txt, "Game Over")
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
