package components

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// SpatialComponent contains spatial data
type SpatialComponent struct {
	Width    float64
	Height   float64
	PrevRect pixel.Rect
	Rect     pixel.Rect
	Shape    *imdraw.IMDraw
}
