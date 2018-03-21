package entities

import (
	"engo.io/ecs"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
)

// Player is an entity made up of components
type Player struct {
	ecs.BasicEntity
	Win *pixelgl.Window
	*components.AppearanceComponent
	*components.SpatialComponent
	*components.MovementComponent
}
