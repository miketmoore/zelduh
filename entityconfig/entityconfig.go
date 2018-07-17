package entityconfig

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
)

// AnimationConfig is a map of animation types to sprite index lists
type AnimationConfig map[string][]int

// MovementConfig is used to configure an entity's Movement component
type MovementConfig struct {
	LastDirection  terraform2d.Direction
	Direction      terraform2d.Direction
	MaxSpeed       float64
	Speed          float64
	MaxMoves       int
	RemainingMoves int
	HitSpeed       float64
	MovingFromHit  bool
	HitBackMoves   int
	PatternName    string
}

// HitboxConfig is used to configure an entity's hitbox
type HitboxConfig struct {
	Box                  *imdraw.IMDraw
	Radius               float64
	CollisionWithRectMod int
}

// DashConfig is used to configure an entity's dash stats
type DashConfig struct {
	Charge, MaxCharge int
	SpeedMod          float64
}

// Config is used to simplify building entities
type Config struct {
	Category                                                      terraform2d.EntityCategory
	Moveable, Animated, Toggleable, Impassable, Invincible, Coins bool
	X, Y, W, H                                                    float64
	SpriteFrames                                                  []int
	WarpToRoomID                                                  terraform2d.RoomID
	PatternName                                                   string
	Toggled                                                       bool
	Health                                                        int
	Expiration                                                    int
	Ignore                                                        bool
	Animation                                                     AnimationConfig
	Hitbox                                                        *HitboxConfig
	Dash                                                          *DashConfig
	Movement                                                      *MovementConfig
}
