package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

func GameStatePause(win *pixelgl.Window, txt *text.Text, currLocaleMsgs map[string]string, gameModel *GameModel) {
	win.Clear(colornames.Darkgray)
	DrawMapBackground(win, MapX, MapY, MapW, MapH, colornames.White)
	DrawCenterText(win, txt, currLocaleMsgs["pauseScreenMessage"], colornames.Black)

	if win.JustPressed(pixelgl.KeyP) {
		gameModel.CurrentState = StateGame
	}
	if win.JustPressed(pixelgl.KeyEscape) {
		gameModel.CurrentState = StateStart
	}
}
