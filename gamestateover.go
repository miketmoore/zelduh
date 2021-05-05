package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateOver(
	ui UISystem,
	currLocaleMsgs LocaleMessagesMap,
	currentState *State,
) error {
	ui.Window.Clear(colornames.Darkgray)
	ui.DrawMapBackground(colornames.Black)
	ui.DrawCenterText(currLocaleMsgs["gameOverScreenMessage"], colornames.White)

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		*currentState = StateStart
	}

	return nil
}
