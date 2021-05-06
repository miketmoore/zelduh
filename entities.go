package zelduh

import "github.com/miketmoore/zelduh/core/entity"

// Entity is used to represent each character and tangable "thing" in the game
type Entity struct {
	id       entity.EntityID
	Category entity.EntityCategory
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

type EntitiesMap map[entity.EntityID]Entity

func NewEntitiesMap() EntitiesMap {
	return EntitiesMap{}
}

// ID returns the entity ID
func (e *Entity) ID() entity.EntityID {
	return e.id
}
