package entities

import (
	"engo.io/ecs"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

// Player is an entity made up of components
type Player struct {
	ecs.BasicEntity
	*systems.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
	*systems.CoinsComponent
	*components.EntityTypeComponent
}
