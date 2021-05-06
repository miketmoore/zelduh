package zelduh

import "github.com/miketmoore/zelduh/core/entity"

// componentHealth contains health data
type componentHealth struct {
	Total int
}

func NewComponentHealth(total int) *componentHealth {
	return &componentHealth{
		Total: total,
	}
}

type healthEntity struct {
	ID entity.EntityID
	*componentHealth
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
		componentHealth: entity.componentHealth,
	})
}

// Hit reduces entity health by d
func (s *HealthSystem) Hit(entityID entity.EntityID, d int) bool {
	for i := 0; i < len(s.entities); i++ {
		entity := s.entities[i]
		if entity.ID == entityID {
			entity.componentHealth.Total -= d
			return entity.componentHealth.Total == 0
		}
	}
	return false
}

func (s *HealthSystem) Health(entityID entity.EntityID) int {
	for i := 0; i < len(s.entities); i++ {
		entity := s.entities[i]
		if entity.ID == entityID {
			return entity.componentHealth.Total
		}
	}
	return 0
}

// Update is a no-op
func (s *HealthSystem) Update() error {
	return nil
}
