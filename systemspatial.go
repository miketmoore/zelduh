package zelduh

import (
	"math/rand"

	"github.com/faiface/pixel"
)

type spatialEntity struct {
	ID EntityID
	*ComponentMovement
	*ComponentSpatial
	*ComponentDash
	TotalMoves  int
	MoveCounter int
}

// SystemSpatial is a custom system
type SystemSpatial struct {
	Rand              *rand.Rand
	player            spatialEntity
	sword             spatialEntity
	arrow             spatialEntity
	enemies           []*spatialEntity
	moveableObstacles []*spatialEntity
}

// AddEntity adds an entity to the system
func (s *SystemSpatial) AddEntity(entity Entity) {
	r := spatialEntity{
		ID:                entity.ID(),
		ComponentSpatial:  entity.ComponentSpatial,
		ComponentMovement: entity.ComponentMovement,
	}
	switch entity.Category {
	case CategoryPlayer:
		r.ComponentDash = entity.ComponentDash
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
func (s *SystemSpatial) Remove(category EntityCategory, id EntityID) {
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
func (s *SystemSpatial) RemoveAll(category EntityCategory) {
	switch category {
	case CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// MovePlayerBack moves the player back
func (s *SystemSpatial) MovePlayerBack() {
	player := s.player
	var v pixel.Vec
	switch player.ComponentMovement.Direction {
	case DirectionUp:
		v = pixel.V(0, -48)
	case DirectionRight:
		v = pixel.V(-48, 0)
	case DirectionDown:
		v = pixel.V(0, 48)
	case DirectionLeft:
		v = pixel.V(48, 0)
	}
	player.ComponentSpatial.Rect = player.ComponentSpatial.PrevRect.Moved(v)
	player.ComponentSpatial.PrevRect = player.ComponentSpatial.Rect
}

// MoveMoveableObstacle moves a moveable obstacle
func (s *SystemSpatial) MoveMoveableObstacle(obstacleID EntityID, dir Direction) bool {
	entity, ok := s.moveableObstacle(obstacleID)
	if ok && !entity.ComponentMovement.MovingFromHit {
		entity.ComponentMovement.MovingFromHit = true
		entity.ComponentMovement.RemainingMoves = entity.ComponentMovement.MaxMoves
		entity.ComponentMovement.Direction = dir
		return true
	}
	return false
}

// UndoEnemyRect resets current rect to previous rect
func (s *SystemSpatial) UndoEnemyRect(enemyID EntityID) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		enemy.ComponentSpatial.Rect = enemy.ComponentSpatial.PrevRect
	}
}

// MoveEnemyBack moves the enemy back
func (s *SystemSpatial) MoveEnemyBack(enemyID EntityID, directionHit Direction) {
	enemy, ok := s.enemy(enemyID)
	if ok && !enemy.ComponentMovement.MovingFromHit {
		enemy.ComponentMovement.MovingFromHit = true
		enemy.ComponentMovement.RemainingMoves = enemy.ComponentMovement.HitBackMoves
		enemy.ComponentMovement.Direction = directionHit
	}
}

// GetEnemySpatial returns the spatial component
func (s *SystemSpatial) GetEnemySpatial(enemyID EntityID) (*ComponentSpatial, bool) {
	for _, enemy := range s.enemies {
		if enemy.ID == enemyID {
			return enemy.ComponentSpatial, true
		}
	}
	return &ComponentSpatial{}, false
}

// EnemyMovingFromHit indicates if the enemy is moving after being hit
func (s *SystemSpatial) EnemyMovingFromHit(enemyID EntityID) bool {
	enemy, ok := s.enemy(enemyID)
	if ok {
		if enemy.ID == enemyID {
			return enemy.ComponentMovement.MovingFromHit == true
		}
	}
	return false
}

// Update changes spatial data based on movement data
func (s *SystemSpatial) Update() {
	s.movePlayer()
	s.moveSword()
	s.moveArrow()

	for i := 0; i < len(s.moveableObstacles); i++ {
		entity := s.moveableObstacles[i]
		s.moveMoveableObstacle(entity)
	}

	for i := 0; i < len(s.enemies); i++ {
		enemy := s.enemies[i]
		switch enemy.ComponentMovement.PatternName {
		case "random":
			s.moveEnemyRandom(enemy)
		case "left-right":
			s.moveEnemyLeftRight(enemy)
		}
	}
}

func (s *SystemSpatial) moveableObstacle(id EntityID) (spatialEntity, bool) {
	for _, e := range s.moveableObstacles {
		if e.ID == id {
			return *e, true
		}
	}
	return spatialEntity{}, false
}

func (s *SystemSpatial) enemy(id EntityID) (spatialEntity, bool) {
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

func (s *SystemSpatial) moveSword() {
	sword := s.sword
	speed := sword.ComponentMovement.Speed
	w := sword.ComponentSpatial.Width
	h := sword.ComponentSpatial.Height
	if speed > 0 {
		sword.ComponentSpatial.PrevRect = sword.ComponentSpatial.Rect
		v := delta(sword.ComponentMovement.Direction, speed+w, speed+h)
		sword.ComponentSpatial.Rect = s.player.ComponentSpatial.Rect.Moved(v)
	} else {
		sword.ComponentSpatial.Rect = s.player.ComponentSpatial.Rect
	}
}

func (s *SystemSpatial) moveArrow() {
	arrow := s.arrow
	speed := arrow.ComponentMovement.Speed
	if arrow.ComponentMovement.RemainingMoves > 0 {
		arrow.ComponentSpatial.PrevRect = arrow.ComponentSpatial.Rect
		v := delta(arrow.ComponentMovement.Direction, speed, speed)
		arrow.ComponentSpatial.Rect = arrow.ComponentSpatial.Rect.Moved(v)
	} else {
		arrow.ComponentSpatial.Rect = s.player.ComponentSpatial.Rect
	}
}

func (s *SystemSpatial) movePlayer() {
	player := s.player
	speed := player.ComponentMovement.Speed
	if player.ComponentDash.Charge == player.ComponentDash.MaxCharge {
		speed += player.ComponentDash.SpeedMod
	}
	if speed > 0 {
		v := delta(player.ComponentMovement.Direction, speed, speed)
		player.ComponentSpatial.PrevRect = player.ComponentSpatial.Rect
		player.ComponentSpatial.Rect = player.ComponentSpatial.Rect.Moved(v)
	}
}

func (s *SystemSpatial) moveMoveableObstacle(entity *spatialEntity) {
	if entity.ComponentMovement.RemainingMoves > 0 {
		speed := entity.ComponentMovement.MaxSpeed
		entity.ComponentSpatial.PrevRect = entity.ComponentSpatial.Rect
		moveVec := delta(entity.ComponentMovement.Direction, speed, speed)
		entity.ComponentSpatial.Rect = entity.ComponentSpatial.Rect.Moved(moveVec)
		entity.ComponentMovement.RemainingMoves--
	} else {
		entity.ComponentMovement.MovingFromHit = false
		entity.ComponentMovement.RemainingMoves = 0
	}
}

func (s *SystemSpatial) moveEnemyRandom(enemy *spatialEntity) {
	if enemy.ComponentMovement.RemainingMoves == 0 {
		enemy.ComponentMovement.MovingFromHit = false
		enemy.ComponentMovement.RemainingMoves = s.Rand.Intn(enemy.ComponentMovement.MaxMoves)
		enemy.ComponentMovement.Direction = RandomDirection(s.Rand)
	} else if enemy.ComponentMovement.RemainingMoves > 0 {
		var speed float64
		if enemy.ComponentMovement.MovingFromHit {
			speed = enemy.ComponentMovement.HitSpeed
		} else {
			speed = enemy.ComponentMovement.MaxSpeed
		}
		enemy.ComponentSpatial.PrevRect = enemy.ComponentSpatial.Rect
		moveVec := delta(enemy.ComponentMovement.Direction, speed, speed)
		enemy.ComponentSpatial.Rect = enemy.ComponentSpatial.Rect.Moved(moveVec)
		enemy.ComponentMovement.RemainingMoves--
	} else {
		enemy.ComponentMovement.MovingFromHit = false
		enemy.ComponentMovement.RemainingMoves = int(enemy.ComponentSpatial.Rect.W())
	}
}

func (s *SystemSpatial) moveEnemyLeftRight(enemy *spatialEntity) {
	if enemy.ComponentMovement.RemainingMoves == 0 {
		enemy.ComponentMovement.MovingFromHit = false
		enemy.ComponentMovement.RemainingMoves = enemy.ComponentMovement.MaxMoves
		switch enemy.ComponentMovement.Direction {
		case DirectionLeft:
			enemy.ComponentMovement.Direction = DirectionRight
		case DirectionRight:
			enemy.ComponentMovement.Direction = DirectionLeft
		}
	} else if enemy.ComponentMovement.RemainingMoves > 0 {
		var speed float64
		if enemy.ComponentMovement.MovingFromHit {
			speed = enemy.ComponentMovement.HitSpeed
		} else {
			speed = enemy.ComponentMovement.MaxSpeed
		}
		enemy.ComponentSpatial.PrevRect = enemy.ComponentSpatial.Rect
		moveVec := delta(enemy.ComponentMovement.Direction, speed, speed)
		enemy.ComponentSpatial.Rect = enemy.ComponentSpatial.Rect.Moved(moveVec)
		enemy.ComponentMovement.RemainingMoves--
	} else {
		enemy.ComponentMovement.MovingFromHit = false
		enemy.ComponentMovement.RemainingMoves = int(enemy.ComponentSpatial.Rect.W())
	}
}
