package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

func GameStateOver(
	ui UISystem,
	currentState *State,
) error {
	ui.DrawGameOverScreen()

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		*currentState = StateStart
	}

	return nil
}
