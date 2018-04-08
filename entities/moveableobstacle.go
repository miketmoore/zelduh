package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"golang.org/x/image/colornames"
)

// MoveableObstacle represents an impassable, but moveable object/tile
type MoveableObstacle struct {
	ID       EntityID
	Category categories.Category
	*components.Appearance
	*components.Spatial
	*components.Movement
	*components.Animation
}

// BuildMoveableObstacle builds a new moveable obstacle
func BuildMoveableObstacle(id EntityID, w, h, x, y float64) MoveableObstacle {
	return MoveableObstacle{
		ID:       id,
		Category: categories.MovableObstacle,
		Appearance: &components.Appearance{
			Color: colornames.Purple,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect:   pixel.R(x, y, x+w, y+h),
			Shape:  imdraw.New(nil),
		},
		Movement: &components.Movement{
			Direction: direction.Down,
			Speed:     1.0,
		},
	}
}
