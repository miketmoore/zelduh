package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

// GameStateStart handles functionality for the game "start" state
func GameStateStart(ui UI, currLocaleMsgs LocaleMessagesMap, currentState *State, activeSpaceRectangle ActiveSpaceRectangle) {
	DrawScreenStart(ui.Window, ui.Text, currLocaleMsgs, activeSpaceRectangle)

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		*currentState = StateGame
	}
}
