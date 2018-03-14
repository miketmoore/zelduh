package collision

import "github.com/faiface/pixel"

// IsColliding determines if two vectors are colliding
func IsColliding(a, b pixel.Vec, spriteSize float64) bool {
	aBottom := a.Y
	aTop := a.Y + spriteSize
	aLeft := a.X
	aRight := a.X + spriteSize

	bBottom := b.Y
	bTop := b.Y + spriteSize
	bLeft := b.X
	bRight := b.X + spriteSize

	notCollidingWithPlayer := aBottom > bTop ||
		bBottom > aTop ||
		aLeft > bRight ||
		bLeft > aRight || false

	return !notCollidingWithPlayer
}
