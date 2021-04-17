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

func NewUI(currLocaleMsgs LocaleMessagesMap) UI {

	// Initialize text
	orig := pixel.V(20, 50)
	txt := text.New(orig, text.Atlas7x13)
	txt.Color = colornames.Black

	// Initialize window
	win, err := pixelgl.NewWindow(
		pixelgl.WindowConfig{
			Title:  currLocaleMsgs["gameTitle"],
			Bounds: pixel.R(WinX, WinY, WinW, WinH),
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

func DrawMapBackground(win *pixelgl.Window, x, y, w, h float64, color color.Color) {
	s := imdraw.New(nil)
	s.Color = color
	s.Push(pixel.V(x, y))
	s.Push(pixel.V(x+w, y+h))
	s.Rectangle(0)
	s.Draw(win)
}

func DrawScreenStart(win *pixelgl.Window, txt *text.Text, currLocaleMsgs LocaleMessagesMap) {
	win.Clear(colornames.Darkgray)
	DrawMapBackground(win, MapX, MapY, MapW, MapH, colornames.White)
	DrawCenterText(win, txt, currLocaleMsgs["gameTitle"], colornames.Black)
}

func DrawMapBackgroundImage(
	win *pixelgl.Window,
	spritesheet map[int]*pixel.Sprite,
	mapDrawData MapDrawData,
	name MapName,
	modX, modY float64,
	tileSize float64,
) {

	d := mapDrawData[name]
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			sprite := spritesheet[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+MapX+modX+tileSize/2,
				vec.Y+MapY+modY+tileSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func AddUICoin(systemsManager *SystemsManager, entityConfigPresetFnManager *EntityConfigPresetFnManager, frameRate int) {
	presetFn := entityConfigPresetFnManager.GetPreset("uiCoin")
	entityConfig := presetFn(4, 14)
	coin := BuildEntityFromConfig(entityConfig, systemsManager.NewEntityID(), frameRate)
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
	entityConfigPresetFnManager *EntityConfigPresetFnManager,
	roomsMap Rooms,
	mapDrawData MapDrawData,
	roomID *RoomID,
	modX, modY float64,
	tileSize float64,
	frameRate int,
) []Entity {
	d := mapDrawData[roomsMap[*roomID].MapName()]
	obstacles := []Entity{}
	mod := 0.5
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+MapX+modX+tileSize/2,
				vec.Y+MapY+modY+tileSize/2,
			)

			if _, ok := NonObstacleSprites[spriteData.SpriteID]; !ok {
				x := movedVec.X/tileSize - mod
				y := movedVec.Y/tileSize - mod
				id := systemsManager.NewEntityID()
				obstacle := BuildEntityFromConfig(entityConfigPresetFnManager.GetPreset("obstacle")(x, y), id, frameRate)
				obstacles = append(obstacles, obstacle)
			}
		}
	}
	return obstacles
}

func DrawMask(win *pixelgl.Window) {
	// top
	s := imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, MapY+MapH))
	s.Push(pixel.V(WinW, MapY+MapH+(WinH-(MapY+MapH))))
	s.Rectangle(0)
	s.Draw(win)

	// bottom
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(WinW, (WinH - (MapY + MapH))))
	s.Rectangle(0)
	s.Draw(win)

	// left
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(0, 0))
	s.Push(pixel.V(0+MapX, WinH))
	s.Rectangle(0)
	s.Draw(win)

	// right
	s = imdraw.New(nil)
	s.Color = colornames.White
	s.Push(pixel.V(MapX+MapW, MapY))
	s.Push(pixel.V(WinW, WinH))
	s.Rectangle(0)
	s.Draw(win)
}
