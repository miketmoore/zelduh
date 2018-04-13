package systems

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/categories"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
	"github.com/miketmoore/zelduh/entities"
)

type renderEntity struct {
	ID       entities.EntityID
	Category categories.Category
	*components.Spatial
	*components.Appearance
	*components.Animation
	*components.Movement
	*components.Ignore
}

// Render is a custom system
type Render struct {
	Win         *pixelgl.Window
	Spritesheet map[int]*pixel.Sprite

	player renderEntity
	sword  renderEntity

	defaultEntities []renderEntity

	generic           []renderEntity
	arrow             renderEntity
	obstacles         []renderEntity
	collisionSwitches []renderEntity
	hearts            []renderEntity
}

// AddEntity adds an entity to the system
func (s *Render) AddEntity(entity entities.Entity) {
	r := renderEntity{
		entity.ID,
		entity.Category,
		entity.Spatial,
		entity.Appearance,
		entity.Animation,
		entity.Movement,
		entity.Ignore,
	}
	switch entity.Category {
	case categories.Player:
		s.player = r
	case categories.Sword:
		s.sword = r
	case categories.Arrow:
		s.arrow = r
	case categories.Explosion:
		s.generic = append(s.generic, r)
	case categories.MovableObstacle:
		s.defaultEntities = append(s.defaultEntities, r)
	case categories.CollisionSwitch:
		s.collisionSwitches = append(s.collisionSwitches, r)
	case categories.Coin:
		s.defaultEntities = append(s.defaultEntities, r)
	case categories.Enemy:
		s.defaultEntities = append(s.defaultEntities, r)
	case categories.Heart:
		s.hearts = append(s.hearts, r)
	}
}

// Remove removes the entity from the system
func (s *Render) Remove(category categories.Category, id entities.EntityID) {
	switch category {
	case categories.Explosion:
		for i := len(s.generic) - 1; i >= 0; i-- {
			generic := s.generic[i]
			if generic.ID == id {
				s.generic = append(s.generic[:i], s.generic[i+1:]...)
			}
		}
	case categories.Coin:
		fallthrough
	case categories.Enemy:
		for i := len(s.defaultEntities) - 1; i >= 0; i-- {
			entity := s.defaultEntities[i]
			if entity.ID == id {
				s.defaultEntities = append(s.defaultEntities[:i], s.defaultEntities[i+1:]...)
			}
		}
	}
}

// RemoveAll removes all entities from one category
func (s *Render) RemoveAll(category categories.Category) {
	switch category {
	case categories.Enemy:
		for i := len(s.defaultEntities) - 1; i >= 0; i-- {
			s.defaultEntities = append(s.defaultEntities[:i], s.defaultEntities[i+1:]...)
		}
	case categories.CollisionSwitch:
		for i := len(s.collisionSwitches) - 1; i >= 0; i-- {
			s.collisionSwitches = append(s.collisionSwitches[:i], s.collisionSwitches[i+1:]...)
		}
	case categories.MovableObstacle:
		for i := len(s.defaultEntities) - 1; i >= 0; i-- {
			s.defaultEntities = append(s.defaultEntities[:i], s.defaultEntities[i+1:]...)
		}
	}
}

// Update changes spatial data based on movement data
func (s *Render) Update() {

	for _, collisionSwitch := range s.collisionSwitches {
		if collisionSwitch.Animation != nil {
			s.animateDefault(collisionSwitch)
		} else {
			// Draw an invisible collision switch
			collisionSwitch.Shape.Clear()
			collisionSwitch.Shape.Color = collisionSwitch.Appearance.Color
			collisionSwitch.Shape.Push(collisionSwitch.Spatial.Rect.Min)
			collisionSwitch.Shape.Push(collisionSwitch.Spatial.Rect.Max)
			collisionSwitch.Shape.Rectangle(0)
		}

	}

	for _, generic := range s.generic {
		s.animateDefault(generic)
		if generic.Animation.Expiration == 0 {
			generic.Animation.OnExpiration()
			s.Remove(generic.Category, generic.ID)
		} else {
			generic.Animation.Expiration--
		}
	}

	player := s.player

	if !s.sword.Ignore.Value {
		s.animateDirections(player.Movement.Direction, s.sword)
	}

	if !s.arrow.Ignore.Value {
		s.animateDirections(player.Movement.Direction, s.arrow)
	}

	if s.sword.Ignore.Value && s.arrow.Ignore.Value {
		s.animateDirections(player.Movement.Direction, player)
	} else {
		s.animateAttackDirection(player.Movement.Direction, player)
	}

	for _, entity := range s.defaultEntities {
		s.animateDefault(entity)
	}

	for _, entity := range s.hearts {
		// fmt.Printf("Heart: %v\n", entity)
		s.animateDefault(entity)
	}

}

func (s *Render) animateDefault(entity renderEntity) {
	if anim := entity.Animation; anim != nil {
		if animData := anim.Default; animData != nil {
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

func (s *Render) animateAttackDirection(dir direction.Name, entity renderEntity) {
	if anim := entity.Animation; anim != nil {
		var animData *components.AnimationData
		switch dir {
		case direction.Up:
			animData = anim.SwordAttackUp
		case direction.Right:
			animData = anim.SwordAttackRight
		case direction.Down:
			animData = anim.SwordAttackDown
		case direction.Left:
			animData = anim.SwordAttackLeft
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

func (s *Render) animateDirections(dir direction.Name, entity renderEntity) {
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
		case direction.Up:
			animData = anim.Up
		case direction.Right:
			animData = anim.Right
		case direction.Down:
			animData = anim.Down
		case direction.Left:
			animData = anim.Left
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
