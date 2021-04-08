package zelduh

// TransitionStyle represents a transition style
type TransitionStyle string

const (
	// TransitionSlide represents a slide transition
	TransitionSlide TransitionStyle = "slide"
	// TransitionWarp represents a warp transition
	TransitionWarp TransitionStyle = "warp"
)

// RoomTransition is used to transition between map room
type RoomTransition struct {
	Active bool
	Side   Bound
	Start  float64
	Timer  int
	Style  TransitionStyle
}
