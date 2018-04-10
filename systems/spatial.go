package systems

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"github.com/miketmoore/zelduh/entities"
)

type spatialEntity struct {
	ID entities.EntityID
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
	moveableObstacles []spatialEntity
}

// Add adds an entity to the system
func (s *Spatial) Add(category categories.Category, id entities.EntityID, spatial *components.Spatial, movement *components.Movement, dash *components.Dash) {
	switch category {
	case categories.Player:
		s.player = spatialEntity{
			Spatial:  spatial,
			Movement: movement,
			Dash:     dash,
		}
	case categories.Sword:
		s.sword = spatialEntity{
			Spatial:  spatial,
			Movement: movement,
		}
	case categories.Arrow:
		s.arrow = spatialEntity{
			Spatial:  spatial,
			Movement: movement,
		}
	case categories.Enemy:
		s.enemies = append(s.enemies, &spatialEntity{
			ID:          id,
			Spatial:     spatial,
			Movement:    movement,
			TotalMoves:  0,
			MoveCounter: 0,
		})
	case categories.MovableObstacle:
		s.moveableObstacles = append(s.moveableObstacles, spatialEntity{
			ID:       id,
			Spatial:  spatial,
			Movement: movement,
		})
	}
}

// Remove removes the entity from the system
func (s *Spatial) Remove(category categories.Category, id entities.EntityID) {
	switch category {
	case categories.Enemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			enemy := s.enemies[i]
			if enemy.ID == id {
				s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
			}
		}
	}
}

// RemoveAll removes all entities from one category
func (s *Spatial) RemoveAll(category categories.Category) {
	switch category {
	case categories.Enemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// MovePlayerBack moves the player back
func (s *Spatial) MovePlayerBack() {
	fmt.Println("Move back")
	player := s.player
	var v pixel.Vec
	switch player.Movement.Direction {
	case direction.Up:
		v = pixel.V(0, -48)
	case direction.Right:
		v = pixel.V(-48, 0)
	case direction.Down:
		v = pixel.V(0, 48)
	case direction.Left:
		v = pixel.V(48, 0)
	}
	player.Spatial.Rect = player.Spatial.PrevRect.Moved(v)
	player.Spatial.PrevRect = player.Spatial.Rect
}

// MoveMoveableObstacle moves a moveable obstacle :P
func (s *Spatial) MoveMoveableObstacle(obstacleID entities.EntityID, dir direction.Name) {
	obstacle, ok := s.moveableObstacle(obstacleID)
	if ok {
		w := obstacle.Spatial.Width
		h := obstacle.Spatial.Height
		var v pixel.Vec
		switch dir {
		case direction.Up:
			v = pixel.V(0, h)
		case direction.Right:
			v = pixel.V(w, 0)
		case direction.Down:
			v = pixel.V(0, -h)
		case direction.Left:
			v = pixel.V(-w, 0)
		}
		obstacle.Spatial.PrevRect = obstacle.Spatial.Rect
		obstacle.Spatial.Rect = obstacle.Spatial.Rect.Moved(v)
	}
}

// UndoEnemyRect resets current rect to previous rect
func (s *Spatial) UndoEnemyRect(enemyID entities.EntityID) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		enemy.Spatial.Rect = enemy.Spatial.PrevRect
	}
}

// MoveEnemyBack moves the enemy back
func (s *Spatial) MoveEnemyBack(enemyID entities.EntityID, directionHit direction.Name, distance float64) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		v := delta(directionHit, distance, distance)
		enemy.Spatial.Rect = enemy.Spatial.PrevRect.Moved(v)
		enemy.Spatial.PrevRect = enemy.Spatial.Rect
	}
}

// GetEnemySpatial returns the spatial component
func (s *Spatial) GetEnemySpatial(enemyID entities.EntityID) (*components.Spatial, bool) {
	for _, enemy := range s.enemies {
		if enemy.ID == enemyID {
			return enemy.Spatial, true
		}
	}
	return &components.Spatial{}, false
}

// Update changes spatial data based on movement data
func (s *Spatial) Update() {
	s.movePlayer()
	s.moveSword()
	s.moveArrow()

	for i := 0; i < len(s.enemies); i++ {
		enemy := s.enemies[i]
		moveEnemy(s, enemy)
	}
}

func (s *Spatial) moveableObstacle(id entities.EntityID) (spatialEntity, bool) {
	for _, e := range s.moveableObstacles {
		if e.ID == id {
			return e, true
		}
	}
	return spatialEntity{}, false
}

func (s *Spatial) enemy(id entities.EntityID) (spatialEntity, bool) {
	for _, e := range s.enemies {
		if e.ID == id {
			return *e, true
		}
	}
	return spatialEntity{}, false
}

func delta(dir direction.Name, modX, modY float64) pixel.Vec {
	switch dir {
	case direction.Up:
		return pixel.V(0, modY)
	case direction.Right:
		return pixel.V(modX, 0)
	case direction.Down:
		return pixel.V(0, -modY)
	case direction.Left:
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

func moveEnemy(s *Spatial, enemy *spatialEntity) {
	if enemy.Movement.RemainingMoves == 0 {
		enemy.Movement.RemainingMoves = s.Rand.Intn(enemy.Movement.MaxMoves)
		enemy.Movement.Direction = direction.Rand()
	} else if enemy.Movement.RemainingMoves > 0 {
		enemy.Spatial.PrevRect = enemy.Spatial.Rect
		speed := enemy.Movement.MaxSpeed
		moveVec := delta(enemy.Movement.Direction, speed, speed)
		enemy.Spatial.Rect = enemy.Spatial.Rect.Moved(moveVec)
		enemy.Movement.RemainingMoves--
	} else {
		enemy.Movement.RemainingMoves = int(enemy.Spatial.Rect.W())
	}
}
