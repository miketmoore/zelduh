package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

func GameStatePause(
	ui UISystem,
	currentState *State,
) error {

	ui.DrawPauseScreen()

	if ui.Window.JustPressed(pixelgl.KeyP) {
		*currentState = StateGame
	}
	if ui.Window.JustPressed(pixelgl.KeyEscape) {
		*currentState = StateStart
	}
	return nil
}
