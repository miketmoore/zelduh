package zelduh

// ComponentTemporary is used to track when an entity should be removed
type ComponentTemporary struct {
	Expiration   int
	OnExpiration func()
}

func NewComponentTemporary(expiration int) *ComponentTemporary {
	return &ComponentTemporary{
		Expiration: expiration,
	}
}

type temporaryEntity struct {
	ID EntityID
	*ComponentTemporary
}

type TemporarySystem struct {
	entityByID map[EntityID]temporaryEntity
	*ComponentIgnore
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
		ComponentTemporary: entity.ComponentTemporary,
	}
}

// Update
func (s *TemporarySystem) Update() error {

	return nil
}

func (s *TemporarySystem) SetExpiration(entityID EntityID, value int, onExpirationHandler func()) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentTemporary.Expiration = value
		entity.ComponentTemporary.OnExpiration = onExpirationHandler
	}
}

func (s *TemporarySystem) IsTemporary(entityID EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return entity.ComponentTemporary != nil
	}
	return false
}

func (s *TemporarySystem) IsExpired(entityID EntityID) bool {
	entity, ok := s.entityByID[entityID]
	if ok {
		return s.IsTemporary(entityID) && entity.ComponentTemporary.Expiration == 0
	}
	return false
}

func (s *TemporarySystem) DecrementExpiration(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentTemporary.Expiration--
	}
}

func (s *TemporarySystem) CallOnExpirationHandler(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentTemporary.OnExpiration()
	}

}
