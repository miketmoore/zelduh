package zelduh

// componentMovement contains data about movement
type componentMovement struct {
	LastDirection       Direction
	Direction           Direction
	MaxSpeed            float64
	Speed               float64
	MaxMoves            int
	RemainingMoves      int
	HitSpeed            float64
	MovingFromHit       bool
	HitBackMoves        int
	MovementPatternName string
}

func NewComponentMovement(
	direction Direction,
	speed, maxSpeed float64,
	maxMoves, remainingMoves int,
	hitSpeed float64,
	movingFromHit bool,
	hitBackMoves int,
	patternName string,

) *componentMovement {
	return &componentMovement{
		Direction:           direction,
		MaxSpeed:            maxSpeed,
		Speed:               speed,
		MaxMoves:            maxMoves,
		RemainingMoves:      remainingMoves,
		HitSpeed:            hitSpeed,
		MovingFromHit:       movingFromHit,
		HitBackMoves:        hitBackMoves,
		MovementPatternName: patternName,
	}
}

type movementEntity struct {
	ID EntityID
	*componentMovement
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
		componentMovement: entity.componentMovement,
	}
	s.entityByID[e.ID] = e
	// switch entity.Category {
	// case CategoryPlayer:
	// 	// r.componentDash = entity.componentDash
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
	// 	switch enemy.componentMovement.MovementPatternName {
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
		entity.componentMovement.Speed = entity.componentMovement.MaxSpeed
	}
}

func (s *MovementSystem) SetZeroSpeed(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentMovement.Speed = 0
	}
}

func (s *MovementSystem) ChangeSpeed(entityID EntityID, speed float64) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentMovement.Speed = speed
	}
}

func (s *MovementSystem) ChangeDirection(entityID EntityID, direction Direction) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentMovement.Direction = direction
	}
}

func (s *MovementSystem) MatchDirectionToPlayer(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	player, playerOk := s.entityByID[s.playerID]
	if ok && playerOk {
		entity.componentMovement.Direction = player.componentMovement.Direction
	}
}

func (s *MovementSystem) SetRemainingMoves(entityID EntityID, value int) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentMovement.RemainingMoves = value
	}
}

func (s *MovementSystem) DecrementRemainingMoves(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentMovement.RemainingMoves--
	}
}

func (s *MovementSystem) RemainingMoves(entityID EntityID) int {
	entity, ok := s.entityByID[entityID]
	if ok {
		return entity.componentMovement.RemainingMoves
	}
	return 0
}

func (s *MovementSystem) UpdateLastDirection(entityID EntityID) {
	entity, ok := s.entityByID[entityID]
	if ok {
		entity.componentMovement.LastDirection = entity.componentMovement.Direction
	}
}
