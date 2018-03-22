package entities

import (
	"engo.io/ecs"
	"github.com/miketmoore/zelduh/components"
)

// Player is an entity made up of components
type Player struct {
	ecs.BasicEntity
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
	*components.CoinsComponent
}
