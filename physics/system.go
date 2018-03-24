package physics

import (
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
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
	player := s.playerEntity

	// Determine speed by forces
	if player.PhysicsComponent.ForceUp > 0 {
		player.MovementComponent.Speed = 1
		player.MovementComponent.Direction = direction.Up
	} else if player.PhysicsComponent.ForceRight > 0 {
		player.MovementComponent.Speed = 1
		player.MovementComponent.Direction = direction.Right
	} else if player.PhysicsComponent.ForceDown > 0 {
		player.MovementComponent.Speed = 1
		player.MovementComponent.Direction = direction.Down
	} else if player.PhysicsComponent.ForceLeft > 0 {
		player.MovementComponent.Speed = 1
		player.MovementComponent.Direction = direction.Left
	} else {
		player.MovementComponent.Speed = 0
	}
}
