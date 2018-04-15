package world

import (
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/entities"
	"github.com/miketmoore/zelduh/systems"
)

// System is an interface
type System interface {
	Update()
	AddEntity(entities.Entity)
}

// World is a world struct
type World struct {
	systems      []System
	SystemsMap   map[string]System
	lastEntityID entities.EntityID
}

// New returns a new World
func New() World {
	return World{
		lastEntityID: 0,
		SystemsMap:   map[string]System{},
	}
}

// AddSystem adds a System to the World
func (w *World) AddSystem(sys System) {
	w.systems = append(w.systems, sys)
}

// Update executes Update on all systems in this World
func (w *World) Update() {
	for _, sys := range w.systems {
		sys.Update()
	}
}

// Systems returns the systems in this World
func (w *World) Systems() []System {
	return w.systems
}

// NewEntityID generates and returns a new Entity ID
func (w *World) NewEntityID() entities.EntityID {
	w.lastEntityID++
	return w.lastEntityID
}

// Remove removes the specific entity from all systems
func (w *World) Remove(category categories.Category, id entities.EntityID) {
	switch category {
	case categories.Coin:
		for _, sys := range w.systems {
			switch sys := sys.(type) {
			case *systems.Collision:
				sys.Remove(categories.Coin, id)
			case *systems.Render:
				sys.RemoveEntity(id)
			}
		}
	}
}

// RemoveEnemy removes the enemy from all system
func (w *World) RemoveEnemy(id entities.EntityID) {
	for _, sys := range w.systems {
		switch sys := sys.(type) {
		case *systems.Spatial:
			sys.Remove(categories.Enemy, id)
		case *systems.Collision:
			sys.Remove(categories.Enemy, id)
		case *systems.Render:
			sys.RemoveEntity(id)
		}
	}
}
