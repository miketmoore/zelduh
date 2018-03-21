package systems

import (
	"fmt"

	"engo.io/ecs"
	"github.com/miketmoore/zelduh/components"
)

type spatialEntity struct {
	ecs.BasicEntity
	*components.SpatialComponent
}

// SpatialSystem updates spatial component data based on physics component data
type SpatialSystem struct {
	entities []spatialEntity
}

// New is the initialisation of the System
func (*SpatialSystem) New(*ecs.World) {
	fmt.Println("SpatialSystem was added to the Scene")
}

// Add adds an entity to the system and specifies required components
func (s *SpatialSystem) Add(
	basic *ecs.BasicEntity,
	space *components.SpatialComponent,
) {
	s.entities = append(s.entities, spatialEntity{
		BasicEntity:      *basic,
		SpatialComponent: space,
	})
}

// Remove removes an entity from the system completely
func (s *SpatialSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

// Update is called from World.Update on every frame
// dt is the time in seconds since the last frame
func (s *SpatialSystem) Update(dt float32) {
	for _, entity := range s.entities {
		fmt.Printf("SpatialSystem Update loop %v\n", entity)
	}
}
