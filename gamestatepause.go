package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStatePause(
	ui UI,
	currLocaleMsgs LocaleMessagesMap,
	currentState *State,
	mapConfig MapConfig,
) {
	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, mapConfig, colornames.White)
	DrawCenterText(ui.Window, ui.Text, currLocaleMsgs["pauseScreenMessage"], colornames.Black)

	if ui.Window.JustPressed(pixelgl.KeyP) {
		*currentState = StateGame
	}
	if ui.Window.JustPressed(pixelgl.KeyEscape) {
		*currentState = StateStart
	}
}
