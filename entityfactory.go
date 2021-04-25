package zelduh

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
	return BuildEntityFromConfig(
		entityConfig,
		entityID,
		frameRate,
	)
}

func (ef *EntityFactory) NewEntity2(entityConfig EntityConfig, frameRate int) Entity {
	return BuildEntityFromConfig(
		entityConfig,
		ef.systemsManager.NewEntityID(),
		frameRate,
	)
}
