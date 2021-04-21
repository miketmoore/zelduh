package zelduh

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// ComponentAppearance contains data about visual appearance
type ComponentAppearance struct {
	Color color.RGBA
}

// ComponentCoins contains info about an entity's coins
type ComponentCoins struct {
	Coins int
}

// ComponentDash indicates that an entity can dash
type ComponentDash struct {
	Charge    int
	MaxCharge int
	SpeedMod  float64
}

// ComponentEnabled is a component for tracking enabled/disabled state of an entity
type ComponentEnabled struct {
	Value bool
}

// ComponentToggler contains information to use when something is toggled
type ComponentToggler struct {
	enabled bool
}

// Enabled determine if the Toggler is enabled or not
func (s *ComponentToggler) Enabled() bool {
	return s.enabled
}

// Toggle handles the switch being toggled
func (s *ComponentToggler) Toggle() {
	s.enabled = !s.enabled
}

// ComponentHealth contains health data
type ComponentHealth struct {
	Total int
}

// ComponentIgnore determines if an entity is ignored by the game, or not
type ComponentIgnore struct {
	Value bool
}

// ComponentInvincible is used to track if an enemy is immune to damage of all kinds
type ComponentInvincible struct {
	Enabled bool
}

// ComponentMovement contains data about movement
type ComponentMovement struct {
	LastDirection  Direction
	Direction      Direction
	MaxSpeed       float64
	Speed          float64
	MaxMoves       int
	RemainingMoves int
	HitSpeed       float64
	MovingFromHit  bool
	HitBackMoves   int
	PatternName    string
}

// ComponentSpatial contains spatial data
type ComponentSpatial struct {
	Width                float64
	Height               float64
	PrevRect             pixel.Rect
	Rect                 pixel.Rect
	Shape                *imdraw.IMDraw
	HitBox               *imdraw.IMDraw
	HitBoxRadius         float64
	CollisionWithRectMod float64
}

// ComponentAnimationData contains data about animating one sequence of sprites
type ComponentAnimationData struct {
	Frames         []int
	Frame          int
	FrameRate      int
	FrameRateCount int
}

// ComponentAnimationMap indexes ComponentAnimationData by use/context
type ComponentAnimationMap map[string]*ComponentAnimationData

// ComponentAnimation contains everything necessary to animate basic characters
type ComponentAnimation struct {
	ComponentAnimationByName ComponentAnimationMap
}

// ComponentTemporary is used to track when an entity should be removed
type ComponentTemporary struct {
	Expiration   int
	OnExpiration func()
}
