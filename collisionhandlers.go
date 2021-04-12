package zelduh

import (
	"github.com/faiface/pixel"
)

// CollisionHandler contains collision handlers
type CollisionHandler struct {
	GameModel         *GameModel
	SystemsManager    *SystemsManager
	SpatialSystem     *SpatialSystem
	HealthSystem      *HealthSystem
	ShouldAddEntities *bool
}

// OnPlayerCollisionWithBounds handles collisions between player and bounds
func (ch *CollisionHandler) OnPlayerCollisionWithBounds(side Bound) {
	if !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Side = side
		ch.GameModel.RoomTransition.Style = TransitionSlide
		ch.GameModel.RoomTransition.Timer = int(ch.GameModel.RoomTransition.Start)
		ch.GameModel.CurrentState = StateMapTransition
		*ch.ShouldAddEntities = true
	}
}

// OnPlayerCollisionWithCoin handles collision between player and coin
func (ch *CollisionHandler) OnPlayerCollisionWithCoin(coinID EntityID) {
	ch.GameModel.Entities.Player.ComponentCoins.Coins++
	ch.SystemsManager.Remove(CategoryCoin, coinID)
}

// OnPlayerCollisionWithEnemy handles collision between player and enemy
func (ch *CollisionHandler) OnPlayerCollisionWithEnemy(enemyID EntityID) {
	// TODO repeat what I did with the enemies
	ch.SpatialSystem.MovePlayerBack()
	ch.GameModel.Entities.Player.ComponentHealth.Total--

	// remove heart entity
	heartIndex := len(ch.GameModel.Entities.Hearts) - 1
	ch.SystemsManager.Remove(CategoryHeart, ch.GameModel.Entities.Hearts[heartIndex].ID())
	ch.GameModel.Entities.Hearts = append(ch.GameModel.Entities.Hearts[:heartIndex], ch.GameModel.Entities.Hearts[heartIndex+1:]...)

	if ch.GameModel.Entities.Player.ComponentHealth.Total == 0 {
		ch.GameModel.CurrentState = StateOver
	}
}

func dropCoin(v pixel.Vec, systemsManager *SystemsManager) {
	coin := BuildEntityFromConfig(GetPreset("coin")(v.X/TileSize, v.Y/TileSize), systemsManager.NewEntityID())
	systemsManager.AddEntity(coin)
}

// OnSwordCollisionWithEnemy handles collision between sword and enemy
func (ch *CollisionHandler) OnSwordCollisionWithEnemy(enemyID EntityID) {
	if !ch.GameModel.Entities.Sword.ComponentIgnore.Value {
		dead := false
		if !ch.SpatialSystem.EnemyMovingFromHit(enemyID) {
			dead = ch.HealthSystem.Hit(enemyID, 1)
			if dead {
				enemySpatial, _ := ch.SpatialSystem.GetEnemySpatial(enemyID)
				ch.GameModel.Entities.Explosion.ComponentTemporary.Expiration = len(ch.GameModel.Entities.Explosion.ComponentAnimation.Map["default"].Frames)
				ch.GameModel.Entities.Explosion.ComponentSpatial = &ComponentSpatial{
					Width:  TileSize,
					Height: TileSize,
					Rect:   enemySpatial.Rect,
				}
				ch.GameModel.Entities.Explosion.ComponentTemporary.OnExpiration = func() {
					dropCoin(ch.GameModel.Entities.Explosion.ComponentSpatial.Rect.Min, ch.SystemsManager)
				}
				ch.SystemsManager.AddEntity(ch.GameModel.Entities.Explosion)
				ch.SystemsManager.RemoveEnemy(enemyID)
			} else {
				ch.SpatialSystem.MoveEnemyBack(enemyID, ch.GameModel.Entities.Player.ComponentMovement.Direction)
			}
		}

	}
}

// OnArrowCollisionWithEnemy handles collision between arrow and enemy
func (ch *CollisionHandler) OnArrowCollisionWithEnemy(enemyID EntityID) {
	if !ch.GameModel.Entities.Arrow.ComponentIgnore.Value {
		dead := ch.HealthSystem.Hit(enemyID, 1)
		ch.GameModel.Entities.Arrow.ComponentIgnore.Value = true
		if dead {
			enemySpatial, _ := ch.SpatialSystem.GetEnemySpatial(enemyID)
			ch.GameModel.Entities.Explosion.ComponentTemporary.Expiration = len(ch.GameModel.Entities.Explosion.ComponentAnimation.Map["default"].Frames)
			ch.GameModel.Entities.Explosion.ComponentSpatial = &ComponentSpatial{
				Width:  TileSize,
				Height: TileSize,
				Rect:   enemySpatial.Rect,
			}
			ch.GameModel.Entities.Explosion.ComponentTemporary.OnExpiration = func() {
				dropCoin(ch.GameModel.Entities.Explosion.ComponentSpatial.Rect.Min, ch.SystemsManager)
			}
			ch.SystemsManager.AddEntity(ch.GameModel.Entities.Explosion)
			ch.SystemsManager.RemoveEnemy(enemyID)
		} else {
			ch.SpatialSystem.MoveEnemyBack(enemyID, ch.GameModel.Entities.Player.ComponentMovement.Direction)
		}
	}
}

// OnArrowCollisionWithObstacle handles collision between arrow and obstacle
func (ch *CollisionHandler) OnArrowCollisionWithObstacle() {
	ch.GameModel.Entities.Arrow.ComponentMovement.RemainingMoves = 0
}

// OnPlayerCollisionWithObstacle handles collision between player and obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithObstacle(obstacleID EntityID) {
	// "Block" by undoing rect
	ch.GameModel.Entities.Player.ComponentSpatial.Rect = ch.GameModel.Entities.Player.ComponentSpatial.PrevRect
	ch.GameModel.Entities.Sword.ComponentSpatial.Rect = ch.GameModel.Entities.Sword.ComponentSpatial.PrevRect
}

// OnPlayerCollisionWithMoveableObstacle handles collision between player and moveable obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithMoveableObstacle(obstacleID EntityID) {
	moved := ch.SpatialSystem.MoveMoveableObstacle(obstacleID, ch.GameModel.Entities.Player.ComponentMovement.Direction)
	if !moved {
		ch.GameModel.Entities.Player.ComponentSpatial.Rect = ch.GameModel.Entities.Player.ComponentSpatial.PrevRect
	}
}

// OnMoveableObstacleCollisionWithSwitch handles collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnMoveableObstacleNoCollisionWithSwitch handles *no* collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleNoCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
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
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && !entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerNoCollisionWithSwitch handles *no* collision between player and switch
func (ch *CollisionHandler) OnPlayerNoCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.GameModel.EntitiesMap {
		if id == collisionSwitchID && entity.ComponentToggler.Enabled() {
			entity.ComponentToggler.Toggle()
		}
	}
}

// OnPlayerCollisionWithWarp handles collision between player and warp
func (ch *CollisionHandler) OnPlayerCollisionWithWarp(warpID EntityID) {
	entityConfig, ok := ch.GameModel.RoomWarps[warpID]
	if ok && !ch.GameModel.RoomTransition.Active {
		ch.GameModel.RoomTransition.Active = true
		ch.GameModel.RoomTransition.Style = TransitionWarp
		ch.GameModel.RoomTransition.Timer = 1
		ch.GameModel.CurrentState = StateMapTransition
		*ch.ShouldAddEntities = true
		ch.GameModel.NextRoomID = entityConfig.WarpToRoomID
	}
}
