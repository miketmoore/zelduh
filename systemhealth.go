package zelduh

type healthEntity struct {
	ID EntityID
	*ComponentHealth
}

// SystemHealth is a custom system for altering character health
type SystemHealth struct {
	entities []healthEntity
}

// AddEntity adds the entity to the system
func (s *SystemHealth) AddEntity(entity Entity) {
	s.entities = append(s.entities, healthEntity{
		ID:              entity.ID(),
		ComponentHealth: entity.ComponentHealth,
	})
}

// Hit reduces entity health by d
func (s *SystemHealth) Hit(entityID EntityID, d int) bool {
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
func (s *SystemHealth) Update() {}
