package entity

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// Entity is a generic (oh no) container for any sprite that can be used to draw it to the screen
// and track it's position
type Entity struct {
	Win    *pixelgl.Window
	Shape  *imdraw.IMDraw
	Size   float64
	Last   pixel.Vec
	Frames []*pixel.Sprite
	// FrameRate is how many frames to render each sprite in the Frames slice
	FrameRate int

	frameRateCount int
	currFrame      int
	totalFrames    int
}

// New returns a new entity
func New(win *pixelgl.Window, size float64, last pixel.Vec, frames []*pixel.Sprite, frameRate int) Entity {
	entity := Entity{
		Win:       win,
		Size:      size,
		Last:      last,
		Frames:    frames,
		FrameRate: frameRate,
		Shape:     imdraw.New(nil),
	}

	entity.frameRateCount = 0
	entity.currFrame = 0
	entity.totalFrames = len(frames)

	return entity
}

// Draw draws the entity's current animation frame
func (entity *Entity) Draw() {
	if entity.frameRateCount < entity.FrameRate {
		entity.frameRateCount++
	} else {
		entity.frameRateCount = 0
		if entity.currFrame < (entity.totalFrames - 1) {
			entity.currFrame++
		} else {
			entity.currFrame = 0
		}
	}

	entity.Shape.Clear()
	entity.Shape.Push(entity.Last)
	entity.Shape.Push(pixel.V(entity.Last.X+entity.Size, entity.Last.Y+entity.Size))
	entity.Shape.Color = colornames.White
	entity.Shape.Rectangle(0)

	matrix := pixel.IM.Moved(pixel.V(entity.Last.X+entity.Size/2, entity.Last.Y+entity.Size/2))
	entity.Frames[entity.currFrame].Draw(entity.Win, matrix)
}
