package zelduh

type InputHandlers struct {
	OnUp, OnRight, OnDown, OnLeft, OnNoDirection,
	OnPrimaryAttack, OnNoPrimaryAttack,
	OnSecondaryAttack, OnNoSecondaryAttack func()
}

// InputSystem is a custom system for detecting collisions and what to do when they occur
type InputSystem struct {
	Input        Input
	inputEnabled bool
	handlers     InputHandlers
}

// NewInputSystem creates a new InputSystem
func NewInputSystem(
	input Input,
	handlers InputHandlers,
) InputSystem {
	return InputSystem{
		Input:    input,
		handlers: handlers,
	}
}

// Disable disables input
func (s *InputSystem) Disable() {
	s.inputEnabled = false
}

// Enable enables  input
func (s *InputSystem) Enable() {
	s.inputEnabled = true
}

// AddEntity adds an entity to the system
func (s *InputSystem) AddEntity(entity Entity) {}

type Input interface {
	Up() bool
	Right() bool
	Down() bool
	Left() bool
	PrimaryAttack() bool
	SecondaryAttack() bool
	Combo() bool
}

// Update checks for player input
func (s InputSystem) Update() error {
	if !s.inputEnabled {
		return nil
	}

	input := s.Input

	if input.Up() {
		s.handlers.OnUp()
	} else if input.Right() {
		s.handlers.OnRight()
	} else if input.Down() {
		s.handlers.OnDown()
	} else if input.Left() {
		s.handlers.OnLeft()
	} else {
		s.handlers.OnNoDirection()
	}

	if s.Input.PrimaryAttack() {
		s.handlers.OnPrimaryAttack()
	} else {
		s.handlers.OnNoPrimaryAttack()
	}

	if s.Input.SecondaryAttack() {
		s.handlers.OnSecondaryAttack()
	} else {
		s.handlers.OnNoSecondaryAttack()
	}

	// s.handleInputDash()

	return nil
}

// func (s InputSystem) handleInputDash() {
// 	input := s.Input

// 	if input.PrimaryAttack() && input.Combo() {
// 		if s.player.ComponentDash.Charge < s.player.ComponentDash.MaxCharge {
// 			s.player.ComponentDash.Charge++
// 			s.sword.ComponentMovement.Speed = 0
// 			s.sword.ComponentIgnore.Value = true
// 		} else {
// 			s.sword.ComponentMovement.Speed = 1.0
// 			s.sword.ComponentIgnore.Value = false
// 		}
// 	} else {
// 		s.player.ComponentDash.Charge = 0
// 	}
// }
