package components

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/direction"
)

// Appearance contains data about visual appearance
type Appearance struct {
	Color color.RGBA
}

// Coins contains info about an entity's coins
type Coins struct {
	Coins int
}

// Dash indicates that an entity can dash
type Dash struct {
	Charge    int
	MaxCharge int
	SpeedMod  float64
}

// Enabled is a component for tracking enabled/disabled state of an entity
type Enabled struct {
	Value bool
}

// Health contains health data
type Health struct {
	Total int
}

// Ignore determines if an entity is ignored by the game, or not
type Ignore struct {
	Value bool
}

// Movement contains data about movement
type Movement struct {
	LastDirection direction.Name
	Direction     direction.Name
	MaxSpeed      float64
	Speed         float64
	MoveCount     int
}

// Spatial contains spatial data
type Spatial struct {
	Width        float64
	Height       float64
	PrevRect     pixel.Rect
	Rect         pixel.Rect
	Shape        *imdraw.IMDraw
	HitBox       *imdraw.IMDraw
	HitBoxRadius float64
}

// AnimationData contains data about animating one sequence of sprites
type AnimationData struct {
	Frames         []pixel.Sprite
	Frame          int
	FrameRate      int
	FrameRateCount int
}

// Animation contains everything necessary to animate basic characters (four directions)
type Animation struct {
	Default          *AnimationData
	SwordAttackDown  *AnimationData
	SwordAttackUp    *AnimationData
	SwordAttackRight *AnimationData
	SwordAttackLeft  *AnimationData
	Up               *AnimationData
	Right            *AnimationData
	Down             *AnimationData
	Left             *AnimationData
}
