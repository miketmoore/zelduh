package zelduh

import (
	"math"

	"github.com/faiface/pixel"
)

type collisionEntity struct {
	ID EntityID
	*ComponentSpatial
	*ComponentInvincible
}

// SystemCollision is a custom system for detecting collisions and what to do when they occur
type SystemCollision struct {
	MapBounds         pixel.Rect
	player            collisionEntity
	sword             collisionEntity
	arrow             collisionEntity
	enemies           []collisionEntity
	coins             []collisionEntity
	obstacles         []collisionEntity
	moveableObstacles []collisionEntity
	collisionSwitches []collisionEntity
	warps             []collisionEntity
	CollisionHandler  CollisionHandler
}

func NewSystemCollision(
	mapBoundsConfig MapBoundsConfig,
	roomTransitionManager *RoomTransitionManager,
	systemsManager *SystemsManager,
	healthSystem *SystemHealth,
	spatialSystem *SystemSpatial,
	entitiesMap EntityByEntityID,
	roomWarps map[EntityID]EntityConfig,
	entities Entities,
) SystemCollision {
	return SystemCollision{
		MapBounds: pixel.R(
			mapBoundsConfig.X,
			mapBoundsConfig.Y,
			mapBoundsConfig.Width,
			mapBoundsConfig.Height,
		),
		CollisionHandler: CollisionHandler{
			RoomTransitionManager: roomTransitionManager,
			SystemsManager:        systemsManager,
			HealthSystem:          healthSystem,
			SpatialSystem:         spatialSystem,
			EntitiesMap:           entitiesMap,
			RoomWarps:             roomWarps,
			Entities:              entities,
		},
	}
}

// AddEntity adds an entity to the system
func (s *SystemCollision) AddEntity(entity Entity) {
	r := collisionEntity{
		ID:               entity.ID(),
		ComponentSpatial: entity.ComponentSpatial,
	}
	switch entity.Category {
	case CategoryPlayer:
		s.player = r
	case CategorySword:
		s.sword = r
	case CategoryArrow:
		s.arrow = r
	case CategoryMovableObstacle:
		s.moveableObstacles = append(s.moveableObstacles, r)
	case CategoryCollisionSwitch:
		s.collisionSwitches = append(s.collisionSwitches, r)
	case CategoryWarp:
		s.warps = append(s.warps, r)
	case CategoryEnemy:
		r.ComponentInvincible = entity.ComponentInvincible
		s.enemies = append(s.enemies, r)
	case CategoryCoin:
		s.coins = append(s.coins, r)
	case CategoryObstacle:
		s.obstacles = append(s.obstacles, r)
	}
}

// Remove removes the entity from the system
func (s *SystemCollision) Remove(category EntityCategory, id EntityID) {
	switch category {
	case CategoryCoin:
		for i := len(s.coins) - 1; i >= 0; i-- {
			coin := s.coins[i]
			if coin.ID == id {
				s.coins = append(s.coins[:i], s.coins[i+1:]...)
			}
		}
	case CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			enemy := s.enemies[i]
			if enemy.ID == id {
				s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
			}
		}
	}
}

// RemoveAll removes all entities from one category
func (s *SystemCollision) RemoveAll(category EntityCategory) {
	switch category {
	case CategoryEnemy:
		for i := len(s.enemies) - 1; i >= 0; i-- {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	case CategoryCollisionSwitch:
		for i := len(s.collisionSwitches) - 1; i >= 0; i-- {
			s.collisionSwitches = append(s.collisionSwitches[:i], s.collisionSwitches[i+1:]...)
		}
	case CategoryMovableObstacle:
		for i := len(s.moveableObstacles) - 1; i >= 0; i-- {
			s.moveableObstacles = append(s.moveableObstacles[:i], s.moveableObstacles[i+1:]...)
		}
	case CategoryObstacle:
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
func (s *SystemCollision) Update() {

	player := s.player
	playerR := s.player.ComponentSpatial.Rect
	mapBounds := s.MapBounds

	// is player at map edge?
	if player.ComponentSpatial.Rect.Min.Y <= mapBounds.Min.Y {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundBottom)
	} else if player.ComponentSpatial.Rect.Min.X <= mapBounds.Min.X {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundLeft)
	} else if player.ComponentSpatial.Rect.Max.X >= mapBounds.Max.X {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundRight)
	} else if player.ComponentSpatial.Rect.Max.Y >= mapBounds.Max.Y {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundTop)
	}

	w, h := player.ComponentSpatial.Width, player.ComponentSpatial.Height
	for _, enemy := range s.enemies {
		enemyR := enemy.ComponentSpatial.Rect

		if isCircleCollision(
			player.ComponentSpatial.HitBoxRadius,
			enemy.ComponentSpatial.HitBoxRadius,
			w, h, playerR, enemyR) {
			s.CollisionHandler.OnPlayerCollisionWithEnemy(enemy.ID)
		}

		if !enemy.ComponentInvincible.Enabled {
			if isCircleCollision(
				s.sword.ComponentSpatial.HitBoxRadius,
				enemy.ComponentSpatial.HitBoxRadius,
				w, h, s.sword.ComponentSpatial.Rect, enemyR) {
				s.CollisionHandler.OnSwordCollisionWithEnemy(enemy.ID)
			}

			if isCircleCollision(
				s.arrow.ComponentSpatial.HitBoxRadius,
				enemy.ComponentSpatial.HitBoxRadius,
				w, h, s.arrow.ComponentSpatial.Rect, enemyR) {
				s.CollisionHandler.OnArrowCollisionWithEnemy(enemy.ID)
			}
		}
	}
	for _, coin := range s.coins {
		if isColliding(coin.ComponentSpatial.Rect, s.player.ComponentSpatial.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithCoin(coin.ID)
		}
	}

	for _, obstacle := range s.obstacles {
		mod := player.ComponentSpatial.CollisionWithRectMod
		if isColliding(obstacle.ComponentSpatial.Rect, pixel.R(
			s.player.ComponentSpatial.Rect.Min.X+mod,
			s.player.ComponentSpatial.Rect.Min.Y+mod,
			s.player.ComponentSpatial.Rect.Max.X-mod,
			s.player.ComponentSpatial.Rect.Max.Y-mod,
		)) {
			s.CollisionHandler.OnPlayerCollisionWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			mod = enemy.ComponentSpatial.CollisionWithRectMod
			if isColliding(obstacle.ComponentSpatial.Rect, pixel.R(
				enemy.ComponentSpatial.Rect.Min.X+mod,
				enemy.ComponentSpatial.Rect.Min.Y+mod,
				enemy.ComponentSpatial.Rect.Max.X-mod,
				enemy.ComponentSpatial.Rect.Max.Y-mod,
			)) {
				s.CollisionHandler.OnEnemyCollisionWithObstacle(enemy.ID, obstacle.ID)
			}
		}

		if isColliding(obstacle.ComponentSpatial.Rect, s.arrow.ComponentSpatial.Rect) {
			s.CollisionHandler.OnArrowCollisionWithObstacle()
		}
	}
	for _, moveableObstacle := range s.moveableObstacles {
		if isColliding(moveableObstacle.ComponentSpatial.Rect, s.player.ComponentSpatial.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithMoveableObstacle(moveableObstacle.ID)
		}

		for _, collisionSwitch := range s.collisionSwitches {
			if isColliding(moveableObstacle.ComponentSpatial.Rect, collisionSwitch.ComponentSpatial.Rect) {
				s.CollisionHandler.OnMoveableObstacleCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnMoveableObstacleNoCollisionWithSwitch(collisionSwitch.ID)
			}
		}

		for _, enemy := range s.enemies {
			if isColliding(moveableObstacle.ComponentSpatial.Rect, enemy.ComponentSpatial.Rect) {
				// s.EnemyCollisionWithMoveableObstacle(enemy.ID)
			}
		}

		for _, obstacle := range s.obstacles {
			if isColliding(moveableObstacle.ComponentSpatial.Rect, obstacle.ComponentSpatial.Rect) {
				// s.MoveableObstacleCollisionWithObstacle(moveableObstacle.ID)
			}
		}

		if isColliding(moveableObstacle.ComponentSpatial.Rect, s.arrow.ComponentSpatial.Rect) {
			s.CollisionHandler.OnArrowCollisionWithObstacle()
		}
	}

	for _, collisionSwitch := range s.collisionSwitches {
		if collisionSwitch.ComponentSpatial.HitBoxRadius > 0 {
			if isCircleCollision(
				s.player.ComponentSpatial.HitBoxRadius,
				collisionSwitch.ComponentSpatial.HitBoxRadius,
				w, h, s.player.ComponentSpatial.Rect, collisionSwitch.ComponentSpatial.Rect) {
				s.CollisionHandler.OnPlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnPlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		} else {
			if isColliding(s.player.ComponentSpatial.Rect, collisionSwitch.ComponentSpatial.Rect) {
				s.CollisionHandler.OnPlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnPlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		}

	}

	for _, warp := range s.warps {
		if isColliding(s.player.ComponentSpatial.Rect, warp.ComponentSpatial.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithWarp(warp.ID)
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
