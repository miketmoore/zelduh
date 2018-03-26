package spatial

import (
	"fmt"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type enemySpatialComponent struct {
	TotalMoves  int
	MoveCounter int
}

type spatialEntity struct {
	ID int
	*components.MovementComponent
	*components.SpatialComponent
	EnemySpatialComponent *enemySpatialComponent
}

// System is a custom system
type System struct {
	Rand         *rand.Rand
	playerEntity spatialEntity
	sword        spatialEntity
	enemies      []spatialEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(spatial *components.SpatialComponent, movement *components.MovementComponent) {
	s.playerEntity = spatialEntity{
		SpatialComponent:  spatial,
		MovementComponent: movement,
	}
}

// AddSword adds the player to the system
func (s *System) AddSword(spatial *components.SpatialComponent, movement *components.MovementComponent) {
	s.sword = spatialEntity{
		SpatialComponent:  spatial,
		MovementComponent: movement,
	}
}

// AddEnemy adds an enemy to the system
func (s *System) AddEnemy(id int, spatial *components.SpatialComponent, movement *components.MovementComponent) {
	s.enemies = append(s.enemies, spatialEntity{
		ID:                    id,
		SpatialComponent:      spatial,
		MovementComponent:     movement,
		EnemySpatialComponent: &enemySpatialComponent{},
	})
}

// MovePlayerBack moves the player back
func (s *System) MovePlayerBack() {
	fmt.Println("Move back")
	player := s.playerEntity
	var v pixel.Vec
	switch player.MovementComponent.Direction {
	case direction.Up:
		v = pixel.V(0, -48)
	case direction.Right:
		v = pixel.V(-48, 0)
	case direction.Down:
		v = pixel.V(0, 48)
	case direction.Left:
		v = pixel.V(48, 0)
	}
	player.SpatialComponent.Rect = player.SpatialComponent.PrevRect.Moved(v)
	player.SpatialComponent.PrevRect = player.SpatialComponent.Rect
}

func (s *System) enemy(id int) (spatialEntity, bool) {
	for _, e := range s.enemies {
		if e.ID == id {
			return e, true
		}
	}
	return spatialEntity{}, false
}

// MoveEnemyBack moves the enemy back
func (s *System) MoveEnemyBack(enemyID int, directionHit direction.Name) {
	fmt.Printf("spatial MoveEnemyBack called enemyID: %d\n", enemyID)
	enemy, ok := s.enemy(enemyID)
	if ok {
		fmt.Printf("spatial MoveEnemyBack OK, found enemy.ID: %d\n", enemy.ID)
		var v pixel.Vec
		fmt.Printf("spatial MoveEnemyBack direction: %v\n", enemy.MovementComponent.Direction)
		switch directionHit {
		case direction.Up:
			v = pixel.V(0, 48)
		case direction.Right:
			v = pixel.V(48, 0)
		case direction.Down:
			v = pixel.V(0, -48)
		case direction.Left:
			v = pixel.V(-48, 0)
		}
		enemy.SpatialComponent.Rect = enemy.SpatialComponent.PrevRect.Moved(v)
		enemy.SpatialComponent.PrevRect = enemy.SpatialComponent.Rect
	}
}

// Update changes spatial data based on movement data
func (s *System) Update() {
	player := s.playerEntity
	speed := player.MovementComponent.Speed
	if speed > 0 {
		var v pixel.Vec

		switch player.MovementComponent.Direction {
		case direction.Up:
			v = pixel.V(0, speed)
		case direction.Right:
			v = pixel.V(speed, 0)
		case direction.Down:
			v = pixel.V(0, -speed)
		case direction.Left:
			v = pixel.V(-speed, 0)
		}
		player.SpatialComponent.PrevRect = player.SpatialComponent.Rect
		player.SpatialComponent.Rect = player.SpatialComponent.Rect.Moved(v)
	}

	sword := s.sword
	speed = sword.MovementComponent.Speed
	swordW := sword.SpatialComponent.Width
	swordH := sword.SpatialComponent.Height
	if speed > 0 {
		var v pixel.Vec

		switch sword.MovementComponent.Direction {
		case direction.Up:
			v = pixel.V(0, speed+swordH)
		case direction.Right:
			v = pixel.V(speed+swordW, 0)
		case direction.Down:
			v = pixel.V(0, -speed-swordH)
		case direction.Left:
			v = pixel.V(-speed-swordW, 0)
		}
		sword.SpatialComponent.Rect = player.SpatialComponent.Rect.Moved(v)
	} else {
		sword.SpatialComponent.Rect = player.SpatialComponent.Rect
	}

	for _, enemy := range s.enemies {
		moveEnemy(s, enemy)
	}

}

// Enemy moves constantly, never stopping.
// Quick orientation changes.
func moveEnemy(s *System, enemy spatialEntity) {

	if enemy.EnemySpatialComponent.TotalMoves == 0 {
		maxMoves := 5
		enemy.EnemySpatialComponent.TotalMoves = s.Rand.Intn(maxMoves)

		directionIndex := s.Rand.Intn(4)
		switch directionIndex {
		case 0:
			enemy.MovementComponent.Direction = direction.Up
		case 1:
			enemy.MovementComponent.Direction = direction.Right
		case 2:
			enemy.MovementComponent.Direction = direction.Down
		case 3:
			enemy.MovementComponent.Direction = direction.Left
		}
	} else {
		if enemy.EnemySpatialComponent.MoveCounter > 0 {
			moveVec := pixel.V(0, 0)

			speed := 1.0
			switch enemy.MovementComponent.Direction {
			case direction.Up:
				moveVec = pixel.V(0, speed)
			case direction.Right:
				moveVec = pixel.V(speed, 0)
			case direction.Down:
				moveVec = pixel.V(0, -speed)
			case direction.Left:
				moveVec = pixel.V(-speed, 0)
			}
			// player.SpatialComponent.PrevRect = player.SpatialComponent.Rect
			// player.SpatialComponent.Rect = player.SpatialComponent.Rect.Moved(v)
			enemy.SpatialComponent.PrevRect = enemy.SpatialComponent.Rect
			enemy.SpatialComponent.Rect = enemy.SpatialComponent.Rect.Moved(moveVec)
			enemy.EnemySpatialComponent.MoveCounter--

		} else {
			enemy.EnemySpatialComponent.TotalMoves--
			enemy.EnemySpatialComponent.MoveCounter = int(enemy.SpatialComponent.Rect.W())
		}
	}
}
