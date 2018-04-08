package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"golang.org/x/image/colornames"
)

// Enemy is an entity made up of components
type Enemy struct {
	ID       int
	Category categories.Category
	*components.Appearance
	*components.Spatial
	*components.Movement
	*components.Health
	*components.Animation
}

// BuildEnemy builds an enemy entity
func BuildEnemy(id int, w, h, x, y, hitRadius float64) Enemy {
	return Enemy{
		ID:       id,
		Category: categories.Enemy,
		Health:   &components.Health{Total: 2},
		Appearance: &components.Appearance{
			Color: colornames.Red,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect: pixel.R(
				x,
				y,
				x+w,
				y+h,
			),
			Shape:                imdraw.New(nil),
			HitBox:               imdraw.New(nil),
			HitBoxRadius:         hitRadius,
			CollisionWithRectMod: 5,
		},
		Movement: &components.Movement{
			Direction: direction.Down,
			Speed:     1.0,
		},
	}
}
