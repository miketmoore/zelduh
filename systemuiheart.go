package zelduh

type uiHeartEntity struct {
	*componentHealth
}

type UIHeartSystem struct {
	systemsManager *SystemsManager
	player         uiHeartEntity
	heartEntities  []Entity
	frameRate      int
	entityFactory  EntityFactory
}

func NewUIHeartSystem(systemsManager *SystemsManager, entityFactory EntityFactory, frameRate int) UIHeartSystem {
	return UIHeartSystem{
		systemsManager: systemsManager,
		heartEntities:  []Entity{},
		frameRate:      frameRate,
		entityFactory:  entityFactory,
	}
}

// AddEntity adds an entity to the system
func (s *UIHeartSystem) AddEntity(entity Entity) {
	if entity.Category == CategoryPlayer {
		s.player = uiHeartEntity{
			componentHealth: entity.componentHealth,
		}
	}
}

// Update checks for collisions
func (s *UIHeartSystem) Update() error {

	if s.player.componentHealth.Total != len(s.heartEntities) {
		s.heartEntities = nil

		x := 1.5
		xMod := 0.65
		y := 14.0
		for i := 0.0; i < float64(s.player.componentHealth.Total); i++ {
			entity := s.entityFactory.NewEntity("heart", NewCoordinates(x+(xMod*i), y), s.frameRate)
			s.heartEntities = append(s.heartEntities, entity)
			s.systemsManager.AddEntities(entity)
		}
	}

	return nil
}
