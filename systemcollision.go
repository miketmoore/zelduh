package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type collisionEntity struct {
	ID EntityID
	*ComponentSpatial
	*ComponentInvincible
}

// CollisionSystem is a custom system for detecting collisions and what to do when they occur
type CollisionSystem struct {
	MapBounds            pixel.Rect
	player               collisionEntity
	sword                collisionEntity
	arrow                collisionEntity
	enemies              []collisionEntity
	coins                []collisionEntity
	obstacles            []collisionEntity
	moveableObstacles    []collisionEntity
	collisionSwitches    []collisionEntity
	warps                []collisionEntity
	CollisionHandler     *CollisionHandler
	ActiveSpaceRectangle ActiveSpaceRectangle
	Win                  *pixelgl.Window
}

func NewCollisionSystem(
	mapBounds pixel.Rect,
	collisionHandler *CollisionHandler,
	activeSpaceRectangle ActiveSpaceRectangle,
	win *pixelgl.Window,
) CollisionSystem {
	return CollisionSystem{
		MapBounds:            mapBounds,
		CollisionHandler:     collisionHandler,
		ActiveSpaceRectangle: activeSpaceRectangle,
		Win:                  win,
	}
}

// AddEntity adds an entity to the system
func (s *CollisionSystem) AddEntity(entity Entity) {
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
func (s *CollisionSystem) Remove(category EntityCategory, id EntityID) {
	switch category {
	case CategoryCoin:
		removeEntityFromCollection(s.coins, id)
	case CategoryEnemy:
		removeEntityFromCollection(s.enemies, id)
	}
}

func removeEntityFromCollection(entities []collisionEntity, entityIDToRemove EntityID) {
	for i := len(entities) - 1; i >= 0; i-- {
		entity := entities[i]
		if entity.ID == entityIDToRemove {
			entities = append(entities[:i], entities[i+1:]...)
		}
	}
}

// RemoveAll removes all entities from one category
func (s *CollisionSystem) RemoveAll(category EntityCategory) {
	switch category {
	case CategoryEnemy:
		removeAllEntities(s.enemies)
	case CategoryCollisionSwitch:
		removeAllEntities(s.collisionSwitches)
	case CategoryMovableObstacle:
		removeAllEntities(s.moveableObstacles)
	case CategoryObstacle:
		removeAllEntities(s.obstacles)
	}
}

func removeAllEntities(entities []collisionEntity) {
	for i := len(entities) - 1; i >= 0; i-- {
		entities = append(entities[:i], entities[i+1:]...)
	}
}

// Update checks for collisions
func (s *CollisionSystem) Update() error {
	s.handlePlayerAtMapEdge()
	s.handleEnemyCollisions()
	s.handleCoinCollisions()
	s.handleObstacleCollisions()
	s.handleMoveableObstacleCollisions()
	s.handleSwitchCollisions()
	s.handleWarpCollisions()
	return nil
}

func (s *CollisionSystem) handlePlayerAtMapEdge() {

	player := s.player
	mapBounds := s.MapBounds

	if player.ComponentSpatial.Rect.Min.Y <= mapBounds.Min.Y {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundBottom)
	} else if player.ComponentSpatial.Rect.Min.X <= mapBounds.Min.X {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundLeft)
	} else if player.ComponentSpatial.Rect.Max.X >= mapBounds.Max.X {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundRight)
	} else if player.ComponentSpatial.Rect.Max.Y >= mapBounds.Max.Y {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundTop)
	}
}

func (s *CollisionSystem) drawHitbox(rect pixel.Rect, radius float64) {

	circle := imdraw.New(nil)
	circle.Color = colornames.Blue
	circle.Push(rect.Center())

	circle.Circle(radius, 5)
	circle.Draw(s.Win)
}

func (s *CollisionSystem) handleEnemyCollisions() {

	player := s.player
	// playerR := s.player.ComponentSpatial.Rect

	w, h := player.ComponentSpatial.Width, player.ComponentSpatial.Height
	for _, enemy := range s.enemies {
		// enemyR := enemy.ComponentSpatial.Rect

		// v := s.buildSpriteVector(player.ComponentSpatial)
		playerRect := player.ComponentSpatial.Rect
		// playerVector := pixel.V(
		// 	playerRect.Min.X+player.ComponentSpatial.Width/2,
		// 	playerRect.Min.Y+player.ComponentSpatial.Height/2,
		// )
		// m := s.buildSpriteMatrix(player.ComponentSpatial, v)
		s.drawHitbox(playerRect, player.HitBoxRadius)

		enemyRect := enemy.ComponentSpatial.Rect
		// enemyVector := pixel.V(
		// 	enemyRect.Min.X+player.ComponentSpatial.Width/2+200,
		// 	enemyRect.Min.Y+player.ComponentSpatial.Height/2+200,
		// )
		// enemyRect.Moved(enemyVector)
		s.drawHitbox(enemyRect, enemy.HitBoxRadius)

		// shape := imdraw.New(nil)
		// shape.Color = colornames.Blue
		// shape.Push(pixel.V(
		// 	enemyRect.Min.X+s.CollisionHandler.TileSize,
		// 	enemyRect.Min.Y+s.CollisionHandler.TileSize,
		// ))
		// shape.Push(pixel.V(
		// 	enemyRect.Max.X+s.CollisionHandler.TileSize,
		// 	enemyRect.Max.Y+s.CollisionHandler.TileSize,
		// ))
		// shape.Rectangle(3)
		// shape.Draw(s.Win)

		// circle := imdraw.New(nil)
		// circle.Color = colornames.Blue
		// circle.Push(pixel.V(
		// 	enemyRect.Min.X+s.CollisionHandler.TileSize,
		// 	enemyRect.Min.Y+s.CollisionHandler.TileSize,
		// ))

		// circle.Circle(enemy.HitBoxRadius, 5)
		// circle.Draw(s.Win)

		// Check if player and enemy are colliding
		if isCircleCollision(
			player.ComponentSpatial.HitBoxRadius,
			enemy.ComponentSpatial.HitBoxRadius,
			w, h, playerRect, enemyRect) {
			s.CollisionHandler.OnPlayerCollisionWithEnemy(enemy.ID)
		}

		if !enemy.ComponentInvincible.Enabled {

			// Check if the player sword is colliding with the enemy
			if isCircleCollision(
				s.sword.ComponentSpatial.HitBoxRadius,
				enemy.ComponentSpatial.HitBoxRadius,
				w, h, s.sword.ComponentSpatial.Rect, enemyRect) {
				s.CollisionHandler.OnSwordCollisionWithEnemy(enemy.ID)
			}

			// Check if the player arrow is colliding with the enemy
			if isCircleCollision(
				s.arrow.ComponentSpatial.HitBoxRadius,
				enemy.ComponentSpatial.HitBoxRadius,
				w, h, s.arrow.ComponentSpatial.Rect, enemyRect) {
				s.CollisionHandler.OnArrowCollisionWithEnemy(enemy.ID)
			}
		}
	}
}

func (s *CollisionSystem) handleCoinCollisions() {
	for _, coin := range s.coins {
		if isColliding(coin.ComponentSpatial.Rect, s.player.ComponentSpatial.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithCoin(coin.ID)
		}
	}
}

func (s *CollisionSystem) handleObstacleCollisions() {
	player := s.player

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
}

func (s *CollisionSystem) handleMoveableObstacleCollisions() {

	player := s.player

	for _, moveableObstacle := range s.moveableObstacles {
		if isColliding(moveableObstacle.ComponentSpatial.Rect, player.ComponentSpatial.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithMoveableObstacle(moveableObstacle.ID)
		}

		for _, collisionSwitch := range s.collisionSwitches {
			if isColliding(moveableObstacle.ComponentSpatial.Rect, collisionSwitch.ComponentSpatial.Rect) {
				s.CollisionHandler.OnMoveableObstacleCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnMoveableObstacleNoCollisionWithSwitch(collisionSwitch.ID)
			}
		}

		// for _, enemy := range s.enemies {
		// 	if isColliding(moveableObstacle.ComponentSpatial.Rect, enemy.ComponentSpatial.Rect) {
		// 		// s.EnemyCollisionWithMoveableObstacle(enemy.ID)
		// 	}
		// }

		// for _, obstacle := range s.obstacles {
		// 	if isColliding(moveableObstacle.ComponentSpatial.Rect, obstacle.ComponentSpatial.Rect) {
		// 		// s.MoveableObstacleCollisionWithObstacle(moveableObstacle.ID)
		// 	}
		// }

		if isColliding(moveableObstacle.ComponentSpatial.Rect, s.arrow.ComponentSpatial.Rect) {
			s.CollisionHandler.OnArrowCollisionWithObstacle()
		}
	}
}

func (s *CollisionSystem) handleSwitchCollisions() {

	player := s.player

	for _, collisionSwitch := range s.collisionSwitches {
		if collisionSwitch.ComponentSpatial.HitBoxRadius > 0 {
			w, h := player.ComponentSpatial.Width, player.ComponentSpatial.Height
			if isCircleCollision(
				s.player.ComponentSpatial.HitBoxRadius,
				collisionSwitch.ComponentSpatial.HitBoxRadius,
				w, h, player.ComponentSpatial.Rect, collisionSwitch.ComponentSpatial.Rect) {
				s.CollisionHandler.OnPlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnPlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		} else {
			if isColliding(player.ComponentSpatial.Rect, collisionSwitch.ComponentSpatial.Rect) {
				s.CollisionHandler.OnPlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnPlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		}

	}
}

func (s *CollisionSystem) handleWarpCollisions() {

	player := s.player

	for _, warp := range s.warps {
		if isColliding(player.ComponentSpatial.Rect, warp.ComponentSpatial.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithWarp(warp.ID)
		}
	}
}
