package npc

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/mvmt"
	"github.com/miketmoore/zelduh/palette"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type stateName string

const (
	stateNameAppear stateName = "appear"
	stateNameActive stateName = "active"
	stateNamePause  stateName = "pause"
)

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
	Shape             *imdraw.IMDraw
	Sprites           map[string]*pixel.Sprite
	WalkCycleCountMax int
	WalkCycleCount    int
	WalkCycleFlag     bool
	// Win is a pointer to the GUI window
	Win *pixelgl.Window

	Stride float64

	AttackPower int

	totalMoves         int
	moveCounter        int
	currentState       stateName
	lastState          stateName
	appearDelay        int
	currentAppearDelay int
}

// NewBlob returns a new blob
func NewBlob(win *pixelgl.Window, size, x, y, stride float64, attackPower int, sprites map[string]*pixel.Sprite) Blob {
	blob := Blob{
		Win:                win,
		Size:               size,
		Start:              pixel.V(x, y),
		Last:               pixel.V(0, 0),
		Shape:              imdraw.New(nil),
		Sprites:            sprites,
		WalkCycleCountMax:  7,
		WalkCycleFlag:      false,
		Stride:             stride,
		AttackPower:        attackPower,
		currentState:       stateNameAppear,
		appearDelay:        20,
		currentAppearDelay: 20,
	}
	blob.WalkCycleCount = blob.WalkCycleCountMax
	return blob
}

// Reset updates all values to starting values
func (blob *Blob) Reset() {
	// TODO reset to random start position
	// need limits defined, probably on blob constructor
	blob.Last = blob.Start
	blob.currentState = stateNameAppear
	blob.currentAppearDelay = blob.appearDelay
}

// Draw the blob's current state
func (blob *Blob) Draw(screenW, screenH float64) {
	blob.Shape.Clear()

	switch blob.currentState {
	case stateNameAppear:
		blob.Shape.Color = palette.Map[palette.Lightest]
		if blob.currentAppearDelay > 0 {
			blob.currentAppearDelay--
		} else {
			blob.currentState = stateNameActive
		}
	case stateNameActive:
		blob.Shape.Color = palette.Map[palette.Darkest]
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
					if blob.Last.Y+blob.Size < screenH {
						blob.Last = pixel.V(blob.Last.X, blob.Last.Y+stride)
					}
				case mvmt.DirectionXPos:
					if blob.Last.X+blob.Size < screenW {
						blob.Last = pixel.V(blob.Last.X+stride, blob.Last.Y)
					}
				case mvmt.DirectionYNeg:
					if blob.Last.Y >= 0 {
						blob.Last = pixel.V(blob.Last.X, blob.Last.Y-stride)
					}
				case mvmt.DirectionXNeg:
					if blob.Last.X >= 0 {
						blob.Last = pixel.V(blob.Last.X-stride, blob.Last.Y)
					}
				}
				blob.moveCounter--
			} else {
				blob.totalMoves--
				blob.moveCounter = int(blob.Size)
			}

		}
	}

	blob.Shape.Push(blob.Last)
	blob.Shape.Push(pixel.V(blob.Last.X+blob.Size, blob.Last.Y+blob.Size))
	blob.Shape.Rectangle(0)
	// blob.Shape.Draw(blob.Win)

	// Do some walk cycle "math"
	if blob.WalkCycleCount > 0 {
		blob.WalkCycleCount--
	} else {
		blob.WalkCycleCount = blob.WalkCycleCountMax
		blob.WalkCycleFlag = !blob.WalkCycleFlag
	}

	// Figure out which walk cycle sprite to use
	var sprite *pixel.Sprite

	if blob.LastDir == mvmt.DirectionYNeg {
		if blob.WalkCycleFlag {
			sprite = blob.Sprites["downA"]
		} else {
			sprite = blob.Sprites["downB"]
		}
	} else if blob.LastDir == mvmt.DirectionYPos {
		if blob.WalkCycleFlag {
			sprite = blob.Sprites["upA"]
		} else {
			sprite = blob.Sprites["upB"]
		}
	} else if blob.LastDir == mvmt.DirectionXPos {
		if blob.WalkCycleFlag {
			sprite = blob.Sprites["rightA"]
		} else {
			sprite = blob.Sprites["rightB"]
		}
	} else if blob.LastDir == mvmt.DirectionXNeg {
		if blob.WalkCycleFlag {
			sprite = blob.Sprites["leftA"]
		} else {
			sprite = blob.Sprites["leftB"]
		}
	} else {
		sprite = blob.Sprites["downA"]
	}

	matrix := pixel.IM.Moved(pixel.V(blob.Last.X+blob.Size/2, blob.Last.Y+blob.Size/2))
	sprite.Draw(blob.Win, matrix)
}
