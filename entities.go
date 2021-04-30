package zelduh

import (
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

func NewEntitiesMap() EntitiesMap {
	return EntitiesMap{}
}

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
		ComponentSpatial: NewComponentSpatial(
			c.Coordinates,
			c.Dimensions,
			c.Color,
		),
		ComponentIgnore: NewComponentIgnore(c.Ignore),
	}

	if c.Expiration > 0 {
		entity.ComponentTemporary = NewComponentTemporary(c.Expiration)
	}

	if c.Category == CategoryWarp {
		entity.ComponentEnabled = NewComponentEnabled(true)
	}

	if c.Health > 0 {
		entity.ComponentHealth = NewComponentHealth(c.Health)
	}

	if c.Hitbox != nil {
		entity.ComponentSpatial.HitBoxRadius = c.Hitbox.Radius
	}

	if c.Transform != nil {
		// How to rotate pixel.Rect?
		// entity.ComponentSpatial.

		// https://github.com/faiface/pixel/wiki/Moving,-scaling-and-rotating-with-Matrix#rotation
		// mat := pixel.IM
		// mat = mat.Moved(win.Bounds().Center())
		// mat = mat.Rotated(win.Bounds().Center(), math.Pi/4)
		// sprite.Draw(win, mat)

		// matrix := pixel.IM
		// vector := entity.ComponentSpatial.Rect.Center()
		// matrix = matrix.Rotated(vector, 90)

		entity.ComponentSpatial.Transform = &ComponentSpatialTransform{
			Rotation: c.Transform.Rotation,
		}
	}

	if c.Toggleable {
		entity.ComponentToggler = NewComponentToggler(c.Toggled)
	}

	entity.ComponentInvincible = NewComponentInvincible(c.Invincible)

	if c.Movement != nil {
		entity.ComponentMovement = NewComponentMovement(
			c.Movement.Direction,
			c.Movement.Speed,
			c.Movement.MaxSpeed,
			c.Movement.MaxMoves,
			c.Movement.RemainingMoves,
			c.Movement.HitSpeed,
			c.Movement.MovingFromHit,
			c.Movement.HitBackMoves,
			c.PatternName,
		)
	}

	if c.Coins {
		entity.ComponentCoins = NewComponentCoins(0)
	}

	if c.Dash != nil {
		entity.ComponentDash = NewComponentDash(
			c.Dash.Charge,
			c.Dash.MaxCharge,
			c.Dash.SpeedMod,
		)
	}

	if c.Animation != nil {
		entity.ComponentAnimation = NewComponentAnimation()
		for key, val := range c.Animation {
			entity.ComponentAnimation.ComponentAnimationByName[key] = NewComponentAnimationData(val, frameRate)
		}
	} else {
		entity.ComponentAppearance = NewComponentAppearance(colornames.Sandybrown)
	}

	return entity
}
