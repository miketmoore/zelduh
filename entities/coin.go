package entities

import (
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

// Coin is a collectible item entity
type Coin struct {
	ID int
	*systems.AppearanceComponent
	*components.SpatialComponent
}
