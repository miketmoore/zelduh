package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"golang.org/x/image/colornames"
)

// Player is an entity made up of components
type Player struct {
	Category categories.Category
	*components.Animation
	*components.Appearance
	*components.Spatial
	*components.Movement
	*components.Coins
	*components.Health
	*components.Dash
}

// BuildPlayer builds a player entity
func BuildPlayer(w, h, x, y float64) Player {
	return Player{
		Category: categories.Player,
		Health: &components.Health{
			Total: 3,
		},
		Appearance: &components.Appearance{
			Color: colornames.Green,
		},
		Spatial: &components.Spatial{
			Width:  w,
			Height: h,
			Rect: pixel.R(
				x,
				y,
				x+w,
				y+h,
			),
			Shape:                imdraw.New(nil),
			HitBox:               imdraw.New(nil),
			HitBoxRadius:         15,
			CollisionWithRectMod: 5,
		},
		Movement: &components.Movement{
			Direction: direction.Down,
			MaxSpeed:  7.0,
			Speed:     0.0,
		},
		Coins: &components.Coins{
			Coins: 0,
		},
		Dash: &components.Dash{
			Charge:    0,
			MaxCharge: 50,
			SpeedMod:  7,
		},
	}
}
