package systems

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/entities"
)

type spatialEntity struct {
	ID terraform2d.EntityID
	*components.Movement
	*components.Spatial
	*components.Dash
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
func (s *Spatial) AddEntity(entity entities.Entity) {
	r := spatialEntity{
		ID:       entity.ID(),
		Spatial:  entity.Spatial,
		Movement: entity.Movement,
	}
	switch entity.Category {
	case zelduh.CategoryPlayer:
		r.Dash = entity.Dash
		s.player = r
	case zelduh.CategorySword:
		s.sword = r
	case zelduh.CategoryArrow:
		s.arrow = r
	case zelduh.CategoryMovableObstacle:
		s.moveableObstacles = append(s.moveableObstacles, &r)
	case zelduh.CategoryEnemy:
		s.enemies = append(s.enemies, &r)
	}
}

// Remove removes the entity from the system
func (s *Spatial) Remove(category terraform2d.EntityCategory, id terraform2d.EntityID) {
	switch category {
	case zelduh.CategoryEnemy:
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
	case zelduh.CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// MovePlayerBack moves the player back
func (s *Spatial) MovePlayerBack() {
	player := s.player
	var v pixel.Vec
	switch player.Movement.Direction {
	case terraform2d.DirectionUp:
		v = pixel.V(0, -48)
	case terraform2d.DirectionRight:
		v = pixel.V(-48, 0)
	case terraform2d.DirectionDown:
		v = pixel.V(0, 48)
	case terraform2d.DirectionLeft:
		v = pixel.V(48, 0)
	}
	player.Spatial.Rect = player.Spatial.PrevRect.Moved(v)
	player.Spatial.PrevRect = player.Spatial.Rect
}

// MoveMoveableObstacle moves a moveable obstacle
func (s *Spatial) MoveMoveableObstacle(obstacleID terraform2d.EntityID, dir terraform2d.Direction) bool {
	entity, ok := s.moveableObstacle(obstacleID)
	if ok && !entity.Movement.MovingFromHit {
		entity.Movement.MovingFromHit = true
		entity.Movement.RemainingMoves = entity.Movement.MaxMoves
		entity.Movement.Direction = dir
		return true
	}
	return false
}

// UndoEnemyRect resets current rect to previous rect
func (s *Spatial) UndoEnemyRect(enemyID terraform2d.EntityID) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		enemy.Spatial.Rect = enemy.Spatial.PrevRect
	}
}

// MoveEnemyBack moves the enemy back
func (s *Spatial) MoveEnemyBack(enemyID terraform2d.EntityID, directionHit terraform2d.Direction) {
	enemy, ok := s.enemy(enemyID)
	if ok && !enemy.Movement.MovingFromHit {
		enemy.Movement.MovingFromHit = true
		enemy.Movement.RemainingMoves = enemy.Movement.HitBackMoves
		enemy.Movement.Direction = directionHit
	}
}

// GetEnemySpatial returns the spatial component
func (s *Spatial) GetEnemySpatial(enemyID terraform2d.EntityID) (*components.Spatial, bool) {
	for _, enemy := range s.enemies {
		if enemy.ID == enemyID {
			return enemy.Spatial, true
		}
	}
	return &components.Spatial{}, false
}

// EnemyMovingFromHit indicates if the enemy is moving after being hit
func (s *Spatial) EnemyMovingFromHit(enemyID terraform2d.EntityID) bool {
	enemy, ok := s.enemy(enemyID)
	if ok {
		if enemy.ID == enemyID {
			return enemy.Movement.MovingFromHit == true
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
		switch enemy.Movement.PatternName {
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
	speed := sword.Movement.Speed
	w := sword.Spatial.Width
	h := sword.Spatial.Height
	if speed > 0 {
		sword.Spatial.PrevRect = sword.Spatial.Rect
		v := delta(sword.Movement.Direction, speed+w, speed+h)
		sword.Spatial.Rect = s.player.Spatial.Rect.Moved(v)
	} else {
		sword.Spatial.Rect = s.player.Spatial.Rect
	}
}

func (s *Spatial) moveArrow() {
	arrow := s.arrow
	speed := arrow.Movement.Speed
	if arrow.Movement.RemainingMoves > 0 {
		arrow.Spatial.PrevRect = arrow.Spatial.Rect
		v := delta(arrow.Movement.Direction, speed, speed)
		arrow.Spatial.Rect = arrow.Spatial.Rect.Moved(v)
	} else {
		arrow.Spatial.Rect = s.player.Spatial.Rect
	}
}

func (s *Spatial) movePlayer() {
	player := s.player
	speed := player.Movement.Speed
	if player.Dash.Charge == player.Dash.MaxCharge {
		speed += player.Dash.SpeedMod
	}
	if speed > 0 {
		v := delta(player.Movement.Direction, speed, speed)
		player.Spatial.PrevRect = player.Spatial.Rect
		player.Spatial.Rect = player.Spatial.Rect.Moved(v)
	}
}

func (s *Spatial) moveMoveableObstacle(entity *spatialEntity) {
	if entity.Movement.RemainingMoves > 0 {
		speed := entity.Movement.MaxSpeed
		entity.Spatial.PrevRect = entity.Spatial.Rect
		moveVec := delta(entity.Movement.Direction, speed, speed)
		entity.Spatial.Rect = entity.Spatial.Rect.Moved(moveVec)
		entity.Movement.RemainingMoves--
	} else {
		entity.Movement.MovingFromHit = false
		entity.Movement.RemainingMoves = 0
	}
}

func (s *Spatial) moveEnemyRandom(enemy *spatialEntity) {
	if enemy.Movement.RemainingMoves == 0 {
		enemy.Movement.MovingFromHit = false
		enemy.Movement.RemainingMoves = s.Rand.Intn(enemy.Movement.MaxMoves)
		enemy.Movement.Direction = terraform2d.RandomDirection(s.Rand)
	} else if enemy.Movement.RemainingMoves > 0 {
		var speed float64
		if enemy.Movement.MovingFromHit {
			speed = enemy.Movement.HitSpeed
		} else {
			speed = enemy.Movement.MaxSpeed
		}
		enemy.Spatial.PrevRect = enemy.Spatial.Rect
		moveVec := delta(enemy.Movement.Direction, speed, speed)
		enemy.Spatial.Rect = enemy.Spatial.Rect.Moved(moveVec)
		enemy.Movement.RemainingMoves--
	} else {
		enemy.Movement.MovingFromHit = false
		enemy.Movement.RemainingMoves = int(enemy.Spatial.Rect.W())
	}
}

func (s *Spatial) moveEnemyLeftRight(enemy *spatialEntity) {
	if enemy.Movement.RemainingMoves == 0 {
		enemy.Movement.MovingFromHit = false
		enemy.Movement.RemainingMoves = enemy.Movement.MaxMoves
		switch enemy.Movement.Direction {
		case terraform2d.DirectionLeft:
			enemy.Movement.Direction = terraform2d.DirectionRight
		case terraform2d.DirectionRight:
			enemy.Movement.Direction = terraform2d.DirectionLeft
		}
	} else if enemy.Movement.RemainingMoves > 0 {
		var speed float64
		if enemy.Movement.MovingFromHit {
			speed = enemy.Movement.HitSpeed
		} else {
			speed = enemy.Movement.MaxSpeed
		}
		enemy.Spatial.PrevRect = enemy.Spatial.Rect
		moveVec := delta(enemy.Movement.Direction, speed, speed)
		enemy.Spatial.Rect = enemy.Spatial.Rect.Moved(moveVec)
		enemy.Movement.RemainingMoves--
	} else {
		enemy.Movement.MovingFromHit = false
		enemy.Movement.RemainingMoves = int(enemy.Spatial.Rect.W())
	}
}
