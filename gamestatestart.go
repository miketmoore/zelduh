package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

// GameStateStart handles functionality for the game "start" state
func GameStateStart(ui UISystem, currentState *State) error {
	ui.DrawScreenStart()

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		*currentState = StateGame
	}

	return nil
}
