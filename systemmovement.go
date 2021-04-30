package zelduh

// ComponentMovement contains data about movement
type ComponentMovement struct {
	LastDirection  Direction
	Direction      Direction
	MaxSpeed       float64
	Speed          float64
	MaxMoves       int
	RemainingMoves int
	HitSpeed       float64
	MovingFromHit  bool
	HitBackMoves   int
	PatternName    string
}

type movementEntity struct {
	ID EntityID
	*ComponentMovement
}

type MovementSystem struct {
	entityByID map[EntityID]movementEntity
	playerID   EntityID
}

func NewMovementSystem(playerID EntityID) MovementSystem {
	return MovementSystem{
		entityByID: map[EntityID]movementEntity{},
		playerID:   playerID,
	}
}

// AddEntity adds an entity to the system
func (s *MovementSystem) AddEntity(entity Entity) {
	e := movementEntity{
		ID:                entity.ID(),
		ComponentMovement: entity.ComponentMovement,
	}
	s.entityByID[e.ID] = e
	// switch entity.Category {
	// case CategoryPlayer:
	// 	// r.ComponentDash = entity.ComponentDash
	// 	// s.player = r
	// 	// case CategorySword:
	// 	// 	s.sword = r
	// 	// case CategoryArrow:
	// 	// 	s.arrow = r
	// 	// case CategoryMovableObstacle:
	// 	// 	s.moveableObstacles = append(s.moveableObstacles, &r)
	// 	// case CategoryEnemy:
	// 	// 	s.enemies = append(s.enemies, &r)
	// }
}

// Update
func (s *MovementSystem) Update() error {
	// s.movePlayer()
	// s.moveSword()
	// s.moveArrow()

	// for i := 0; i < len(s.moveableObstacles); i++ {
	// 	entity := s.moveableObstacles[i]
	// 	s.moveMoveableObstacle(entity)
	// }

	// for i := 0; i < len(s.enemies); i++ {
	// 	enemy := s.enemies[i]
	// 	switch enemy.ComponentMovement.PatternName {
	// 	case "random":
	// 		s.moveEnemyRandom(enemy)
	// 	case "left-right":
	// 		s.moveEnemyLeftRight(enemy)
	// 	}
	// }
	return nil
}

func (s *MovementSystem) SetMaxSpeed(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentMovement.Speed = entity.ComponentMovement.MaxSpeed
	}
}

func (s *MovementSystem) SetZeroSpeed(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentMovement.Speed = 0
	}
}

func (s *MovementSystem) ChangeSpeed(entityID EntityID, speed float64) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentMovement.Speed = speed
	}
}

func (s *MovementSystem) ChangeDirection(entityID EntityID, direction Direction) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentMovement.Direction = direction
	}
}

func (s *MovementSystem) MatchDirectionToPlayer(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	player, playerOk := s.entityByID[s.playerID]
	if ok && playerOk {
		entity.ComponentMovement.Direction = player.ComponentMovement.Direction
	}
}

func (s *MovementSystem) SetRemainingMoves(entityID EntityID, value int) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentMovement.RemainingMoves = value
	}
}

func (s *MovementSystem) DecrementRemainingMoves(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.ComponentMovement.RemainingMoves--
	}
}
