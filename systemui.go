package zelduh

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miketmoore/zelduh/core/tmx"
	"golang.org/x/image/colornames"
)

type UISystem struct {
	Window               *pixelgl.Window
	Text                 *text.Text
	activeSpaceRectangle ActiveSpaceRectangle
	spriteMap            SpriteMap
	mapDrawData          MapDrawData
	currLocaleMsgs       LocaleMessagesMap
	tileSize             float64
	frameRate            int
	windowConfig         WindowConfig
	nonObstacleSprites   map[int]bool
	levelManager         *LevelManager
	entityFactory        *EntityFactory
}

func NewUISystem(
	currLocaleMsgs LocaleMessagesMap,
	windowConfig WindowConfig,
	activeSpaceRectangle ActiveSpaceRectangle,
	spriteMap SpriteMap,
	mapDrawData MapDrawData,
	tileSize float64,
	frameRate int,
	levelManager *LevelManager,
	nonObstacleSprites map[int]bool,
	entityFactory *EntityFactory,
) UISystem {

	// Initialize text
	orig := pixel.V(20, 50)
	txt := text.New(orig, text.Atlas7x13)
	txt.Color = colornames.Black

	// Initialize window
	win, err := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  currLocaleMsgs["gameTitle"],
			Bounds: pixel.R(windowConfig.X, windowConfig.Y, windowConfig.Width, windowConfig.Height),
			VSync:  true,
		},
	)
	if err != nil {
		fmt.Println("Initializing GUI window failed:")
		fmt.Println(err)
		os.Exit(1)
	}

	return UISystem{
		Window:               win,
		Text:                 txt,
		activeSpaceRectangle: activeSpaceRectangle,
		spriteMap:            spriteMap,
		mapDrawData:          mapDrawData,
		currLocaleMsgs:       currLocaleMsgs,
		tileSize:             tileSize,
		frameRate:            frameRate,
		windowConfig:         windowConfig,
		levelManager:         levelManager,
		nonObstacleSprites:   nonObstacleSprites,
		entityFactory:        entityFactory,
	}
}

func (s *UISystem) DrawCenterText(str string, c color.RGBA) {
	s.Text.Clear()
	s.Text.Color = c
	fmt.Fprintln(s.Text, str)
	s.Text.Draw(s.Window, pixel.IM.Moved(s.Window.Bounds().Center().Sub(s.Text.Bounds().Center())))
}

func (s *UISystem) DrawMapBackground(color color.Color) {
	shape := imdraw.New(nil)
	shape.Color = color
	shape.Push(pixel.V(s.activeSpaceRectangle.X, s.activeSpaceRectangle.Y))
	shape.Push(pixel.V(s.activeSpaceRectangle.X+s.activeSpaceRectangle.Width, s.activeSpaceRectangle.Y+s.activeSpaceRectangle.Height))
	shape.Rectangle(0)
	shape.Draw(s.Window)
}

func (s *UISystem) DrawScreenStart() {
	s.Window.Clear(colornames.Darkgray)
	s.DrawMapBackground(colornames.White)
	s.DrawCenterText(s.currLocaleMsgs["gameTitle"], colornames.Black)
}

func (s *UISystem) DrawMapBackgroundImage(
	name tmx.TMXFileName,
	modX, modY float64,
) error {

	data, dataOk := s.mapDrawData[name]
	if strings.TrimSpace(string(name)) == "" {
		return fmt.Errorf("background image name is an empty string")
	}
	if !dataOk {
		return fmt.Errorf("DrawMapBackgroundImage: tmx file not found in map by name=%s", name)
	}
	for _, spriteData := range data.Data {
		if spriteData.SpriteID != 0 {
			sprite := s.spriteMap[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+s.activeSpaceRectangle.X+modX+s.tileSize/2,
				vec.Y+s.activeSpaceRectangle.Y+modY+s.tileSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(s.Window, matrix)
		}
	}
	return nil
}

func (s *UISystem) DrawMask() {
	// top
	shape := imdraw.New(nil)
	shape.Color = colornames.White
	shape.Push(pixel.V(0, s.activeSpaceRectangle.Y+s.activeSpaceRectangle.Height))
	shape.Push(pixel.V(s.windowConfig.Width, s.activeSpaceRectangle.Y+s.activeSpaceRectangle.Height+(s.windowConfig.Height-(s.activeSpaceRectangle.Y+s.activeSpaceRectangle.Height))))
	shape.Rectangle(0)
	shape.Draw(s.Window)

	// bottom
	shape = imdraw.New(nil)
	shape.Color = colornames.White
	shape.Push(pixel.V(0, 0))
	shape.Push(pixel.V(s.windowConfig.Width, (s.windowConfig.Height - (s.activeSpaceRectangle.Y + s.activeSpaceRectangle.Height))))
	shape.Rectangle(0)
	shape.Draw(s.Window)

	// left
	shape = imdraw.New(nil)
	shape.Color = colornames.White
	shape.Push(pixel.V(0, 0))
	shape.Push(pixel.V(0+s.activeSpaceRectangle.X, s.windowConfig.Height))
	shape.Rectangle(0)
	shape.Draw(s.Window)

	// right
	shape = imdraw.New(nil)
	shape.Color = colornames.White
	shape.Push(pixel.V(s.activeSpaceRectangle.X+s.activeSpaceRectangle.Width, s.activeSpaceRectangle.Y))
	shape.Push(pixel.V(s.windowConfig.Width, s.windowConfig.Height))
	shape.Rectangle(0)
	shape.Draw(s.Window)
}

func (s *UISystem) DrawPauseScreen() {
	s.Window.Clear(colornames.Darkgray)
	s.DrawMapBackground(colornames.White)
	s.DrawCenterText(s.currLocaleMsgs["pauseScreenMessage"], colornames.Black)
}

func (s *UISystem) DrawGameOverScreen() {
	s.Window.Clear(colornames.Darkgray)
	s.DrawMapBackground(colornames.Black)
	s.DrawCenterText(s.currLocaleMsgs["gameOverScreenMessage"], colornames.White)
}

func (s *UISystem) DrawLevelBackground(roomName tmx.TMXFileName) error {
	s.Window.Clear(colornames.Darkgray)
	s.DrawMapBackground(colornames.White)

	err := s.DrawMapBackgroundImage(
		roomName,
		0, 0,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
