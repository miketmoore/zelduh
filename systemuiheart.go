package zelduh

type UIHeartSystem struct {
	systemsManager     *SystemsManager
	player             Entity
	frameRate          int
	entityFactory      EntityFactory
	totalHeartEntities int
}

func NewUIHeartSystem(systemsManager *SystemsManager, entityFactory EntityFactory, frameRate int) UIHeartSystem {
	return UIHeartSystem{
		systemsManager:     systemsManager,
		frameRate:          frameRate,
		entityFactory:      entityFactory,
		totalHeartEntities: 0,
	}
}

// AddEntity adds an entity to the system
func (s *UIHeartSystem) AddEntity(entity Entity) {
	if entity.Category == CategoryPlayer {
		s.player = entity
	}
}

// Update checks for collisions
func (s *UIHeartSystem) Update() error {

	if s.totalHeartEntities != s.player.componentHealth.Total {
		s.totalHeartEntities = 0

		x := 1.5
		y := 14.0

		xDistanceBetweenHearts := 0.65

		for i := 0.0; i < float64(s.player.componentHealth.Total); i++ {

			thisX := x + (xDistanceBetweenHearts * i)

			entity := s.entityFactory.NewEntity("heart", NewCoordinates(thisX, y), s.frameRate)
			s.totalHeartEntities++

			// add heart entity so that it can be rendered
			s.systemsManager.AddEntities(entity)
		}
	}

	return nil
}
