package zelduh

type ignoreEntity struct {
	ID EntityID
	*ComponentIgnore
}

type IgnoreSystem struct {
	entityByID map[EntityID]ignoreEntity
	*ComponentIgnore
}

func NewIgnoreSystem() IgnoreSystem {
	return IgnoreSystem{
		entityByID: map[EntityID]ignoreEntity{},
	}
}

// AddEntity adds an entity to the system
func (s *IgnoreSystem) AddEntity(entity Entity) {
	s.entityByID[entity.ID()] = ignoreEntity{
		ID:              entity.ID(),
		ComponentIgnore: entity.ComponentIgnore,
	}
}

// Update
func (s *IgnoreSystem) Update() error {

	return nil
}

func (s *IgnoreSystem) Ignore(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentIgnore.Value = true
	}
}

func (s *IgnoreSystem) DoNotIgnore(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentIgnore.Value = false
	}
}
