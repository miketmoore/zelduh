package entities

import (
	"github.com/miketmoore/zelduh/components"
)

// Coin is a collectible item entity
type Coin struct {
	ID int
	*components.Appearance
	*components.Spatial
	*components.Animation
}
