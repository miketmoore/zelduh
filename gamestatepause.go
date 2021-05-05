package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

func (g *GameStateManager) statePause() error {

	g.UI.DrawPauseScreen()

	if g.UI.Window.JustPressed(pixelgl.KeyP) {
		*g.CurrentState = StateGame
	}
	if g.UI.Window.JustPressed(pixelgl.KeyEscape) {
		*g.CurrentState = StateStart
	}
	return nil
}
