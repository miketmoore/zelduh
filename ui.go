package zelduh

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func DrawCenterText(win *pixelgl.Window, txt *text.Text, s string, c color.RGBA) {
	txt.Clear()
	txt.Color = c
	fmt.Fprintln(txt, s)
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
}

func DrawMapBackground(win *pixelgl.Window, x, y, w, h float64, color color.Color) {
	s := imdraw.New(nil)
	s.Color = color
	s.Push(pixel.V(x, y))
	s.Push(pixel.V(x+w, y+h))
	s.Rectangle(0)
	s.Draw(win)
}

func DrawScreenStart(win *pixelgl.Window, txt *text.Text, currLocaleMsgs map[string]string) {
	win.Clear(colornames.Darkgray)
	DrawMapBackground(win, MapX, MapY, MapW, MapH, colornames.White)
	DrawCenterText(win, txt, currLocaleMsgs["gameTitle"], colornames.Black)
}
