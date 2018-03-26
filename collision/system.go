package collision

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/components"
)

type collisionEntity struct {
	ID int
	*components.SpatialComponent
}

// System is a custom system for detecting collisions and what to do when they occur
type System struct {
	playerEntity                collisionEntity
	sword                       collisionEntity
	enemies                     []collisionEntity
	coins                       []collisionEntity
	obstacles                   []collisionEntity
	PlayerCollisionWithCoin     func(int)
	PlayerCollisionWithEnemy    func(int)
	SwordCollisionWithEnemy     func(int)
	PlayerCollisionWithObstacle func(int)
	EnemyCollisionWithObstacle  func(int, int)
	playerwithenemy             int
	swordhitenemy               int
}

// AddPlayer adds the player to the system
func (s *System) AddPlayer(spatial *components.SpatialComponent) {
	s.playerEntity = collisionEntity{
		SpatialComponent: spatial,
	}
}

// AddSword adds the player to the system
func (s *System) AddSword(spatial *components.SpatialComponent) {
	s.sword = collisionEntity{
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

// RemoveEnemy removes the specified enemy from the system
func (s *System) RemoveEnemy(id int) {
	for i := len(s.enemies) - 1; i >= 0; i-- {
		enemy := s.enemies[i]
		if enemy.ID == id {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

func isColliding(r1, r2 pixel.Rect) bool {
	return r1.Min.X < r2.Max.X &&
		r1.Max.X > r2.Min.X &&
		r1.Min.Y < r2.Max.Y &&
		r1.Max.Y > r2.Min.Y
}

// Update checks for collisions
func (s *System) Update() {
	for _, enemy := range s.enemies {
		enemyR := enemy.SpatialComponent.Rect
		if isColliding(enemyR, s.playerEntity.SpatialComponent.Rect) {
			s.playerwithenemy++
			s.PlayerCollisionWithEnemy(enemy.ID)
		}
		if isColliding(enemyR, s.sword.SpatialComponent.Rect) {
			s.SwordCollisionWithEnemy(enemy.ID)
		}
	}
	for _, coin := range s.coins {
		if isColliding(coin.SpatialComponent.Rect, s.playerEntity.SpatialComponent.Rect) {
			s.PlayerCollisionWithCoin(coin.ID)
		}
	}
	for _, obstacle := range s.obstacles {
		if isColliding(obstacle.SpatialComponent.Rect, s.playerEntity.SpatialComponent.Rect) {
			s.PlayerCollisionWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			if isColliding(obstacle.SpatialComponent.Rect, enemy.SpatialComponent.Rect) {
				s.EnemyCollisionWithObstacle(enemy.ID, obstacle.ID)
			}
		}
	}
}
