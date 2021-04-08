package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

// GameStateStart handles functionality for the game "start" state
func GameStateStart(win *pixelgl.Window, txt *text.Text, currLocaleMsgs map[string]string, gameModel *GameModel) {
	DrawScreenStart(win, txt, currLocaleMsgs)

	if win.JustPressed(pixelgl.KeyEnter) {
		gameModel.CurrentState = StateGame
	}
}
