package entities

import "github.com/miketmoore/zelduh/components"

// Generic is a generic entity
type Generic struct {
	ID int
	*components.Spatial
	*components.Animation
}
