package world

import (
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/systems"
)

// System is an interface
type System interface {
	Update()
}

// World is a world struct
type World struct {
	systems      []System
	lastEntityID int
}

// New returns a new World
func New() World {
	return World{
		lastEntityID: 0,
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
func (w *World) NewEntityID() int {
	w.lastEntityID++
	return w.lastEntityID
}

// Remove removes the specific entity from all systems
func (w *World) Remove(category categories.Category, id int) {
	switch category {
	case categories.Coin:
		for _, sys := range w.systems {
			switch sys := sys.(type) {
			case *systems.Collision:
				sys.Remove(categories.Coin, id)
			case *systems.Render:
				sys.Remove(categories.Coin, id)
			}
		}
	}
}

// RemoveEnemy removes the enemy from all system
func (w *World) RemoveEnemy(id int) {
	for _, sys := range w.systems {
		switch sys := sys.(type) {
		case *systems.Spatial:
			sys.RemoveEnemy(id)
		case *systems.Collision:
			sys.Remove(categories.Enemy, id)
		case *systems.Render:
			sys.Remove(categories.Enemy, id)
		}
	}
}
