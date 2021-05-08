package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/core/entity"
	"golang.org/x/image/colornames"
)

type entityConfigPresetFnByNameMap map[PresetName]EntityConfigPresetFn

type EntityFactory struct {
	systemsManager                *SystemsManager
	entityConfigPresetFnByNameMap entityConfigPresetFnByNameMap
	temporarySystem               *TemporarySystem
	movementSystem                *MovementSystem
	tileSize                      float64
	frameRate                     int
}

func NewEntityFactory(
	systemsManager *SystemsManager,
	entityConfigPresetFnByNameMap entityConfigPresetFnByNameMap,
	temporarySystem *TemporarySystem,
	movementSystem *MovementSystem,
	tileSize float64,
	frameRate int,
) EntityFactory {
	return EntityFactory{
		systemsManager:                systemsManager,
		temporarySystem:               temporarySystem,
		movementSystem:                movementSystem,
		entityConfigPresetFnByNameMap: entityConfigPresetFnByNameMap,
		tileSize:                      tileSize,
		frameRate:                     frameRate,
	}
}

type PresetName string

func (ef *EntityFactory) GetPreset(presetName PresetName) EntityConfigPresetFn {
	return ef.entityConfigPresetFnByNameMap[presetName]
}

func (ef *EntityFactory) NewEntityFromPresetName(presetName PresetName, coordinates Coordinates, frameRate int) Entity {
	presetFn := ef.entityConfigPresetFnByNameMap[presetName]
	entityConfig := presetFn(coordinates)
	entityID := ef.systemsManager.NewEntityID()
	return ef.buildEntityFromConfig(
		entityConfig,
		entityID,
		frameRate,
	)
}

func (ef *EntityFactory) NewEntityFromConfig(entityConfig EntityConfig, frameRate int) Entity {
	return ef.buildEntityFromConfig(
		entityConfig,
		ef.systemsManager.NewEntityID(),
		frameRate,
	)
}

// BuildEntitiesFromConfigs builds and returns a batch of entities
func (ef *EntityFactory) buildEntitiesFromConfigs(newEntityID func() entity.EntityID, frameRate int, configs ...EntityConfig) []Entity {
	batch := []Entity{}
	for _, config := range configs {
		entity := ef.buildEntityFromConfig(config, newEntityID(), frameRate)
		batch = append(batch, entity)
	}
	return batch
}

func (ef *EntityFactory) buildEntityFromConfig(c EntityConfig, id entity.EntityID, frameRate int) Entity {
	entity := Entity{
		id:                   id,
		Category:             c.Category,
		componentIgnore:      NewComponentIgnore(c.Ignore),
		componentCoordinates: NewComponentCoordinates(c.Coordinates.X, c.Coordinates.Y),
		componentDimensions:  NewComponentDimensions(c.Dimensions.Width, c.Dimensions.Height),
		componentRectangle: NewComponentRectangle(
			c.Coordinates.X,
			c.Coordinates.Y,
			c.Dimensions.Width,
			c.Dimensions.Height,
		),

		// Create default shape and color
		// useful for debugging
		// might want to remove this later... not sure if creating
		// shapes that aren't being used increases heap memory
		componentShape: NewComponentShape(),
		componentColor: NewComponentColor(colornames.Greenyellow),
	}

	if c.Expiration > 0 {
		entity.componentTemporary = NewComponentTemporary(c.Expiration)
	}

	if c.Category == CategoryWarp {
		entity.componentEnabled = NewComponentEnabled(true)
	}

	if c.Health > 0 {
		entity.componentHealth = NewComponentHealth(c.Health)
	}

	if c.Hitbox != nil {
		entity.componentHitbox = NewComponentHitbox(c.Hitbox.Radius, float64(c.Hitbox.CollisionWithRectMod))
	}

	if c.Transform != nil {
		entity.componentRotation = NewComponentRotation(c.Transform.Rotation)
	}

	if c.Toggleable {
		entity.componentToggler = NewComponentToggler(c.Toggled)
	}

	entity.componentInvincible = NewComponentInvincible(c.Invincible)

	if c.Movement != nil {
		entity.componentMovement = NewComponentMovement(
			c.Movement.Direction,
			c.Movement.Speed,
			c.Movement.MaxSpeed,
			c.Movement.MaxMoves,
			c.Movement.RemainingMoves,
			c.Movement.HitSpeed,
			c.Movement.MovingFromHit,
			c.Movement.HitBackMoves,
			c.MovementPatternName,
		)
	}

	if c.Coins {
		entity.componentCoins = NewComponentCoins(0)
	}

	if c.Dash != nil {
		entity.componentDash = NewComponentDash(
			c.Dash.Charge,
			c.Dash.MaxCharge,
			c.Dash.SpeedMod,
		)
	}

	// An animation is a sprite graphic that may have one or more frames
	// so technically it might not be an animation
	if c.Animation != nil {
		entity.componentAnimation = NewComponentAnimation(c.Animation, frameRate)
	}

	return entity
}

func (ef *EntityFactory) CreateCoin(
	v pixel.Vec,
) {
	coordinates := Coordinates{
		X: v.X / ef.tileSize,
		Y: v.Y / ef.tileSize,
	}
	coin := ef.NewEntityFromPresetName("coin", coordinates, ef.frameRate)
	ef.systemsManager.AddEntity(coin)
}

func (ef *EntityFactory) CreateUICoin() {
	presetFn := ef.entityConfigPresetFnByNameMap["uiCoin"]
	entityConfig := presetFn(Coordinates{X: 4, Y: 14})
	coin := ef.NewEntityFromConfig(entityConfig, ef.frameRate)
	ef.systemsManager.AddEntity(coin)
}

func (ef *EntityFactory) CreateExplosion(
	entityID entity.EntityID,
) {
	explosion := ef.NewEntityFromPresetName("explosion", NewCoordinates(0, 0), ef.frameRate)

	ef.temporarySystem.SetExpiration(
		explosion.ID(),
		len(explosion.componentAnimation.ComponentAnimationByName["default"].Frames),
		func() {
			ef.CreateCoin(explosion.componentRectangle.Rect.Min)
		},
	)

	explosion.componentDimensions = NewComponentDimensions(ef.tileSize, ef.tileSize)
	enemyComponentRectangle, _ := ef.movementSystem.ComponentRectangle(entityID)
	explosion.componentRectangle = &componentRectangle{
		Rect: enemyComponentRectangle.Rect,
	}

	ef.systemsManager.AddEntity(explosion)
}

func (ef *EntityFactory) buildWarpFnFactory(
	dimensions Dimensions,
) BuildWarpFn {

	return func(
		warpToRoomID RoomID,
		coordinates Coordinates,
		hitboxRadius float64,
	) EntityConfig {
		return EntityConfig{
			Category:     CategoryWarp,
			WarpToRoomID: warpToRoomID,
			Dimensions:   dimensions,
			Coordinates:  coordinates,
			Hitbox: &HitboxConfig{
				Radius: hitboxRadius,
			},
		}
	}
}

type BuildWarpStoneFn func(
	WarpToRoomID RoomID,
	coordinates Coordinates,
	HitBoxRadius float64,
) EntityConfig

func (ef *EntityFactory) buildWarpStoneFnFactory() BuildWarpStoneFn {
	return func(
		WarpToRoomID RoomID,
		coordinates Coordinates,
		HitBoxRadius float64,
	) EntityConfig {
		presetFn := ef.entityConfigPresetFnByNameMap["warpStone"]
		e := presetFn(Coordinates{X: coordinates.X, Y: coordinates.Y})
		e.WarpToRoomID = 6
		e.Hitbox.Radius = 5
		return e
	}
}
