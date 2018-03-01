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

type Direction string

const (
	DirectionXPos Direction = "xPositive"
	DirectionXNeg Direction = "xNegative"
	DirectionYPos Direction = "yPositive"
	DirectionYNeg Direction = "yNegative"
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

	var playerSize float64 = 8
	var playerStartX float64 = (screenW / 2) - playerSize
	var playerStartY float64 = (screenH / 2) - playerSize
	var playerLastX float64 = playerStartX
	var playerLastY float64 = playerStartY

	var playerSwordSize float64 = 8

	var pcStrid float64 = playerSize

	var npcSize float64 = 8
	var npcStartX float64 = 0
	var npcStartY float64 = 0
	var npcLastX float64 = npcStartX
	var npcLastY float64 = npcStartY

	currentState := GameStateStart

	player := imdraw.New(nil)
	var playerLastDir Direction
	playerSword := imdraw.New(nil)
	npc := imdraw.New(nil)

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

			playerLastX = playerStartX
			playerLastY = playerStartY

			npcLastX = npcStartX
			npcLastY = npcStartY

			if win.JustPressed(pixelgl.KeyEnter) {
				fmt.Println("Transition from state %s to %s\n", currentState, GameStateGame)
				currentState = GameStateGame
			}
		case GameStateGame:
			win.Clear(colornames.Darkgreen)
			txt.Clear()
			npc.Color = colornames.Darkblue
			npc.Push(pixel.V(npcLastX, npcLastY))
			npc.Push(pixel.V(npcLastX+npcSize, npcLastY+npcSize))
			npc.Rectangle(0)
			npc.Draw(win)

			player.Clear()
			player.Color = colornames.White
			player.Push(pixel.V(playerLastX, playerLastY))
			player.Push(pixel.V(playerLastX+playerSize, playerLastY+playerSize))
			player.Rectangle(0)
			player.Draw(win)

			// Detect edge of window
			if win.JustPressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
				if playerLastY+pcStrid < screenH {
					playerLastY += pcStrid
					playerLastDir = DirectionYPos
				}
			} else if win.JustPressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
				if playerLastY-pcStrid >= 0 {
					playerLastY -= pcStrid
					playerLastDir = DirectionYNeg
				}
			} else if win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
				if playerLastX+pcStrid < screenW {
					playerLastX += pcStrid
					playerLastDir = DirectionXPos
				}
			} else if win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
				if playerLastX-pcStrid >= 0 {
					playerLastX -= pcStrid
					playerLastDir = DirectionXNeg
				}
			}

			if win.JustPressed(pixelgl.KeyP) {
				currentState = GameStatePause
			}

			if win.JustPressed(pixelgl.KeyX) {
				currentState = GameStateOver
			}

			if win.JustPressed(pixelgl.KeySpace) {
				// Attack with sword
				fmt.Printf("Sword attack direction: %s\n", playerLastDir)

				playerSword.Clear()
				playerSword.Color = colornames.Darkgray

				switch playerLastDir {
				case DirectionXPos:
					playerSword.Push(pixel.V(playerLastX+playerSwordSize, playerLastY))
					playerSword.Push(pixel.V(playerLastX+(playerSwordSize*2), playerLastY+playerSwordSize))
				case DirectionXNeg:
					playerSword.Push(pixel.V(playerLastX-playerSwordSize, playerLastY))
					playerSword.Push(pixel.V(playerLastX+playerSwordSize, playerLastY+playerSwordSize))
				case DirectionYPos:
					playerSword.Push(pixel.V(playerLastX, playerLastY+playerSwordSize))
					playerSword.Push(pixel.V(playerLastX+playerSwordSize, playerLastY+(playerSwordSize*2)))
				case DirectionYNeg:
					playerSword.Push(pixel.V(playerLastX, playerLastY-playerSwordSize))
					playerSword.Push(pixel.V(playerLastX+playerSwordSize, playerLastY+playerSwordSize))
				}

				playerSword.Rectangle(0)
				playerSword.Draw(win)
			}
		case GameStatePause:
			win.Clear(colornames.Darkblue)
			txt.Clear()
			fmt.Fprintln(txt, "Pause")
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyP) {
				currentState = GameStateGame
			}
			if win.JustPressed(pixelgl.KeyEscape) {
				currentState = GameStateStart
			}
		case GameStateOver:
			win.Clear(colornames.Black)
			txt.Clear()
			fmt.Fprintln(txt, "Game Over")
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			if win.JustPressed(pixelgl.KeyEnter) {
				currentState = GameStateStart
			}
		}

		win.Update()

	}
}

func main() {
	pixelgl.Run(run)
}
