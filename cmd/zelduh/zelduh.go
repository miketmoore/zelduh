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
		Bounds: pixel.R(0, 0, screenW, screenH),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Draw player character
	pc := imdraw.New(nil)
	pc.Color = colornames.White

	var size float64 = 8
	var lastX float64 = screenW - size
	var lastY float64 = screenH - size

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

		// Detect edge of window
		if win.JustPressed(pixelgl.KeyUp) || win.Repeated(pixelgl.KeyUp) {
			if lastY+stride < screenH {
				lastY += stride
				drawPC = true
			}
		} else if win.JustPressed(pixelgl.KeyDown) || win.Repeated(pixelgl.KeyDown) {
			if lastY-stride >= 0 {
				lastY -= stride
				drawPC = true
			}
		} else if win.JustPressed(pixelgl.KeyRight) || win.Repeated(pixelgl.KeyRight) {
			if lastX+stride < screenW {
				lastX += stride
				drawPC = true
			}
		} else if win.JustPressed(pixelgl.KeyLeft) || win.Repeated(pixelgl.KeyLeft) {
			if lastX-stride >= 0 {
				lastX -= stride
				drawPC = true
			}
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
