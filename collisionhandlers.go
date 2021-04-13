package zelduh

import (
	"github.com/faiface/pixel"
)

// CollisionHandler contains collision handlers
type CollisionHandler struct {
	SystemsManager                  *SystemsManager
	SpatialSystem                   *SpatialSystem
	HealthSystem                    *HealthSystem
	ShouldAddEntities               *bool
	NextRoomID                      *RoomID
	CurrentState                    *State
	RoomTransition                  *RoomTransition
	EntitiesMap                     EntitiesMap
	Player, Sword, Explosion, Arrow *Entity
	Hearts                          []Entity
	RoomWarps                       RoomWarps
}

func NewCollisionHandler(
	systemsManager *SystemsManager,
	spatialSystem *SpatialSystem,
	healthSystem *HealthSystem,
	shouldAddEntities *bool,
	nextRoomID *RoomID,
	currentState *State,
	roomTransition *RoomTransition,
	entitiesMap EntitiesMap,
	player, sword, explosion, arrow *Entity,
	hearts []Entity,
	roomWarps RoomWarps,
) CollisionHandler {
	return CollisionHandler{
		SystemsManager:    systemsManager,
		SpatialSystem:     spatialSystem,
		HealthSystem:      healthSystem,
		ShouldAddEntities: shouldAddEntities,
		NextRoomID:        nextRoomID,
		CurrentState:      currentState,
		RoomTransition:    roomTransition,
		EntitiesMap:       entitiesMap,
		Player:            player,
		Sword:             sword,
		Explosion:         explosion,
		Arrow:             arrow,
		Hearts:            hearts,
		RoomWarps:         roomWarps,
	}
}

// OnPlayerCollisionWithBounds handles collisions between player and bounds
func (ch *CollisionHandler) OnPlayerCollisionWithBounds(side Bound) {
	if !ch.RoomTransition.Active {
		ch.RoomTransition.Active = true
		ch.RoomTransition.Side = side
		ch.RoomTransition.Style = TransitionSlide
		ch.RoomTransition.Timer = int(ch.RoomTransition.Start)
		*ch.CurrentState = StateMapTransition
		*ch.ShouldAddEntities = true
	}
}

// OnPlayerCollisionWithCoin handles collision between player and coin
func (ch *CollisionHandler) OnPlayerCollisionWithCoin(coinID EntityID) {
	ch.Player.ComponentCoins.Coins++
	ch.SystemsManager.Remove(CategoryCoin, coinID)
}

// OnPlayerCollisionWithEnemy handles collision between player and enemy
func (ch *CollisionHandler) OnPlayerCollisionWithEnemy(enemyID EntityID) {
	// TODO repeat what I did with the enemies
	ch.SpatialSystem.MovePlayerBack()
	ch.Player.ComponentHealth.Total--

	// remove heart entity
	heartIndex := len(ch.Hearts) - 1
	ch.SystemsManager.Remove(CategoryHeart, ch.Hearts[heartIndex].ID())
	ch.Hearts = append(ch.Hearts[:heartIndex], ch.Hearts[heartIndex+1:]...)

	if ch.Player.ComponentHealth.Total == 0 {
		*ch.CurrentState = StateOver
	}
}

func dropCoin(v pixel.Vec, systemsManager *SystemsManager) {
	coin := BuildEntityFromConfig(GetPreset("coin")(v.X/TileSize, v.Y/TileSize), systemsManager.NewEntityID())
	systemsManager.AddEntity(coin)
}

// OnSwordCollisionWithEnemy handles collision between sword and enemy
func (ch *CollisionHandler) OnSwordCollisionWithEnemy(enemyID EntityID) {
	if !ch.Sword.ComponentIgnore.Value {
		dead := false
		if !ch.SpatialSystem.EnemyMovingFromHit(enemyID) {
			dead = ch.HealthSystem.Hit(enemyID, 1)
			if dead {
				enemySpatial, _ := ch.SpatialSystem.GetEnemySpatial(enemyID)
				ch.Explosion.ComponentTemporary.Expiration = len(ch.Explosion.ComponentAnimation.Map["default"].Frames)
				ch.Explosion.ComponentSpatial = &ComponentSpatial{
					Width:  TileSize,
					Height: TileSize,
					Rect:   enemySpatial.Rect,
				}
				ch.Explosion.ComponentTemporary.OnExpiration = func() {
					dropCoin(ch.Explosion.ComponentSpatial.Rect.Min, ch.SystemsManager)
				}
				ch.SystemsManager.AddEntity(*ch.Explosion)
				ch.SystemsManager.RemoveEnemy(enemyID)
			} else {
				ch.SpatialSystem.MoveEnemyBack(enemyID, ch.Player.ComponentMovement.Direction)
			}
		}

	}
}

// OnArrowCollisionWithEnemy handles collision between arrow and enemy
func (ch *CollisionHandler) OnArrowCollisionWithEnemy(enemyID EntityID) {
	if !ch.Arrow.ComponentIgnore.Value {
		dead := ch.HealthSystem.Hit(enemyID, 1)
		ch.Arrow.ComponentIgnore.Value = true
		if dead {
			enemySpatial, _ := ch.SpatialSystem.GetEnemySpatial(enemyID)
			ch.Explosion.ComponentTemporary.Expiration = len(ch.Explosion.ComponentAnimation.Map["default"].Frames)
			ch.Explosion.ComponentSpatial = &ComponentSpatial{
				Width:  TileSize,
				Height: TileSize,
				Rect:   enemySpatial.Rect,
			}
			ch.Explosion.ComponentTemporary.OnExpiration = func() {
				dropCoin(ch.Explosion.ComponentSpatial.Rect.Min, ch.SystemsManager)
			}
			ch.SystemsManager.AddEntity(*ch.Explosion)
			ch.SystemsManager.RemoveEnemy(enemyID)
		} else {
			ch.SpatialSystem.MoveEnemyBack(enemyID, ch.Player.ComponentMovement.Direction)
		}
	}
}

// OnArrowCollisionWithObstacle handles collision between arrow and obstacle
func (ch *CollisionHandler) OnArrowCollisionWithObstacle() {
	ch.Arrow.ComponentMovement.RemainingMoves = 0
}

// OnPlayerCollisionWithObstacle handles collision between player and obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithObstacle(obstacleID EntityID) {
	// "Block" by undoing rect
	ch.Player.ComponentSpatial.Rect = ch.Player.ComponentSpatial.PrevRect
	ch.Sword.ComponentSpatial.Rect = ch.Sword.ComponentSpatial.PrevRect
}

// OnPlayerCollisionWithMoveableObstacle handles collision between player and moveable obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithMoveableObstacle(obstacleID EntityID) {
	moved := ch.SpatialSystem.MoveMoveableObstacle(obstacleID, ch.Player.ComponentMovement.Direction)
	if !moved {
		ch.Player.ComponentSpatial.Rect = ch.Player.ComponentSpatial.PrevRect
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
	if ok && !ch.RoomTransition.Active {
		ch.RoomTransition.Active = true
		ch.RoomTransition.Style = TransitionWarp
		ch.RoomTransition.Timer = 1
		*ch.CurrentState = StateMapTransition
		*ch.ShouldAddEntities = true
		*ch.NextRoomID = entityConfig.WarpToRoomID
	}
}
