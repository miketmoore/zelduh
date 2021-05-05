package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

// GameStateStart handles functionality for the game "start" state
func (g *GameStateManager) stateStart() error {
	g.UI.DrawScreenStart()

	if g.UI.Window.JustPressed(pixelgl.KeyEnter) {
		g.setCurrentState(StateGame)
	}

	return nil
}
