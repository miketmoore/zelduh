package zelduh

// ComponentCoins contains info about an entity's coins
type ComponentCoins struct {
	Coins int
}

func NewComponentCoins(coins int) *ComponentCoins {
	return &ComponentCoins{
		Coins: coins,
	}
}

// ComponentEnabled is a component for tracking enabled/disabled state of an entity
type ComponentEnabled struct {
	Value bool
}

func NewComponentEnabled(enabled bool) *ComponentEnabled {
	return &ComponentEnabled{
		Value: enabled,
	}
}

// ComponentToggler contains information to use when something is toggled
type ComponentToggler struct {
	enabled bool
}

func NewComponentToggler(toggled bool) *ComponentToggler {
	component := &ComponentToggler{}
	if toggled {
		component.Toggle()
	}
	return component
}

// Enabled determine if the Toggler is enabled or not
func (s *ComponentToggler) Enabled() bool {
	return s.enabled
}

// Toggle handles the switch being toggled
func (s *ComponentToggler) Toggle() {
	s.enabled = !s.enabled
}
