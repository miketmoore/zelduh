package spatial

import (
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
	enemies      []spatialEntity
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(spatial *components.SpatialComponent, movement *components.MovementComponent) {
	s.playerEntity = spatialEntity{
		SpatialComponent:  spatial,
		MovementComponent: movement,
	}
}

// AddEnemy adds an enemy to the system
func (s *System) AddEnemy(id int, spatial *components.SpatialComponent, movement *components.MovementComponent) {
	s.enemies = append(s.enemies, spatialEntity{
		SpatialComponent:      spatial,
		MovementComponent:     movement,
		EnemySpatialComponent: &enemySpatialComponent{},
	})
}

// Update changes spatial data based on movement data
func (s *System) Update() {
	player := s.playerEntity
	if player.MovementComponent.Moving {
		var v pixel.Vec
		speed := player.MovementComponent.Speed
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
		newRect := player.SpatialComponent.Rect.Moved(v)
		if player.SpatialComponent.BoundsRect.Contains(newRect.Min) &&
			player.SpatialComponent.BoundsRect.Contains(newRect.Max) {
			player.SpatialComponent.Rect = newRect
		}
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
			enemy.SpatialComponent.LastDir = direction.Up
		case 1:
			enemy.SpatialComponent.LastDir = direction.Right
		case 2:
			enemy.SpatialComponent.LastDir = direction.Down
		case 3:
			enemy.SpatialComponent.LastDir = direction.Left
		}
	} else {
		if enemy.EnemySpatialComponent.MoveCounter > 0 {
			moveVec := pixel.V(0, 0)
			if enemy.SpatialComponent.BoundsRect.Contains(enemy.SpatialComponent.Rect.Min) &&
				enemy.SpatialComponent.BoundsRect.Contains(enemy.SpatialComponent.Rect.Max) {
				// speed := enemy.MovementComponent.Speed
				speed := 1.0
				switch enemy.SpatialComponent.LastDir {
				case direction.Up:
					// fmt.Printf("Up\n")
					moveVec = pixel.V(0, speed)
				case direction.Right:
					// fmt.Printf("Right\n")
					moveVec = pixel.V(speed, 0)
				case direction.Down:
					// fmt.Printf("Down\n")
					moveVec = pixel.V(0, -speed)
				case direction.Left:
					// fmt.Printf("Left\n")
					moveVec = pixel.V(-speed, 0)
				}
				enemy.SpatialComponent.Rect = enemy.SpatialComponent.Rect.Moved(moveVec)
				enemy.EnemySpatialComponent.MoveCounter--
			}
		} else {
			enemy.EnemySpatialComponent.TotalMoves--
			enemy.EnemySpatialComponent.MoveCounter = int(enemy.SpatialComponent.Rect.W())
		}
	}
}
