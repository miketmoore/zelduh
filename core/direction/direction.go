package direction

import "math/rand"

// Direction is the type of direction
type Direction string

const (
	// DirectionUp indicates oriented up
	DirectionUp Direction = "up"
	// DirectionRight indicates oriented right
	DirectionRight Direction = "right"
	// DirectionDown indicates oriented down
	DirectionDown Direction = "down"
	// DirectionLeft indicates oriented left
	DirectionLeft Direction = "left"
)

// RandomDirection returns a random direction
func RandomDirection(r *rand.Rand) Direction {
	i := r.Intn(4)
	switch i {
	case 0:
		return DirectionUp
	case 1:
		return DirectionRight
	case 2:
		return DirectionDown
	case 3:
		return DirectionLeft
	default:
		return DirectionUp
	}
}
