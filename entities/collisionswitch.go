package entities

import "github.com/miketmoore/zelduh/components"

// CollisionSwitch is a switch that is trigger by collision
type CollisionSwitch struct {
	ID      int
	Enabled bool
	*components.AppearanceComponent
	*components.SpatialComponent
}
