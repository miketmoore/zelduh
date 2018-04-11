package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
)

// Bomb is an entity
type Bomb struct {
	Category categories.Category
	*components.Spatial
	*components.Animation
}

// BuildBomb builds a bomb entity
func BuildBomb(w, h, x, y float64) Bomb {
	return Bomb{
		Category: categories.Bomb,
		Spatial: &components.Spatial{
			Width:        w,
			Height:       h,
			Rect:         pixel.R(x, y, x+w, y+h),
			Shape:        imdraw.New(nil),
			HitBox:       imdraw.New(nil),
			HitBoxRadius: 50,
		},
	}
}
