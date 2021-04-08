package zelduh

import (
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

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

func DrawScreenStart(win *pixelgl.Window, txt *text.Text, currLocaleMsgs map[string]string) {
	win.Clear(colornames.Darkgray)
	DrawMapBackground(win, MapX, MapY, MapW, MapH, colornames.White)
	DrawCenterText(win, txt, currLocaleMsgs["gameTitle"], colornames.Black)
}

func DrawMapBackgroundImage(
	win *pixelgl.Window,
	spritesheet map[int]*pixel.Sprite,
	allMapDrawData map[string]MapData,
	name string,
	modX, modY float64) {

	d := allMapDrawData[name]
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			sprite := spritesheet[spriteData.SpriteID]

			vec := spriteData.Rect.Min

			movedVec := pixel.V(
				vec.X+MapX+modX+TileSize/2,
				vec.Y+MapY+modY+TileSize/2,
			)
			matrix := pixel.IM.Moved(movedVec)
			sprite.Draw(win, matrix)
		}
	}
}

func AddUICoin(gameWorld *World) {
	coin := BuildEntityFromConfig(GetPreset("uiCoin")(4, 14), gameWorld.NewEntityID())
	gameWorld.AddEntity(coin)
}

// make sure only correct number of hearts exists in systems
// so, if health is reduced, need to remove a heart entity from the systems,
// the correct one... last one
func AddUIHearts(gameWorld *World, hearts []Entity, health int) {
	for i, entity := range hearts {
		if i < health {
			gameWorld.AddEntity(entity)
		}
	}
}

func DrawObstaclesPerMapTiles(gameWorld *World, roomsMap Rooms, allMapDrawData map[string]MapData, roomID RoomID, modX, modY float64) []Entity {
	d := allMapDrawData[roomsMap[roomID].MapName()]
	obstacles := []Entity{}
	mod := 0.5
	for _, spriteData := range d.Data {
		if spriteData.SpriteID != 0 {
			vec := spriteData.Rect.Min
			movedVec := pixel.V(
				vec.X+MapX+modX+TileSize/2,
				vec.Y+MapY+modY+TileSize/2,
			)

			if _, ok := NonObstacleSprites[spriteData.SpriteID]; !ok {
				x := movedVec.X/TileSize - mod
				y := movedVec.Y/TileSize - mod
				id := gameWorld.NewEntityID()
				obstacle := BuildEntityFromConfig(GetPreset("obstacle")(x, y), id)
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
