package render

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/systems"
)

type renderEntity struct {
	ID int
	*components.SpatialComponent
	*systems.AppearanceComponent
}

// System is a custom system
type System struct {
	Win          *pixelgl.Window
	playerEntity renderEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(appearance *systems.AppearanceComponent, spatial *components.SpatialComponent) {
	s.playerEntity = renderEntity{
		AppearanceComponent: appearance,
		SpatialComponent:    spatial,
	}
}

// Update changes spatial data based on movement data
func (s *System) Update() {
	player := s.playerEntity
	player.Shape.Clear()
	player.Shape.Color = player.AppearanceComponent.Color
	player.Shape.Push(player.SpatialComponent.Rect.Min)
	player.Shape.Push(player.SpatialComponent.Rect.Max)
	player.Shape.Rectangle(0)
	player.Shape.Draw(s.Win)
}
