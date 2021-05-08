package zelduh

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/core/entity"
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
	ID entity.EntityID
	*componentInvincible
	*componentHitbox
	*componentDimensions
	*componentRectangle
}

type OnCollisionHandlers struct {
	PlayerWithEnemy                       func(entity.EntityID)
	PlayerWithCoin                        func(entity.EntityID)
	SwordWithEnemy                        func(entity.EntityID)
	ArrowWithEnemy                        func(entity.EntityID)
	PlayerWithObstacle                    func(entity.EntityID)
	ArrowWithObstacle                     func(entity.EntityID)
	EnemyWithObstacle                     func(entity.EntityID)
	PlayerWithMoveableObstacle            func(entity.EntityID)
	MoveableObstacleWithSwitch            func(entity.EntityID)
	MoveableObstacleWithSwitchNoCollision func(entity.EntityID)
	PlayerWithSwitch                      func(entity.EntityID)
	PlayerWithSwitchNoCollision           func(entity.EntityID)
	PlayerWithWarp                        func(entity.EntityID)
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
	ActiveSpaceRectangle ActiveSpaceRectangle
	Win                  *pixelgl.Window
	onCollisionHandlers  *OnCollisionHandlers
}

func NewCollisionSystem(
	mapBounds pixel.Rect,
	activeSpaceRectangle ActiveSpaceRectangle,
	win *pixelgl.Window,
	onCollisionHandlers *OnCollisionHandlers,
) CollisionSystem {
	return CollisionSystem{
		MapBounds:            mapBounds,
		ActiveSpaceRectangle: activeSpaceRectangle,
		Win:                  win,
		onCollisionHandlers:  onCollisionHandlers,
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
func (s *CollisionSystem) Remove(category entity.EntityCategory, id entity.EntityID) {
	switch category {
	case CategoryCoin:
		removeEntityFromCollection(s.coins, id)
	case CategoryEnemy:
		removeEntityFromCollection(s.enemies, id)
	}
}

func removeEntityFromCollection(entities []collisionEntity, entityIDToRemove entity.EntityID) {
	for i := len(entities) - 1; i >= 0; i-- {
		entity := entities[i]
		if entity.ID == entityIDToRemove {
			entities = append(entities[:i], entities[i+1:]...)
		}
	}
}

// RemoveAll removes all entities from one category
func (s *CollisionSystem) RemoveAll(category entity.EntityCategory) {
	switch category {
	case CategoryEnemy:
		removeAllEntities(s.enemies)
	case CategoryCollisionSwitch:
		removeAllEntities(s.collisionSwitches)
	case CategoryMovableObstacle:
		removeAllEntities(s.moveableObstacles)
	case CategoryObstacle:
		// removeAllEntities(s.obstacles)

		for i := len(s.obstacles) - 1; i >= 0; i-- {
			s.obstacles = append(s.obstacles[:i], s.obstacles[i+1:]...)
		}
	}
}

func removeAllEntities(entities []collisionEntity) {
	for i := len(entities) - 1; i >= 0; i-- {
		entities = append(entities[:i], entities[i+1:]...)
	}
}

// Update checks for collisions
func (s *CollisionSystem) Update() error {
	s.handleEnemyCollisions()
	s.handleCoinCollisions()
	s.handleObstacleCollisions()
	s.handleMoveableObstacleCollisions()
	s.handleSwitchCollisions()
	s.handleWarpCollisions()
	return nil
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
	// fmt.Printf("total enemies %d\n", len(s.enemies))
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
			s.onCollisionHandlers.PlayerWithEnemy(enemy.ID)
		}

		if !enemy.componentInvincible.Enabled {

			// Check if the player sword is colliding with the enemy
			if isCircleCollision(
				s.sword.componentHitbox.HitBoxRadius,
				enemy.componentHitbox.HitBoxRadius,
				w, h, s.sword.componentRectangle.Rect, enemyRect) {
				s.onCollisionHandlers.SwordWithEnemy(enemy.ID)
			}

			// Check if the player arrow is colliding with the enemy
			if isCircleCollision(
				s.arrow.componentHitbox.HitBoxRadius,
				enemy.componentHitbox.HitBoxRadius,
				w, h, s.arrow.componentRectangle.Rect, enemyRect) {
				s.onCollisionHandlers.ArrowWithEnemy(enemy.ID)
			}
		}
	}
}

func (s *CollisionSystem) handleCoinCollisions() {
	for _, coin := range s.coins {
		if isColliding(coin.componentRectangle.Rect, s.player.componentRectangle.Rect) {
			s.onCollisionHandlers.PlayerWithCoin(coin.ID)
		}
	}
}

func (s *CollisionSystem) handleObstacleCollisions() {
	player := s.player

	fmt.Printf("CollisionSystem total obstacles %d\n", len(s.obstacles))
	for _, obstacle := range s.obstacles {
		mod := player.componentHitbox.CollisionWithRectMod
		if isColliding(obstacle.componentRectangle.Rect, pixel.R(
			s.player.componentRectangle.Rect.Min.X+mod,
			s.player.componentRectangle.Rect.Min.Y+mod,
			s.player.componentRectangle.Rect.Max.X-mod,
			s.player.componentRectangle.Rect.Max.Y-mod,
		)) {
			s.onCollisionHandlers.PlayerWithObstacle(obstacle.ID)
		}

		for _, enemy := range s.enemies {
			mod = enemy.componentHitbox.CollisionWithRectMod
			if isColliding(obstacle.componentRectangle.Rect, pixel.R(
				enemy.componentRectangle.Rect.Min.X+mod,
				enemy.componentRectangle.Rect.Min.Y+mod,
				enemy.componentRectangle.Rect.Max.X-mod,
				enemy.componentRectangle.Rect.Max.Y-mod,
			)) {
				s.onCollisionHandlers.EnemyWithObstacle(enemy.ID)
			}
		}

		if isColliding(obstacle.componentRectangle.Rect, s.arrow.componentRectangle.Rect) {
			s.onCollisionHandlers.ArrowWithEnemy(s.arrow.ID)
		}
	}
}

func (s *CollisionSystem) handleMoveableObstacleCollisions() {

	player := s.player

	for _, moveableObstacle := range s.moveableObstacles {
		if isColliding(moveableObstacle.componentRectangle.Rect, player.componentRectangle.Rect) {
			s.onCollisionHandlers.PlayerWithMoveableObstacle(moveableObstacle.ID)
		}

		for _, collisionSwitch := range s.collisionSwitches {
			if isColliding(moveableObstacle.componentRectangle.Rect, collisionSwitch.componentRectangle.Rect) {
				s.onCollisionHandlers.MoveableObstacleWithSwitch(collisionSwitch.ID)
			} else {
				s.onCollisionHandlers.MoveableObstacleWithSwitchNoCollision(collisionSwitch.ID)
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
			s.onCollisionHandlers.ArrowWithObstacle(s.arrow.ID)
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
				s.onCollisionHandlers.PlayerWithSwitch(collisionSwitch.ID)
			} else {
				s.onCollisionHandlers.PlayerWithSwitchNoCollision(collisionSwitch.ID)
			}
		} else {
			if isColliding(player.componentRectangle.Rect, collisionSwitch.componentRectangle.Rect) {
				s.onCollisionHandlers.PlayerWithSwitch(collisionSwitch.ID)
			} else {
				s.onCollisionHandlers.PlayerWithSwitchNoCollision(collisionSwitch.ID)
			}
		}

	}
}

func (s *CollisionSystem) handleWarpCollisions() {

	player := s.player

	for _, warp := range s.warps {
		if isColliding(player.componentRectangle.Rect, warp.componentRectangle.Rect) {
			s.onCollisionHandlers.PlayerWithWarp(warp.ID)
		}
	}
}
