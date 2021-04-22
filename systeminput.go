package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

type inputEntity struct {
	*ComponentMovement
	*ComponentIgnore
	*ComponentDash
}

// InputSystem is a custom system for detecting collisions and what to do when they occur
type InputSystem struct {
	Win           *pixelgl.Window
	playerEntity  inputEntity
	playerEnabled bool
	sword         inputEntity
	arrow         inputEntity
}

func NewInputSystem(window *pixelgl.Window) InputSystem {
	return InputSystem{Win: window}
}

// DisablePlayer disables player input
func (s *InputSystem) DisablePlayer() {
	s.playerEnabled = false
}

// EnablePlayer enables player input
func (s *InputSystem) EnablePlayer() {
	s.playerEnabled = true
}

// AddEntity adds an entity to the system
func (s *InputSystem) AddEntity(entity Entity) {
	r := inputEntity{
		ComponentMovement: entity.ComponentMovement,
		ComponentDash:     entity.ComponentDash,
		ComponentIgnore:   entity.ComponentIgnore,
	}
	switch entity.Category {
	case CategoryPlayer:
		s.playerEntity = r
	case CategorySword:
		s.sword = r
	case CategoryArrow:
		s.arrow = r
	}
}

// Update checks for player input
func (s InputSystem) Update() {
	if !s.playerEnabled {
		return
	}

	win := s.Win
	player := s.playerEntity

	movingSpeed := player.ComponentMovement.MaxSpeed

	player.ComponentMovement.LastDirection = player.ComponentMovement.Direction
	if win.Pressed(pixelgl.KeyUp) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionUp
	} else if win.Pressed(pixelgl.KeyRight) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionRight
	} else if win.Pressed(pixelgl.KeyDown) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionDown
	} else if win.Pressed(pixelgl.KeyLeft) {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionLeft
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
