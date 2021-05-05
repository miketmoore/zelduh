package zelduh

// TransitionStyle represents a transition style
type TransitionStyle string

const (
	// RoomTransitionSlide represents a slide transition
	RoomTransitionSlide TransitionStyle = "slide"
	// RoomTransitionWarp represents a warp transition
	RoomTransitionWarp TransitionStyle = "warp"
)

// RoomTransition is used to transition between map room
type RoomTransition struct {
	Active bool
	Side   Bound
	Start  float64
	Timer  int
	Style  TransitionStyle
}

func NewRoomTransition(start float64) RoomTransition {
	return RoomTransition{Start: start}
}
