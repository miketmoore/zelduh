package playerinput

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
)

type playerInputEntity struct {
	ID int
	*components.PhysicsComponent
}

// System is a custom system for detecting collisions and what to do when they occur
type System struct {
	Win          *pixelgl.Window
	playerEntity playerInputEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(physics *components.PhysicsComponent) {
	s.playerEntity = playerInputEntity{
		PhysicsComponent: physics,
	}
}

// Update checks for collisions
func (s *System) Update() {
	win := s.Win
	player := s.playerEntity

	player.PhysicsComponent.ForceUp = 0
	player.PhysicsComponent.ForceRight = 0
	player.PhysicsComponent.ForceDown = 0
	player.PhysicsComponent.ForceLeft = 0

	if win.Pressed(pixelgl.KeyUp) {
		player.PhysicsComponent.ForceUp = 1
	} else if win.Pressed(pixelgl.KeyRight) {
		player.PhysicsComponent.ForceRight = 1
	} else if win.Pressed(pixelgl.KeyDown) {
		player.PhysicsComponent.ForceDown = 1
	} else if win.Pressed(pixelgl.KeyLeft) {
		player.PhysicsComponent.ForceLeft = 1
	}
}
