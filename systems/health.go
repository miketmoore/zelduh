package systems

import (
	"fmt"

	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/entities"
)

type healthEntity struct {
	ID entities.EntityID
	*components.Health
}

// Health is a custom system for altering character health
type Health struct {
	entities []healthEntity
}

// Add adds the health entity to the system
func (s *Health) Add(id entities.EntityID, health *components.Health) {
	s.entities = append(s.entities, healthEntity{
		ID:     id,
		Health: health,
	})
}

// Hit reduces entity health by d
func (s *Health) Hit(entityID entities.EntityID, d int) bool {
	for i := 0; i < len(s.entities); i++ {
		entity := s.entities[i]
		if entity.ID == entityID {
			entity.Health.Total -= d
			fmt.Printf("Entity health reduced by %d\n", d)
			return entity.Health.Total == 0
		}
	}
	return false
}

// Update is a no-op
func (s *Health) Update() {}
