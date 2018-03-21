package components

import "github.com/miketmoore/zelduh/direction"

// MovementComponent contains data about visual appearance
type MovementComponent struct {
	Direction direction.Name
	Moving    bool
}
