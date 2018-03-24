package world

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
