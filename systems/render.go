package systems

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/terraform2d"
	"github.com/miketmoore/zelduh"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/entities"
)

type renderEntity struct {
	ID       terraform2d.EntityID
	Category terraform2d.EntityCategory
	*components.Spatial
	*components.Appearance
	*components.Animation
	*components.Movement
	*components.Ignore
	*components.Toggler
	*components.Temporary
}

// Render is a custom system
type Render struct {
	Win         *pixelgl.Window
	Spritesheet map[int]*pixel.Sprite

	player renderEntity
	arrow  renderEntity
	sword  renderEntity

	entities  []renderEntity
	obstacles []renderEntity
}

// AddEntity adds an entity to the system
func (s *Render) AddEntity(entity entities.Entity) {
	r := renderEntity{
		ID:        entity.ID(),
		Category:  entity.Category,
		Spatial:   entity.Spatial,
		Animation: entity.Animation,
		Movement:  entity.Movement,
		Temporary: entity.Temporary,
		Ignore:    entity.Ignore,
	}
	switch entity.Category {
	case zelduh.CategoryPlayer:
		s.player = r
	case zelduh.CategoryArrow:
		s.arrow = r
	case zelduh.CategorySword:
		s.sword = r
	case zelduh.CategoryExplosion:
		fallthrough
	case zelduh.CategoryHeart:
		fallthrough
	case zelduh.CategoryEnemy:
		fallthrough
	case zelduh.CategoryCollisionSwitch:
		fallthrough
	case zelduh.CategoryMovableObstacle:
		fallthrough
	case zelduh.CategoryWarp:
		fallthrough
	case zelduh.CategoryCoin:
		fallthrough
	default:
		if entity.Toggler != nil {
			r.Toggler = entity.Toggler
		}
		s.entities = append(s.entities, r)
	}
}

// RemoveAll removes all entities from one category
func (s *Render) RemoveAll(category terraform2d.EntityCategory) {
	switch category {
	case zelduh.CategoryEnemy:
		for i := len(s.entities) - 1; i >= 0; i-- {
			if s.entities[i].Category == zelduh.CategoryEnemy {
				s.entities = append(s.entities[:i], s.entities[i+1:]...)
			}
		}
	}
}

// RemoveEntity removes an entity by ID
func (s *Render) RemoveEntity(id terraform2d.EntityID) {
	for i := len(s.entities) - 1; i >= 0; i-- {
		if s.entities[i].ID == id {
			s.entities = append(s.entities[:i], s.entities[i+1:]...)
		}
	}
}

// RemoveAllEntities removes all entities
func (s *Render) RemoveAllEntities() {
	for i := len(s.entities) - 1; i >= 0; i-- {
		s.entities = append(s.entities[:i], s.entities[i+1:]...)
	}
}

// Update changes spatial data based on movement data
func (s *Render) Update() {

	for _, entity := range s.entities {
		if entity.Ignore != nil && !entity.Ignore.Value {
			if entity.Temporary != nil {
				if entity.Temporary.Expiration == 0 {
					entity.Temporary.OnExpiration()
					s.RemoveEntity(entity.ID)
				} else {
					entity.Temporary.Expiration--
				}
			}
			if entity.Toggler != nil {
				s.animateToggleFrame(entity)
			} else {
				s.animateDefault(entity)
			}
		}
	}

	player := s.player
	arrow := s.arrow
	sword := s.sword

	if !sword.Ignore.Value {
		s.animateDirections(player.Movement.Direction, sword)
	}

	if !arrow.Ignore.Value {
		s.animateDirections(player.Movement.Direction, arrow)
	}

	if sword.Ignore != nil && sword.Ignore.Value && arrow.Ignore.Value {
		s.animateDirections(player.Movement.Direction, player)
	} else {
		s.animateAttackDirection(player.Movement.Direction, player)
	}

}

func (s *Render) animateToggleFrame(entity renderEntity) {
	if anim := entity.Animation; anim != nil {
		if animData := anim.Map["default"]; animData != nil {
			var frameIndex int
			if !entity.Toggler.Enabled() {
				frameIndex = animData.Frames[0]
			} else {
				frameIndex = animData.Frames[1]
			}
			frame := s.Spritesheet[frameIndex]

			v := pixel.V(
				entity.Spatial.Rect.Min.X+entity.Spatial.Width/2,
				entity.Spatial.Rect.Min.Y+entity.Spatial.Height/2,
			)
			frame.Draw(s.Win, pixel.IM.Moved(v))
		}
	}
}

func (s *Render) animateDefault(entity renderEntity) {
	if anim := entity.Animation; anim != nil {
		if animData := anim.Map["default"]; animData != nil {
			rate := animData.FrameRateCount
			if rate < animData.FrameRate {
				rate++
			} else {
				rate = 0
			}
			animData.FrameRateCount = rate

			frameNum := animData.Frame
			if rate == animData.FrameRate {
				if frameNum < len(animData.Frames)-1 {
					frameNum++
				} else {
					frameNum = 0
				}
				animData.Frame = frameNum
			}

			frameIndex := animData.Frames[frameNum]
			frame := s.Spritesheet[frameIndex]

			v := pixel.V(
				entity.Spatial.Rect.Min.X+entity.Spatial.Width/2,
				entity.Spatial.Rect.Min.Y+entity.Spatial.Height/2,
			)
			frame.Draw(s.Win, pixel.IM.Moved(v))
		}
	}
}

func (s *Render) animateAttackDirection(dir terraform2d.Direction, entity renderEntity) {
	if anim := entity.Animation; anim != nil {
		var animData *components.AnimationData
		switch dir {
		case terraform2d.DirectionUp:
			animData = anim.Map["swordAttackUp"]
		case terraform2d.DirectionRight:
			animData = anim.Map["swordAttackRight"]
		case terraform2d.DirectionDown:
			animData = anim.Map["swordAttackDown"]
		case terraform2d.DirectionLeft:
			animData = anim.Map["swordAttackLeft"]
		}

		rate := animData.FrameRateCount
		if rate < animData.FrameRate {
			rate++
		} else {
			rate = 0
		}
		animData.FrameRateCount = rate

		frameNum := animData.Frame
		if rate == animData.FrameRate {
			if frameNum < len(animData.Frames)-1 {
				frameNum++
			} else {
				frameNum = 0
			}
			animData.Frame = frameNum
		}

		frameIndex := animData.Frames[frameNum]
		frame := s.Spritesheet[frameIndex]

		rect := entity.Spatial.Rect
		v := pixel.V(
			rect.Min.X+entity.Spatial.Width/2,
			rect.Min.Y+entity.Spatial.Height/2,
		)

		frame.Draw(s.Win, pixel.IM.Moved(v))
	}
}

func (s *Render) animateDirections(dir terraform2d.Direction, entity renderEntity) {
	// if entity.Spatial.HitBoxRadius > 0 {
	// 	shape := entity.Spatial.Shape
	// 	shape.Clear()
	// 	shape.Color = colornames.Yellow
	// 	shape.Push(pixel.V(
	// 		entity.Spatial.Rect.Min.X+entity.Spatial.Width/2,
	// 		entity.Spatial.Rect.Min.Y+entity.Spatial.Height/2,
	// 	))
	// 	// s.Push(entity.Spatial.Rect.Max)
	// 	shape.Circle(entity.Spatial.HitBoxRadius, 0)
	// 	shape.Draw(s.Win)
	// }
	if anim := entity.Animation; anim != nil {
		var animData *components.AnimationData
		switch dir {
		case terraform2d.DirectionUp:
			animData = anim.Map["up"]
		case terraform2d.DirectionRight:
			animData = anim.Map["right"]
		case terraform2d.DirectionDown:
			animData = anim.Map["down"]
		case terraform2d.DirectionLeft:
			animData = anim.Map["left"]
		}

		rate := animData.FrameRateCount
		if rate < animData.FrameRate {
			rate++
		} else {
			rate = 0
		}
		animData.FrameRateCount = rate

		frameNum := animData.Frame
		if rate == animData.FrameRate {
			if frameNum < len(animData.Frames)-1 {
				frameNum++
			} else {
				frameNum = 0
			}
			animData.Frame = frameNum
		}

		frameIndex := animData.Frames[frameNum]
		frame := s.Spritesheet[frameIndex]

		rect := entity.Spatial.Rect
		v := pixel.V(
			rect.Min.X+entity.Spatial.Width/2,
			rect.Min.Y+entity.Spatial.Height/2,
		)

		frame.Draw(s.Win, pixel.IM.Moved(v))
	}
}
