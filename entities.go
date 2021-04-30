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

	width := c.Dimensions.Width
	height := c.Dimensions.Height

	x := c.Coordinates.X
	y := c.Coordinates.Y

	spatialRectangle := pixel.R(x, y, x+width, y+height)

	spatialComponent := &ComponentSpatial{
		Width:  width,
		Height: height,
		Rect:   spatialRectangle,
		Shape:  imdraw.New(nil),
		HitBox: imdraw.New(nil),
		Color:  c.Color,
	}

	entity := Entity{
		id:               id,
		Category:         c.Category,
		ComponentSpatial: spatialComponent,
		ComponentIgnore: &ComponentIgnore{
			Value: c.Ignore,
		},
	}

	if c.Expiration > 0 {
		entity.ComponentTemporary = NewComponentTemporary(c.Expiration)
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
			ComponentAnimationByName: ComponentAnimationMap{},
		}
		for key, val := range c.Animation {
			entity.ComponentAnimation.ComponentAnimationByName[key] = &ComponentAnimationData{
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
