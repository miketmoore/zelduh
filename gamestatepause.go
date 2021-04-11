package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStatePause(ui UI, currLocaleMsgs LocaleMessagesMap, gameStateManager *GameStateManager, mapConfig MapConfig) {
	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, mapConfig, colornames.White)
	DrawCenterText(ui.Window, ui.Text, currLocaleMsgs["pauseScreenMessage"], colornames.Black)

	if ui.Window.JustPressed(pixelgl.KeyP) {
		gameStateManager.CurrentState = StateGame
	}
	if ui.Window.JustPressed(pixelgl.KeyEscape) {
		gameStateManager.CurrentState = StateStart
	}
}
