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

type enemySpatialComponent struct {
	TotalMoves  int
	MoveCounter int
}

type spatialEntity struct {
	ID entities.EntityID
	*components.Movement
	*components.Spatial
	*components.Dash
	EnemySpatialComponent *enemySpatialComponent
}

// Spatial is a custom system
type Spatial struct {
	Rand              *rand.Rand
	playerEntity      spatialEntity
	sword             spatialEntity
	arrow             spatialEntity
	enemies           []spatialEntity
	moveableObstacles []spatialEntity
}

// AddPlayer adds the player to the system
func (s *Spatial) AddPlayer(
	spatial *components.Spatial,
	movement *components.Movement,
	dash *components.Dash,
) {
	s.playerEntity = spatialEntity{
		Spatial:  spatial,
		Movement: movement,
		Dash:     dash,
	}
}

// AddSword adds the sword to the system
func (s *Spatial) AddSword(spatial *components.Spatial, movement *components.Movement) {
	s.sword = spatialEntity{
		Spatial:  spatial,
		Movement: movement,
	}
}

// AddArrow adds the arrow to the system
func (s *Spatial) AddArrow(spatial *components.Spatial, movement *components.Movement) {
	s.arrow = spatialEntity{
		Spatial:  spatial,
		Movement: movement,
	}
}

// AddEnemy adds an enemy to the system
func (s *Spatial) AddEnemy(id entities.EntityID, spatial *components.Spatial, movement *components.Movement) {
	s.enemies = append(s.enemies, spatialEntity{
		ID:                    id,
		Spatial:               spatial,
		Movement:              movement,
		EnemySpatialComponent: &enemySpatialComponent{},
	})
}

// AddMoveableObstacle adds a moveable obstacle to the system
func (s *Spatial) AddMoveableObstacle(id entities.EntityID, spatial *components.Spatial, movement *components.Movement) {
	s.moveableObstacles = append(s.moveableObstacles, spatialEntity{
		ID:       id,
		Spatial:  spatial,
		Movement: movement,
	})
}

// RemoveEnemy removes the specified enemy from the system
func (s *Spatial) RemoveEnemy(id entities.EntityID) {
	for i := len(s.enemies) - 1; i >= 0; i-- {
		enemy := s.enemies[i]
		if enemy.ID == id {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
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
	player := s.playerEntity
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
			return e, true
		}
	}
	return spatialEntity{}, false
}

// UndoEnemyRect resets current rect to previous rect
func (s *Spatial) UndoEnemyRect(enemyID entities.EntityID) {
	enemy, ok := s.enemy(enemyID)
	if ok {
		enemy.Spatial.Rect = enemy.Spatial.PrevRect
	}
}

// MoveEnemyBack moves the enemy back
// TODO how to prevent entity from passing through obstacles, map boundaries?
func (s *Spatial) MoveEnemyBack(enemyID entities.EntityID, directionHit direction.Name, distance float64) {
	// fmt.Printf("spatial MoveEnemyBack called enemyID: %d\n", enemyID)
	enemy, ok := s.enemy(enemyID)
	if ok {
		// fmt.Printf("spatial MoveEnemyBack OK, found enemy.ID: %d\n", enemy.ID)
		var v pixel.Vec
		// fmt.Printf("spatial MoveEnemyBack direction: %v\n", enemy.Movement.Direction)
		switch directionHit {
		case direction.Up:
			v = pixel.V(0, distance)
		case direction.Right:
			v = pixel.V(distance, 0)
		case direction.Down:
			v = pixel.V(0, -distance)
		case direction.Left:
			v = pixel.V(-distance, 0)
		}
		enemy.Spatial.Rect = enemy.Spatial.PrevRect.Moved(v)
		enemy.Spatial.PrevRect = enemy.Spatial.Rect
	}
}

// Update changes spatial data based on movement data
func (s *Spatial) Update() {
	player := s.playerEntity
	speed := player.Movement.Speed
	if player.Dash.Charge == player.Dash.MaxCharge {
		speed += player.Dash.SpeedMod
	}
	if speed > 0 {
		var v pixel.Vec

		switch player.Movement.Direction {
		case direction.Up:
			v = pixel.V(0, speed)
		case direction.Right:
			v = pixel.V(speed, 0)
		case direction.Down:
			v = pixel.V(0, -speed)
		case direction.Left:
			v = pixel.V(-speed, 0)
		}
		player.Spatial.PrevRect = player.Spatial.Rect
		player.Spatial.Rect = player.Spatial.Rect.Moved(v)
	}

	sword := s.sword
	speed = sword.Movement.Speed
	swordW := sword.Spatial.Width
	swordH := sword.Spatial.Height
	if speed > 0 {
		var v pixel.Vec

		switch sword.Movement.Direction {
		case direction.Up:
			v = pixel.V(0, speed+swordH)
		case direction.Right:
			v = pixel.V(speed+swordW, 0)
		case direction.Down:
			v = pixel.V(0, -speed-swordH)
		case direction.Left:
			v = pixel.V(-speed-swordW, 0)
		}
		player.Spatial.PrevRect = player.Spatial.Rect
		sword.Spatial.Rect = player.Spatial.Rect.Moved(v)
	} else {
		sword.Spatial.Rect = player.Spatial.Rect
	}

	arrow := s.arrow
	speed = arrow.Movement.Speed
	if arrow.Movement.MoveCount > 0 {
		var v pixel.Vec

		switch arrow.Movement.Direction {
		case direction.Up:
			v = pixel.V(0, speed)
		case direction.Right:
			v = pixel.V(speed, 0)
		case direction.Down:
			v = pixel.V(0, -speed)
		case direction.Left:
			v = pixel.V(-speed, 0)
		}
		arrow.Spatial.PrevRect = arrow.Spatial.Rect
		arrow.Spatial.Rect = arrow.Spatial.Rect.Moved(v)
	} else {
		arrow.Spatial.Rect = player.Spatial.Rect
	}

	for _, enemy := range s.enemies {
		moveEnemy(s, enemy)
	}

}

// Enemy moves constantly, never stopping.
// Quick orientation changes.
func moveEnemy(s *Spatial, enemy spatialEntity) {

	if enemy.EnemySpatialComponent.TotalMoves == 0 {
		maxMoves := 5
		enemy.EnemySpatialComponent.TotalMoves = s.Rand.Intn(maxMoves)

		directionIndex := s.Rand.Intn(4)
		switch directionIndex {
		case 0:
			enemy.Movement.Direction = direction.Up
		case 1:
			enemy.Movement.Direction = direction.Right
		case 2:
			enemy.Movement.Direction = direction.Down
		case 3:
			enemy.Movement.Direction = direction.Left
		}
	} else {
		if enemy.EnemySpatialComponent.MoveCounter > 0 {
			moveVec := pixel.V(0, 0)

			speed := 1.0
			switch enemy.Movement.Direction {
			case direction.Up:
				moveVec = pixel.V(0, speed)
			case direction.Right:
				moveVec = pixel.V(speed, 0)
			case direction.Down:
				moveVec = pixel.V(0, -speed)
			case direction.Left:
				moveVec = pixel.V(-speed, 0)
			}
			// player.Spatial.PrevRect = player.Spatial.Rect
			// player.Spatial.Rect = player.Spatial.Rect.Moved(v)
			enemy.Spatial.PrevRect = enemy.Spatial.Rect
			enemy.Spatial.Rect = enemy.Spatial.Rect.Moved(moveVec)
			enemy.EnemySpatialComponent.MoveCounter--

		} else {
			enemy.EnemySpatialComponent.TotalMoves--
			enemy.EnemySpatialComponent.MoveCounter = int(enemy.Spatial.Rect.W())
		}
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
