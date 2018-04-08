package entities

import (
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
)

// Obstacle represents an impassable object/tile
type Obstacle struct {
	ID       int
	Category categories.Category
	*components.Appearance
	*components.Spatial
}
