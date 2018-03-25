package collision

import (
	"fmt"

	"github.com/miketmoore/zelduh/components"
)

type collisionEntity struct {
	ID int
	*components.SpatialComponent
}

// System is a custom system for detecting collisions and what to do when they occur
type System struct {
	playerEntity                collisionEntity
	enemies                     []collisionEntity
	coins                       []collisionEntity
	obstacles                   []collisionEntity
	PlayerCollisionWithCoin     func(int)
	PlayerCollisionWithEnemy    func(int)
	PlayerCollisionWithObstacle func(int)
	EnemyCollisionWithObstacle  func(int, int)
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(spatial *components.SpatialComponent) {
	s.playerEntity = collisionEntity{
		SpatialComponent: spatial,
	}
}

// AddEnemy adds an enemy to the system
func (s *System) AddEnemy(id int, spatial *components.SpatialComponent) {
	s.enemies = append(s.enemies, collisionEntity{
		ID:               id,
		SpatialComponent: spatial,
	})
}

// AddObstacle adds an obstacle entity to the system
func (s *System) AddObstacle(id int, spatial *components.SpatialComponent) {
	s.obstacles = append(s.obstacles, collisionEntity{
		ID:               id,
		SpatialComponent: spatial,
	})
}

// AddCoin adds a coin to the system
func (s *System) AddCoin(id int, spatial *components.SpatialComponent) {
	fmt.Printf("collision.System.AddCoin() id %d\n", id)
	s.coins = append(s.coins, collisionEntity{
		ID:               id,
		SpatialComponent: spatial,
	})
}

// RemoveCoin removes the specified coin from the system
func (s *System) RemoveCoin(id int) {
	for i := len(s.coins) - 1; i >= 0; i-- {
		coin := s.coins[i]
		if coin.ID == id {
			s.coins = append(s.coins[:i], s.coins[i+1:]...)
		}
	}
}

// Update checks for collisions
func (s *System) Update() {
	for _, enemy := range s.enemies {
		intersection := enemy.SpatialComponent.Rect.Intersect(s.playerEntity.SpatialComponent.Rect)
		if intersection.Area() > 0 {
			s.PlayerCollisionWithEnemy(enemy.ID)
		}
	}
	for _, coin := range s.coins {
		intersection := coin.SpatialComponent.Rect.Intersect(s.playerEntity.SpatialComponent.Rect)
		if intersection.Area() > 0 {
			s.PlayerCollisionWithCoin(coin.ID)
		}
	}
	for _, obstacle := range s.obstacles {
		intersectionWithPlayer := obstacle.SpatialComponent.Rect.Intersect(s.playerEntity.SpatialComponent.Rect)
		if intersectionWithPlayer.Area() > 0 {
			s.PlayerCollisionWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			intersectionWithEnemy := obstacle.SpatialComponent.Rect.Intersect(enemy.SpatialComponent.Rect)
			if intersectionWithEnemy.Area() > 0 {
				s.EnemyCollisionWithObstacle(enemy.ID, obstacle.ID)
			}
		}
	}
}
