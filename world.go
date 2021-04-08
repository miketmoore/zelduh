package zelduh

// System is an interface
type System interface {
	Update()
	AddEntity(Entity)
}

// SystemsManager is a world struct
type SystemsManager struct {
	systems      []System
	SystemsMap   map[string]System
	lastEntityID EntityID
}

// NewSystemsManager returns a new SystemsManager
func NewSystemsManager() SystemsManager {
	return SystemsManager{
		lastEntityID: 0,
		SystemsMap:   map[string]System{},
	}
}

// AddSystem adds a System to the SystemsManager
func (w *SystemsManager) AddSystem(sys System) {
	w.systems = append(w.systems, sys)
}

// AddSystems adds a batch of systems to the world
func (w *SystemsManager) AddSystems(all ...System) {
	for _, sys := range all {
		w.AddSystem(sys)
	}
}

// Update executes Update on all systems in this SystemsManager
func (w *SystemsManager) Update() {
	for _, sys := range w.systems {
		sys.Update()
	}
}

// Systems returns the systems in this SystemsManager
func (w *SystemsManager) Systems() []System {
	return w.systems
}

// NewEntityID generates and returns a new Entity ID
func (w *SystemsManager) NewEntityID() EntityID {
	w.lastEntityID++
	return w.lastEntityID
}

// AddEntity adds the entity to it's system
func (w *SystemsManager) AddEntity(entity Entity) {
	for _, system := range w.Systems() {
		system.AddEntity(entity)
	}
}

// AddEntities adds the entities to their system
func (w *SystemsManager) AddEntities(all ...Entity) {
	for _, entity := range all {
		w.AddEntity(entity)
	}
}

// Remove removes the specific entity from all systems
func (w *SystemsManager) Remove(category EntityCategory, id EntityID) {
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
func (w *SystemsManager) RemoveEnemy(id EntityID) {
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
func (w *SystemsManager) RemoveAllEnemies() {
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
func (w *SystemsManager) RemoveAllEntities() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *SystemRender:
			sys.RemoveAllEntities()
		}
	}
}

// RemoveAllMoveableObstacles removes all moveable obstacles from systems
func (w *SystemsManager) RemoveAllMoveableObstacles() {
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
func (w *SystemsManager) RemoveAllCollisionSwitches() {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *SystemCollision:
			sys.RemoveAll(CategoryCollisionSwitch)
		}
	}
}
