package zelduh

import "github.com/miketmoore/zelduh/core/entity"

// componentTemporary is used to track when an entity should be removed
type componentTemporary struct {
	Expiration   int
	OnExpiration func()
}

func NewComponentTemporary(expiration int) *componentTemporary {
	return &componentTemporary{
		Expiration: expiration,
	}
}

type temporaryEntity struct {
	ID entity.EntityID
	*componentTemporary
}

type TemporarySystem struct {
	entityByID map[entity.EntityID]temporaryEntity
	*componentIgnore
}

func NewTemporarySystem() TemporarySystem {
	return TemporarySystem{
		entityByID: map[entity.EntityID]temporaryEntity{},
	}
}

// AddEntity adds an entity to the system
func (s *TemporarySystem) AddEntity(entity Entity) {
	s.entityByID[entity.ID()] = temporaryEntity{
		ID:                 entity.ID(),
		componentTemporary: entity.componentTemporary,
	}
}

// Update
func (s *TemporarySystem) Update() error {

	return nil
}

func (s *TemporarySystem) SetExpiration(entityID entity.EntityID, value int, onExpirationHandler func()) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentTemporary.Expiration = value
		entity.componentTemporary.OnExpiration = onExpirationHandler
	}
}

func (s *TemporarySystem) IsTemporary(entityID entity.EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return entity.componentTemporary != nil
	}
	return false
}

func (s *TemporarySystem) IsExpired(entityID entity.EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return s.IsTemporary(entityID) && entity.componentTemporary.Expiration == 0
	}
	return false
}

func (s *TemporarySystem) DecrementExpiration(entityID entity.EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentTemporary.Expiration--
	}
}

func (s *TemporarySystem) CallOnExpirationHandler(entityID entity.EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentTemporary.OnExpiration()
	}

}
