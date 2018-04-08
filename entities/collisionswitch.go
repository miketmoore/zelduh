package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"golang.org/x/image/colornames"
)

// CollisionSwitch is a switch that is trigger by collision
type CollisionSwitch struct {
	ID       int
	Category categories.Category
	Enabled  bool
	*components.Appearance
	*components.Spatial
	*components.Animation
}

// BuildCollisionSwitch builds a collision switch entity
func BuildCollisionSwitch(id int, w, h, x, y float64) CollisionSwitch {
	return CollisionSwitch{
		ID:       id,
		Category: categories.CollisionSwitch,
		Enabled:  false,
		Appearance: &components.Appearance{
			Color: colornames.Sandybrown,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect:   pixel.R(x, y, x+w, y+h),
			Shape:  imdraw.New(nil),
		},
	}
}
