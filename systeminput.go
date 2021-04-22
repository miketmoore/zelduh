package zelduh

type inputEntity struct {
	*ComponentMovement
	*ComponentIgnore
	*ComponentDash
}

// InputSystem is a custom system for detecting collisions and what to do when they occur
type InputSystem struct {
	Input        Input
	playerEntity inputEntity
	inputEnabled bool
	sword        inputEntity
	arrow        inputEntity
}

func NewInputSystem(input Input) InputSystem {
	return InputSystem{Input: input}
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

type Input interface {
	Up() bool
	Right() bool
	Down() bool
	Left() bool
	PrimaryAttack() bool
	SecondaryAttack() bool
	Combo() bool
}

func (s InputSystem) handleInputMovement() {
	input := s.Input
	player := s.playerEntity
	movingSpeed := player.ComponentMovement.MaxSpeed

	if input.Up() {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionUp
	} else if input.Right() {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionRight
	} else if input.Down() {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionDown
	} else if input.Left() {
		player.ComponentMovement.Speed = movingSpeed
		player.ComponentMovement.Direction = DirectionLeft
	} else {
		player.ComponentMovement.Speed = 0
	}
}

func (s InputSystem) handleInputSword() {
	input := s.Input
	player := s.playerEntity

	// attack with sword
	s.sword.ComponentMovement.Direction = player.ComponentMovement.Direction
	if input.PrimaryAttack() {
		s.sword.ComponentMovement.Speed = 1.0
		s.sword.ComponentIgnore.Value = false
	} else {
		s.sword.ComponentMovement.Speed = 0
		s.sword.ComponentIgnore.Value = true
	}

}

func (s InputSystem) handleInputArrow() {
	input := s.Input
	player := s.playerEntity

	// fire arrow
	if s.arrow.ComponentMovement.RemainingMoves == 0 {
		s.arrow.ComponentMovement.Direction = player.ComponentMovement.Direction
		if input.SecondaryAttack() {
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
}

func (s InputSystem) handleInputDash() {
	input := s.Input

	// dashing
	if !input.PrimaryAttack() && input.Combo() {
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

// Update checks for player input
func (s InputSystem) Update() {
	if !s.inputEnabled {
		return
	}

	s.playerEntity.ComponentMovement.LastDirection = s.playerEntity.ComponentMovement.Direction

	s.handleInputMovement()
	s.handleInputSword()
	s.handleInputArrow()
	s.handleInputDash()

}
