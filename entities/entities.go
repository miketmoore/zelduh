package entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/config"
	"github.com/miketmoore/zelduh/rooms"
	"golang.org/x/image/colornames"
)

// Entity is used to represent each character and tangable "thing" in the game
type Entity struct {
	id       terraform2d.EntityID
	Category terraform2d.EntityCategory
	*components.Invincible
	*components.Animation
	*components.Appearance
	*components.Coins
	*components.Dash
	*components.Enabled
	*components.Toggler
	*components.Health
	*components.Ignore
	*components.Movement
	*components.Spatial
	*components.Temporary
}

// ID returns the entity ID
func (e *Entity) ID() terraform2d.EntityID {
	return e.id
}

// BuildEntitiesFromConfigs builds and returns a batch of entities
func BuildEntitiesFromConfigs(newEntityID func() terraform2d.EntityID, configs ...rooms.EntityConfig) []Entity {
	batch := []Entity{}
	for _, config := range configs {
		entity := BuildEntityFromConfig(config, newEntityID())
		batch = append(batch, entity)
	}
	return batch
}

// BuildEntityFromConfig builds an entity from a configuration
func BuildEntityFromConfig(c rooms.EntityConfig, id terraform2d.EntityID) Entity {
	entity := Entity{
		id:       id,
		Category: c.Category,
		Spatial: &components.Spatial{
			Width:  c.W,
			Height: c.H,
			Rect:   pixel.R(c.X, c.Y, c.X+c.W, c.Y+c.H),
			Shape:  imdraw.New(nil),
			HitBox: imdraw.New(nil),
		},
		Ignore: &components.Ignore{
			Value: c.Ignore,
		},
	}

	if c.Expiration > 0 {
		entity.Temporary = &components.Temporary{
			Expiration: c.Expiration,
		}
	}

	if c.Category == categories.Warp {
		entity.Enabled = &components.Enabled{
			Value: true,
		}
	}

	if c.Health > 0 {
		entity.Health = &components.Health{
			Total: c.Health,
		}
	}

	if c.Hitbox != nil {
		entity.Spatial.HitBoxRadius = c.Hitbox.Radius
	}

	if c.Toggleable {
		entity.Toggler = &components.Toggler{}
		if c.Toggled {
			entity.Toggler.Toggle()
		}
	}

	if c.Invincible {
		entity.Invincible = &components.Invincible{
			Enabled: true,
		}
	} else {
		entity.Invincible = &components.Invincible{
			Enabled: false,
		}
	}

	if c.Movement != nil {
		entity.Movement = &components.Movement{
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
		entity.Coins = &components.Coins{
			Coins: 0,
		}
	}

	if c.Dash != nil {
		entity.Dash = &components.Dash{
			Charge:    c.Dash.Charge,
			MaxCharge: c.Dash.MaxCharge,
			SpeedMod:  c.Dash.SpeedMod,
		}
	}

	if c.Animation != nil {
		entity.Animation = &components.Animation{
			Map: components.AnimationMap{},
		}
		for key, val := range c.Animation {
			entity.Animation.Map[key] = &components.AnimationData{
				Frames:    val,
				FrameRate: config.FrameRate,
			}
		}
	} else {
		entity.Appearance = &components.Appearance{
			Color: colornames.Sandybrown,
		}
	}

	return entity
}
