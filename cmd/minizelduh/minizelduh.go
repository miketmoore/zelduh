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
		Bounds: pixel.R(0, 0, screenW*3, screenH*3),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Draw player character
	pc := imdraw.New(nil)
	pc.Color = colornames.White

	var lastX float64 = 100
	var lastY float64 = 100
	var size float64 = 32

	pc.Push(pixel.V(lastX, lastY))
	pc.Push(pixel.V(lastX+size, lastY+size))
	pc.Rectangle(0)

	drawPC := false

	var stride float64 = size

	fmt.Fprintln(txt, title)
	for !win.Closed() {

		if win.JustPressed(pixelgl.KeyQ) {
			os.Exit(1)
		}

		win.Clear(colornames.Darkgreen)

		txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))

		// Get mouse position and log to screen
		mpos := win.MousePosition()
		coordDebugTxt.Clear()
		fmt.Fprintln(coordDebugTxt, fmt.Sprintf("%d, %d", int(math.Ceil(mpos.X)), int(math.Ceil(mpos.Y))))
		coordDebugTxt.Draw(win, pixel.IM.Moved(coordDebugTxtOrig))
		pc.Draw(win)
		win.Update()

		if win.JustPressed(pixelgl.KeyUp) {
			lastY += stride
			drawPC = true
		} else if win.JustPressed(pixelgl.KeyDown) {
			lastY -= stride
			drawPC = true
		} else if win.JustPressed(pixelgl.KeyRight) {
			lastX += stride
			drawPC = true
		} else if win.JustPressed(pixelgl.KeyLeft) {
			lastX -= stride
			drawPC = true
		}

		if drawPC {
			pc.Clear()
			pc.Color = colornames.White
			pc.Push(pixel.V(lastX, lastY))
			pc.Push(pixel.V(lastX+size, lastY+size))
			pc.Rectangle(0)
			drawPC = false
		}

	}
}

func main() {
	pixelgl.Run(run)
}
