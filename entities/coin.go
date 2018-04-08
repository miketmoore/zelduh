package entities

import (
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
)

// Coin is a collectible item entity
type Coin struct {
	ID       EntityID
	Category categories.Category
	*components.Appearance
	*components.Spatial
	*components.Animation
}
