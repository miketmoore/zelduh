package systems

import (
	"fmt"

	"engo.io/ecs"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
)

type renderEntity struct {
	ecs.BasicEntity
	*components.SpatialComponent
	*components.AppearanceComponent
}

// RenderSystem translates player input to vehicle input
type RenderSystem struct {
	entities []renderEntity
	Win      *pixelgl.Window
}

// New is called by World when the system is added (I think)
func (*RenderSystem) New(*ecs.World) {
	fmt.Println("RenderSystem was added to the Scene")
}

// Add defines which components are required for an entity in this system and adds it
func (s *RenderSystem) Add(
	basic *ecs.BasicEntity,
	spatial *components.SpatialComponent,
	appearance *components.AppearanceComponent,
) {
	s.entities = append(s.entities, renderEntity{
		BasicEntity:         *basic,
		SpatialComponent:    spatial,
		AppearanceComponent: appearance,
	})
}

// Remove removes an entity from the system completely
func (s *RenderSystem) Remove(basic ecs.BasicEntity) {
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
// This is where we use components and alter component data
func (s *RenderSystem) Update(dt float32) {
	for _, entity := range s.entities {

		entity.Shape.Clear()
		entity.Shape.Color = entity.AppearanceComponent.Color
		entity.Shape.Push(entity.SpatialComponent.Rect.Min)
		entity.Shape.Push(entity.SpatialComponent.Rect.Max)
		entity.Shape.Rectangle(0)
		entity.Shape.Draw(s.Win)
	}
}
