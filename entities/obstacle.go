package entities

import "github.com/miketmoore/zelduh/components"

// Obstacle represents an impassable object/tile
type Obstacle struct {
	ID int
	*components.Appearance
	*components.Spatial
}
