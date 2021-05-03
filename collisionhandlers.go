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
