package zelduh

// CollisionHandler contains collision handlers
type CollisionHandler struct {
	SystemsManager                  *SystemsManager
	MovementSystem                  *MovementSystem
	HealthSystem                    *HealthSystem
	TemporarySystem                 *TemporarySystem
	EntityCreator                   *EntityCreator
	ShouldAddEntities               *bool
	NextRoomID                      *RoomID
	CurrentState                    *State
	RoomTransition                  *RoomTransition
	EntitiesMap                     EntitiesMap
	Player, Sword, Explosion, Arrow *Entity
	Hearts                          []Entity
	RoomWarps                       RoomWarps
	EntityConfigPresetFnManager     *EntityConfigPresetFnManager
	TileSize                        float64
	FrameRate                       int
}

func NewCollisionHandler(
	systemsManager *SystemsManager,
	movementSystem *MovementSystem,
	healthSystem *HealthSystem,
	temporarySystem *TemporarySystem,
	entityCreator *EntityCreator,
	shouldAddEntities *bool,
	nextRoomID *RoomID,
	currentState *State,
	roomTransition *RoomTransition,
	entitiesMap EntitiesMap,
	player, sword, explosion, arrow *Entity,
	hearts []Entity,
	roomWarps RoomWarps,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
	frameRate int,
) CollisionHandler {
	return CollisionHandler{
		SystemsManager:              systemsManager,
		MovementSystem:              movementSystem,
		HealthSystem:                healthSystem,
		TemporarySystem:             temporarySystem,
		EntityCreator:               entityCreator,
		ShouldAddEntities:           shouldAddEntities,
		NextRoomID:                  nextRoomID,
		CurrentState:                currentState,
		RoomTransition:              roomTransition,
		EntitiesMap:                 entitiesMap,
		Player:                      player,
		Sword:                       sword,
		Explosion:                   explosion,
		Arrow:                       arrow,
		Hearts:                      hearts,
		RoomWarps:                   roomWarps,
		EntityConfigPresetFnManager: entityConfigPresetFnManager,
		TileSize:                    tileSize,
		FrameRate:                   frameRate,
	}
}

// OnPlayerCollisionWithBounds handles collisions between player and bounds
func (ch *CollisionHandler) OnPlayerCollisionWithBounds(side Bound) {
	// TODO prevent room transition if no room exists on this side

	if !ch.RoomTransition.Active && (*ch.NextRoomID) > 0 {

		ch.RoomTransition.Active = true
		ch.RoomTransition.Side = side
		ch.RoomTransition.Style = TransitionSlide
		ch.RoomTransition.Timer = int(ch.RoomTransition.Start)
		*ch.CurrentState = StateMapTransition
		*ch.ShouldAddEntities = true
	} else {
		ch.MovementSystem.SetZeroSpeed(ch.Player.ID())
	}
}

// OnPlayerCollisionWithCoin handles collision between player and coin
func (ch *CollisionHandler) OnPlayerCollisionWithCoin(coinID EntityID) {
	ch.Player.componentCoins.Coins++
	ch.SystemsManager.Remove(CategoryCoin, coinID)
}

// OnPlayerCollisionWithEnemy handles collision between player and enemy
func (ch *CollisionHandler) OnPlayerCollisionWithEnemy(enemyID EntityID) {
	// TODO repeat what I did with the enemies
	ch.MovementSystem.MovePlayerBack()
	ch.Player.componentHealth.Total--

	// remove heart entity
	heartIndex := len(ch.Hearts) - 1
	ch.SystemsManager.Remove(CategoryHeart, ch.Hearts[heartIndex].ID())
	ch.Hearts = append(ch.Hearts[:heartIndex], ch.Hearts[heartIndex+1:]...)

	if ch.Player.componentHealth.Total == 0 {
		*ch.CurrentState = StateOver
	}
}

// OnSwordCollisionWithEnemy handles collision between sword and enemy
func (ch *CollisionHandler) OnSwordCollisionWithEnemy(enemyID EntityID) {
	if !ch.Sword.componentIgnore.Value {
		dead := false
		if !ch.MovementSystem.EnemyMovingFromHit(enemyID) {
			dead = ch.HealthSystem.Hit(enemyID, 1)
			if dead {
				ch.EntityCreator.CreateExplosion(enemyID)
				ch.SystemsManager.RemoveEnemy(enemyID)
			} else {
				ch.MovementSystem.MoveEnemyBack(enemyID, ch.Player.componentMovement.Direction)
			}
		}

	}
}

// OnArrowCollisionWithEnemy handles collision between arrow and enemy
func (ch *CollisionHandler) OnArrowCollisionWithEnemy(enemyID EntityID) {
	if !ch.Arrow.componentIgnore.Value {
		dead := ch.HealthSystem.Hit(enemyID, 1)
		ch.Arrow.componentIgnore.Value = true
		if dead {
			ch.EntityCreator.CreateExplosion(enemyID)
			ch.SystemsManager.RemoveEnemy(enemyID)
		} else {
			ch.MovementSystem.MoveEnemyBack(enemyID, ch.Player.componentMovement.Direction)
		}
	}
}

// OnArrowCollisionWithObstacle handles collision between arrow and obstacle
func (ch *CollisionHandler) OnArrowCollisionWithObstacle() {
	ch.Arrow.componentMovement.RemainingMoves = 0
}

// OnPlayerCollisionWithObstacle handles collision between player and obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithObstacle(obstacleID EntityID) {
	// "Block" by undoing rect
	ch.Player.componentRectangle.Rect = ch.Player.componentRectangle.PrevRect
	ch.Sword.componentRectangle.Rect = ch.Sword.componentRectangle.PrevRect
}

// OnPlayerCollisionWithMoveableObstacle handles collision between player and moveable obstacle
func (ch *CollisionHandler) OnPlayerCollisionWithMoveableObstacle(obstacleID EntityID) {
	moved := ch.MovementSystem.MoveMoveableObstacle(obstacleID, ch.Player.componentMovement.Direction)
	if !moved {
		ch.Player.componentRectangle.Rect = ch.Player.componentRectangle.PrevRect
	}
}

// OnMoveableObstacleCollisionWithSwitch handles collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && !entity.componentToggler.Enabled() {
			entity.componentToggler.Toggle()
		}
	}
}

// OnMoveableObstacleNoCollisionWithSwitch handles *no* collision between moveable obstacle and switch
func (ch *CollisionHandler) OnMoveableObstacleNoCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && entity.componentToggler.Enabled() {
			entity.componentToggler.Toggle()
		}
	}
}

// OnEnemyCollisionWithObstacle handles collision between enemy and obstacle
func (ch *CollisionHandler) OnEnemyCollisionWithObstacle(enemyID, obstacleID EntityID) {
	// Block enemy within the spatial system by reseting current rect to previous rect
	ch.MovementSystem.UndoEnemyRect(enemyID)
}

// OnPlayerCollisionWithSwitch handles collision between player and switch
func (ch *CollisionHandler) OnPlayerCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && !entity.componentToggler.Enabled() {
			entity.componentToggler.Toggle()
		}
	}
}

// OnPlayerNoCollisionWithSwitch handles *no* collision between player and switch
func (ch *CollisionHandler) OnPlayerNoCollisionWithSwitch(collisionSwitchID EntityID) {
	for id, entity := range ch.EntitiesMap {
		if id == collisionSwitchID && entity.componentToggler.Enabled() {
			entity.componentToggler.Toggle()
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
