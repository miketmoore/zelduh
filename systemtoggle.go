package zelduh

import "github.com/miketmoore/zelduh/core/entity"

type toggleEntity struct {
	*componentToggler
}

// ToggleSystem is a custom system for detecting collisions and what to do when they occur
type ToggleSystem struct {
	entityByID map[entity.EntityID]toggleEntity
}

func NewToggleSystem() ToggleSystem {
	return ToggleSystem{
		entityByID: map[entity.EntityID]toggleEntity{},
	}
}

// AddEntity adds an entity to the system
func (s *ToggleSystem) AddEntity(entity Entity) {
	s.entityByID[entity.ID()] = toggleEntity{
		componentToggler: entity.componentToggler,
	}
}

// Update checks for collisions
func (s *ToggleSystem) Update() error {
	return nil
}

func (s *ToggleSystem) Enabled(entityID entity.EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return entity.componentToggler.enabled
	}
	return false
}

func (s *ToggleSystem) Toggle(entityID entity.EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentToggler.Toggle()
	}
}
