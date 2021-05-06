package zelduh

import "github.com/miketmoore/zelduh/core/entity"

// componentIgnore determines if an entity is ignored by the game, or not
type componentIgnore struct {
	Value bool
}

func NewComponentIgnore(ignore bool) *componentIgnore {
	return &componentIgnore{
		Value: ignore,
	}
}

type ignoreEntity struct {
	ID entity.EntityID
	*componentIgnore
}

type IgnoreSystem struct {
	entityByID map[entity.EntityID]ignoreEntity
	*componentIgnore
}

func NewIgnoreSystem() IgnoreSystem {
	return IgnoreSystem{
		entityByID: map[entity.EntityID]ignoreEntity{},
	}
}

// AddEntity adds an entity to the system
func (s *IgnoreSystem) AddEntity(entity Entity) {
	s.entityByID[entity.ID()] = ignoreEntity{
		ID:              entity.ID(),
		componentIgnore: entity.componentIgnore,
	}
}

// Update
func (s *IgnoreSystem) Update() error {

	return nil
}

func (s *IgnoreSystem) Ignore(entityID entity.EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentIgnore.Value = true
	}
}

func (s *IgnoreSystem) DoNotIgnore(entityID entity.EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentIgnore.Value = false
	}
}

func (s *IgnoreSystem) IsCurrentlyIgnored(entityID entity.EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return entity.componentIgnore.Value
	}
	return false
}
