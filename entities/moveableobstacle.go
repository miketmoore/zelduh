package entities

import "github.com/miketmoore/zelduh/components"

// MoveableObstacle represents an impassable, but moveable object/tile
type MoveableObstacle struct {
	ID int
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
}
