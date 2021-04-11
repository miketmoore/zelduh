package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateOver(ui UI, currLocaleMsgs LocaleMessagesMap, gameStateManager *GameStateManager, mapConfig MapConfig) {
	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, mapConfig, colornames.Black)
	DrawCenterText(ui.Window, ui.Text, currLocaleMsgs["gameOverScreenMessage"], colornames.White)

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		gameStateManager.CurrentState = StateStart
	}
}
