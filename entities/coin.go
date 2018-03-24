package entities

import (
	"engo.io/ecs"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

// Coin is a collectible item entity
type Coin struct {
	ID int
	ecs.BasicEntity
	*systems.AppearanceComponent
	*components.SpatialComponent
	*components.EntityTypeComponent
}
