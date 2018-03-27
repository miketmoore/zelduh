package components

import "github.com/miketmoore/zelduh/direction"

// MovementComponent contains data about movement
type MovementComponent struct {
	LastDirection direction.Name
	Direction     direction.Name
	Speed         float64
	MoveCount     int
}
