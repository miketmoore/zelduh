package systems

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh/bounds"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/entities"
)

type collisionEntity struct {
	ID terraform2d.EntityID
	*components.Spatial
	*components.Invincible
}

// Collision is a custom system for detecting collisions and what to do when they occur
type Collision struct {
	MapBounds                               pixel.Rect
	player                                  collisionEntity
	sword                                   collisionEntity
	arrow                                   collisionEntity
	enemies                                 []collisionEntity
	coins                                   []collisionEntity
	obstacles                               []collisionEntity
	moveableObstacles                       []collisionEntity
	collisionSwitches                       []collisionEntity
	warps                                   []collisionEntity
	OnPlayerCollisionWithCoin               func(terraform2d.EntityID)
	OnPlayerCollisionWithEnemy              func(terraform2d.EntityID)
	OnSwordCollisionWithEnemy               func(terraform2d.EntityID)
	OnArrowCollisionWithEnemy               func(terraform2d.EntityID)
	OnArrowCollisionWithObstacle            func()
	OnPlayerCollisionWithObstacle           func(terraform2d.EntityID)
	OnPlayerCollisionWithMoveableObstacle   func(terraform2d.EntityID)
	OnEnemyCollisionWithObstacle            func(terraform2d.EntityID, terraform2d.EntityID)
	OnEnemyCollisionWithMoveableObstacle    func(terraform2d.EntityID)
	OnMoveableObstacleCollisionWithObstacle func(terraform2d.EntityID)
	OnPlayerCollisionWithSwitch             func(terraform2d.EntityID)
	OnPlayerNoCollisionWithSwitch           func(terraform2d.EntityID)
	OnPlayerCollisionWithBounds             func(bounds.Bound)
	OnMoveableObstacleCollisionWithSwitch   func(terraform2d.EntityID)
	OnMoveableObstacleNoCollisionWithSwitch func(terraform2d.EntityID)
	OnPlayerCollisionWithWarp               func(terraform2d.EntityID)
}

// AddEntity adds an entity to the system
func (s *Collision) AddEntity(entity entities.Entity) {
	r := collisionEntity{
		ID:      entity.ID(),
		Spatial: entity.Spatial,
	}
	switch entity.Category {
	case categories.Player:
		s.player = r
	case categories.Sword:
		s.sword = r
	case categories.Arrow:
		s.arrow = r
	case categories.MovableObstacle:
		s.moveableObstacles = append(s.moveableObstacles, r)
	case categories.CollisionSwitch:
		s.collisionSwitches = append(s.collisionSwitches, r)
	case categories.Warp:
		s.warps = append(s.warps, r)
	case categories.Enemy:
		r.Invincible = entity.Invincible
		s.enemies = append(s.enemies, r)
	case categories.Coin:
		s.coins = append(s.coins, r)
	case categories.Obstacle:
		s.obstacles = append(s.obstacles, r)
	}
}

// Remove removes the entity from the system
func (s *Collision) Remove(category categories.Category, id terraform2d.EntityID) {
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
		s.OnPlayerCollisionWithBounds(bounds.Bottom)
	} else if player.Spatial.Rect.Min.X <= mapBounds.Min.X {
		s.OnPlayerCollisionWithBounds(bounds.Left)
	} else if player.Spatial.Rect.Max.X >= mapBounds.Max.X {
		s.OnPlayerCollisionWithBounds(bounds.Right)
	} else if player.Spatial.Rect.Max.Y >= mapBounds.Max.Y {
		s.OnPlayerCollisionWithBounds(bounds.Top)
	}

	w, h := player.Spatial.Width, player.Spatial.Height
	for _, enemy := range s.enemies {
		enemyR := enemy.Spatial.Rect

		if isCircleCollision(
			player.Spatial.HitBoxRadius,
			enemy.Spatial.HitBoxRadius,
			w, h, playerR, enemyR) {
			s.OnPlayerCollisionWithEnemy(enemy.ID)
		}

		if !enemy.Invincible.Enabled {
			if isCircleCollision(
				s.sword.Spatial.HitBoxRadius,
				enemy.Spatial.HitBoxRadius,
				w, h, s.sword.Spatial.Rect, enemyR) {
				s.OnSwordCollisionWithEnemy(enemy.ID)
			}

			if isCircleCollision(
				s.arrow.Spatial.HitBoxRadius,
				enemy.Spatial.HitBoxRadius,
				w, h, s.arrow.Spatial.Rect, enemyR) {
				s.OnArrowCollisionWithEnemy(enemy.ID)
			}
		}
	}
	for _, coin := range s.coins {
		if isColliding(coin.Spatial.Rect, s.player.Spatial.Rect) {
			s.OnPlayerCollisionWithCoin(coin.ID)
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
			s.OnPlayerCollisionWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			mod = enemy.Spatial.CollisionWithRectMod
			if isColliding(obstacle.Spatial.Rect, pixel.R(
				enemy.Spatial.Rect.Min.X+mod,
				enemy.Spatial.Rect.Min.Y+mod,
				enemy.Spatial.Rect.Max.X-mod,
				enemy.Spatial.Rect.Max.Y-mod,
			)) {
				s.OnEnemyCollisionWithObstacle(enemy.ID, obstacle.ID)
			}
		}

		if isColliding(obstacle.Spatial.Rect, s.arrow.Spatial.Rect) {
			s.OnArrowCollisionWithObstacle()
		}
	}
	for _, moveableObstacle := range s.moveableObstacles {
		if isColliding(moveableObstacle.Spatial.Rect, s.player.Spatial.Rect) {
			s.OnPlayerCollisionWithMoveableObstacle(moveableObstacle.ID)
		}

		for _, collisionSwitch := range s.collisionSwitches {
			if isColliding(moveableObstacle.Spatial.Rect, collisionSwitch.Spatial.Rect) {
				s.OnMoveableObstacleCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.OnMoveableObstacleNoCollisionWithSwitch(collisionSwitch.ID)
			}
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
			s.OnArrowCollisionWithObstacle()
		}
	}

	for _, collisionSwitch := range s.collisionSwitches {
		if collisionSwitch.Spatial.HitBoxRadius > 0 {
			if isCircleCollision(
				s.player.Spatial.HitBoxRadius,
				collisionSwitch.Spatial.HitBoxRadius,
				w, h, s.player.Spatial.Rect, collisionSwitch.Spatial.Rect) {
				s.OnPlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.OnPlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		} else {
			if isColliding(s.player.Spatial.Rect, collisionSwitch.Spatial.Rect) {
				s.OnPlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.OnPlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		}

	}

	for _, warp := range s.warps {
		if isColliding(s.player.Spatial.Rect, warp.Spatial.Rect) {
			s.OnPlayerCollisionWithWarp(warp.ID)
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
