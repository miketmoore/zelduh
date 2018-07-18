package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/entityconfig"
	"golang.org/x/image/colornames"
)

// Entity is used to represent each character and tangable "thing" in the game
type Entity struct {
	id       terraform2d.EntityID
	Category terraform2d.EntityCategory
	*zelduh.ComponentInvincible
	*zelduh.ComponentAnimation
	*zelduh.ComponentAppearance
	*zelduh.ComponentCoins
	*zelduh.ComponentDash
	*zelduh.ComponentEnabled
	*zelduh.ComponentToggler
	*zelduh.ComponentHealth
	*zelduh.ComponentIgnore
	*zelduh.ComponentMovement
	*zelduh.ComponentSpatial
	*zelduh.ComponentTemporary
}

// ID returns the entity ID
func (e *Entity) ID() terraform2d.EntityID {
	return e.id
}

// BuildEntitiesFromConfigs builds and returns a batch of entities
func BuildEntitiesFromConfigs(newEntityID func() terraform2d.EntityID, configs ...entityconfig.Config) []Entity {
	batch := []Entity{}
	for _, config := range configs {
		entity := BuildEntityFromConfig(config, newEntityID())
		batch = append(batch, entity)
	}
	return batch
}

// BuildEntityFromConfig builds an entity from a configuration
func BuildEntityFromConfig(c entityconfig.Config, id terraform2d.EntityID) Entity {
	entity := Entity{
		id:       id,
		Category: c.Category,
		ComponentSpatial: &zelduh.ComponentSpatial{
			Width:  c.W,
			Height: c.H,
			Rect:   pixel.R(c.X, c.Y, c.X+c.W, c.Y+c.H),
			Shape:  imdraw.New(nil),
			HitBox: imdraw.New(nil),
		},
		ComponentIgnore: &zelduh.ComponentIgnore{
			Value: c.Ignore,
		},
	}

	if c.Expiration > 0 {
		entity.ComponentTemporary = &zelduh.ComponentTemporary{
			Expiration: c.Expiration,
		}
	}

	if c.Category == zelduh.CategoryWarp {
		entity.ComponentEnabled = &zelduh.ComponentEnabled{
			Value: true,
		}
	}

	if c.Health > 0 {
		entity.ComponentHealth = &zelduh.ComponentHealth{
			Total: c.Health,
		}
	}

	if c.Hitbox != nil {
		entity.ComponentSpatial.HitBoxRadius = c.Hitbox.Radius
	}

	if c.Toggleable {
		entity.ComponentToggler = &zelduh.ComponentToggler{}
		if c.Toggled {
			entity.ComponentToggler.Toggle()
		}
	}

	if c.Invincible {
		entity.ComponentInvincible = &zelduh.ComponentInvincible{
			Enabled: true,
		}
	} else {
		entity.ComponentInvincible = &zelduh.ComponentInvincible{
			Enabled: false,
		}
	}

	if c.Movement != nil {
		entity.ComponentMovement = &zelduh.ComponentMovement{
			Direction:      c.Movement.Direction,
			MaxSpeed:       c.Movement.MaxSpeed,
			Speed:          c.Movement.Speed,
			MaxMoves:       c.Movement.MaxMoves,
			RemainingMoves: c.Movement.RemainingMoves,
			HitSpeed:       c.Movement.HitSpeed,
			MovingFromHit:  c.Movement.MovingFromHit,
			HitBackMoves:   c.Movement.HitBackMoves,
			PatternName:    c.Movement.PatternName,
		}
	}

	if c.Coins {
		entity.ComponentCoins = &zelduh.ComponentCoins{
			Coins: 0,
		}
	}

	if c.Dash != nil {
		entity.ComponentDash = &zelduh.ComponentDash{
			Charge:    c.Dash.Charge,
			MaxCharge: c.Dash.MaxCharge,
			SpeedMod:  c.Dash.SpeedMod,
		}
	}

	if c.Animation != nil {
		entity.ComponentAnimation = &zelduh.ComponentAnimation{
			Map: zelduh.ComponentAnimationMap{},
		}
		for key, val := range c.Animation {
			entity.ComponentAnimation.Map[key] = &zelduh.ComponentAnimationData{
				Frames:    val,
				FrameRate: zelduh.FrameRate,
			}
		}
	} else {
		entity.ComponentAppearance = &zelduh.ComponentAppearance{
			Color: colornames.Sandybrown,
		}
	}

	return entity
}
