package zelduh

import "github.com/miketmoore/zelduh/core/entity"

type coinsEntity struct {
	*componentCoins
}

// CoinsSystem is a custom system for detecting collisions and what to do when they occur
type CoinsSystem struct {
	entityByID map[entity.EntityID]coinsEntity
}

func NewCoinsSystem() CoinsSystem {
	return CoinsSystem{
		entityByID: map[entity.EntityID]coinsEntity{},
	}
}

// AddEntity adds an entity to the system
func (s *CoinsSystem) AddEntity(entity Entity) {
	s.entityByID[entity.ID()] = coinsEntity{
		componentCoins: entity.componentCoins,
	}
}

// Update checks for collisions
func (s *CoinsSystem) Update() error {
	return nil
}

func (s *CoinsSystem) AddCoins(entityID entity.EntityID, value int) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentCoins.Coins = entity.componentCoins.Coins + value
	}
}
