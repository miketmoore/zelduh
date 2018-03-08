package equipment

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Sword represents one sword item
type Sword struct {
	// Size is the dimensions (square)
	Size float64
	// Shape represents the sword shape that is rendered
	Shape *imdraw.IMDraw
	// Last is the last vector
	Last pixel.Vec
}

// NewSword returns a new sword
func NewSword(size float64) Sword {
	return Sword{
		Size:  size,
		Shape: imdraw.New(nil),
	}
}
