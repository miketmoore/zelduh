package entities

import (
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
)

// Generic is a generic entity
type Generic struct {
	ID       EntityID
	Category categories.Category
	*components.Spatial
	*components.Animation
}
