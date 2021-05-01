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
	*componentInvincible
	*componentHitbox
	*componentAnimation
	*componentColor
	*componentCoins
	*componentDash
	*componentEnabled
	*componentToggler
	*componentHealth
	*componentIgnore
	*componentMovement
	*componentSpatial
	*componentTemporary
	*componentRotation
	*componentDimensions
	*componentCoordinates
	*componentRectangle
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
		componentSpatial: NewComponentSpatial(
			c.Color,
		),
		componentIgnore:      NewComponentIgnore(c.Ignore),
		componentCoordinates: NewComponentCoordinates(c.Coordinates.X, c.Coordinates.Y),
		componentDimensions:  NewComponentDimensions(c.Dimensions.Width, c.Dimensions.Height),
		componentRectangle: NewComponentRectangle(
			c.Coordinates.X,
			c.Coordinates.Y,
			c.Dimensions.Width,
			c.Dimensions.Height,
		),
	}

	if c.Expiration > 0 {
		entity.componentTemporary = NewComponentTemporary(c.Expiration)
	}

	if c.Category == CategoryWarp {
		entity.componentEnabled = NewComponentEnabled(true)
	}

	if c.Health > 0 {
		entity.componentHealth = NewComponentHealth(c.Health)
	}

	if c.Hitbox != nil {
		entity.componentHitbox = NewComponentHitbox(c.Hitbox.Radius, float64(c.Hitbox.CollisionWithRectMod))
	}

	if c.Transform != nil {
		// How to rotate pixel.Rect?
		// entity.componentSpatial.

		// https://github.com/faiface/pixel/wiki/Moving,-scaling-and-rotating-with-Matrix#rotation
		// mat := pixel.IM
		// mat = mat.Moved(win.Bounds().Center())
		// mat = mat.Rotated(win.Bounds().Center(), math.Pi/4)
		// sprite.Draw(win, mat)

		// matrix := pixel.IM
		// vector := entity.componentSpatial.Rect.Center()
		// matrix = matrix.Rotated(vector, 90)

		entity.componentRotation = NewComponentRotation(c.Transform.Rotation)

	}

	if c.Toggleable {
		entity.componentToggler = NewComponentToggler(c.Toggled)
	}

	entity.componentInvincible = NewComponentInvincible(c.Invincible)

	if c.Movement != nil {
		entity.componentMovement = NewComponentMovement(
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
		entity.componentCoins = NewComponentCoins(0)
	}

	if c.Dash != nil {
		entity.componentDash = NewComponentDash(
			c.Dash.Charge,
			c.Dash.MaxCharge,
			c.Dash.SpeedMod,
		)
	}

	// An animation is a sprite graphic that may have one or more frames
	// so technically it might not be an animation
	if c.Animation != nil {
		entity.componentAnimation = NewComponentAnimation(c.Animation, frameRate)
	} else {
		// If the "animation" configuration is not set, then this is not a sprite graphic
		// instead, it is a simple shape
		entity.componentColor = NewComponentColor(colornames.Sandybrown)
	}

	return entity
}
