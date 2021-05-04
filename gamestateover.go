package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateOver(
	ui UISystem,
	currLocaleMsgs LocaleMessagesMap,
	currentState *State,
	activeSpaceRectangle ActiveSpaceRectangle,
) error {
	ui.Window.Clear(colornames.Darkgray)
	ui.DrawMapBackground(colornames.Black)
	DrawCenterText(ui.Window, ui.Text, currLocaleMsgs["gameOverScreenMessage"], colornames.White)

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		*currentState = StateStart
	}

	return nil
}
