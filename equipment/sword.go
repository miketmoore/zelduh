package equipment

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Sword represents one sword item
type Sword struct {
	// Size is the dimensions (square)
	Size float64
	// Shape represents the sword shape that is rendered
	Shape *imdraw.IMDraw
	// Last is the last vector
	Last pixel.Vec
	// Win is the GUI window
	Win *pixelgl.Window
	// Sprite is the graphic of the sword
	Sprite *pixel.Sprite
}

// NewSword returns a new sword
func NewSword(win *pixelgl.Window, size float64, sprite *pixel.Sprite) Sword {
	return Sword{
		Win:    win,
		Size:   size,
		Shape:  imdraw.New(nil),
		Sprite: sprite,
	}
}

// Draw renders the sword in the correct location on the window/map
func (sword *Sword) Draw() {
	sword.Shape.Clear()
	sword.Shape.Color = colornames.White

	sword.Shape.Push(sword.Last)
	sword.Shape.Push(pixel.V(sword.Last.X+sword.Size, sword.Last.Y+sword.Size))

	matrix := pixel.IM.Moved(pixel.V(sword.Last.X+sword.Size/2, sword.Last.Y+sword.Size/2))
	sword.Sprite.Draw(sword.Win, matrix)

}
