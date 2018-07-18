package zelduh

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/miketmoore/terraform2d"
)

type spatialEntity struct {
	ID terraform2d.EntityID
	*ComponentMovement
	*ComponentSpatial
	*ComponentDash
	TotalMoves  int
	MoveCounter int
}

// Spatial is a custom system
type Spatial struct {
	Rand              *rand.Rand
	player            spatialEntity
	sword             spatialEntity
	arrow             spatialEntity
	enemies           []*spatialEntity
	moveableObstacles []*spatialEntity
}

// AddEntity adds an entity to the system
func (s *Spatial) AddEntity(entity Entity) {
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
func (s *Spatial) Remove(category terraform2d.EntityCategory, id terraform2d.EntityID) {
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
func (s *Spatial) RemoveAll(category terraform2d.EntityCategory) {
	switch category {
	case CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// MovePlayerBack moves the player back
func (s *Spatial) MovePlayerBack() {
	player := s.player
	var v pixel.Vec
	switch player.ComponentMovement.Direction {
	case terraform2d.DirectionUp:
		v = pixel.V(0, -48)
	case terraform2d.DirectionRight:
		v = pixel.V(-48, 0)
	case terraform2d.DirectionDown:
		v = pixel.V(0, 48)
	case terraform2d.DirectionLeft:
		v = pixel.V(48, 0)
	}
	player.ComponentSpatial.Rect = player.ComponentSpatial.PrevRect.Moved(v)
	player.ComponentSpatial.PrevRect = player.ComponentSpatial.Rect
}

// MoveMoveableObstacle moves a moveable obstacle
func (s *Spatial) MoveMoveableObstacle(obstacleID terraform2d.EntityID, dir terraform2d.Direction) bool {
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
func (s *Spatial) UndoEnemyRect(enemyID terraform2d.EntityID) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		enemy.ComponentSpatial.Rect = enemy.ComponentSpatial.PrevRect
	}
}

// MoveEnemyBack moves the enemy back
func (s *Spatial) MoveEnemyBack(enemyID terraform2d.EntityID, directionHit terraform2d.Direction) {
	enemy, ok := s.enemy(enemyID)
	if ok && !enemy.ComponentMovement.MovingFromHit {
		enemy.ComponentMovement.MovingFromHit = true
		enemy.ComponentMovement.RemainingMoves = enemy.ComponentMovement.HitBackMoves
		enemy.ComponentMovement.Direction = directionHit
	}
}

// GetEnemySpatial returns the spatial component
func (s *Spatial) GetEnemySpatial(enemyID terraform2d.EntityID) (*ComponentSpatial, bool) {
	for _, enemy := range s.enemies {
		if enemy.ID == enemyID {
			return enemy.ComponentSpatial, true
		}
	}
	return &ComponentSpatial{}, false
}

// EnemyMovingFromHit indicates if the enemy is moving after being hit
func (s *Spatial) EnemyMovingFromHit(enemyID terraform2d.EntityID) bool {
	enemy, ok := s.enemy(enemyID)
	if ok {
		if enemy.ID == enemyID {
			return enemy.ComponentMovement.MovingFromHit == true
		}
	}
	return false
}

// Update changes spatial data based on movement data
func (s *Spatial) Update() {
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

func (s *Spatial) moveableObstacle(id terraform2d.EntityID) (spatialEntity, bool) {
	for _, e := range s.moveableObstacles {
		if e.ID == id {
			return *e, true
		}
	}
	return spatialEntity{}, false
}

func (s *Spatial) enemy(id terraform2d.EntityID) (spatialEntity, bool) {
	for _, e := range s.enemies {
		if e.ID == id {
			return *e, true
		}
	}
	return spatialEntity{}, false
}

func delta(dir terraform2d.Direction, modX, modY float64) pixel.Vec {
	switch dir {
	case terraform2d.DirectionUp:
		return pixel.V(0, modY)
	case terraform2d.DirectionRight:
		return pixel.V(modX, 0)
	case terraform2d.DirectionDown:
		return pixel.V(0, -modY)
	case terraform2d.DirectionLeft:
		return pixel.V(-modX, 0)
	default:
		return pixel.V(0, 0)
	}
}

func (s *Spatial) moveSword() {
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

func (s *Spatial) moveArrow() {
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

func (s *Spatial) movePlayer() {
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

func (s *Spatial) moveMoveableObstacle(entity *spatialEntity) {
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

func (s *Spatial) moveEnemyRandom(enemy *spatialEntity) {
	if enemy.ComponentMovement.RemainingMoves == 0 {
		enemy.ComponentMovement.MovingFromHit = false
		enemy.ComponentMovement.RemainingMoves = s.Rand.Intn(enemy.ComponentMovement.MaxMoves)
		enemy.ComponentMovement.Direction = terraform2d.RandomDirection(s.Rand)
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

func (s *Spatial) moveEnemyLeftRight(enemy *spatialEntity) {
	if enemy.ComponentMovement.RemainingMoves == 0 {
		enemy.ComponentMovement.MovingFromHit = false
		enemy.ComponentMovement.RemainingMoves = enemy.ComponentMovement.MaxMoves
		switch enemy.ComponentMovement.Direction {
		case terraform2d.DirectionLeft:
			enemy.ComponentMovement.Direction = terraform2d.DirectionRight
		case terraform2d.DirectionRight:
			enemy.ComponentMovement.Direction = terraform2d.DirectionLeft
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
