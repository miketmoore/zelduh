package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

func (g *GameStateManager) stateOver() error {
	g.UI.DrawGameOverScreen()

	if g.UI.Window.JustPressed(pixelgl.KeyEnter) {
		g.setCurrentState(StateStart)
	}

	return nil
}
