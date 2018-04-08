package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"golang.org/x/image/colornames"
)

// Arrow is an entity
type Arrow struct {
	Category categories.Category
	*components.Ignore
	*components.Appearance
	*components.Spatial
	*components.Movement
	*components.Animation
}

// BuildArrow builds an arrow entity
func BuildArrow(w, h, x, y float64, dir direction.Name) Arrow {
	return Arrow{
		Category: categories.Arrow,
		Ignore: &components.Ignore{
			Value: true,
		},
		Appearance: &components.Appearance{
			Color: colornames.Deeppink,
		},
		Spatial: &components.Spatial{
			Width:        w,
			Height:       h,
			Rect:         pixel.R(x, y, x+w, y+h),
			Shape:        imdraw.New(nil),
			HitBox:       imdraw.New(nil),
			HitBoxRadius: 5,
		},
		Movement: &components.Movement{
			Direction: dir,
			Speed:     0.0,
		},
	}
}
