package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type boundsCollisionEntity struct {
	*componentRectangle
}

type onPlayerCollisionWithBounds func(side Bound)

type BoundsCollisionSystem struct {
	player                      boundsCollisionEntity
	window                      *pixelgl.Window
	mapBounds                   pixel.Rect
	onPlayerCollisionWithBounds func(side Bound)
}

func NewBoundsCollisionSystem(
	window *pixelgl.Window,
	mapBounds pixel.Rect,
	onPlayerCollisionWithBounds onPlayerCollisionWithBounds,
) BoundsCollisionSystem {
	return BoundsCollisionSystem{
		window:                      window,
		mapBounds:                   mapBounds,
		onPlayerCollisionWithBounds: onPlayerCollisionWithBounds,
	}
}

// AddEntity adds an entity to the system
func (s *BoundsCollisionSystem) AddEntity(entity Entity) {

	systemEntity := boundsCollisionEntity{
		componentRectangle: entity.componentRectangle,
	}

	switch entity.Category {
	case CategoryPlayer:
		s.player = systemEntity
	}
}

// Update checks for collisions
func (s *BoundsCollisionSystem) Update() error {

	s.handlePlayerAtMapEdge()

	return nil
}

func (s *BoundsCollisionSystem) handlePlayerAtMapEdge() {
	// DrawActiveSpace(s.Win, ActiveSpaceRectangle{
	// 	X:      s.MapBounds.Min.X,
	// 	Y:      s.MapBounds.Min.Y,
	// 	Width:  s.MapBounds.W(),
	// 	Height: s.MapBounds.H(),
	// })
	// DrawRect(s.window, s.mapBounds)
	// DrawRect(s.window, s.player.componentRectangle.Rect)

	player := s.player
	mapBounds := s.mapBounds

	if player.componentRectangle.Rect.Min.Y <= mapBounds.Min.Y {
		s.onPlayerCollisionWithBounds(BoundBottom)
	} else if player.componentRectangle.Rect.Min.X <= mapBounds.Min.X {
		s.onPlayerCollisionWithBounds(BoundLeft)
	} else if player.componentRectangle.Rect.Max.X >= mapBounds.Max.X {
		s.onPlayerCollisionWithBounds(BoundRight)
	} else if player.componentRectangle.Rect.Max.Y >= mapBounds.Max.Y {
		s.onPlayerCollisionWithBounds(BoundTop)
	}
}
