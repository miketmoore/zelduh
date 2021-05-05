package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStatePause(
	ui UISystem,
	currLocaleMsgs LocaleMessagesMap,
	currentState *State,
) error {
	ui.Window.Clear(colornames.Darkgray)
	ui.DrawMapBackground(colornames.White)
	ui.DrawCenterText(currLocaleMsgs["pauseScreenMessage"], colornames.Black)

	if ui.Window.JustPressed(pixelgl.KeyP) {
		*currentState = StateGame
	}
	if ui.Window.JustPressed(pixelgl.KeyEscape) {
		*currentState = StateStart
	}
	return nil
}
