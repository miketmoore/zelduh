package zelduh

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
	*componentTemporary
	*componentRotation
	*componentDimensions
	*componentCoordinates
	*componentRectangle
	*componentShape
}

type EntitiesMap map[EntityID]Entity

func NewEntitiesMap() EntitiesMap {
	return EntitiesMap{}
}

// ID returns the entity ID
func (e *Entity) ID() EntityID {
	return e.id
}
