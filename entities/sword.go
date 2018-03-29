package entities

import "github.com/miketmoore/zelduh/components"

// Sword is an entity
type Sword struct {
	*components.Ignore
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
}