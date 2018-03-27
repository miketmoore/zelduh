package input

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type inputEntity struct {
	*components.MovementComponent
	*components.Ignore
	*components.Dash
}

// System is a custom system for detecting collisions and what to do when they occur
type System struct {
	Win          *pixelgl.Window
	playerEntity inputEntity
	sword        inputEntity
	arrow        inputEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(movement *components.MovementComponent, dash *components.Dash) {
	s.playerEntity = inputEntity{
		MovementComponent: movement,
		Dash:              dash,
	}
}

// AddSword adds the sword entity to the sytem
func (s *System) AddSword(movement *components.MovementComponent, ignore *components.Ignore) {
	s.sword = inputEntity{
		MovementComponent: movement,
		Ignore:            ignore,
	}
}

// AddArrow adds the arrow entity to the sytem
func (s *System) AddArrow(movement *components.MovementComponent, ignore *components.Ignore) {
	s.arrow = inputEntity{
		MovementComponent: movement,
		Ignore:            ignore,
	}
}

// Update checks for player input
func (s *System) Update() {
	win := s.Win
	player := s.playerEntity

	movingSpeed := 3.0

	player.MovementComponent.LastDirection = player.MovementComponent.Direction
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

	// attack with sword
	s.sword.MovementComponent.Direction = player.MovementComponent.Direction
	if win.Pressed(pixelgl.KeyF) {
		s.sword.MovementComponent.Speed = 1.0
		s.sword.Ignore.Value = false
	} else {
		s.sword.MovementComponent.Speed = 0
		s.sword.Ignore.Value = true
	}

	// fire arrow
	if s.arrow.MovementComponent.MoveCount == 0 {
		s.arrow.MovementComponent.Direction = player.MovementComponent.Direction
		if win.Pressed(pixelgl.KeyG) {
			s.arrow.MovementComponent.Speed = 7.0
			s.arrow.MovementComponent.MoveCount = 100
			s.arrow.Ignore.Value = false
		} else {
			s.arrow.MovementComponent.Speed = 0
			s.arrow.MovementComponent.MoveCount = 0
			s.arrow.Ignore.Value = true
		}
	} else {
		s.arrow.MovementComponent.MoveCount--
	}

	// dashing
	if win.Pressed(pixelgl.KeySpace) {
		if s.playerEntity.Dash.Charge < s.playerEntity.Dash.MaxCharge {
			s.playerEntity.Dash.Charge++
		}
	} else {
		s.playerEntity.Dash.Charge = 0
	}
}
