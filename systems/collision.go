package systems

import (
	"fmt"

	"engo.io/ecs"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/message"
)

// CollisionMessage is sent whenever a collision is detected by the CollisionSystem.
type CollisionMessage struct {
	Entity     collisionEntity
	EntityType string
	To         collisionEntity
	ToType     string
	// Groups CollisionGroup
}

// Type implements the zelduh.Message interface
func (CollisionMessage) Type() string { return "CollisionMessage" }

type collisionEntity struct {
	ecs.BasicEntity
	*components.SpatialComponent
	*components.EntityTypeComponent
}

// CollisionSystem determines effect of vehicle input on vehicle physics
type CollisionSystem struct {
	entities []collisionEntity
	Mailbox  *message.Manager
}

// New is called by World when the system is added (I think)
func (s *CollisionSystem) New(*ecs.World) {
	fmt.Println("CollisionSystem was added to the Scene")
}

// Add defines which components are required for an entity in this system and adds it
func (s *CollisionSystem) Add(
	basic *ecs.BasicEntity,
	spatial *components.SpatialComponent,
	entityType *components.EntityTypeComponent,
) {
	s.entities = append(s.entities, collisionEntity{
		BasicEntity:         *basic,
		SpatialComponent:    spatial,
		EntityTypeComponent: entityType,
	})
}

// Remove removes an entity from the system completely
func (s *CollisionSystem) Remove(basic ecs.BasicEntity) {
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
func (s *CollisionSystem) Update(dt float32) {
	for _, entity := range s.entities {
		for _, entityB := range s.entities {
			if entity.ID() != entityB.ID() {
				if entity.SpatialComponent.Rect.Contains(entityB.SpatialComponent.Rect.Min) ||
					entity.SpatialComponent.Rect.Contains(entityB.SpatialComponent.Rect.Max) {
					fmt.Println("Collision!")
					// TODO Determine what to do
					// How do I know which is the player and which is the coin?
					// If one is a coin and the other is a player...
					s.Mailbox.Dispatch(CollisionMessage{
						Entity:     entity,
						EntityType: entity.EntityTypeComponent.Type,
						To:         entityB,
						ToType:     entity.EntityTypeComponent.Type,
					})
				}
			}
		}
	}
}
