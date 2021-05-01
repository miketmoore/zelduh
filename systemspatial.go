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

type spatialEntity struct {
	ID EntityID
	*componentMovement
	*componentDash
	*componentDimensions
	*componentRectangle
	TotalMoves  int
	MoveCounter int
}

// SpatialSystem is a custom system
type SpatialSystem struct {
	Rand              *rand.Rand
	player            spatialEntity
	sword             spatialEntity
	arrow             spatialEntity
	enemies           []*spatialEntity
	moveableObstacles []*spatialEntity
}

// AddEntity adds an entity to the system
func (s *SpatialSystem) AddEntity(entity Entity) {
	r := spatialEntity{
		ID:                  entity.ID(),
		componentMovement:   entity.componentMovement,
		componentDimensions: entity.componentDimensions,
		componentRectangle:  entity.componentRectangle,
	}
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
func (s *SpatialSystem) Remove(category EntityCategory, id EntityID) {
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
func (s *SpatialSystem) RemoveAll(category EntityCategory) {
	switch category {
	case CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// MovePlayerBack moves the player back
func (s *SpatialSystem) MovePlayerBack() {
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
func (s *SpatialSystem) MoveMoveableObstacle(obstacleID EntityID, dir Direction) bool {
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
func (s *SpatialSystem) UndoEnemyRect(enemyID EntityID) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		enemy.componentRectangle.Rect = enemy.componentRectangle.PrevRect
	}
}

// MoveEnemyBack moves the enemy back
func (s *SpatialSystem) MoveEnemyBack(enemyID EntityID, directionHit Direction) {
	enemy, ok := s.enemy(enemyID)
	if ok && !enemy.componentMovement.MovingFromHit {
		enemy.componentMovement.MovingFromHit = true
		enemy.componentMovement.RemainingMoves = enemy.componentMovement.HitBackMoves
		enemy.componentMovement.Direction = directionHit
	}
}

// ComponentRectangle returns the ComponentRectangle for the entity
func (s *SpatialSystem) ComponentRectangle(entityID EntityID) (*componentRectangle, bool) {
	for _, entity := range s.enemies {
		if entity.ID == entityID {
			return entity.componentRectangle, true
		}
	}
	return &componentRectangle{}, false
}

// EnemyMovingFromHit indicates if the enemy is moving after being hit
func (s *SpatialSystem) EnemyMovingFromHit(enemyID EntityID) bool {
	enemy, ok := s.enemy(enemyID)
	if ok {
		if enemy.ID == enemyID {
			return enemy.componentMovement.MovingFromHit == true
		}
	}
	return false
}

// Update changes spatial data based on movement data
func (s *SpatialSystem) Update() error {
	s.movePlayer()
	s.moveSword()
	s.moveArrow()

	for i := 0; i < len(s.moveableObstacles); i++ {
		entity := s.moveableObstacles[i]
		s.moveMoveableObstacle(entity)
	}

	// for i := 0; i < len(s.enemies); i++ {
	// 	enemy := s.enemies[i]
	// 	switch enemy.componentMovement.PatternName {
	// 	case "random":
	// 		s.moveEnemyRandom(enemy)
	// 	case "left-right":
	// 		s.moveEnemyLeftRight(enemy)
	// 	}
	// }
	return nil
}

func (s *SpatialSystem) moveableObstacle(id EntityID) (spatialEntity, bool) {
	for _, e := range s.moveableObstacles {
		if e.ID == id {
			return *e, true
		}
	}
	return spatialEntity{}, false
}

func (s *SpatialSystem) enemy(id EntityID) (spatialEntity, bool) {
	for _, e := range s.enemies {
		if e.ID == id {
			return *e, true
		}
	}
	return spatialEntity{}, false
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

func (s *SpatialSystem) moveSword() {
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

func (s *SpatialSystem) moveArrow() {
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

func (s *SpatialSystem) movePlayer() {
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

func (s *SpatialSystem) moveMoveableObstacle(entity *spatialEntity) {
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

func (s *SpatialSystem) moveEnemyRandom(enemy *spatialEntity) {
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

func (s *SpatialSystem) moveEnemyLeftRight(enemy *spatialEntity) {
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
