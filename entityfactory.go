package zelduh

import "golang.org/x/image/colornames"

type EntityFactory struct {
	systemsManager              *SystemsManager
	entityConfigPresetFnManager *EntityConfigPresetFnManager
}

func NewEntityFactory(
	systemsManager *SystemsManager,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
) EntityFactory {
	return EntityFactory{
		systemsManager:              systemsManager,
		entityConfigPresetFnManager: entityConfigPresetFnManager,
	}
}

type PresetName string

func (ef *EntityFactory) NewEntity(presetName PresetName, coordinates Coordinates, frameRate int) Entity {
	presetFn := ef.entityConfigPresetFnManager.GetPreset(presetName)
	entityConfig := presetFn(coordinates)
	entityID := ef.systemsManager.NewEntityID()
	return ef.buildEntityFromConfig(
		entityConfig,
		entityID,
		frameRate,
	)
}

func (ef *EntityFactory) NewEntity2(entityConfig EntityConfig, frameRate int) Entity {
	return ef.buildEntityFromConfig(
		entityConfig,
		ef.systemsManager.NewEntityID(),
		frameRate,
	)
}

// BuildEntitiesFromConfigs builds and returns a batch of entities
func (ef *EntityFactory) buildEntitiesFromConfigs(newEntityID func() EntityID, frameRate int, configs ...EntityConfig) []Entity {
	batch := []Entity{}
	for _, config := range configs {
		entity := ef.buildEntityFromConfig(config, newEntityID(), frameRate)
		batch = append(batch, entity)
	}
	return batch
}

func (ef *EntityFactory) buildEntityFromConfig(c EntityConfig, id EntityID, frameRate int) Entity {
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
