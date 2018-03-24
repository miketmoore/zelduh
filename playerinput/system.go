package playerinput

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type playerInputEntity struct {
	ID int
	*components.MovementComponent
}

// System is a custom system for detecting collisions and what to do when they occur
type System struct {
	Win          *pixelgl.Window
	playerEntity playerInputEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(movement *components.MovementComponent) {
	s.playerEntity = playerInputEntity{
		MovementComponent: movement,
	}
}

// Update checks for collisions
func (s *System) Update() {
	win := s.Win
	player := s.playerEntity

	player.MovementComponent.Moving = true
	if win.Pressed(pixelgl.KeyUp) {
		player.MovementComponent.Direction = direction.Up
	} else if win.Pressed(pixelgl.KeyRight) {
		player.MovementComponent.Direction = direction.Right
	} else if win.Pressed(pixelgl.KeyDown) {
		player.MovementComponent.Direction = direction.Down
	} else if win.Pressed(pixelgl.KeyLeft) {
		player.MovementComponent.Direction = direction.Left
	} else {
		player.MovementComponent.Moving = false
	}

	fmt.Printf("%v\n", player.MovementComponent)
}
