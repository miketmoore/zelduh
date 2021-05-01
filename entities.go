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
	*ComponentHitbox
	*ComponentAnimation
	*ComponentColor
	*ComponentCoins
	*ComponentDash
	*ComponentEnabled
	*ComponentToggler
	*ComponentHealth
	*ComponentIgnore
	*ComponentMovement
	*ComponentSpatial
	*ComponentTemporary
	*ComponentRotation
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
		entity.ComponentHitbox = NewComponentHitbox(c.Hitbox.Radius, float64(c.Hitbox.CollisionWithRectMod))
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

		entity.ComponentRotation = NewComponentRotation(c.Transform.Rotation)
		// entity.ComponentSpatial.Transform = &ComponentSpatialTransform{
		// 	Rotation: c.Transform.Rotation,
		// }
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

	// An animation is a sprite graphic that may have one or more frames
	// so technically it might not be an animation
	if c.Animation != nil {
		entity.ComponentAnimation = NewComponentAnimation(c.Animation, frameRate)
	} else {
		// If the "animation" configuration is not set, then this is not a sprite graphic
		// instead, it is a simple shape
		entity.ComponentColor = NewComponentColor(colornames.Sandybrown)
	}

	return entity
}
