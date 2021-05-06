package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

type GameStateGameOver struct {
	context  *GameStateContext
	uiSystem *UISystem
}

func NewGameStateGameOver(context *GameStateContext, uiSystem *UISystem) GameState {
	return GameStateStart{
		context:  context,
		uiSystem: uiSystem,
	}
}

func (g GameStateGameOver) Update() error {
	g.uiSystem.DrawGameOverScreen()

	if g.uiSystem.Window.JustPressed(pixelgl.KeyEnter) {
		g.context.SetState(GameStateNameStart)
	}

	return nil
}
