package zelduh

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
	ID EntityID
	*componentTemporary
}

type TemporarySystem struct {
	entityByID map[EntityID]temporaryEntity
	*componentIgnore
}

func NewTemporarySystem() TemporarySystem {
	return TemporarySystem{
		entityByID: map[EntityID]temporaryEntity{},
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

func (s *TemporarySystem) SetExpiration(entityID EntityID, value int, onExpirationHandler func()) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentTemporary.Expiration = value
		entity.componentTemporary.OnExpiration = onExpirationHandler
	}
}

func (s *TemporarySystem) IsTemporary(entityID EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return entity.componentTemporary != nil
	}
	return false
}

func (s *TemporarySystem) IsExpired(entityID EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return s.IsTemporary(entityID) && entity.componentTemporary.Expiration == 0
	}
	return false
}

func (s *TemporarySystem) DecrementExpiration(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentTemporary.Expiration--
	}
}

func (s *TemporarySystem) CallOnExpirationHandler(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentTemporary.OnExpiration()
	}

}
