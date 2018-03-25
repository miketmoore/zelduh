package components

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/direction"
)

// SpatialComponent contains spatial data
type SpatialComponent struct {
	Width    float64
	Height   float64
	PrevRect pixel.Rect
	Rect     pixel.Rect
	Shape    *imdraw.IMDraw
	LastDir  direction.Name
}
