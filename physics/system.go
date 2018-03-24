package physics

import (
	"github.com/miketmoore/zelduh/components"
)

type physicsEntity struct {
	ID int
	*components.MovementComponent
	*components.PhysicsComponent
}

// System is a custom system
type System struct {
	playerEntity physicsEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(physics *components.PhysicsComponent, movement *components.MovementComponent) {
	s.playerEntity = physicsEntity{
		PhysicsComponent:  physics,
		MovementComponent: movement,
	}
}

// Update changes spatial data based on movement data
func (s *System) Update() {
	// player := s.playerEntity

	// TODO
}
