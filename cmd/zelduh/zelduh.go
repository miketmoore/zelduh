package main

import (
	"fmt"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const screenW = 160
const screenH = 144

var title string = "Zelduh"

type GameState string

const (
	GameStateStart GameState = "start"
	GameStateGame  GameState = "game"
	GameStatePause GameState = "pause"
	GameStateOver  GameState = "over"
)

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

	// Draw player character

	var size float64 = 8
	var lastX float64 = screenW - size
	var lastY float64 = screenH - size

	var pcStrid float64 = size

	var npcSize float64 = 8
	var npcLastX float64 = 0
	var npcLastY float64 = 0

	currentState := GameStateStart

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
		case GameStateStart:
			win.Clear(colornames.Darkgreen)
			txt.Clear()
			fmt.Fprintln(txt, title)
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyEnter) {
				fmt.Println("Transition from state %s to %s\n", currentState, GameStateGame)
				currentState = GameStateGame
			}
		case GameStateGame:
			win.Clear(colornames.Darkgreen)
			txt.Clear()
			npc := imdraw.New(nil)
			npc.Color = colornames.Darkblue
			npc.Push(pixel.V(npcLastX, npcLastY))
			npc.Push(pixel.V(npcLastX+npcSize, npcLastY+npcSize))
			npc.Rectangle(0)
			npc.Draw(win)

			pc := imdraw.New(nil)
			pc.Color = colornames.White
			pc.Push(pixel.V(lastX, lastY))
			pc.Push(pixel.V(lastX+size, lastY+size))
			pc.Rectangle(0)
			pc.Draw(win)

			// Detect edge of window
			if win.JustPressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
				if lastY+pcStrid < screenH {
					lastY += pcStrid
				}
			} else if win.JustPressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
				if lastY-pcStrid >= 0 {
					lastY -= pcStrid
				}
			} else if win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
				if lastX+pcStrid < screenW {
					lastX += pcStrid
				}
			} else if win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
				if lastX-pcStrid >= 0 {
					lastX -= pcStrid
				}
			}

		}

		win.Update()

	}
}

func main() {
	pixelgl.Run(run)
}
