package systems

import (
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/entities"
)

type healthEntity struct {
	ID terraform2d.EntityID
	*zelduh.ComponentHealth
}

// Health is a custom system for altering character health
type Health struct {
	entities []healthEntity
}

// AddEntity adds the entity to the system
func (s *Health) AddEntity(entity entities.Entity) {
	s.entities = append(s.entities, healthEntity{
		ID:              entity.ID(),
		ComponentHealth: entity.ComponentHealth,
	})
}

// Hit reduces entity health by d
func (s *Health) Hit(entityID terraform2d.EntityID, d int) bool {
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
func (s *Health) Update() {}
