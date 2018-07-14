package sprites

import (
	"math"

	"github.com/faiface/pixel"
)

// BuildSpritesheet this is a map of pixel engine sprites
func BuildSpritesheet(pic pixel.Picture, s float64) map[int]*pixel.Sprite {
	cols := pic.Bounds().W() / s
	rows := pic.Bounds().H() / s

	maxIndex := (rows * cols) - 1.0

	index := maxIndex
	id := maxIndex + 1
	spritesheet := map[int]*pixel.Sprite{}
	for row := (rows - 1); row >= 0; row-- {
		for col := (cols - 1); col >= 0; col-- {
			x := col
			y := math.Abs(rows-row) - 1
			spritesheet[int(id)] = pixel.NewSprite(pic, pixel.R(
				x*s,
				y*s,
				x*s+s,
				y*s+s,
			))
			index--
			id--
		}
	}
	return spritesheet
}
