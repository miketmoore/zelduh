package systems

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/components"
)

type collisionEntity struct {
	ID int
	*components.Spatial
}

// Collision is a custom system for detecting collisions and what to do when they occur
type Collision struct {
	MapBounds                             pixel.Rect
	player                                collisionEntity
	sword                                 collisionEntity
	arrow                                 collisionEntity
	enemies                               []collisionEntity
	coins                                 []collisionEntity
	obstacles                             []collisionEntity
	moveableObstacles                     []collisionEntity
	collisionSwitches                     []collisionEntity
	PlayerCollisionWithCoin               func(int)
	PlayerCollisionWithEnemy              func(int)
	SwordCollisionWithEnemy               func(int)
	ArrowCollisionWithEnemy               func(int)
	ArrowCollisionWithObstacle            func()
	PlayerCollisionWithObstacle           func(int)
	PlayerCollisionWithMoveableObstacle   func(int)
	EnemyCollisionWithObstacle            func(int, int)
	EnemyCollisionWithMoveableObstacle    func(int)
	MoveableObstacleCollisionWithObstacle func(int)
	PlayerCollisionWithSwitch             func(int)
	PlayerNoCollisionWithSwitch           func(int)
	PlayerCollisionWithBounds             func(string)
}

// AddPlayer adds the player to the system
func (s *Collision) AddPlayer(spatial *components.Spatial) {
	s.player = collisionEntity{
		Spatial: spatial,
	}
}

// AddSword adds the sword to the system
func (s *Collision) AddSword(spatial *components.Spatial) {
	s.sword = collisionEntity{
		Spatial: spatial,
	}
}

// AddArrow adds the arrow to the system
func (s *Collision) AddArrow(spatial *components.Spatial) {
	s.arrow = collisionEntity{
		Spatial: spatial,
	}
}

// AddEnemy adds an enemy to the system
func (s *Collision) AddEnemy(id int, spatial *components.Spatial) {
	s.enemies = append(s.enemies, collisionEntity{
		ID:      id,
		Spatial: spatial,
	})
}

// AddObstacle adds an obstacle entity to the system
func (s *Collision) AddObstacle(id int, spatial *components.Spatial) {
	s.obstacles = append(s.obstacles, collisionEntity{
		ID:      id,
		Spatial: spatial,
	})
}

// AddMoveableObstacle adds an obstacle entity to the system
func (s *Collision) AddMoveableObstacle(id int, spatial *components.Spatial) {
	s.moveableObstacles = append(s.moveableObstacles, collisionEntity{
		ID:      id,
		Spatial: spatial,
	})
}

// AddCollisionSwitch adds a collision switch entity to the system
func (s *Collision) AddCollisionSwitch(id int, spatial *components.Spatial) {
	s.collisionSwitches = append(s.collisionSwitches, collisionEntity{
		ID:      id,
		Spatial: spatial,
	})
}

// AddCoin adds a coin to the system
func (s *Collision) AddCoin(id int, spatial *components.Spatial) {
	fmt.Printf("collision.Collision.AddCoin() id %d\n", id)
	s.coins = append(s.coins, collisionEntity{
		ID:      id,
		Spatial: spatial,
	})
}

// RemoveCoin removes the specified coin from the system
func (s *Collision) RemoveCoin(id int) {
	for i := len(s.coins) - 1; i >= 0; i-- {
		coin := s.coins[i]
		if coin.ID == id {
			s.coins = append(s.coins[:i], s.coins[i+1:]...)
		}
	}
}

// RemoveEnemy removes the specified enemy from the system
func (s *Collision) RemoveEnemy(id int) {
	for i := len(s.enemies) - 1; i >= 0; i-- {
		enemy := s.enemies[i]
		if enemy.ID == id {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// RemoveAllEnemies removes all enemy entities from the system
func (s *Collision) RemoveAllEnemies() {
	for i := len(s.enemies) - 1; i >= 0; i-- {
		s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
	}
}

// RemoveAllCollisionSwitches removes all collision switches
func (s *Collision) RemoveAllCollisionSwitches() {
	for i := len(s.collisionSwitches) - 1; i >= 0; i-- {
		s.collisionSwitches = append(s.collisionSwitches[:i], s.collisionSwitches[i+1:]...)
	}
}

// RemoveObstacles removes all obstacles from the system
func (s *Collision) RemoveObstacles() {
	for i := len(s.obstacles) - 1; i >= 0; i-- {
		s.obstacles = append(s.obstacles[:i], s.obstacles[i+1:]...)
	}
}

func isColliding(r1, r2 pixel.Rect) bool {
	return r1.Min.X < r2.Max.X &&
		r1.Max.X > r2.Min.X &&
		r1.Min.Y < r2.Max.Y &&
		r1.Max.Y > r2.Min.Y
}

// Update checks for collisions
func (s *Collision) Update() {

	player := s.player
	playerR := s.player.Spatial.Rect
	mapBounds := s.MapBounds

	// is player at map edge?
	if player.Spatial.Rect.Min.Y <= mapBounds.Min.Y {
		s.PlayerCollisionWithBounds("bottom")
	} else if player.Spatial.Rect.Min.X <= mapBounds.Min.X {
		s.PlayerCollisionWithBounds("left")
	} else if player.Spatial.Rect.Max.X >= mapBounds.Max.X {
		s.PlayerCollisionWithBounds("right")
	} else if player.Spatial.Rect.Max.Y >= mapBounds.Max.Y {
		s.PlayerCollisionWithBounds("top")
	}

	w, h := player.Spatial.Width, player.Spatial.Height
	for _, enemy := range s.enemies {
		enemyR := enemy.Spatial.Rect

		if isCircleCollision(
			player.Spatial.HitBoxRadius,
			enemy.Spatial.HitBoxRadius,
			w, h, playerR, enemyR) {
			s.PlayerCollisionWithEnemy(enemy.ID)
		}

		if isCircleCollision(
			s.sword.Spatial.HitBoxRadius,
			enemy.Spatial.HitBoxRadius,
			w, h, s.sword.Spatial.Rect, enemyR) {
			s.SwordCollisionWithEnemy(enemy.ID)
		}

		if isCircleCollision(
			s.arrow.Spatial.HitBoxRadius,
			enemy.Spatial.HitBoxRadius,
			w, h, s.arrow.Spatial.Rect, enemyR) {
			s.ArrowCollisionWithEnemy(enemy.ID)
		}
	}
	for _, coin := range s.coins {
		if isColliding(coin.Spatial.Rect, s.player.Spatial.Rect) {
			s.PlayerCollisionWithCoin(coin.ID)
		}
	}
	// fmt.Printf("Total obstacles in collision system %d\n", len(s.obstacles))
	for _, obstacle := range s.obstacles {
		if isColliding(obstacle.Spatial.Rect, s.player.Spatial.Rect) {
			s.PlayerCollisionWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			if isColliding(obstacle.Spatial.Rect, enemy.Spatial.Rect) {
				s.EnemyCollisionWithObstacle(enemy.ID, obstacle.ID)
			}
		}

		if isColliding(obstacle.Spatial.Rect, s.arrow.Spatial.Rect) {
			s.ArrowCollisionWithObstacle()
		}
	}
	for _, moveableObstacle := range s.moveableObstacles {
		if isColliding(moveableObstacle.Spatial.Rect, s.player.Spatial.Rect) {
			s.PlayerCollisionWithMoveableObstacle(moveableObstacle.ID)
		}

		for _, enemy := range s.enemies {
			if isColliding(moveableObstacle.Spatial.Rect, enemy.Spatial.Rect) {
				// s.EnemyCollisionWithMoveableObstacle(enemy.ID)
			}
		}

		for _, obstacle := range s.obstacles {
			if isColliding(moveableObstacle.Spatial.Rect, obstacle.Spatial.Rect) {
				// s.MoveableObstacleCollisionWithObstacle(moveableObstacle.ID)
			}
		}

		if isColliding(moveableObstacle.Spatial.Rect, s.arrow.Spatial.Rect) {
			s.ArrowCollisionWithObstacle()
		}
	}

	for _, collisionSwitch := range s.collisionSwitches {
		if collisionSwitch.Spatial.HitBoxRadius > 0 {
			if isCircleCollision(
				s.player.Spatial.HitBoxRadius,
				collisionSwitch.Spatial.HitBoxRadius,
				w, h, s.player.Spatial.Rect, collisionSwitch.Spatial.Rect) {
				s.PlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.PlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		} else {
			if isColliding(s.player.Spatial.Rect, collisionSwitch.Spatial.Rect) {
				s.PlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.PlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		}

	}
}

func isCircleCollision(radius1, radius2, w, h float64, rect1, rect2 pixel.Rect) bool {
	x1 := rect1.Min.X + (w / 2)
	y1 := rect1.Min.Y + (h / 2)

	x2 := rect2.Min.X + (w / 2)
	y2 := rect2.Min.Y + (h / 2)

	dx := x1 - x2
	dy := y1 - y2

	distance := math.Sqrt(dx*dx + dy*dy)

	return distance < (radius1 + radius2)
}
