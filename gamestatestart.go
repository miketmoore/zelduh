package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

// GameStateStart handles functionality for the game "start" state
func GameStateStart(ui UI, currLocaleMsgs LocaleMessagesMap, gameStateManager *GameStateManager) {
	DrawScreenStart(ui.Window, ui.Text, currLocaleMsgs)

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		gameStateManager.CurrentState = StateGame
	}
}
