package zelduh

import (
	"fmt"
	"image/color"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type UISystem struct {
	Window               *pixelgl.Window
	Text                 *text.Text
	activeSpaceRectangle ActiveSpaceRectangle
}

func NewUISystem(
	currLocaleMsgs LocaleMessagesMap,
	windowConfig WindowConfig,
	activeSpaceRectangle ActiveSpaceRectangle,
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
	}
}

func DrawCenterText(win *pixelgl.Window, txt *text.Text, s string, c color.RGBA) {
	txt.Clear()
	txt.Color = c
	fmt.Fprintln(txt, s)
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
}

func (s *UISystem) DrawMapBackground(color color.Color) {
	shape := imdraw.New(nil)
	shape.Color = color
	shape.Push(pixel.V(s.activeSpaceRectangle.X, s.activeSpaceRectangle.Y))
	shape.Push(pixel.V(s.activeSpaceRectangle.X+s.activeSpaceRectangle.Width, s.activeSpaceRectangle.Y+s.activeSpaceRectangle.Height))
	shape.Rectangle(0)
	shape.Draw(s.Window)
}

func (s *UISystem) DrawScreenStart(win *pixelgl.Window, txt *text.Text, currLocaleMsgs LocaleMessagesMap, activeSpaceRectangle ActiveSpaceRectangle) {
	win.Clear(colornames.Darkgray)
	s.DrawMapBackground(colornames.White)
	DrawCenterText(win, txt, currLocaleMsgs["gameTitle"], colornames.Black)
}

func DrawMapBackgroundImage(
	win *pixelgl.Window,
	spriteMap SpriteMap,
	mapDrawData MapDrawData,
	name RoomName,
	modX, modY float64,
	tileSize float64,
	activeSpaceRectangle ActiveSpaceRectangle,
) {

	d := mapDrawData[name]
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			sprite := spriteMap[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+activeSpaceRectangle.X+modX+tileSize/2,
				vec.Y+activeSpaceRectangle.Y+modY+tileSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func AddUICoin(systemsManager *SystemsManager, entityConfigPresetFnManager *EntityConfigPresetFnManager, frameRate int) {
	presetFn := entityConfigPresetFnManager.GetPreset("uiCoin")
	entityConfig := presetFn(Coordinates{X: 4, Y: 14})
	coin := BuildEntityFromConfig(entityConfig, systemsManager.NewEntityID(), frameRate)
	systemsManager.AddEntity(coin)
}

func DrawObstaclesPerMapTiles(
	systemsManager *SystemsManager,
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	roomByIDMap RoomByIDMap,
	mapDrawData MapDrawData,
	roomID *RoomID,
	modX, modY float64,
	tileSize float64,
	frameRate int,
	nonObstacleSprites map[int]bool,
	activeSpaceRectangle ActiveSpaceRectangle,
) []Entity {
	d := mapDrawData[roomByIDMap[*roomID].Name]
	obstacles := []Entity{}
	mod := 0.5
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+activeSpaceRectangle.X+modX+tileSize/2,
				vec.Y+activeSpaceRectangle.Y+modY+tileSize/2,
			)

			if _, ok := nonObstacleSprites[spriteData.SpriteID]; !ok {
				coordinates := Coordinates{
					X: movedVec.X/tileSize - mod,
					Y: movedVec.Y/tileSize - mod,
				}
				id := systemsManager.NewEntityID()
				obstacle := BuildEntityFromConfig(entityConfigPresetFnManager.GetPreset("obstacle")(coordinates), id, frameRate)
				obstacles = append(obstacles, obstacle)
			}
		}
	}
	return obstacles
}

func DrawMask(win *pixelgl.Window, windowConfig WindowConfig, activeSpaceRectangle ActiveSpaceRectangle) {
	// top
	s := imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, activeSpaceRectangle.Y+activeSpaceRectangle.Height))
	s.Push(pixel.V(windowConfig.Width, activeSpaceRectangle.Y+activeSpaceRectangle.Height+(windowConfig.Height-(activeSpaceRectangle.Y+activeSpaceRectangle.Height))))
	s.Rectangle(0)
	s.Draw(win)

	// bottom
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(windowConfig.Width, (windowConfig.Height - (activeSpaceRectangle.Y + activeSpaceRectangle.Height))))
	s.Rectangle(0)
	s.Draw(win)

	// left
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(0+activeSpaceRectangle.X, windowConfig.Height))
	s.Rectangle(0)
	s.Draw(win)

	// right
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(activeSpaceRectangle.X+activeSpaceRectangle.Width, activeSpaceRectangle.Y))
	s.Push(pixel.V(windowConfig.Width, windowConfig.Height))
	s.Rectangle(0)
	s.Draw(win)
}
