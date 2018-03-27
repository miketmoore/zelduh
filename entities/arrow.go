package entities

import "github.com/miketmoore/zelduh/components"

// Arrow is an entity
type Arrow struct {
	*components.Ignore
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
}
