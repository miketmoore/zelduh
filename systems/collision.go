package systems

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/bounds"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/entities"
)

type collisionEntity struct {
	ID entities.EntityID
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
	PlayerCollisionWithCoin               func(entities.EntityID)
	PlayerCollisionWithEnemy              func(entities.EntityID)
	SwordCollisionWithEnemy               func(entities.EntityID)
	ArrowCollisionWithEnemy               func(entities.EntityID)
	ArrowCollisionWithObstacle            func()
	PlayerCollisionWithObstacle           func(entities.EntityID)
	PlayerCollisionWithMoveableObstacle   func(entities.EntityID)
	EnemyCollisionWithObstacle            func(entities.EntityID, entities.EntityID)
	EnemyCollisionWithMoveableObstacle    func(entities.EntityID)
	MoveableObstacleCollisionWithObstacle func(entities.EntityID)
	PlayerCollisionWithSwitch             func(entities.EntityID)
	PlayerNoCollisionWithSwitch           func(entities.EntityID)
	PlayerCollisionWithBounds             func(bounds.Bound)
}

// AddEntity adds an entity to the system
func (s *Collision) AddEntity(entity entities.Entity) {
	switch entity.Category {
	case categories.Player:
		s.player = collisionEntity{
			Spatial: entity.Spatial,
		}
	case categories.Sword:
		s.sword = collisionEntity{
			Spatial: entity.Spatial,
		}
	}
}

// Add adds the entity to the system
func (s *Collision) Add(category categories.Category, id entities.EntityID, spatial *components.Spatial) {
	switch category {
	case categories.Arrow:
		s.arrow = collisionEntity{
			Spatial: spatial,
		}
	case categories.Enemy:
		s.enemies = append(s.enemies, collisionEntity{
			ID:      id,
			Spatial: spatial,
		})
	case categories.Obstacle:
		s.obstacles = append(s.obstacles, collisionEntity{
			ID:      id,
			Spatial: spatial,
		})
	case categories.MovableObstacle:
		s.moveableObstacles = append(s.moveableObstacles, collisionEntity{
			ID:      id,
			Spatial: spatial,
		})
	case categories.CollisionSwitch:
		s.collisionSwitches = append(s.collisionSwitches, collisionEntity{
			ID:      id,
			Spatial: spatial,
		})
	case categories.Coin:
		s.coins = append(s.coins, collisionEntity{
			ID:      id,
			Spatial: spatial,
		})
	}
}

// Remove removes the entity from the system
func (s *Collision) Remove(category categories.Category, id entities.EntityID) {
	switch category {
	case categories.Coin:
		for i := len(s.coins) - 1; i >= 0; i-- {
			coin := s.coins[i]
			if coin.ID == id {
				s.coins = append(s.coins[:i], s.coins[i+1:]...)
			}
		}
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
func (s *Collision) RemoveAll(category categories.Category) {
	switch category {
	case categories.Enemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	case categories.CollisionSwitch:
		for i := len(s.collisionSwitches) - 1; i >= 0; i-- {
			s.collisionSwitches = append(s.collisionSwitches[:i], s.collisionSwitches[i+1:]...)
		}
	case categories.MovableObstacle:
		for i := len(s.moveableObstacles) - 1; i >= 0; i-- {
			s.moveableObstacles = append(s.moveableObstacles[:i], s.moveableObstacles[i+1:]...)
		}
	case categories.Obstacle:
		for i := len(s.obstacles) - 1; i >= 0; i-- {
			s.obstacles = append(s.obstacles[:i], s.obstacles[i+1:]...)
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
func (s *Collision) Update() {

	player := s.player
	playerR := s.player.Spatial.Rect
	mapBounds := s.MapBounds

	// is player at map edge?
	if player.Spatial.Rect.Min.Y <= mapBounds.Min.Y {
		s.PlayerCollisionWithBounds(bounds.Bottom)
	} else if player.Spatial.Rect.Min.X <= mapBounds.Min.X {
		s.PlayerCollisionWithBounds(bounds.Left)
	} else if player.Spatial.Rect.Max.X >= mapBounds.Max.X {
		s.PlayerCollisionWithBounds(bounds.Right)
	} else if player.Spatial.Rect.Max.Y >= mapBounds.Max.Y {
		s.PlayerCollisionWithBounds(bounds.Top)
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

	for _, obstacle := range s.obstacles {
		mod := player.Spatial.CollisionWithRectMod
		if isColliding(obstacle.Spatial.Rect, pixel.R(
			s.player.Spatial.Rect.Min.X+mod,
			s.player.Spatial.Rect.Min.Y+mod,
			s.player.Spatial.Rect.Max.X-mod,
			s.player.Spatial.Rect.Max.Y-mod,
		)) {
			s.PlayerCollisionWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			mod = enemy.Spatial.CollisionWithRectMod
			if isColliding(obstacle.Spatial.Rect, pixel.R(
				enemy.Spatial.Rect.Min.X+mod,
				enemy.Spatial.Rect.Min.Y+mod,
				enemy.Spatial.Rect.Max.X-mod,
				enemy.Spatial.Rect.Max.Y-mod,
			)) {
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
