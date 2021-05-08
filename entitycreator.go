package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/miketmoore/zelduh/core/entity"
)

type EntityCreator struct {
	systemsManager              *SystemsManager
	temporarySystem             *TemporarySystem
	movementSystem              *MovementSystem
	entityFactory               *EntityFactory
	entityConfigPresetFnManager *EntityConfigPresetFnManager
	tileSize                    float64
	frameRate                   int
}

func NewEntityCreator(
	systemsManager *SystemsManager,
	temporarySystem *TemporarySystem,
	movementSystem *MovementSystem,
	entityFactory *EntityFactory,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	tileSize float64,
	frameRate int,
) EntityCreator {
	return EntityCreator{
		systemsManager:              systemsManager,
		temporarySystem:             temporarySystem,
		movementSystem:              movementSystem,
		entityFactory:               entityFactory,
		entityConfigPresetFnManager: entityConfigPresetFnManager,
		tileSize:                    tileSize,
		frameRate:                   frameRate,
	}
}

func (ec *EntityCreator) CreateCoin(
	v pixel.Vec,
) {
	coordinates := Coordinates{
		X: v.X / ec.tileSize,
		Y: v.Y / ec.tileSize,
	}
	coin := ec.entityFactory.NewEntityFromPresetName("coin", coordinates, ec.frameRate)
	ec.systemsManager.AddEntity(coin)
}

func (ec *EntityCreator) CreateUICoin() {
	presetFn := ec.entityConfigPresetFnManager.GetPreset("uiCoin")
	entityConfig := presetFn(Coordinates{X: 4, Y: 14})
	coin := ec.entityFactory.NewEntityFromConfig(entityConfig, ec.frameRate)
	ec.systemsManager.AddEntity(coin)
}

func (ec *EntityCreator) CreateExplosion(
	entityID entity.EntityID,
) {
	explosion := ec.entityFactory.NewEntityFromPresetName("explosion", NewCoordinates(0, 0), ec.frameRate)

	ec.temporarySystem.SetExpiration(
		explosion.ID(),
		len(explosion.componentAnimation.ComponentAnimationByName["default"].Frames),
		func() {
			ec.CreateCoin(explosion.componentRectangle.Rect.Min)
		},
	)

	explosion.componentDimensions = NewComponentDimensions(ec.tileSize, ec.tileSize)
	enemyComponentRectangle, _ := ec.movementSystem.ComponentRectangle(entityID)
	explosion.componentRectangle = &componentRectangle{
		Rect: enemyComponentRectangle.Rect,
	}

	ec.systemsManager.AddEntity(explosion)
}
