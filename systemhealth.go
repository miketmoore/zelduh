package zelduh

type healthEntity struct {
	ID EntityID
	*ComponentHealth
}

// HealthSystem is a custom system for altering character health
type HealthSystem struct {
	entities []healthEntity
}

func NewHealthSystem() HealthSystem {
	return HealthSystem{}
}

// AddEntity adds the entity to the system
func (s *HealthSystem) AddEntity(entity Entity) {
	s.entities = append(s.entities, healthEntity{
		ID:              entity.ID(),
		ComponentHealth: entity.ComponentHealth,
	})
}

// Hit reduces entity health by d
func (s *HealthSystem) Hit(entityID EntityID, d int) bool {
	for i := 0; i < len(s.entities); i++ {
		entity := s.entities[i]
		if entity.ID == entityID {
			entity.ComponentHealth.Total -= d
			return entity.ComponentHealth.Total == 0
		}
	}
	return false
}

// Update is a no-op
func (s *HealthSystem) Update() error {
	return nil
}
