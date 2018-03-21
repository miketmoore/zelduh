package systems

import (
	"fmt"

	"engo.io/ecs"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type playerInputEntity struct {
	ecs.BasicEntity
	*components.MovementComponent
}

// PlayerInputSystem translates player input to vehicle input
type PlayerInputSystem struct {
	entities []playerInputEntity
	Win      *pixelgl.Window
}

// New is called by World when the system is added (I think)
func (*PlayerInputSystem) New(*ecs.World) {
	fmt.Println("PlayerInputSystem was added to the Scene")
}

// Add defines which components are required for an entity in this system and adds it
func (s *PlayerInputSystem) Add(
	basic *ecs.BasicEntity,
	movement *components.MovementComponent,
) {
	s.entities = append(s.entities, playerInputEntity{
		BasicEntity:       *basic,
		MovementComponent: movement,
	})
}

// Remove removes an entity from the system completely
func (s *PlayerInputSystem) Remove(basic ecs.BasicEntity) {
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
func (s *PlayerInputSystem) Update(dt float32) {
	for _, entity := range s.entities {
		win := s.Win

		entity.MovementComponent.Moving = true
		if win.Pressed(pixelgl.KeyUp) {
			entity.MovementComponent.Direction = direction.Up
		} else if win.Pressed(pixelgl.KeyRight) {
			entity.MovementComponent.Direction = direction.Right
		} else if win.Pressed(pixelgl.KeyDown) {
			entity.MovementComponent.Direction = direction.Down
		} else if win.Pressed(pixelgl.KeyLeft) {
			entity.MovementComponent.Direction = direction.Left
		} else {
			entity.MovementComponent.Moving = false
		}
	}
}
