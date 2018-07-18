package world

import (
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/systems"
)

// System is an interface
type System interface {
	Update()
	AddEntity(zelduh.Entity)
}

// World is a world struct
type World struct {
	systems      []System
	SystemsMap   map[string]System
	lastEntityID terraform2d.EntityID
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

// AddSystems adds a batch of systems to the world
func (w *World) AddSystems(all ...System) {
	for _, sys := range all {
		w.AddSystem(sys)
	}
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
func (w *World) NewEntityID() terraform2d.EntityID {
	w.lastEntityID++
	return w.lastEntityID
}

// AddEntity adds the entity to it's system
func (w *World) AddEntity(entity zelduh.Entity) {
	for _, system := range w.Systems() {
		system.AddEntity(entity)
	}
}

// AddEntities adds the terraform2d to their system
func (w *World) AddEntities(all ...zelduh.Entity) {
	for _, entity := range all {
		w.AddEntity(entity)
	}
}

// Remove removes the specific entity from all systems
func (w *World) Remove(category terraform2d.EntityCategory, id terraform2d.EntityID) {
	switch category {
	case zelduh.CategoryCoin:
		for _, sys := range w.systems {
			switch sys := sys.(type) {
			case *systems.Collision:
				sys.Remove(zelduh.CategoryCoin, id)
			case *systems.Render:
				sys.RemoveEntity(id)
			}
		}
	case zelduh.CategoryHeart:
		for _, sys := range w.systems {
			switch sys := sys.(type) {
			case *systems.Render:
				sys.RemoveEntity(id)
			}
		}
	}
}

// RemoveEnemy removes the enemy from all system
func (w *World) RemoveEnemy(id terraform2d.EntityID) {
	for _, sys := range w.systems {
		switch sys := sys.(type) {
		case *systems.Spatial:
			sys.Remove(zelduh.CategoryEnemy, id)
		case *systems.Collision:
			sys.Remove(zelduh.CategoryEnemy, id)
		case *systems.Render:
			sys.RemoveEntity(id)
		}
	}
}

// RemoveAllEnemies removes all enemies from all systems
func (w *World) RemoveAllEnemies() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *systems.Spatial:
			sys.RemoveAll(zelduh.CategoryEnemy)
		case *systems.Collision:
			sys.RemoveAll(zelduh.CategoryEnemy)
		case *systems.Render:
			sys.RemoveAll(zelduh.CategoryEnemy)
		}
	}
}

// RemoveAllEntities removes all entities from systems
func (w *World) RemoveAllEntities() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *systems.Render:
			sys.RemoveAllEntities()
		}
	}
}

// RemoveAllMoveableObstacles removes all moveable obstacles from systems
func (w *World) RemoveAllMoveableObstacles() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
			sys.RemoveAll(zelduh.CategoryMovableObstacle)
		case *systems.Render:
			sys.RemoveAll(zelduh.CategoryMovableObstacle)
		}
	}
}

// RemoveAllCollisionSwitches removes all collision switches from systems
func (w *World) RemoveAllCollisionSwitches() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *systems.Collision:
			sys.RemoveAll(zelduh.CategoryCollisionSwitch)
		}
	}
}
