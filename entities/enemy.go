package entities

import (
	"github.com/miketmoore/zelduh/components"
)

// Enemy is an entity made up of components
type Enemy struct {
	ID int
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
	*components.PhysicsComponent
}
