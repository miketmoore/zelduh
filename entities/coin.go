package entities

import (
	"engo.io/ecs"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
)

// Coin is a collectible item entity
type Coin struct {
	ecs.BasicEntity
	Win *pixelgl.Window
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.EntityTypeComponent
}
