package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"golang.org/x/image/colornames"
)

// Sword is an entity
type Sword struct {
	Category categories.Category
	*components.Ignore
	*components.Appearance
	*components.Spatial
	*components.Movement
	*components.Animation
}

// BuildSword builds a sword entity
func BuildSword(w, h float64, dir direction.Name) Sword {
	return Sword{
		Category: categories.Sword,
		Ignore: &components.Ignore{
			Value: true,
		},
		Appearance: &components.Appearance{
			Color: colornames.Deeppink,
		},
		Spatial: &components.Spatial{
			Width:        w,
			Height:       h,
			Rect:         pixel.R(0, 0, 0, 0),
			Shape:        imdraw.New(nil),
			HitBox:       imdraw.New(nil),
			HitBoxRadius: 20,
		},
		Movement: &components.Movement{
			Direction: dir,
			Speed:     0.0,
		},
	}
}
