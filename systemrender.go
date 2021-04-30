package zelduh

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// ComponentAnimationData contains data about animating one sequence of sprites
type ComponentAnimationData struct {
	Frames         []int
	Frame          int
	FrameRate      int
	FrameRateCount int
}

// ComponentAnimationMap indexes ComponentAnimationData by use/context
type ComponentAnimationMap map[string]*ComponentAnimationData

// ComponentAnimation contains everything necessary to animate basic characters
type ComponentAnimation struct {
	ComponentAnimationByName ComponentAnimationMap
}

// ComponentAppearance contains data about visual appearance
type ComponentAppearance struct {
	Color color.RGBA
}

type renderEntity struct {
	ID       EntityID
	Category EntityCategory
	*ComponentSpatial
	*ComponentAppearance
	*ComponentAnimation
	*ComponentMovement
	*ComponentIgnore
	*ComponentToggler
}

func (entity *renderEntity) shouldNotIgnore() bool {
	return entity.ComponentIgnore == nil || (entity.ComponentIgnore != nil && !entity.ComponentIgnore.Value)
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

	TemporarySystem *TemporarySystem
}

func NewRenderSystem(
	window *pixelgl.Window,
	spriteMap SpriteMap,
	activeSpaceRectangle ActiveSpaceRectangle,
	tileSize float64,
	temporarySystem *TemporarySystem,
) RenderSystem {
	return RenderSystem{
		Win:                  window,
		SpriteMap:            spriteMap,
		ActiveSpaceRectangle: activeSpaceRectangle,
		TileSize:             tileSize,
		TemporarySystem:      temporarySystem,
	}
}

// AddEntity adds an entity to the system
func (s *RenderSystem) AddEntity(entity Entity) {
	r := renderEntity{
		ID:                 entity.ID(),
		Category:           entity.Category,
		ComponentSpatial:   entity.ComponentSpatial,
		ComponentAnimation: entity.ComponentAnimation,
		ComponentMovement:  entity.ComponentMovement,
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
func (s *RenderSystem) Update() error {

	// DrawActiveSpace(s.Win, s.ActiveSpaceRectangle)

	for _, entity := range s.entities {

		if entity.shouldNotIgnore() {

			if entity.Category == CategoryRectangle {
				s.drawRectangle(entity)
			}

			if s.TemporarySystem.IsTemporary(entity.ID) {
				if s.TemporarySystem.IsExpired(entity.ID) {
					s.TemporarySystem.CallOnExpirationHandler(entity.ID)
					s.RemoveEntity(entity.ID)
				} else {
					s.TemporarySystem.DecrementExpiration(entity.ID)
				}
			}

			s.drawDefaultFrame(entity)

		}
	}

	player := s.player
	arrow := s.arrow
	sword := s.sword

	if !sword.ComponentIgnore.Value {
		s.drawByPlayerDirection(sword)
	}

	if !arrow.ComponentIgnore.Value {
		s.drawByPlayerDirection(arrow)
	}

	if sword.ComponentIgnore != nil && sword.ComponentIgnore.Value && arrow.ComponentIgnore.Value {
		s.drawByPlayerDirection(player)
	} else {
		animDataKey := swordComponentAnimationByDirection[player.ComponentMovement.Direction]
		componentAnimationData, ok := getComponentAnimationByName(player, animDataKey)
		if ok {
			s.drawSprite(componentAnimationData, player)
		}
	}

	return nil
}

func (s *RenderSystem) drawDefaultFrame(entity renderEntity) {
	componentAnimationData, ok := getComponentAnimationByName(entity, "default")
	if ok {
		s.drawSprite(componentAnimationData, entity)
	}
}

func (s *RenderSystem) drawByPlayerDirection(entity renderEntity) {
	animDataKey := playerComponentAnimationByDirection[s.player.ComponentMovement.Direction]
	componentAnimationData, ok := getComponentAnimationByName(entity, animDataKey)
	if ok {
		s.drawSprite(componentAnimationData, entity)
	}
}

var playerComponentAnimationByDirection = map[Direction]string{
	DirectionUp:    "up",
	DirectionRight: "right",
	DirectionDown:  "down",
	DirectionLeft:  "left",
}

var swordComponentAnimationByDirection = map[Direction]string{
	DirectionUp:    "swordAttackUp",
	DirectionRight: "swordAttackRight",
	DirectionDown:  "swordAttackDown",
	DirectionLeft:  "swordAttackLeft",
}

func getComponentAnimationByName(entity renderEntity, name string) (*ComponentAnimationData, bool) {
	componentAnimation := entity.ComponentAnimation
	if componentAnimation == nil {
		return nil, false
	}

	animationData, ok := componentAnimation.ComponentAnimationByName[name]

	return animationData, ok
}

func (s *RenderSystem) drawSprite(
	animData *ComponentAnimationData,
	entity renderEntity,
) {
	frame, _, matrix := s.getSpriteDrawData(animData, entity.ComponentSpatial)
	frame.Draw(s.Win, matrix)
}

func (s *RenderSystem) getSpriteDrawData(
	animData *ComponentAnimationData,
	spatialComponent *ComponentSpatial,
) (*pixel.Sprite, pixel.Vec, pixel.Matrix) {
	rate := determineFrameRate(animData)

	animData.FrameRateCount = rate

	frameNum := determineFrameNumber(animData)

	frameIndex := animData.Frames[frameNum]

	frame := s.SpriteMap[frameIndex]

	vector := buildSpriteVector(spatialComponent, s.ActiveSpaceRectangle)
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

func DrawActiveSpace(window *pixelgl.Window, activeSpaceRectangle ActiveSpaceRectangle) {
	rect := imdraw.New(nil)
	rect.Color = colornames.Blue

	vectorX := activeSpaceRectangle.X
	vectorY := activeSpaceRectangle.Y
	point := pixel.V(vectorX, vectorY)

	rect.Push(point)

	point2 := pixel.V(
		point.X+(activeSpaceRectangle.Width),
		point.Y+(activeSpaceRectangle.Height),
	)
	rect.Push(point2)

	rect.Rectangle(5)

	rect.Draw(window)
}
