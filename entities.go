package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// Entity is an interface for implementing concrete "things" in the game. TODO rename...
type Entityer interface {
	ID() EntityID
	Category() EntityCategory
}

// EntityCategory is used to group entities
type EntityCategory uint

// EntityID represents an entity ID
type EntityID int

// Entity is used to represent each character and tangable "thing" in the game
type Entity struct {
	id       EntityID
	Category EntityCategory
	*ComponentInvincible
	*ComponentAnimation
	*ComponentAppearance
	*ComponentCoins
	*ComponentDash
	*ComponentEnabled
	*ComponentToggler
	*ComponentHealth
	*ComponentIgnore
	*ComponentMovement
	*ComponentSpatial
	*ComponentTemporary
}

type EntitiesMap map[EntityID]Entity

// ID returns the entity ID
func (e *Entity) ID() EntityID {
	return e.id
}

// BuildEntitiesFromConfigs builds and returns a batch of entities
func BuildEntitiesFromConfigs(newEntityID func() EntityID, frameRate int, configs ...EntityConfig) []Entity {
	batch := []Entity{}
	for _, config := range configs {
		entity := BuildEntityFromConfig(config, newEntityID(), frameRate)
		batch = append(batch, entity)
	}
	return batch
}

// BuildEntityFromConfig builds an entity from a configuration
func BuildEntityFromConfig(c EntityConfig, id EntityID, frameRate int) Entity {
	entity := Entity{
		id:       id,
		Category: c.Category,
		ComponentSpatial: &ComponentSpatial{
			Width:  c.W,
			Height: c.H,
			Rect:   pixel.R(c.X, c.Y, c.X+c.W, c.Y+c.H),
			Shape:  imdraw.New(nil),
			HitBox: imdraw.New(nil),
		},
		ComponentIgnore: &ComponentIgnore{
			Value: c.Ignore,
		},
	}

	if c.Expiration > 0 {
		entity.ComponentTemporary = &ComponentTemporary{
			Expiration: c.Expiration,
		}
	}

	if c.Category == CategoryWarp {
		entity.ComponentEnabled = &ComponentEnabled{
			Value: true,
		}
	}

	if c.Health > 0 {
		entity.ComponentHealth = &ComponentHealth{
			Total: c.Health,
		}
	}

	if c.Hitbox != nil {
		entity.ComponentSpatial.HitBoxRadius = c.Hitbox.Radius
	}

	if c.Toggleable {
		entity.ComponentToggler = &ComponentToggler{}
		if c.Toggled {
			entity.ComponentToggler.Toggle()
		}
	}

	if c.Invincible {
		entity.ComponentInvincible = &ComponentInvincible{
			Enabled: true,
		}
	} else {
		entity.ComponentInvincible = &ComponentInvincible{
			Enabled: false,
		}
	}

	if c.Movement != nil {
		entity.ComponentMovement = &ComponentMovement{
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
		entity.ComponentCoins = &ComponentCoins{
			Coins: 0,
		}
	}

	if c.Dash != nil {
		entity.ComponentDash = &ComponentDash{
			Charge:    c.Dash.Charge,
			MaxCharge: c.Dash.MaxCharge,
			SpeedMod:  c.Dash.SpeedMod,
		}
	}

	if c.Animation != nil {
		entity.ComponentAnimation = &ComponentAnimation{
			Map: ComponentAnimationMap{},
		}
		for key, val := range c.Animation {
			entity.ComponentAnimation.Map[key] = &ComponentAnimationData{
				Frames:    val,
				FrameRate: frameRate,
			}
		}
	} else {
		entity.ComponentAppearance = &ComponentAppearance{
			Color: colornames.Sandybrown,
		}
	}

	return entity
}
