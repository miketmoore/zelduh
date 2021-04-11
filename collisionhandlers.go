package zelduh

import (
	"github.com/faiface/pixel"
)

// CollisionHandler contains collision handlers
type CollisionHandler struct {
	RoomTransitionManager *RoomTransitionManager
	SystemsManager        *SystemsManager
	HealthSystem          *SystemHealth
	SpatialSystem         *SystemSpatial
	EntitiesMap           EntityByEntityID
	RoomWarps             map[EntityID]EntityConfig
	Entities              Entities
	GameStateManager      GameStateManager
	RoomData              RoomData
	FrameRate             int
}

// OnPlayerCollisionWithBounds handles collisions between player and bounds
func (ch *CollisionHandler) OnPlayerCollisionWithBounds(side Bound) {
	if !ch.RoomTransitionManager.Active() {
		ch.RoomTransitionManager.SetSlideStart(side)
		ch.GameStateManager.CurrentState = StateMapTransition
		ch.SystemsManager.SetShouldAddEntities(true)
	}
}

// OnPlayerCollisionWithCoin handles collision between player and coin
func (ch *CollisionHandler) OnPlayerCollisionWithCoin(coinID EntityID) {
	ch.Entities.Player.ComponentCoins.Coins++
	ch.SystemsManager.Remove(CategoryCoin, coinID)
}

// OnPlayerCollisionWithEnemy handles collision between player and enemy
func (ch *CollisionHandler) OnPlayerCollisionWithEnemy(enemyID EntityID) {
	// TODO repeat what I did with the enemies
	ch.SpatialSystem.MovePlayerBack()
	ch.Entities.Player.ComponentHealth.Total--

	// remove heart entity
	heartIndex := len(ch.Entities.Hearts) - 1
	ch.SystemsManager.Remove(CategoryHeart, ch.Entities.Hearts[heartIndex].ID())
	ch.Entities.Hearts = append(ch.Entities.Hearts[:heartIndex], ch.Entities.Hearts[heartIndex+1:]...)

	if ch.Entities.Player.ComponentHealth.Total == 0 {
		ch.GameStateManager.CurrentState = StateOver
	}
}

func dropCoin(v pixel.Vec, systemsManager *SystemsManager, frameRate int) {
	coin := BuildEntityFromConfig(GetPreset("coin")(v.X/TileSize, v.Y/TileSize), systemsManager.NewEntityID(), frameRate)
	systemsManager.AddEntity(coin)
}

// OnSwordCollisionWithEnemy handles collision between sword and enemy
func (ch *CollisionHandler) OnSwordCollisionWithEnemy(enemyID EntityID) {
	if !ch.Entities.Sword.ComponentIgnore.Value {
		dead := false
		if !ch.SpatialSystem.EnemyMovingFromHit(enemyID) {
			dead = ch.HealthSystem.Hit(enemyID, 1)
			if dead {
				enemySpatial, _ := ch.SpatialSystem.GetEnemySpatial(enemyID)
				ch.Entities.Explosion.ComponentTemporary.Expiration = len(ch.Entities.Explosion.ComponentAnimation.Map["default"].Frames)
				ch.Entities.Explosion.ComponentSpatial = &ComponentSpatial{
					Width:  TileSize,
					Height: TileSize,
					Rect:   enemySpatial.Rect,
				}
				ch.Entities.Explosion.ComponentTemporary.OnExpiration = func() {
					dropCoin(ch.Entities.Explosion.ComponentSpatial.Rect.Min, ch.SystemsManager, ch.FrameRate)
				}
				ch.SystemsManager.AddEntity(ch.Entities.Explosion)
				ch.SystemsManager.RemoveEnemy(enemyID)
			} else {
				ch.SpatialSystem.MoveEnemyBack(enemyID, ch.Entities.Player.ComponentMovement.Direction)
			}
		}

	}
}

// OnArrowCollisionWithEnemy handles collision between arrow and enemy
func (ch *CollisionHandler) OnArrowCollisionWithEnemy(enemyID EntityID) {
	if !ch.Entities.Arrow.ComponentIgnore.Value {
		dead := ch.HealthSystem.Hit(enemyID, 1)
		ch.Entities.Arrow.ComponentIgnore.Value = true
		if dead {
			enemySpatial, _ := ch.SpatialSystem.GetEnemySpatial(enemyID)
			ch.Entities.Explosion.ComponentTemporary.Expiration = len(ch.Entities.Explosion.ComponentAnimation.Map["default"].Frames)
			ch.Entities.Explosion.ComponentSpatial = &ComponentSpatial{
				Width:  TileSize,
				Height: TileSize,
				Rect:   enemySpatial.Rect,
			}
			ch.Entities.Explosion.ComponentTemporary.OnExpiration = func() {
				dropCoin(ch.Entities.Explosion.ComponentSpatial.Rect.Min, ch.SystemsManager, ch.FrameRate)
			}
			ch.SystemsManager.AddEntity(ch.Entities.Explosion)
			ch.SystemsManager.RemoveEnemy(enemyID)
		} else {
			ch.SpatialSystem.MoveEnemyBack(enemyID, ch.Entities.Player.ComponentMovement.Direction)
		}
	}
}

// OnArrowCollisionWithObstacle handles collision between arrow and obstacle
func (ch *CollisionHandler) OnArrowCollisionWithObstacle() {
	ch.Entities.Arrow.ComponentMovement.RemainingMoves = 0
}

// OnPlayerCollisionWithObstacle handles collision between player and obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithObstacle(obstacleID EntityID) {
	// "Block" by undoing rect
	ch.Entities.Player.ComponentSpatial.Rect = ch.Entities.Player.ComponentSpatial.PrevRect
	ch.Entities.Sword.ComponentSpatial.Rect = ch.Entities.Sword.ComponentSpatial.PrevRect
}

// OnPlayerCollisionWithMoveableObstacle handles collision between player and moveable obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithMoveableObstacle(obstacleID EntityID) {
	moved := ch.SpatialSystem.MoveMoveableObstacle(obstacleID, ch.Entities.Player.ComponentMovement.Direction)
	if !moved {
		ch.Entities.Player.ComponentSpatial.Rect = ch.Entities.Player.ComponentSpatial.PrevRect
	}
}

// OnMoveableObstacleCollisionWithSwitch handles collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnMoveableObstacleNoCollisionWithSwitch handles *no* collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleNoCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnEnemyCollisionWithObstacle handles collision between enemy and obstacle
func (ch *CollisionHandler) OnEnemyCollisionWithObstacle(enemyID, obstacleID EntityID) {
	// Block enemy within the spatial system by reseting current rect to previous rect
	ch.SpatialSystem.UndoEnemyRect(enemyID)
}

// OnPlayerCollisionWithSwitch handles collision between player and switch
func (ch *CollisionHandler) OnPlayerCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerNoCollisionWithSwitch handles *no* collision between player and switch
func (ch *CollisionHandler) OnPlayerNoCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerCollisionWithWarp handles collision between player and warp
func (ch *CollisionHandler) OnPlayerCollisionWithWarp(warpID EntityID) {
	entityConfig, ok := ch.RoomWarps[warpID]
	if ok && !ch.RoomTransitionManager.Active() {
		ch.RoomTransitionManager.SetWarp()
		ch.GameStateManager.CurrentState = StateMapTransition
		ch.SystemsManager.SetShouldAddEntities(true)
		ch.RoomData.NextRoomID = entityConfig.WarpToRoomID
	}
}
