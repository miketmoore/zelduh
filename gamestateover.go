package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func GameStateOver(win *pixelgl.Window, txt *text.Text, currLocaleMsgs map[string]string, gameModel *GameModel) {
	win.Clear(colornames.Darkgray)
	DrawMapBackground(win, MapX, MapY, MapW, MapH, colornames.Black)
	DrawCenterText(win, txt, currLocaleMsgs["gameOverScreenMessage"], colornames.White)

	if win.JustPressed(pixelgl.KeyEnter) {
		gameModel.CurrentState = StateStart
	}
}
