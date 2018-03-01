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

// Player represents the player character
type Player struct {
	// Size is the dimensions (square)
	Size float64
	// Start is the starting vector
	Start pixel.Vec
	// Last is the last vector
	Last pixel.Vec
	// LastDir is the last direction the player was headed in
	LastDir Direction
	// Shape is the view
	Shape *imdraw.IMDraw
	// Win is a pointer to the GUI window
	Win *pixelgl.Window
	// SwordSize is the dimensions of the sword
	SwordSize float64
	// Stride is how many tiles character can move in one "step"
	Stride float64
}

func (player *Player) Draw() {
	shape := player.Shape
	shape.Clear()
	shape.Color = colornames.White
	shape.Push(pixel.V(player.Last.X, player.Last.Y))
	shape.Push(pixel.V(player.Last.X+player.Size, player.Last.Y+player.Size))
	shape.Rectangle(0)
	shape.Draw(player.Win)
}

// NPC represents one non-player character
type NPC struct {
	// Size is the dimensions (square)
	Size float64
	// Start is the starting vector
	Start pixel.Vec
	// Last is the last vector
	Last pixel.Vec
	// LastDir is the last direction the NPC was headed in
	LastDir Direction
	// Shape is the view
	Shape *imdraw.IMDraw
	// Win is a pointer to the GUI window
	Win *pixelgl.Window
}

func (npc *NPC) Draw(playerLast pixel.Vec) {
	// Move npc with AI
	// velocity := 1
	// max_velocity = 2
	// desired_velocity = normalize(target - position) * max_velocity
	// steering = desired_velocity - velocity
	// velocity := float64(1)    // move one space per update
	npc.Shape.Clear()
	maxVelocity := float64(0.75) // move three spaces per update
	normalized := normalize(playerLast, npc.Last)
	desiredVelocity := pixel.V(normalized.X*maxVelocity, normalized.Y*maxVelocity)
	fmt.Printf("desiredVelocity: %f, %f\n", desiredVelocity.X, desiredVelocity.Y)

	npc.Shape.Color = colornames.Darkblue
	// npc.Push(pixel.V(npcLastX, npcLastY))
	// npc.Push(pixel.V(npcLastX+npcSize, npcLastY+npcSize))
	npc.Shape.Push(desiredVelocity)
	npc.Shape.Push(pixel.V(desiredVelocity.X+npc.Size, desiredVelocity.Y+npc.Size))
	npc.Shape.Rectangle(0)
	npc.Shape.Draw(npc.Win)
}

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
	player := Player{
		Win:       win,
		Size:      8,
		Shape:     imdraw.New(nil),
		SwordSize: 8,
	}
	player.Start = pixel.V((screenW/2)-player.Size, (screenH/2)-player.Size)
	player.Last = player.Start
	player.Stride = player.Size

	// Init non-player character
	var npc = NPC{
		Win:   win,
		Size:  8,
		Start: pixel.V(0, 0),
		Last:  pixel.V(0, 0),
		Shape: imdraw.New(nil),
	}

	currentState := GameStateStart

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
		case GameStateStart:
			win.Clear(colornames.Darkgreen)
			txt.Clear()
			fmt.Fprintln(txt, title)
			txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

			// player.Last.X = playerStartX
			// player.Last.Y = playerStartY
			player.Last = player.Start

			// npcLastX = npcStartX
			// npcLastY = npcStartY
			npc.Last = npc.Start

			if win.JustPressed(pixelgl.KeyEnter) {
				fmt.Println("Transition from state %s to %s\n", currentState, GameStateGame)
				currentState = GameStateGame
			}
		case GameStateGame:
			win.Clear(colornames.Darkgreen)
			txt.Clear()

			npc.Draw(pixel.V(player.Last.X, player.Last.Y))
			player.Draw()

			// Detect edge of window
			if win.JustPressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
				if player.Last.Y+player.Stride < screenH {
					player.Last.Y += player.Stride
					player.LastDir = DirectionYPos
				}
			} else if win.JustPressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
				if player.Last.Y-player.Stride >= 0 {
					player.Last.Y -= player.Stride
					player.LastDir = DirectionYNeg
				}
			} else if win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
				if player.Last.X+player.Stride < screenW {
					player.Last.X += player.Stride
					player.LastDir = DirectionXPos
				}
			} else if win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
				if player.Last.X-player.Stride >= 0 {
					player.Last.X -= player.Stride
					player.LastDir = DirectionXNeg
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
				fmt.Printf("Sword attack direction: %s\n", player.LastDir)

				playerSword.Clear()
				playerSword.Color = colornames.Darkgray

				// Attack in direction player last moved
				switch player.LastDir {
				case DirectionXPos:
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y))
					playerSword.Push(pixel.V(player.Last.X+(player.SwordSize*2), player.Last.Y+player.SwordSize))
				case DirectionXNeg:
					playerSword.Push(pixel.V(player.Last.X-player.SwordSize, player.Last.Y))
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+player.SwordSize))
				case DirectionYPos:
					playerSword.Push(pixel.V(player.Last.X, player.Last.Y+player.SwordSize))
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+(player.SwordSize*2)))
				case DirectionYNeg:
					playerSword.Push(pixel.V(player.Last.X, player.Last.Y-player.SwordSize))
					playerSword.Push(pixel.V(player.Last.X+player.SwordSize, player.Last.Y+player.SwordSize))
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

func normalize(target, position pixel.Vec) pixel.Vec {
	return pixel.V(target.X-position.X, target.Y-position.Y)
}
