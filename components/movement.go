package components

import "github.com/miketmoore/zelduh/direction"

// MovementComponent contains data about movement
type MovementComponent struct {
	Direction direction.Name
	Speed     float64
}
