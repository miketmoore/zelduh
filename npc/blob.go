package npc

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/mvmt"
	"golang.org/x/image/colornames"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// Blob represents one non-player character
type Blob struct {
	// Size is the dimensions (square)
	Size float64
	// Start is the starting vector
	Start pixel.Vec
	// Last is the last vector
	Last pixel.Vec
	// LastDir is the last direction the Blob was headed in
	LastDir mvmt.Direction
	// Shape is the view
	Shape *imdraw.IMDraw
	// Win is a pointer to the GUI window
	Win *pixelgl.Window

	Stride float64

	AttackPower int

	totalMoves int

	moveCounter int
}

// Draw the blob's current state
func (blob *Blob) Draw(screenW, screenH float64) {
	// v1 will just move randomly

	// move smoothly into next tile, then orient for next move
	stride := blob.Stride

	if blob.totalMoves == 0 {
		// number of tiles to move until next change in movement
		maxMoves := 10
		blob.totalMoves = r.Intn(maxMoves)

		// which direction to move
		directionIndex := r.Intn(4)
		switch directionIndex {
		case 0:
			blob.LastDir = mvmt.DirectionYPos
		case 1:
			blob.LastDir = mvmt.DirectionXPos
		case 2:
			blob.LastDir = mvmt.DirectionYNeg
		case 3:
			blob.LastDir = mvmt.DirectionXNeg
		}

	} else {
		if blob.moveCounter > 0 {
			switch blob.LastDir {
			case mvmt.DirectionYPos:
				if blob.Last.Y+stride < screenH {
					blob.Last = pixel.V(blob.Last.X, blob.Last.Y+stride)
				}
			case mvmt.DirectionXPos:
				if blob.Last.X+stride < screenW {
					blob.Last = pixel.V(blob.Last.X+stride, blob.Last.Y)
				}
			case mvmt.DirectionYNeg:
				if blob.Last.Y-stride >= 0 {
					blob.Last = pixel.V(blob.Last.X, blob.Last.Y-stride)
				}
			case mvmt.DirectionXNeg:
				if blob.Last.X-stride >= 0 {
					blob.Last = pixel.V(blob.Last.X-stride, blob.Last.Y)
				}
			}
			blob.moveCounter--
		} else {
			blob.totalMoves--
			blob.moveCounter = int(blob.Size)
		}

	}

	blob.Shape.Clear()
	blob.Shape.Color = colornames.Darkblue
	blob.Shape.Push(blob.Last)
	blob.Shape.Push(pixel.V(blob.Last.X+blob.Size, blob.Last.Y+blob.Size))
	blob.Shape.Rectangle(0)
	blob.Shape.Draw(blob.Win)
}
