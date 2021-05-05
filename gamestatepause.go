package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

func (g *GameStateManager) statePause() error {

	g.UI.DrawPauseScreen()

	if g.UI.Window.JustPressed(pixelgl.KeyP) {
		g.setCurrentState(StateGame)
	}
	if g.UI.Window.JustPressed(pixelgl.KeyEscape) {
		g.setCurrentState(StateStart)
	}
	return nil
}
