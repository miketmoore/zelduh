package systems

import (
	"fmt"

	"github.com/miketmoore/zelduh/components"
)

type healthEntity struct {
	ID int
	*components.Health
}

// Health is a custom system for altering character health
type Health struct {
	entities []healthEntity
}

// AddEntity adds the health entity to the system
func (s *Health) AddEntity(id int, health *components.Health) {
	s.entities = append(s.entities, healthEntity{
		ID:     id,
		Health: health,
	})
}

// Hit reduces entity health by d
func (s *Health) Hit(entityID, d int) bool {
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
