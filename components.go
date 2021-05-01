package zelduh

// componentCoins contains info about an entity's coins
type componentCoins struct {
	Coins int
}

func NewComponentCoins(coins int) *componentCoins {
	return &componentCoins{
		Coins: coins,
	}
}

// componentEnabled is a component for tracking enabled/disabled state of an entity
type componentEnabled struct {
	Value bool
}

func NewComponentEnabled(enabled bool) *componentEnabled {
	return &componentEnabled{
		Value: enabled,
	}
}

// componentToggler contains information to use when something is toggled
type componentToggler struct {
	enabled bool
}

func NewComponentToggler(toggled bool) *componentToggler {
	component := &componentToggler{}
	if toggled {
		component.Toggle()
	}
	return component
}

// Enabled determine if the Toggler is enabled or not
func (s *componentToggler) Enabled() bool {
	return s.enabled
}

// Toggle handles the switch being toggled
func (s *componentToggler) Toggle() {
	s.enabled = !s.enabled
}

type componentDimensions struct {
	Width, Height float64
}

func NewComponentDimensions(width, height float64) *componentDimensions {
	return &componentDimensions{width, height}
}

type componentCoordinates struct {
	X, Y float64
}

func NewComponentCoordinates(x, y float64) *componentCoordinates {
	return &componentCoordinates{x, y}
}
