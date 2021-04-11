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

type UI struct {
	Window *pixelgl.Window
	Text   *text.Text
}

func NewUI(currLocaleMsgs LocaleMessagesMap, windowConfig WindowConfig) UI {

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

	return UI{
		Window: win,
		Text:   txt,
	}
}

func DrawCenterText(win *pixelgl.Window, txt *text.Text, s string, c color.RGBA) {
	txt.Clear()
	txt.Color = c
	fmt.Fprintln(txt, s)
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center())))
}

func DrawMapBackground(win *pixelgl.Window, mapConfig MapConfig, color color.Color) {
	s := imdraw.New(nil)
	s.Color = color
	s.Push(pixel.V(mapConfig.X, mapConfig.Y))
	s.Push(pixel.V(mapConfig.X+mapConfig.Width, mapConfig.Y+mapConfig.Height))
	s.Rectangle(0)
	s.Draw(win)
}

func DrawScreenStart(win *pixelgl.Window, txt *text.Text, currLocaleMsgs LocaleMessagesMap, mapConfig MapConfig) {
	win.Clear(colornames.Darkgray)
	DrawMapBackground(win, mapConfig, colornames.White)
	DrawCenterText(win, txt, currLocaleMsgs["gameTitle"], colornames.Black)
}

func DrawMapBackgroundImage(
	win *pixelgl.Window,
	spritesheet map[int]*pixel.Sprite,
	allMapDrawData map[string]MapData,
	name string,
	modX, modY float64,
	mapConfig MapConfig,
) {

	d := allMapDrawData[name]
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			sprite := spritesheet[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+mapConfig.X+modX+TileSize/2,
				vec.Y+mapConfig.Y+modY+TileSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func AddUICoin(systemsManager *SystemsManager) {
	coin := BuildEntityFromConfig(GetPreset("uiCoin")(4, 14), systemsManager.NewEntityID())
	systemsManager.AddEntity(coin)
}

// make sure only correct number of hearts exists in systems
// so, if health is reduced, need to remove a heart entity from the systems,
// the correct one... last one
func AddUIHearts(systemsManager *SystemsManager, hearts []Entity, health int) {
	for i, entity := range hearts {
		if i < health {
			systemsManager.AddEntity(entity)
		}
	}
}

func DrawObstaclesPerMapTiles(
	systemsManager *SystemsManager,
	roomsMap Rooms,
	allMapDrawData map[string]MapData,
	roomID RoomID,
	modX,
	modY float64,
	mapConfig MapConfig,
) []Entity {
	d := allMapDrawData[roomsMap[roomID].MapName()]
	obstacles := []Entity{}
	mod := 0.5
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+mapConfig.X+modX+TileSize/2,
				vec.Y+mapConfig.Y+modY+TileSize/2,
			)

			if _, ok := NonObstacleSprites[spriteData.SpriteID]; !ok {
				x := movedVec.X/TileSize - mod
				y := movedVec.Y/TileSize - mod
				id := systemsManager.NewEntityID()
				obstacle := BuildEntityFromConfig(GetPreset("obstacle")(x, y), id)
				obstacles = append(obstacles, obstacle)
			}
		}
	}
	return obstacles
}

func DrawMask(win *pixelgl.Window, windowConfig WindowConfig, mapConfig MapConfig) {
	// top
	s := imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, mapConfig.Y+mapConfig.Height))
	s.Push(pixel.V(windowConfig.Width, mapConfig.Y+mapConfig.Height+(windowConfig.Height-(mapConfig.Y+mapConfig.Height))))
	s.Rectangle(0)
	s.Draw(win)

	// bottom
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(windowConfig.Width, (windowConfig.Height - (mapConfig.Y + mapConfig.Height))))
	s.Rectangle(0)
	s.Draw(win)

	// left
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(0+mapConfig.X, windowConfig.Height))
	s.Rectangle(0)
	s.Draw(win)

	// right
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(mapConfig.X+mapConfig.Width, mapConfig.Y))
	s.Push(pixel.V(windowConfig.Width, windowConfig.Height))
	s.Rectangle(0)
	s.Draw(win)
}
