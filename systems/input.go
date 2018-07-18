package systems

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/entities"
)

type inputEntity struct {
	*zelduh.ComponentMovement
	*zelduh.ComponentIgnore
	*zelduh.ComponentDash
}

// Input is a custom system for detecting collisions and what to do when they occur
type Input struct {
	Win           *pixelgl.Window
	playerEntity  inputEntity
	playerEnabled bool
	sword         inputEntity
	arrow         inputEntity
}

// DisablePlayer disables player input
func (s *Input) DisablePlayer() {
	s.playerEnabled = false
}

// EnablePlayer enables player input
func (s *Input) EnablePlayer() {
	s.playerEnabled = true
}

// AddEntity adds an entity to the system
func (s *Input) AddEntity(entity entities.Entity) {
	r := inputEntity{
		ComponentMovement: entity.ComponentMovement,
		ComponentDash:     entity.ComponentDash,
		ComponentIgnore:   entity.ComponentIgnore,
	}
	switch entity.Category {
	case zelduh.CategoryPlayer:
		s.playerEntity = r
	case zelduh.CategorySword:
		s.sword = r
	case zelduh.CategoryArrow:
		s.arrow = r
	}
}

// Update checks for player input
func (s Input) Update() {
	if !s.playerEnabled {
		return
	}

	win := s.Win
	player := s.playerEntity

	movingSpeed := player.ComponentMovement.MaxSpeed

	player.ComponentMovement.LastDirection = player.ComponentMovement.Direction
	if win.Pressed(pixelgl.KeyUp) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = terraform2d.DirectionUp
	} else if win.Pressed(pixelgl.KeyRight) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = terraform2d.DirectionRight
	} else if win.Pressed(pixelgl.KeyDown) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = terraform2d.DirectionDown
	} else if win.Pressed(pixelgl.KeyLeft) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = terraform2d.DirectionLeft
	} else {
		player.ComponentMovement.Speed = 0
	}

	// attack with sword
	s.sword.ComponentMovement.Direction = player.ComponentMovement.Direction
	if win.Pressed(pixelgl.KeyF) {
		s.sword.ComponentMovement.Speed = 1.0
		s.sword.ComponentIgnore.Value = false
	} else {
		s.sword.ComponentMovement.Speed = 0
		s.sword.ComponentIgnore.Value = true
	}

	// fire arrow
	if s.arrow.ComponentMovement.RemainingMoves == 0 {
		s.arrow.ComponentMovement.Direction = player.ComponentMovement.Direction
		if win.Pressed(pixelgl.KeyG) {
			s.arrow.ComponentMovement.Speed = 7.0
			s.arrow.ComponentMovement.RemainingMoves = 100
			s.arrow.ComponentIgnore.Value = false
		} else {
			s.arrow.ComponentMovement.Speed = 0
			s.arrow.ComponentMovement.RemainingMoves = 0
			s.arrow.ComponentIgnore.Value = true
		}
	} else {
		s.arrow.ComponentMovement.RemainingMoves--
	}

	// dashing
	if !win.Pressed(pixelgl.KeyF) && win.Pressed(pixelgl.KeySpace) {
		if s.playerEntity.ComponentDash.Charge < s.playerEntity.ComponentDash.MaxCharge {
			s.playerEntity.ComponentDash.Charge++
			s.sword.ComponentMovement.Speed = 0
			s.sword.ComponentIgnore.Value = true
		} else {
			s.sword.ComponentMovement.Speed = 1.0
			s.sword.ComponentIgnore.Value = false
		}
	} else {
		s.playerEntity.ComponentDash.Charge = 0
	}
}
