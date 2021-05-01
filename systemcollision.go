package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type componentHitbox struct {
	HitBox               *imdraw.IMDraw
	HitBoxRadius         float64
	CollisionWithRectMod float64
}

func NewComponentHitbox(radius, collisionWithRectMod float64) *componentHitbox {
	return &componentHitbox{
		HitBox:               imdraw.New(nil),
		HitBoxRadius:         radius,
		CollisionWithRectMod: collisionWithRectMod,
	}
}

// componentInvincible is used to track if an enemy is immune to damage of all kinds
type componentInvincible struct {
	Enabled bool
}

func NewComponentInvincible(enabled bool) *componentInvincible {
	return &componentInvincible{
		Enabled: enabled,
	}
}

type collisionEntity struct {
	ID EntityID
	*componentInvincible
	*componentHitbox
	*componentDimensions
	*componentRectangle
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
		ID:                  entity.ID(),
		componentHitbox:     entity.componentHitbox,
		componentInvincible: entity.componentInvincible,
		componentDimensions: entity.componentDimensions,
		componentRectangle:  entity.componentRectangle,
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
		r.componentInvincible = entity.componentInvincible
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
	// DrawActiveSpace(s.Win, ActiveSpaceRectangle{
	// 	X:      s.MapBounds.Min.X,
	// 	Y:      s.MapBounds.Min.Y,
	// 	Width:  s.MapBounds.W(),
	// 	Height: s.MapBounds.H(),
	// })

	player := s.player
	mapBounds := s.MapBounds

	if player.componentRectangle.Rect.Min.Y <= mapBounds.Min.Y {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundBottom)
	} else if player.componentRectangle.Rect.Min.X <= mapBounds.Min.X {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundLeft)
	} else if player.componentRectangle.Rect.Max.X >= mapBounds.Max.X {
		s.CollisionHandler.OnPlayerCollisionWithBounds(BoundRight)
	} else if player.componentRectangle.Rect.Max.Y >= mapBounds.Max.Y {
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
	// playerR := s.player.componentRectangle.Rect

	w, h := player.componentDimensions.Width, player.componentDimensions.Height
	for _, enemy := range s.enemies {

		playerRect := player.componentRectangle.Rect

		s.drawHitbox(playerRect, player.componentHitbox.HitBoxRadius)

		enemyRect := enemy.componentRectangle.Rect

		s.drawHitbox(enemyRect, enemy.HitBoxRadius)

		// Check if player and enemy are colliding
		if isCircleCollision(
			player.componentHitbox.HitBoxRadius,
			enemy.componentHitbox.HitBoxRadius,
			w, h, playerRect, enemyRect) {
			s.CollisionHandler.OnPlayerCollisionWithEnemy(enemy.ID)
		}

		if !enemy.componentInvincible.Enabled {

			// Check if the player sword is colliding with the enemy
			if isCircleCollision(
				s.sword.componentHitbox.HitBoxRadius,
				enemy.componentHitbox.HitBoxRadius,
				w, h, s.sword.componentRectangle.Rect, enemyRect) {
				s.CollisionHandler.OnSwordCollisionWithEnemy(enemy.ID)
			}

			// Check if the player arrow is colliding with the enemy
			if isCircleCollision(
				s.arrow.componentHitbox.HitBoxRadius,
				enemy.componentHitbox.HitBoxRadius,
				w, h, s.arrow.componentRectangle.Rect, enemyRect) {
				s.CollisionHandler.OnArrowCollisionWithEnemy(enemy.ID)
			}
		}
	}
}

func (s *CollisionSystem) handleCoinCollisions() {
	for _, coin := range s.coins {
		if isColliding(coin.componentRectangle.Rect, s.player.componentRectangle.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithCoin(coin.ID)
		}
	}
}

func (s *CollisionSystem) handleObstacleCollisions() {
	player := s.player

	for _, obstacle := range s.obstacles {
		mod := player.componentHitbox.CollisionWithRectMod
		if isColliding(obstacle.componentRectangle.Rect, pixel.R(
			s.player.componentRectangle.Rect.Min.X+mod,
			s.player.componentRectangle.Rect.Min.Y+mod,
			s.player.componentRectangle.Rect.Max.X-mod,
			s.player.componentRectangle.Rect.Max.Y-mod,
		)) {
			s.CollisionHandler.OnPlayerCollisionWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			mod = enemy.componentHitbox.CollisionWithRectMod
			if isColliding(obstacle.componentRectangle.Rect, pixel.R(
				enemy.componentRectangle.Rect.Min.X+mod,
				enemy.componentRectangle.Rect.Min.Y+mod,
				enemy.componentRectangle.Rect.Max.X-mod,
				enemy.componentRectangle.Rect.Max.Y-mod,
			)) {
				s.CollisionHandler.OnEnemyCollisionWithObstacle(enemy.ID, obstacle.ID)
			}
		}

		if isColliding(obstacle.componentRectangle.Rect, s.arrow.componentRectangle.Rect) {
			s.CollisionHandler.OnArrowCollisionWithObstacle()
		}
	}
}

func (s *CollisionSystem) handleMoveableObstacleCollisions() {

	player := s.player

	for _, moveableObstacle := range s.moveableObstacles {
		if isColliding(moveableObstacle.componentRectangle.Rect, player.componentRectangle.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithMoveableObstacle(moveableObstacle.ID)
		}

		for _, collisionSwitch := range s.collisionSwitches {
			if isColliding(moveableObstacle.componentRectangle.Rect, collisionSwitch.componentRectangle.Rect) {
				s.CollisionHandler.OnMoveableObstacleCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnMoveableObstacleNoCollisionWithSwitch(collisionSwitch.ID)
			}
		}

		// for _, enemy := range s.enemies {
		// 	if isColliding(moveableObstacle.componentRectangle.Rect, enemy.componentRectangle.Rect) {
		// 		// s.EnemyCollisionWithMoveableObstacle(enemy.ID)
		// 	}
		// }

		// for _, obstacle := range s.obstacles {
		// 	if isColliding(moveableObstacle.componentRectangle.Rect, obstacle.componentRectangle.Rect) {
		// 		// s.MoveableObstacleCollisionWithObstacle(moveableObstacle.ID)
		// 	}
		// }

		if isColliding(moveableObstacle.componentRectangle.Rect, s.arrow.componentRectangle.Rect) {
			s.CollisionHandler.OnArrowCollisionWithObstacle()
		}
	}
}

func (s *CollisionSystem) handleSwitchCollisions() {

	player := s.player

	for _, collisionSwitch := range s.collisionSwitches {
		if collisionSwitch.componentHitbox.HitBoxRadius > 0 {
			w, h := player.componentDimensions.Width, player.componentDimensions.Height
			if isCircleCollision(
				s.player.componentHitbox.HitBoxRadius,
				collisionSwitch.componentHitbox.HitBoxRadius,
				w, h, player.componentRectangle.Rect, collisionSwitch.componentRectangle.Rect) {
				s.CollisionHandler.OnPlayerCollisionWithSwitch(collisionSwitch.ID)
			} else {
				s.CollisionHandler.OnPlayerNoCollisionWithSwitch(collisionSwitch.ID)
			}
		} else {
			if isColliding(player.componentRectangle.Rect, collisionSwitch.componentRectangle.Rect) {
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
		if isColliding(player.componentRectangle.Rect, warp.componentRectangle.Rect) {
			s.CollisionHandler.OnPlayerCollisionWithWarp(warp.ID)
		}
	}
}
