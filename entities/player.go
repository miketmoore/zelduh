package entities

import (
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

// Player is an entity made up of components
type Player struct {
	*systems.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
	*components.CoinsComponent
}
