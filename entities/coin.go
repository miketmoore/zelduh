package entities

import (
	"engo.io/ecs"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

// Coin is a collectible item entity
type Coin struct {
	ecs.BasicEntity
	Win *pixelgl.Window
	*systems.AppearanceComponent
	*components.SpatialComponent
	*components.EntityTypeComponent
}
