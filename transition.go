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

type RoomTransitionManager struct {
	transition *RoomTransition
}

func NewRoomTransitionManager() RoomTransitionManager {
	return RoomTransitionManager{
		transition: &RoomTransition{
			Start: float64(TileSize),
		},
	}
}

func (r *RoomTransitionManager) Active() bool {
	return r.transition.Active
}

func (r *RoomTransitionManager) SetActive(value bool) {
	r.transition.Active = value
}

func (r *RoomTransitionManager) Style() TransitionStyle {
	return r.transition.Style
}

func (r *RoomTransitionManager) Timer() int {
	return r.transition.Timer
}

func (r *RoomTransitionManager) DecrementTimer() {
	r.transition.Timer--
}

func (r *RoomTransitionManager) Start() float64 {
	return r.transition.Start
}

func (r *RoomTransitionManager) Side() Bound {
	return r.transition.Side
}

func (r *RoomTransitionManager) SetSlideStart(side Bound) {
	r.transition.Active = true
	r.transition.Side = side
	r.transition.Style = TransitionSlide
	r.transition.Timer = int(r.transition.Start)
}

func (r *RoomTransitionManager) SetWarp() {
	r.transition.Active = true
	r.transition.Style = TransitionWarp
	r.transition.Timer = 1
}
