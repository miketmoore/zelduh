package entities

import (
	"github.com/miketmoore/zelduh/components"
)

// Player is an entity made up of components
type Player struct {
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
	*components.CoinsComponent
	*components.Health
}
