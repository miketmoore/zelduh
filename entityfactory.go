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

func (ef *EntityFactory) NewEntity(presetName string, xTiles, yTiles float64, frameRate int) Entity {
	presetFn := ef.entityConfigPresetFnManager.GetPreset(presetName)
	entityConfig := presetFn(xTiles, yTiles)
	entityID := ef.systemsManager.NewEntityID()
	return BuildEntityFromConfig(
		entityConfig,
		entityID,
		frameRate,
	)
}
