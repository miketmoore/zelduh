package equipment

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/entity"
	"github.com/miketmoore/zelduh/mvmt"
	"golang.org/x/image/colornames"
)

type state string

const (
	stateAttacking state = "attacking"
	stateSheathed  state = "sheathed"
)

// Sword represents one sword item
type Sword struct {
	// Size is the dimensions (square)
	Size float64
	// Shape represents the sword shape that is rendered
	Shape *imdraw.IMDraw
	// Last is the last vector
	Last pixel.Vec
	// LastDir is the last direction of the sword
	LastDir mvmt.Direction
	// Win is the GUI window
	Win *pixelgl.Window
	// Sprite is the graphic of the sword
	Sprite *pixel.Sprite

	swordEntity    *entity.Entity
	state          state
	attackingCount int
}

// NewSword returns a new sword
func NewSword(win *pixelgl.Window, size float64, sprite *pixel.Sprite) Sword {
	return Sword{
		Win:            win,
		Size:           size,
		Shape:          imdraw.New(nil),
		Sprite:         sprite,
		state:          stateSheathed,
		attackingCount: 0,
	}
}

// Draw renders the sword in the correct location on the window/map
func (sword *Sword) Draw() {
	if sword.state == stateAttacking {
		if sword.attackingCount < 10 {
			sword.attackingCount++

			sword.Shape.Clear()
			sword.Shape.Color = colornames.White

			sword.Shape.Push(sword.Last)
			sword.Shape.Push(pixel.V(sword.Last.X+sword.Size, sword.Last.Y+sword.Size))

			mat := pixel.IM
			mat = mat.Moved(pixel.V(sword.Last.X+sword.Size/2, sword.Last.Y+sword.Size/2))

			degrees := 0.0
			if sword.LastDir == mvmt.DirectionXPos {
				degrees = math.Pi * 3 / 2
			} else if sword.LastDir == mvmt.DirectionYNeg {
				degrees = math.Pi
			} else if sword.LastDir == mvmt.DirectionXNeg {
				degrees = math.Pi / 2
			}
			mat = mat.Rotated(pixel.V(sword.Last.X+sword.Size/2, sword.Last.Y+sword.Size/2), degrees)
			sword.Sprite.Draw(sword.Win, mat)
		} else {
			sword.attackingCount = 0
			sword.state = stateSheathed
		}
	}

}

// Attack starts the attack sequence for the sword
func (sword *Sword) Attack() {
	sword.state = stateAttacking
}

// IsAttacking indicates if the sword is currently attacking or no
func (sword *Sword) IsAttacking() bool {
	return sword.state == stateAttacking
}
