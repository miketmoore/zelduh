package entities

import "github.com/miketmoore/zelduh/components"

// Generic is a generic entity
type Generic struct {
	ID EntityID
	*components.Spatial
	*components.Animation
}
