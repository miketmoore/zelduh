package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

// GameStateStart handles functionality for the game "start" state
func GameStateStart(ui UI, currLocaleMsgs LocaleMessagesMap, gameStateManager *GameStateManager, mapConfig MapConfig) {
	DrawScreenStart(ui.Window, ui.Text, currLocaleMsgs, mapConfig)

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		gameStateManager.CurrentState = StateGame
	}
}
