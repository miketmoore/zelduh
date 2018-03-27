package input

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type inputEntity struct {
	*components.MovementComponent
	*components.Ignore
}

// System is a custom system for detecting collisions and what to do when they occur
type System struct {
	Win          *pixelgl.Window
	playerEntity inputEntity
	sword        inputEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(movement *components.MovementComponent) {
	s.playerEntity = inputEntity{
		MovementComponent: movement,
	}
}

// AddSword adds the sword entity to the sytem
func (s *System) AddSword(movement *components.MovementComponent, ignore *components.Ignore) {
	s.sword = inputEntity{
		MovementComponent: movement,
		Ignore:            ignore,
	}
}

// Update checks for player input
func (s *System) Update() {
	win := s.Win
	player := s.playerEntity

	movingSpeed := 3.0

	if win.Pressed(pixelgl.KeyUp) {
		player.MovementComponent.Speed = movingSpeed
		player.MovementComponent.Direction = direction.Up
	} else if win.Pressed(pixelgl.KeyRight) {
		player.MovementComponent.Speed = movingSpeed
		player.MovementComponent.Direction = direction.Right
	} else if win.Pressed(pixelgl.KeyDown) {
		player.MovementComponent.Speed = movingSpeed
		player.MovementComponent.Direction = direction.Down
	} else if win.Pressed(pixelgl.KeyLeft) {
		player.MovementComponent.Speed = movingSpeed
		player.MovementComponent.Direction = direction.Left
	} else {
		player.MovementComponent.Speed = 0
	}

	s.sword.MovementComponent.Direction = player.MovementComponent.Direction
	if win.Pressed(pixelgl.KeySpace) {
		s.sword.MovementComponent.Speed = 1.0
		s.sword.Ignore.Value = false
	} else {
		s.sword.MovementComponent.Speed = 0
		s.sword.Ignore.Value = true
	}
}
