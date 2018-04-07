package systems

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miketmoore/zelduh/components"
	"github.com/miketmoore/zelduh/direction"
)

type renderEntity struct {
	ID int
	*components.Spatial
	*components.Appearance
	*components.Animation
	*components.Movement
	*components.Ignore
}

// Render is a custom system
type Render struct {
	Win               *pixelgl.Window
	playerEntity      renderEntity
	sword             renderEntity
	arrow             renderEntity
	coins             []renderEntity
	enemies           []renderEntity
	obstacles         []renderEntity
	moveableObstacles []renderEntity
	collisionSwitches []renderEntity
}

// AddPlayer adds the player to the system
func (s *Render) AddPlayer(
	appearance *components.Appearance,
	spatial *components.Spatial,
	animation *components.Animation,
	movement *components.Movement,
) {
	s.playerEntity = renderEntity{
		Appearance: appearance,
		Spatial:    spatial,
		Animation:  animation,
		Movement:   movement,
	}
}

// AddSword adds the sword to the system
func (s *Render) AddSword(
	appearance *components.Appearance,
	spatial *components.Spatial,
	ignore *components.Ignore,
	animation *components.Animation,
) {
	s.sword = renderEntity{
		Appearance: appearance,
		Spatial:    spatial,
		Ignore:     ignore,
		Animation:  animation,
	}
}

// AddArrow adds the sword to the system
func (s *Render) AddArrow(
	appearance *components.Appearance,
	spatial *components.Spatial,
	ignore *components.Ignore,
	animation *components.Animation,
) {
	s.arrow = renderEntity{
		Appearance: appearance,
		Spatial:    spatial,
		Ignore:     ignore,
		Animation:  animation,
	}
}

// AddMoveableObstacle adds a moveable obstacle to the system
func (s *Render) AddMoveableObstacle(id int, appearance *components.Appearance, spatial *components.Spatial) {
	s.moveableObstacles = append(s.moveableObstacles, renderEntity{
		ID:         id,
		Appearance: appearance,
		Spatial:    spatial,
	})
}

// AddCollisionSwitch adds a collision switch to the system
func (s *Render) AddCollisionSwitch(
	appearance *components.Appearance,
	spatial *components.Spatial,
	animation *components.Animation,
) {
	s.collisionSwitches = append(s.collisionSwitches, renderEntity{
		Appearance: appearance,
		Spatial:    spatial,
		Animation:  animation,
	})
}

// AddCoin adds the player to the system
func (s *Render) AddCoin(
	id int,
	appearance *components.Appearance,
	spatial *components.Spatial,
	animation *components.Animation,
) {
	s.coins = append(s.coins, renderEntity{
		ID:         id,
		Appearance: appearance,
		Spatial:    spatial,
		Animation:  animation,
	})
}

// AddEnemy adds an enemy to the system
func (s *Render) AddEnemy(id int, appearance *components.Appearance, spatial *components.Spatial, animation *components.Animation) {
	s.enemies = append(s.enemies, renderEntity{
		ID:         id,
		Appearance: appearance,
		Spatial:    spatial,
		Animation:  animation,
	})
}

// RemoveCoin removes the specified coin from the system
func (s *Render) RemoveCoin(id int) {
	for i := len(s.coins) - 1; i >= 0; i-- {
		coin := s.coins[i]
		if coin.ID == id {
			s.coins = append(s.coins[:i], s.coins[i+1:]...)
		}
	}
}

// RemoveEnemy removes the specified enemy from the system
func (s *Render) RemoveEnemy(id int) {
	for i := len(s.enemies) - 1; i >= 0; i-- {
		enemy := s.enemies[i]
		if enemy.ID == id {
			s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
		}
	}
}

// RemoveAllEnemies removes all enemy entities from the system
func (s *Render) RemoveAllEnemies() {
	for i := len(s.enemies) - 1; i >= 0; i-- {
		s.enemies = append(s.enemies[:i], s.enemies[i+1:]...)
	}
}

// RemoveAllCollisionSwitches removes all collision switches
func (s *Render) RemoveAllCollisionSwitches() {
	for i := len(s.collisionSwitches) - 1; i >= 0; i-- {
		s.collisionSwitches = append(s.collisionSwitches[:i], s.collisionSwitches[i+1:]...)
	}
}

// Update changes spatial data based on movement data
func (s *Render) Update() {

	for _, collisionSwitch := range s.collisionSwitches {
		if collisionSwitch.Animation != nil {
			s.animateDefault(collisionSwitch)
		} else {
			collisionSwitch.Shape.Clear()
			collisionSwitch.Shape.Color = collisionSwitch.Appearance.Color
			collisionSwitch.Shape.Push(collisionSwitch.Spatial.Rect.Min)
			collisionSwitch.Shape.Push(collisionSwitch.Spatial.Rect.Max)
			collisionSwitch.Shape.Rectangle(0)
			// collisionSwitch.Shape.Draw(s.Win)
		}

	}

	if !s.sword.Ignore.Value {
		s.animateDirections(s.playerEntity.Movement.Direction, s.sword)
	}

	if !s.arrow.Ignore.Value {
		s.animateDirections(s.playerEntity.Movement.Direction, s.arrow)
	}

	if s.sword.Ignore.Value && s.arrow.Ignore.Value {
		s.animateDirections(s.playerEntity.Movement.Direction, s.playerEntity)
	} else {
		s.animateAttackDirection(s.playerEntity.Movement.Direction, s.playerEntity)
	}

	for _, enemy := range s.enemies {
		s.animateDefault(enemy)
	}

	for _, coin := range s.coins {
		s.animateDefault(coin)
	}

	for _, obstacle := range s.moveableObstacles {
		obstacle.Shape.Clear()
		obstacle.Shape.Color = obstacle.Appearance.Color
		obstacle.Shape.Push(obstacle.Spatial.Rect.Min)
		obstacle.Shape.Push(obstacle.Spatial.Rect.Max)
		obstacle.Shape.Rectangle(0)
		obstacle.Shape.Draw(s.Win)
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

			frame := animData.Frames[frameNum]

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

		frame := animData.Frames[frameNum]

		rect := entity.Spatial.Rect
		v := pixel.V(
			rect.Min.X+entity.Spatial.Width/2,
			rect.Min.Y+entity.Spatial.Height/2,
		)

		frame.Draw(s.Win, pixel.IM.Moved(v))
	}
}

func (s *Render) animateDirections(dir direction.Name, entity renderEntity) {
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

		frame := animData.Frames[frameNum]

		rect := entity.Spatial.Rect
		v := pixel.V(
			rect.Min.X+entity.Spatial.Width/2,
			rect.Min.Y+entity.Spatial.Height/2,
		)

		frame.Draw(s.Win, pixel.IM.Moved(v))
	}
}
