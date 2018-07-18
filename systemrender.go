package zelduh

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/terraform2d"
)

type renderEntity struct {
	ID       terraform2d.EntityID
	Category terraform2d.EntityCategory
	*ComponentSpatial
	*ComponentAppearance
	*ComponentAnimation
	*ComponentMovement
	*ComponentIgnore
	*ComponentToggler
	*ComponentTemporary
}

// SystemRender is a custom system
type SystemRender struct {
	Win         *pixelgl.Window
	Spritesheet map[int]*pixel.Sprite

	player renderEntity
	arrow  renderEntity
	sword  renderEntity

	entities  []renderEntity
	obstacles []renderEntity
}

// AddEntity adds an entity to the system
func (s *SystemRender) AddEntity(entity Entity) {
	r := renderEntity{
		ID:                 entity.ID(),
		Category:           entity.Category,
		ComponentSpatial:   entity.ComponentSpatial,
		ComponentAnimation: entity.ComponentAnimation,
		ComponentMovement:  entity.ComponentMovement,
		ComponentTemporary: entity.ComponentTemporary,
		ComponentIgnore:    entity.ComponentIgnore,
	}
	switch entity.Category {
	case CategoryPlayer:
		s.player = r
	case CategoryArrow:
		s.arrow = r
	case CategorySword:
		s.sword = r
	case CategoryExplosion:
		fallthrough
	case CategoryHeart:
		fallthrough
	case CategoryEnemy:
		fallthrough
	case CategoryCollisionSwitch:
		fallthrough
	case CategoryMovableObstacle:
		fallthrough
	case CategoryWarp:
		fallthrough
	case CategoryCoin:
		fallthrough
	default:
		if entity.ComponentToggler != nil {
			r.ComponentToggler = entity.ComponentToggler
		}
		s.entities = append(s.entities, r)
	}
}

// RemoveAll removes all entities from one category
func (s *SystemRender) RemoveAll(category terraform2d.EntityCategory) {
	switch category {
	case CategoryEnemy:
		for i := len(s.entities) - 1; i >= 0; i-- {
			if s.entities[i].Category == CategoryEnemy {
				s.entities = append(s.entities[:i], s.entities[i+1:]...)
			}
		}
	}
}

// RemoveEntity removes an entity by ID
func (s *SystemRender) RemoveEntity(id terraform2d.EntityID) {
	for i := len(s.entities) - 1; i >= 0; i-- {
		if s.entities[i].ID == id {
			s.entities = append(s.entities[:i], s.entities[i+1:]...)
		}
	}
}

// RemoveAllEntities removes all entities
func (s *SystemRender) RemoveAllEntities() {
	for i := len(s.entities) - 1; i >= 0; i-- {
		s.entities = append(s.entities[:i], s.entities[i+1:]...)
	}
}

// Update changes spatial data based on movement data
func (s *SystemRender) Update() {

	for _, entity := range s.entities {
		if entity.ComponentIgnore != nil && !entity.ComponentIgnore.Value {
			if entity.ComponentTemporary != nil {
				if entity.ComponentTemporary.Expiration == 0 {
					entity.ComponentTemporary.OnExpiration()
					s.RemoveEntity(entity.ID)
				} else {
					entity.ComponentTemporary.Expiration--
				}
			}
			if entity.ComponentToggler != nil {
				s.animateToggleFrame(entity)
			} else {
				s.animateDefault(entity)
			}
		}
	}

	player := s.player
	arrow := s.arrow
	sword := s.sword

	if !sword.ComponentIgnore.Value {
		s.animateDirections(player.ComponentMovement.Direction, sword)
	}

	if !arrow.ComponentIgnore.Value {
		s.animateDirections(player.ComponentMovement.Direction, arrow)
	}

	if sword.ComponentIgnore != nil && sword.ComponentIgnore.Value && arrow.ComponentIgnore.Value {
		s.animateDirections(player.ComponentMovement.Direction, player)
	} else {
		s.animateAttackDirection(player.ComponentMovement.Direction, player)
	}

}

func (s *SystemRender) animateToggleFrame(entity renderEntity) {
	if anim := entity.ComponentAnimation; anim != nil {
		if animData := anim.Map["default"]; animData != nil {
			var frameIndex int
			if !entity.ComponentToggler.Enabled() {
				frameIndex = animData.Frames[0]
			} else {
				frameIndex = animData.Frames[1]
			}
			frame := s.Spritesheet[frameIndex]

			v := pixel.V(
				entity.ComponentSpatial.Rect.Min.X+entity.ComponentSpatial.Width/2,
				entity.ComponentSpatial.Rect.Min.Y+entity.ComponentSpatial.Height/2,
			)
			frame.Draw(s.Win, pixel.IM.Moved(v))
		}
	}
}

func (s *SystemRender) animateDefault(entity renderEntity) {
	if anim := entity.ComponentAnimation; anim != nil {
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
				entity.ComponentSpatial.Rect.Min.X+entity.ComponentSpatial.Width/2,
				entity.ComponentSpatial.Rect.Min.Y+entity.ComponentSpatial.Height/2,
			)
			frame.Draw(s.Win, pixel.IM.Moved(v))
		}
	}
}

func (s *SystemRender) animateAttackDirection(dir terraform2d.Direction, entity renderEntity) {
	if anim := entity.ComponentAnimation; anim != nil {
		var animData *ComponentAnimationData
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

		rect := entity.ComponentSpatial.Rect
		v := pixel.V(
			rect.Min.X+entity.ComponentSpatial.Width/2,
			rect.Min.Y+entity.ComponentSpatial.Height/2,
		)

		frame.Draw(s.Win, pixel.IM.Moved(v))
	}
}

func (s *SystemRender) animateDirections(dir terraform2d.Direction, entity renderEntity) {
	// if entity.ComponentSpatial.HitBoxRadius > 0 {
	// 	shape := entity.ComponentSpatial.Shape
	// 	shape.Clear()
	// 	shape.Color = colornames.Yellow
	// 	shape.Push(pixel.V(
	// 		entity.ComponentSpatial.Rect.Min.X+entity.ComponentSpatial.Width/2,
	// 		entity.ComponentSpatial.Rect.Min.Y+entity.ComponentSpatial.Height/2,
	// 	))
	// 	// s.Push(entity.ComponentSpatial.Rect.Max)
	// 	shape.Circle(entity.ComponentSpatial.HitBoxRadius, 0)
	// 	shape.Draw(s.Win)
	// }
	if anim := entity.ComponentAnimation; anim != nil {
		var animData *ComponentAnimationData
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

		rect := entity.ComponentSpatial.Rect
		v := pixel.V(
			rect.Min.X+entity.ComponentSpatial.Width/2,
			rect.Min.Y+entity.ComponentSpatial.Height/2,
		)

		frame.Draw(s.Win, pixel.IM.Moved(v))
	}
}
