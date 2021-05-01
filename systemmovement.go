package zelduh

import (
	"math/rand"

	"github.com/faiface/pixel"
)

// componentDash indicates that an entity can dash
type componentDash struct {
	Charge    int
	MaxCharge int
	SpeedMod  float64
}

func NewComponentDash(
	charge, maxCharge int, speedMod float64,
) *componentDash {
	return &componentDash{
		Charge:    charge,
		MaxCharge: maxCharge,
		SpeedMod:  speedMod,
	}
}

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
	*componentDash
	*componentDimensions
	*componentRectangle
	TotalMoves  int
	MoveCounter int
}

type MovementSystem struct {
	entityByID        map[EntityID]movementEntity
	Rand              *rand.Rand
	player            movementEntity
	sword             movementEntity
	arrow             movementEntity
	enemies           []*movementEntity
	moveableObstacles []*movementEntity
}

func NewMovementSystem(random *rand.Rand) MovementSystem {
	return MovementSystem{
		entityByID: map[EntityID]movementEntity{},
		Rand:       random,
	}
}

// AddEntity adds an entity to the system
func (s *MovementSystem) AddEntity(entity Entity) {
	// e := movementEntity{
	// 	ID:                entity.ID(),
	// 	componentMovement: entity.componentMovement,
	// }
	r := movementEntity{
		ID:                  entity.ID(),
		componentMovement:   entity.componentMovement,
		componentDimensions: entity.componentDimensions,
		componentRectangle:  entity.componentRectangle,
	}
	s.entityByID[entity.ID()] = r
	switch entity.Category {
	case CategoryPlayer:
		r.componentDash = entity.componentDash
		s.player = r
	case CategorySword:
		s.sword = r
	case CategoryArrow:
		s.arrow = r
	case CategoryMovableObstacle:
		s.moveableObstacles = append(s.moveableObstacles, &r)
	case CategoryEnemy:
		s.enemies = append(s.enemies, &r)
	}
}

// Remove removes the entity from the system
func (s *MovementSystem) Remove(category EntityCategory, id EntityID) {
	switch category {
	case CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			enemy := s.enemies[i]
			if enemy.ID == id {
				s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
			}
		}
	}
}

// RemoveAll removes all entities from one category
func (s *MovementSystem) RemoveAll(category EntityCategory) {
	switch category {
	case CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// MovePlayerBack moves the player back
func (s *MovementSystem) MovePlayerBack() {
	player := s.player
	var v pixel.Vec
	switch player.componentMovement.Direction {
	case DirectionUp:
		v = pixel.V(0, -48)
	case DirectionRight:
		v = pixel.V(-48, 0)
	case DirectionDown:
		v = pixel.V(0, 48)
	case DirectionLeft:
		v = pixel.V(48, 0)
	}
	player.componentRectangle.Rect = player.componentRectangle.PrevRect.Moved(v)
	player.componentRectangle.PrevRect = player.componentRectangle.Rect
}

// MoveMoveableObstacle moves a moveable obstacle
func (s *MovementSystem) MoveMoveableObstacle(obstacleID EntityID, dir Direction) bool {
	entity, ok := s.moveableObstacle(obstacleID)
	if ok && !entity.componentMovement.MovingFromHit {
		entity.componentMovement.MovingFromHit = true
		entity.componentMovement.RemainingMoves = entity.componentMovement.MaxMoves
		entity.componentMovement.Direction = dir
		return true
	}
	return false
}

// UndoEnemyRect resets current rect to previous rect
func (s *MovementSystem) UndoEnemyRect(enemyID EntityID) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		enemy.componentRectangle.Rect = enemy.componentRectangle.PrevRect
	}
}

// MoveEnemyBack moves the enemy back
func (s *MovementSystem) MoveEnemyBack(enemyID EntityID, directionHit Direction) {
	enemy, ok := s.enemy(enemyID)
	if ok && !enemy.componentMovement.MovingFromHit {
		enemy.componentMovement.MovingFromHit = true
		enemy.componentMovement.RemainingMoves = enemy.componentMovement.HitBackMoves
		enemy.componentMovement.Direction = directionHit
	}
}

// ComponentRectangle returns the ComponentRectangle for the entity
func (s *MovementSystem) ComponentRectangle(entityID EntityID) (*componentRectangle, bool) {
	for _, entity := range s.enemies {
		if entity.ID == entityID {
			return entity.componentRectangle, true
		}
	}
	return &componentRectangle{}, false
}

// EnemyMovingFromHit indicates if the enemy is moving after being hit
func (s *MovementSystem) EnemyMovingFromHit(enemyID EntityID) bool {
	enemy, ok := s.enemy(enemyID)
	if ok {
		if enemy.ID == enemyID {
			return enemy.componentMovement.MovingFromHit == true
		}
	}
	return false
}

// Update
func (s *MovementSystem) Update() error {
	s.movePlayer()
	s.moveSword()
	s.moveArrow()

	for i := 0; i < len(s.moveableObstacles); i++ {
		entity := s.moveableObstacles[i]
		s.moveMoveableObstacle(entity)
	}

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
	if ok {
		entity.componentMovement.Direction = s.player.componentMovement.Direction
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

func (s *MovementSystem) moveSword() {
	sword := s.sword
	speed := sword.componentMovement.Speed
	w := sword.componentDimensions.Width
	h := sword.componentDimensions.Height
	if speed > 0 {
		sword.componentRectangle.PrevRect = sword.componentRectangle.Rect
		v := delta(sword.componentMovement.Direction, speed+w, speed+h)
		sword.componentRectangle.Rect = s.player.componentRectangle.Rect.Moved(v)
	} else {
		sword.componentRectangle.Rect = s.player.componentRectangle.Rect
	}
}

func (s *MovementSystem) moveArrow() {
	arrow := s.arrow
	speed := arrow.componentMovement.Speed
	if arrow.componentMovement.RemainingMoves > 0 {
		arrow.componentRectangle.PrevRect = arrow.componentRectangle.Rect
		v := delta(arrow.componentMovement.Direction, speed, speed)
		arrow.componentRectangle.Rect = arrow.componentRectangle.Rect.Moved(v)
	} else {
		arrow.componentRectangle.Rect = s.player.componentRectangle.Rect
	}
}

func (s *MovementSystem) movePlayer() {
	player := s.player
	speed := player.componentMovement.Speed
	if player.componentDash.Charge == player.componentDash.MaxCharge {
		speed += player.componentDash.SpeedMod
	}
	if speed > 0 {
		v := delta(player.componentMovement.Direction, speed, speed)
		player.componentRectangle.PrevRect = player.componentRectangle.Rect
		player.componentRectangle.Rect = player.componentRectangle.Rect.Moved(v)
	}
}

func (s *MovementSystem) moveableObstacle(id EntityID) (movementEntity, bool) {
	for _, e := range s.moveableObstacles {
		if e.ID == id {
			return *e, true
		}
	}
	return movementEntity{}, false
}

func (s *MovementSystem) enemy(id EntityID) (movementEntity, bool) {
	for _, e := range s.enemies {
		if e.ID == id {
			return *e, true
		}
	}
	return movementEntity{}, false
}

func delta(dir Direction, modX, modY float64) pixel.Vec {
	switch dir {
	case DirectionUp:
		return pixel.V(0, modY)
	case DirectionRight:
		return pixel.V(modX, 0)
	case DirectionDown:
		return pixel.V(0, -modY)
	case DirectionLeft:
		return pixel.V(-modX, 0)
	default:
		return pixel.V(0, 0)
	}
}

func (s *MovementSystem) moveMoveableObstacle(entity *movementEntity) {
	if entity.componentMovement.RemainingMoves > 0 {
		speed := entity.componentMovement.MaxSpeed
		entity.componentRectangle.PrevRect = entity.componentRectangle.Rect
		moveVec := delta(entity.componentMovement.Direction, speed, speed)
		entity.componentRectangle.Rect = entity.componentRectangle.Rect.Moved(moveVec)
		entity.componentMovement.RemainingMoves--
	} else {
		entity.componentMovement.MovingFromHit = false
		entity.componentMovement.RemainingMoves = 0
	}
}

func (s *MovementSystem) moveEnemyRandom(enemy *movementEntity) {
	if enemy.componentMovement.RemainingMoves == 0 {
		enemy.componentMovement.MovingFromHit = false
		enemy.componentMovement.RemainingMoves = s.Rand.Intn(enemy.componentMovement.MaxMoves)
		enemy.componentMovement.Direction = RandomDirection(s.Rand)
	} else if enemy.componentMovement.RemainingMoves > 0 {
		var speed float64
		if enemy.componentMovement.MovingFromHit {
			speed = enemy.componentMovement.HitSpeed
		} else {
			speed = enemy.componentMovement.MaxSpeed
		}
		enemy.componentRectangle.PrevRect = enemy.componentRectangle.Rect
		moveVec := delta(enemy.componentMovement.Direction, speed, speed)
		enemy.componentRectangle.Rect = enemy.componentRectangle.Rect.Moved(moveVec)
		enemy.componentMovement.RemainingMoves--
	} else {
		enemy.componentMovement.MovingFromHit = false
		enemy.componentMovement.RemainingMoves = int(enemy.componentRectangle.Rect.W())
	}
}

func (s *MovementSystem) moveEnemyLeftRight(enemy *movementEntity) {
	if enemy.componentMovement.RemainingMoves == 0 {
		enemy.componentMovement.MovingFromHit = false
		enemy.componentMovement.RemainingMoves = enemy.componentMovement.MaxMoves
		switch enemy.componentMovement.Direction {
		case DirectionLeft:
			enemy.componentMovement.Direction = DirectionRight
		case DirectionRight:
			enemy.componentMovement.Direction = DirectionLeft
		}
	} else if enemy.componentMovement.RemainingMoves > 0 {
		var speed float64
		if enemy.componentMovement.MovingFromHit {
			speed = enemy.componentMovement.HitSpeed
		} else {
			speed = enemy.componentMovement.MaxSpeed
		}
		enemy.componentRectangle.PrevRect = enemy.componentRectangle.Rect
		moveVec := delta(enemy.componentMovement.Direction, speed, speed)
		enemy.componentRectangle.Rect = enemy.componentRectangle.Rect.Moved(moveVec)
		enemy.componentMovement.RemainingMoves--
	} else {
		enemy.componentMovement.MovingFromHit = false
		enemy.componentMovement.RemainingMoves = int(enemy.componentRectangle.Rect.W())
	}
}
