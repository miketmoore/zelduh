package zelduh

import (
	"github.com/miketmoore/terraform2d"
)

// System is an interface
type System interface {
	Update()
	AddEntity(Entity)
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
func (w *World) AddEntity(entity Entity) {
	for _, system := range w.Systems() {
		system.AddEntity(entity)
	}
}

// AddEntities adds the terraform2d to their system
func (w *World) AddEntities(all ...Entity) {
	for _, entity := range all {
		w.AddEntity(entity)
	}
}

// Remove removes the specific entity from all systems
func (w *World) Remove(category terraform2d.EntityCategory, id terraform2d.EntityID) {
	switch category {
	case CategoryCoin:
		for _, sys := range w.systems {
			switch sys := sys.(type) {
			case *SystemCollision:
				sys.Remove(CategoryCoin, id)
			case *SystemRender:
				sys.RemoveEntity(id)
			}
		}
	case CategoryHeart:
		for _, sys := range w.systems {
			switch sys := sys.(type) {
			case *SystemRender:
				sys.RemoveEntity(id)
			}
		}
	}
}

// RemoveEnemy removes the enemy from all system
func (w *World) RemoveEnemy(id terraform2d.EntityID) {
	for _, sys := range w.systems {
		switch sys := sys.(type) {
		case *SystemSpatial:
			sys.Remove(CategoryEnemy, id)
		case *SystemCollision:
			sys.Remove(CategoryEnemy, id)
		case *SystemRender:
			sys.RemoveEntity(id)
		}
	}
}

// RemoveAllEnemies removes all enemies from all systems
func (w *World) RemoveAllEnemies() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *SystemSpatial:
			sys.RemoveAll(CategoryEnemy)
		case *SystemCollision:
			sys.RemoveAll(CategoryEnemy)
		case *SystemRender:
			sys.RemoveAll(CategoryEnemy)
		}
	}
}

// RemoveAllEntities removes all entities from systems
func (w *World) RemoveAllEntities() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *SystemRender:
			sys.RemoveAllEntities()
		}
	}
}

// RemoveAllMoveableObstacles removes all moveable obstacles from systems
func (w *World) RemoveAllMoveableObstacles() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *SystemCollision:
			sys.RemoveAll(CategoryMovableObstacle)
		case *SystemRender:
			sys.RemoveAll(CategoryMovableObstacle)
		}
	}
}

// RemoveAllCollisionSwitches removes all collision switches from systems
func (w *World) RemoveAllCollisionSwitches() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *SystemCollision:
			sys.RemoveAll(CategoryCollisionSwitch)
		}
	}
}
