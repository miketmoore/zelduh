package npc

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/mvmt"
	"golang.org/x/image/colornames"
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
	Shape *imdraw.IMDraw
	// Win is a pointer to the GUI window
	Win *pixelgl.Window
}

func (blob *Blob) Draw(playerLast pixel.Vec) {
	// Move blob with AI
	// velocity := 1
	// max_velocity = 2
	// desired_velocity = normalize(target - position) * max_velocity
	// steering = desired_velocity - velocity
	// velocity := float64(1)    // move one space per update
	blob.Shape.Clear()
	maxVelocity := float64(0.75) // move three spaces per update
	normalized := mvmt.Normalize(playerLast, blob.Last)
	desiredVelocity := pixel.V(normalized.X*maxVelocity, normalized.Y*maxVelocity)
	fmt.Printf("desiredVelocity: %f, %f\n", desiredVelocity.X, desiredVelocity.Y)

	blob.Shape.Color = colornames.Darkblue
	// blob.Push(pixel.V(blobLastX, blobLastY))
	// blob.Push(pixel.V(blobLastX+blobSize, blobLastY+blobSize))
	blob.Shape.Push(desiredVelocity)
	blob.Shape.Push(pixel.V(desiredVelocity.X+blob.Size, desiredVelocity.Y+blob.Size))
	blob.Shape.Rectangle(0)
	blob.Shape.Draw(blob.Win)
}
