package systems

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type inputEntity struct {
	*components.Movement
	*components.Ignore
	*components.Dash
}

// Input is a custom system for detecting collisions and what to do when they occur
type Input struct {
	Win           *pixelgl.Window
	playerEntity  inputEntity
	playerEnabled bool
	sword         inputEntity
	arrow         inputEntity
}

// AddPlayer adds the player to the system
func (s *Input) AddPlayer(movement *components.Movement, dash *components.Dash) {
	s.playerEntity = inputEntity{
		Movement: movement,
		Dash:     dash,
	}
	s.playerEnabled = true
}

// DisablePlayer disables player input
func (s *Input) DisablePlayer() {
	s.playerEnabled = false
}

// EnablePlayer enables player input
func (s *Input) EnablePlayer() {
	s.playerEnabled = true
}

// AddSword adds the sword entity to the sytem
func (s *Input) AddSword(movement *components.Movement, ignore *components.Ignore) {
	s.sword = inputEntity{
		Movement: movement,
		Ignore:   ignore,
	}
}

// AddArrow adds the arrow entity to the sytem
func (s *Input) AddArrow(movement *components.Movement, ignore *components.Ignore) {
	s.arrow = inputEntity{
		Movement: movement,
		Ignore:   ignore,
	}
}

// Update checks for player input
func (s *Input) Update() {
	if !s.playerEnabled {
		return
	}

	win := s.Win
	player := s.playerEntity

	movingSpeed := player.Movement.MaxSpeed

	player.Movement.LastDirection = player.Movement.Direction
	if win.Pressed(pixelgl.KeyUp) {
		player.Movement.Speed = movingSpeed
		player.Movement.Direction = direction.Up
	} else if win.Pressed(pixelgl.KeyRight) {
		player.Movement.Speed = movingSpeed
		player.Movement.Direction = direction.Right
	} else if win.Pressed(pixelgl.KeyDown) {
		player.Movement.Speed = movingSpeed
		player.Movement.Direction = direction.Down
	} else if win.Pressed(pixelgl.KeyLeft) {
		player.Movement.Speed = movingSpeed
		player.Movement.Direction = direction.Left
	} else {
		player.Movement.Speed = 0
	}

	// attack with sword
	s.sword.Movement.Direction = player.Movement.Direction
	if win.Pressed(pixelgl.KeyF) {
		s.sword.Movement.Speed = 1.0
		s.sword.Ignore.Value = false
	} else {
		s.sword.Movement.Speed = 0
		s.sword.Ignore.Value = true
	}

	// fire arrow
	if s.arrow.Movement.RemainingMoves == 0 {
		s.arrow.Movement.Direction = player.Movement.Direction
		if win.Pressed(pixelgl.KeyG) {
			s.arrow.Movement.Speed = 7.0
			s.arrow.Movement.RemainingMoves = 100
			s.arrow.Ignore.Value = false
		} else {
			s.arrow.Movement.Speed = 0
			s.arrow.Movement.RemainingMoves = 0
			s.arrow.Ignore.Value = true
		}
	} else {
		s.arrow.Movement.RemainingMoves--
	}

	// dashing
	if !win.Pressed(pixelgl.KeyF) && win.Pressed(pixelgl.KeySpace) {
		if s.playerEntity.Dash.Charge < s.playerEntity.Dash.MaxCharge {
			s.playerEntity.Dash.Charge++
			s.sword.Movement.Speed = 0
			s.sword.Ignore.Value = true
		} else {
			s.sword.Movement.Speed = 1.0
			s.sword.Ignore.Value = false
		}
	} else {
		s.playerEntity.Dash.Charge = 0
	}
}
