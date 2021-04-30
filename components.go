package zelduh

// ComponentCoins contains info about an entity's coins
type ComponentCoins struct {
	Coins int
}

// ComponentEnabled is a component for tracking enabled/disabled state of an entity
type ComponentEnabled struct {
	Value bool
}

// ComponentToggler contains information to use when something is toggled
type ComponentToggler struct {
	enabled bool
}

// Enabled determine if the Toggler is enabled or not
func (s *ComponentToggler) Enabled() bool {
	return s.enabled
}

// Toggle handles the switch being toggled
func (s *ComponentToggler) Toggle() {
	s.enabled = !s.enabled
}
