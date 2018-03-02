package pc

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/mvmt"
	"golang.org/x/image/colornames"
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
	Shape *imdraw.IMDraw
	// Win is a pointer to the GUI window
	Win *pixelgl.Window
	// SwordSize is the dimensions of the sword
	SwordSize float64
	// Stride is how many tiles character can move in one "step"
	Stride float64

	Health int
}

// Draw renders the current state of the player character
func (player *Player) Draw() {
	shape := player.Shape
	shape.Clear()
	shape.Color = colornames.White
	shape.Push(pixel.V(player.Last.X, player.Last.Y))
	shape.Push(pixel.V(player.Last.X+player.Size, player.Last.Y+player.Size))
	shape.Rectangle(0)
	shape.Draw(player.Win)
}

// Hit handles what happens when something with attack power gits the character
func (player *Player) Hit(attackPower int) {
	player.Health -= attackPower
}

// IsDead returns a bool indicating if the player is dead or not
func (player *Player) IsDead() bool {
	return player.Health == 0
}
