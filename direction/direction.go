package direction

import (
	"math/rand"
	"time"
)

// Name is the type of direction
type Name string

const (
	// Up indicates oriented up
	Up Name = "up"
	// Right indicates oriented right
	Right Name = "right"
	// Down indicates oriented down
	Down Name = "down"
	// Left indicates oriented left
	Left Name = "left"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// Rand returns a random direction
func Rand() Name {
	i := r.Intn(4)
	switch i {
	case 0:
		return Up
	case 1:
		return Right
	case 2:
		return Down
	case 3:
		return Left
	default:
		return Up
	}
}
