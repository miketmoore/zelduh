package zelduh

type RoomTransitionManager struct {
	transition RoomTransition
}

func NewRoomTransitionManager(
	tileSize float64,
) RoomTransitionManager {
	return RoomTransitionManager{
		transition: NewRoomTransition(tileSize),
	}
}

func (r *RoomTransitionManager) Style() TransitionStyle {
	return r.transition.Style
}

func (r *RoomTransitionManager) Timer() int {
	return r.transition.Timer
}

func (r *RoomTransitionManager) SetTimer(timer int) {
	r.transition.Timer = timer
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

func (r *RoomTransitionManager) SetSide(side Bound) {
	r.transition.Side = side
}

func (r *RoomTransitionManager) Disable() {
	r.transition.Active = false
}

func (r *RoomTransitionManager) Enable() {
	r.transition.Active = true
}

func (r *RoomTransitionManager) Active() bool {
	return r.transition.Active
}

func (r *RoomTransitionManager) SetSlide() {
	r.transition.Style = TransitionSlide
}

func (r *RoomTransitionManager) SetWarp() {
	r.transition.Style = TransitionWarp
}

func (r *RoomTransitionManager) ResetTimer() {
	r.transition.Timer = int(r.transition.Start)
}
