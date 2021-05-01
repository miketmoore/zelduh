package zelduh

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
	ID EntityID
	*componentIgnore
}

type IgnoreSystem struct {
	entityByID map[EntityID]ignoreEntity
	*componentIgnore
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
		componentIgnore: entity.componentIgnore,
	}
}

// Update
func (s *IgnoreSystem) Update() error {

	return nil
}

func (s *IgnoreSystem) Ignore(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentIgnore.Value = true
	}
}

func (s *IgnoreSystem) DoNotIgnore(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentIgnore.Value = false
	}
}
