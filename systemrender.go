package zelduh

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type renderEntity struct {
	ID       EntityID
	Category EntityCategory
	*ComponentSpatial
	*ComponentAppearance
	*ComponentAnimation
	*ComponentMovement
	*ComponentIgnore
	*ComponentToggler
	*ComponentTemporary
}

// RenderSystem is a custom system
type RenderSystem struct {
	Win       *pixelgl.Window
	SpriteMap SpriteMap
	TileSize  float64

	player renderEntity
	arrow  renderEntity
	sword  renderEntity

	entities             []renderEntity
	obstacles            []renderEntity
	ActiveSpaceRectangle ActiveSpaceRectangle
}

// AddEntity adds an entity to the system
func (s *RenderSystem) AddEntity(entity Entity) {
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
	case CategoryIgnore:
		fallthrough
	default:
		if entity.ComponentToggler != nil {
			r.ComponentToggler = entity.ComponentToggler
		}
		s.entities = append(s.entities, r)
	}
}

// RemoveAll removes all entities from one category
func (s *RenderSystem) RemoveAll(category EntityCategory) {
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
func (s *RenderSystem) RemoveEntity(id EntityID) {
	for i := len(s.entities) - 1; i >= 0; i-- {
		if s.entities[i].ID == id {
			s.entities = append(s.entities[:i], s.entities[i+1:]...)
		}
	}
}

// RemoveAllEntities removes all entities
func (s *RenderSystem) RemoveAllEntities() {
	for i := len(s.entities) - 1; i >= 0; i-- {
		s.entities = append(s.entities[:i], s.entities[i+1:]...)
	}
}

// Update changes spatial data based on movement data
func (s *RenderSystem) Update() {

	for _, entity := range s.entities {

		if entity.Category == CategoryRectangle {
			s.drawRectangle(entity)
		}

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

func (s *RenderSystem) drawSprite(
	animData *ComponentAnimationData,
	entity renderEntity,
) {
	frame, _, matrix := GetSpriteDrawData(animData, entity.ComponentSpatial, s.ActiveSpaceRectangle, s.SpriteMap)
	frame.Draw(s.Win, matrix)
}

func (s *RenderSystem) animateToggleFrame(entity renderEntity) {
	if anim := entity.ComponentAnimation; anim != nil {
		if animData := anim.ComponentAnimationByName["default"]; animData != nil {
			s.drawSprite(animData, entity)
		}
	}
}

func (s *RenderSystem) animateDefault(entity renderEntity) {
	if anim := entity.ComponentAnimation; anim != nil {
		if animData := anim.ComponentAnimationByName["default"]; animData != nil {
			s.drawSprite(animData, entity)
		}
	}
}

func (s *RenderSystem) animateAttackDirection(dir Direction, entity renderEntity) {
	if anim := entity.ComponentAnimation; anim != nil {
		animData := getAttackAnimationDataByDirection(anim, dir)
		s.drawSprite(animData, entity)
	}
}

func (s *RenderSystem) animateDirections(dir Direction, entity renderEntity) {
	if anim := entity.ComponentAnimation; anim != nil {
		animData := getDirectionAnimationDataByDirection(anim, dir)
		s.drawSprite(animData, entity)
	}
}

func GetSpriteDrawData(
	animData *ComponentAnimationData,
	spatialComponent *ComponentSpatial,
	activeSpaceRectangle ActiveSpaceRectangle,
	spriteMap SpriteMap,
) (*pixel.Sprite, pixel.Vec, pixel.Matrix) {
	rate := determineFrameRate(animData)

	animData.FrameRateCount = rate

	frameNum := determineFrameNumber(animData)

	frameIndex := animData.Frames[frameNum]

	frame := spriteMap[frameIndex]

	vector := buildSpriteVector(spatialComponent, activeSpaceRectangle)
	matrix := buildSpriteMatrix(spatialComponent, vector)

	return frame, vector, matrix
}

func determineFrameRate(animData *ComponentAnimationData) int {
	rate := animData.FrameRateCount
	if rate < animData.FrameRate {
		return rate + 1
	}
	return 0
}

func determineFrameNumber(animData *ComponentAnimationData) int {
	rate := animData.FrameRateCount
	frameNum := animData.Frame
	if rate == animData.FrameRate {
		if frameNum < len(animData.Frames)-1 {
			frameNum++
		} else {
			frameNum = 0
		}
		animData.Frame = frameNum
	}
	return frameNum
}

func buildSpriteVector(spatialComponent *ComponentSpatial, activeSpaceRectangle ActiveSpaceRectangle) pixel.Vec {
	vectorX := spatialComponent.Rect.Center().X + activeSpaceRectangle.X
	vectorY := spatialComponent.Rect.Center().Y + activeSpaceRectangle.Y
	return pixel.V(vectorX, vectorY)
}

func buildSpriteMatrix(spatialComponent *ComponentSpatial, vector pixel.Vec) pixel.Matrix {

	matrix := pixel.IM.Moved(vector)

	if spatialComponent.Transform != nil {
		// Transform
		degrees := spatialComponent.Transform.Rotation
		radians := degrees * math.Pi / 180
		matrix = matrix.Rotated(vector, radians)
	}

	return matrix
}

func getAttackAnimationDataByDirection(anim *ComponentAnimation, dir Direction) *ComponentAnimationData {
	var animData *ComponentAnimationData
	switch dir {
	case DirectionUp:
		animData = anim.ComponentAnimationByName["swordAttackUp"]
	case DirectionRight:
		animData = anim.ComponentAnimationByName["swordAttackRight"]
	case DirectionDown:
		animData = anim.ComponentAnimationByName["swordAttackDown"]
	case DirectionLeft:
		animData = anim.ComponentAnimationByName["swordAttackLeft"]
	}
	return animData
}

func getDirectionAnimationDataByDirection(anim *ComponentAnimation, dir Direction) *ComponentAnimationData {
	var animData *ComponentAnimationData
	switch dir {
	case DirectionUp:
		animData = anim.ComponentAnimationByName["up"]
	case DirectionRight:
		animData = anim.ComponentAnimationByName["right"]
	case DirectionDown:
		animData = anim.ComponentAnimationByName["down"]
	case DirectionLeft:
		animData = anim.ComponentAnimationByName["left"]
	}
	return animData
}

// drawRectangle draws a rectangle of any dimensions
func (s *RenderSystem) drawRectangle(entity renderEntity) {

	spatialData := entity.ComponentSpatial

	rect := spatialData.Shape
	rect.Color = spatialData.Color

	vectorX := s.ActiveSpaceRectangle.X + (s.TileSize * entity.ComponentSpatial.Rect.Min.X)
	vectorY := s.ActiveSpaceRectangle.Y + (s.TileSize * entity.ComponentSpatial.Rect.Min.Y)
	point := pixel.V(vectorX, vectorY)

	rect.Push(point)

	point2 := pixel.V(
		point.X+(spatialData.Width*48),
		point.Y+(spatialData.Height*48),
	)
	rect.Push(point2)

	rect.Rectangle(0)

	rect.Draw(s.Win)
}
