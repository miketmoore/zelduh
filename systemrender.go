package zelduh

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type componentShape struct {
	Shape *imdraw.IMDraw
}

func NewComponentShape() *componentShape {
	return &componentShape{
		Shape: imdraw.New(nil),
	}
}

type componentRotation struct {
	Degrees float64
}

func NewComponentRotation(degrees float64) *componentRotation {
	return &componentRotation{Degrees: degrees}
}

func NewComponentAnimation(animationConfig AnimationConfig, frameRate int) *componentAnimation {
	component := componentAnimation{
		ComponentAnimationByName: componentAnimationMap{},
	}

	for key, val := range animationConfig {
		component.ComponentAnimationByName[key] = NewComponentAnimationData(val, frameRate)
	}

	return &component
}

// componentAnimationData contains data about animating one sequence of sprites
type componentAnimationData struct {
	Frames         []int
	Frame          int
	FrameRate      int
	FrameRateCount int
}

func NewComponentAnimationData(frames []int, frameRate int) *componentAnimationData {
	return &componentAnimationData{
		Frames:    frames,
		FrameRate: frameRate,
	}
}

// componentAnimationMap indexes componentAnimationData by use/context
type componentAnimationMap map[string]*componentAnimationData

// componentAnimation contains everything necessary to animate basic characters
type componentAnimation struct {
	ComponentAnimationByName componentAnimationMap
}

type componentColor struct {
	Color color.RGBA
}

func NewComponentColor(color color.RGBA) *componentColor {
	return &componentColor{Color: color}
}

type renderEntity struct {
	ID       EntityID
	Category EntityCategory
	*componentRotation
	*componentColor
	*componentAnimation
	*componentMovement
	*componentIgnore
	*componentToggler
	*componentDimensions
	*componentRectangle
	*componentShape
}

func newRenderEntity(entity Entity) renderEntity {
	return renderEntity{
		ID:                  entity.ID(),
		Category:            entity.Category,
		componentRotation:   entity.componentRotation,
		componentAnimation:  entity.componentAnimation,
		componentMovement:   entity.componentMovement,
		componentIgnore:     entity.componentIgnore,
		componentColor:      entity.componentColor,
		componentDimensions: entity.componentDimensions,
		componentRectangle:  entity.componentRectangle,
		componentShape:      entity.componentShape,
	}
}

func (entity *renderEntity) shouldNotIgnore() bool {
	return entity.componentIgnore == nil || (entity.componentIgnore != nil && !entity.componentIgnore.Value)
}

// RenderSystem is a custom system
type RenderSystem struct {
	Win       *pixelgl.Window
	SpriteMap SpriteMap
	TileSize  float64

	player renderEntity
	arrow  renderEntity
	sword  renderEntity

	entities []renderEntity
	// obstacles            []renderEntity
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
	r := newRenderEntity(entity)
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
		if entity.componentToggler != nil {
			r.componentToggler = entity.componentToggler
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

	if !sword.componentIgnore.Value {
		s.drawByPlayerDirection(sword)
	}

	if !arrow.componentIgnore.Value {
		s.drawByPlayerDirection(arrow)
	}

	if sword.componentIgnore != nil && sword.componentIgnore.Value && arrow.componentIgnore.Value {
		s.drawByPlayerDirection(player)
	} else {
		animDataKey := swordComponentAnimationByDirection[player.componentMovement.Direction]
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
	animDataKey := playerComponentAnimationByDirection[s.player.componentMovement.Direction]
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

func getComponentAnimationByName(entity renderEntity, name string) (*componentAnimationData, bool) {
	componentAnimation := entity.componentAnimation
	if componentAnimation == nil {
		return nil, false
	}

	animationData, ok := componentAnimation.ComponentAnimationByName[name]

	return animationData, ok
}

func (s *RenderSystem) drawSprite(
	animData *componentAnimationData,
	entity renderEntity,
) {
	frame, _, matrix := s.getSpriteDrawData(animData, entity.componentRectangle, entity.componentRotation)
	frame.Draw(s.Win, matrix)
}

func (s *RenderSystem) getSpriteDrawData(
	animData *componentAnimationData,
	componentRectangle *componentRectangle,
	rotationComponent *componentRotation,
) (*pixel.Sprite, pixel.Vec, pixel.Matrix) {
	rate := determineFrameRate(animData)

	animData.FrameRateCount = rate

	frameNum := determineFrameNumber(animData)

	frameIndex := animData.Frames[frameNum]

	frame := s.SpriteMap[frameIndex]

	vector := buildSpriteVector(componentRectangle, s.ActiveSpaceRectangle)
	matrix := buildSpriteMatrix(rotationComponent, vector)

	return frame, vector, matrix
}

func determineFrameRate(animData *componentAnimationData) int {
	rate := animData.FrameRateCount
	if rate < animData.FrameRate {
		return rate + 1
	}
	return 0
}

func determineFrameNumber(animData *componentAnimationData) int {
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

func buildSpriteVector(componentRectangle *componentRectangle, activeSpaceRectangle ActiveSpaceRectangle) pixel.Vec {
	vectorX := componentRectangle.Rect.Center().X + activeSpaceRectangle.X
	vectorY := componentRectangle.Rect.Center().Y + activeSpaceRectangle.Y
	return pixel.V(vectorX, vectorY)
}

func buildSpriteMatrix(rotationComponent *componentRotation, vector pixel.Vec) pixel.Matrix {

	matrix := pixel.IM.Moved(vector)

	if rotationComponent != nil {
		radians := rotationComponent.Degrees * math.Pi / 180
		matrix = matrix.Rotated(vector, radians)
	}

	return matrix
}

// drawRectangle draws a rectangle of any dimensions
func (s *RenderSystem) drawRectangle(entity renderEntity) {

	rect := entity.componentShape.Shape
	rect.Color = entity.componentColor.Color

	vectorX := s.ActiveSpaceRectangle.X + (s.TileSize * entity.componentRectangle.Rect.Min.X)
	vectorY := s.ActiveSpaceRectangle.Y + (s.TileSize * entity.componentRectangle.Rect.Min.Y)
	point := pixel.V(vectorX, vectorY)

	rect.Push(point)

	point2 := pixel.V(
		point.X+(entity.componentDimensions.Width*s.TileSize),
		point.Y+(entity.componentDimensions.Height*s.TileSize),
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
