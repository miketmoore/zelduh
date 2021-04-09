package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStatePause(ui UI, currLocaleMsgs LocaleMessagesMap, gameModel *GameModel) {
	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.White)
	DrawCenterText(ui.Window, ui.Text, currLocaleMsgs["pauseScreenMessage"], colornames.Black)

	if ui.Window.JustPressed(pixelgl.KeyP) {
		gameModel.CurrentState = StateGame
	}
	if ui.Window.JustPressed(pixelgl.KeyEscape) {
		gameModel.CurrentState = StateStart
	}
}
