package zelduh

type EntityFactory struct {
	systemsManager *SystemsManager
}

func NewEntityFactory(
	systemsManager *SystemsManager,
) EntityFactory {
	return EntityFactory{
		systemsManager: systemsManager,
	}
}

func (ef *EntityFactory) NewEntity(presetName string, xTiles, yTiles float64) Entity {
	presetFn := GetPreset(presetName)
	entityConfig := presetFn(xTiles, yTiles)
	entityID := ef.systemsManager.NewEntityID()
	return BuildEntityFromConfig(
		entityConfig,
		entityID,
	)
}
