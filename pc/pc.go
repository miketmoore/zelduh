package pc

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/mvmt"
	"github.com/miketmoore/zelduh/palette"
)

// Player represents the player character
type Player struct {
	// Size is the dimensions (square)
	Size float64
	// Start is the starting vector
	Start pixel.Vec
	// Last is the last vector
	Last pixel.Vec
	// LastDir is the last direction the player was headed in
	LastDir mvmt.Direction
	// Shape is the view
	Shape             *imdraw.IMDraw
	Sprites           map[string]*pixel.Sprite
	WalkCycleCountMax int
	WalkCycleCount    int
	WalkCycleFlag     bool
	// Win is a pointer to the GUI window
	Win *pixelgl.Window
	// SwordSize is the dimensions of the sword
	SwordSize float64
	// Stride is how many tiles character can move in one "step"
	Stride float64

	Health int

	MaxHealth   int
	AttackPower int
	Money       int
}

// New returns a new Player instance
func New(win *pixelgl.Window, size, stride float64, health, maxHealth, attackPower int, sprites map[string]*pixel.Sprite, start pixel.Vec) Player {
	player := Player{
		Win:               win,
		Size:              size,
		Shape:             imdraw.New(nil),
		Sprites:           sprites,
		WalkCycleCountMax: 7,
		WalkCycleFlag:     false,
		SwordSize:         size,
		Health:            health,
		MaxHealth:         maxHealth,
		AttackPower:       attackPower,
		Money:             0,
	}
	player.WalkCycleCount = player.WalkCycleCountMax
	// player.Start = pixel.V((win.Bounds().W()/2.0)-player.Size, (win.Bounds().H()/2.0)-player.Size)
	player.Start = start
	player.Last = player.Start
	player.Stride = stride
	return player
}

// Draw renders the current state of the player character
func (player *Player) Draw() {
	// Create a shape that we won't draw
	shape := player.Shape
	shape.Clear()
	shape.Color = palette.Map[palette.Lightest]
	shape.Push(pixel.V(player.Last.X, player.Last.Y))
	shape.Push(pixel.V(player.Last.X+player.Size, player.Last.Y+player.Size))
	shape.Rectangle(0)

	// Do some walk cycle "math"
	if player.WalkCycleCount > 0 {
		player.WalkCycleCount--
	} else {
		player.WalkCycleCount = player.WalkCycleCountMax
		player.WalkCycleFlag = !player.WalkCycleFlag
	}

	// Figure out which walk cycle sprite to use
	var sprite *pixel.Sprite

	if player.LastDir == mvmt.DirectionYNeg {
		if player.WalkCycleFlag {
			sprite = player.Sprites["downA"]
		} else {
			sprite = player.Sprites["downB"]
		}
	} else if player.LastDir == mvmt.DirectionYPos {
		if player.WalkCycleFlag {
			sprite = player.Sprites["upA"]
		} else {
			sprite = player.Sprites["upB"]
		}
	} else if player.LastDir == mvmt.DirectionXPos {
		if player.WalkCycleFlag {
			sprite = player.Sprites["rightA"]
		} else {
			sprite = player.Sprites["rightB"]
		}
	} else if player.LastDir == mvmt.DirectionXNeg {
		if player.WalkCycleFlag {
			sprite = player.Sprites["leftA"]
		} else {
			sprite = player.Sprites["leftB"]
		}
	} else {
		sprite = player.Sprites["downA"]
	}

	matrix := pixel.IM.Moved(pixel.V(player.Last.X+player.Size/2, player.Last.Y+player.Size/2))
	sprite.Draw(player.Win, matrix)

}

// Hit handles what happens when something with attack power gits the character
func (player *Player) Hit(attackPower int) {
	player.Health -= attackPower
}

// IsDead returns a bool indicating if the player is dead or not
func (player *Player) IsDead() bool {
	return player.Health == 0
}

// Reset updates all values to starting values
func (player *Player) Reset() {
	player.Last = player.Start
	player.Health = player.MaxHealth
}

// Deposit puts coin in your bag!
func (player *Player) Deposit(n int) {
	player.Money += n
	fmt.Printf("Player money: %d (+%d)\n", player.Money, n)
}

// Move moves the player in the specified direction
func (player *Player) Move(dir mvmt.Direction, yPosBound, yNegBound, xPosBound, xNegBound float64) {

	switch dir {
	case mvmt.DirectionYPos:
		if player.Last.Y+player.Stride <= (yPosBound - player.Size) {
			player.Last.Y += player.Stride
			player.LastDir = mvmt.DirectionYPos
		}
	case mvmt.DirectionYNeg:
		if player.Last.Y-player.Stride >= yNegBound {
			player.Last.Y -= player.Stride
			player.LastDir = mvmt.DirectionYNeg
		}
	case mvmt.DirectionXPos:
		if player.Last.X+player.Stride <= (xPosBound - player.Size) {
			player.Last.X += player.Stride
			player.LastDir = mvmt.DirectionXPos
		}
	case mvmt.DirectionXNeg:
		if player.Last.X-player.Stride >= xNegBound {
			player.Last.X -= player.Stride
			player.LastDir = mvmt.DirectionXNeg
		}
	}

}
