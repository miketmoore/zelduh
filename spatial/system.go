package spatial

import (
	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type spatialEntity struct {
	ID int
	*components.MovementComponent
	*components.SpatialComponent
}

// System is a custom system
type System struct {
	playerEntity spatialEntity
	enemies      []spatialEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(spatial *components.SpatialComponent, movement *components.MovementComponent) {
	s.playerEntity = spatialEntity{
		SpatialComponent:  spatial,
		MovementComponent: movement,
	}
}

// AddEnemy adds an enemy to the system
func (s *System) AddEnemy(id int, spatial *components.SpatialComponent, movement *components.MovementComponent) {
	s.enemies = append(s.enemies, spatialEntity{
		SpatialComponent:  spatial,
		MovementComponent: movement,
	})
}

// Update changes spatial data based on movement data
func (s *System) Update() {
	player := s.playerEntity
	if player.MovementComponent.Moving {
		var v pixel.Vec
		speed := player.MovementComponent.Speed
		switch player.MovementComponent.Direction {
		case direction.Up:
			v = pixel.V(0, speed)
		case direction.Right:
			v = pixel.V(speed, 0)
		case direction.Down:
			v = pixel.V(0, -speed)
		case direction.Left:
			v = pixel.V(-speed, 0)
		}
		newRect := player.SpatialComponent.Rect.Moved(v)
		if player.SpatialComponent.BoundsRect.Contains(newRect.Min) &&
			player.SpatialComponent.BoundsRect.Contains(newRect.Max) {
			player.SpatialComponent.Rect = newRect
		}
	}

}
