package entities

import (
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
)

// EntityID represents an entity ID
type EntityID int

// Entity is used to represent each character and tangable "thing" in the game
type Entity struct {
	ID       EntityID
	Category categories.Category
	*components.Animation
	*components.Appearance
	*components.Coins
	*components.Dash
	*components.Enabled
	*components.Health
	*components.Ignore
	*components.Movement
	*components.Spatial
}
