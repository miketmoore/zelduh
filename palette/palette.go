package palette

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type Key string

const (
	Darkest  Key = "darkest"
	Dark     Key = "dark"
	Light    Key = "light"
	Lightest Key = "lightest"
)

var Map = map[Key]color.RGBA{
	Darkest:  colornames.Black,
	Dark:     colornames.Darkgray,
	Light:    colornames.Lightgray,
	Lightest: colornames.White,
}
