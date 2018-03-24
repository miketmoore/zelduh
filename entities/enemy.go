package entities

import (
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

// Enemy is an entity made up of components
type Enemy struct {
	ID int
	*systems.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
}
