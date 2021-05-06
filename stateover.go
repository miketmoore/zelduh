package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

type StateGameOver struct {
	context  *StateContext
	uiSystem *UISystem
}

func NewStateGameOver(context *StateContext, uiSystem *UISystem) State {
	return StateStart{
		context:  context,
		uiSystem: uiSystem,
	}
}

func (g StateGameOver) Update() error {
	g.uiSystem.DrawGameOverScreen()

	if g.uiSystem.Window.JustPressed(pixelgl.KeyEnter) {
		err := g.context.SetState(StateNameStart)
		if err != nil {
			return err
		}
	}

	return nil
}
