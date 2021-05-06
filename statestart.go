package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

type StateStart struct {
	context  *StateContext
	uiSystem *UISystem
}

func NewStateStart(context *StateContext, uiSystem *UISystem) State {
	return StateStart{
		context:  context,
		uiSystem: uiSystem,
	}
}

func (g StateStart) Update() error {
	g.uiSystem.DrawScreenStart()

	if g.uiSystem.Window.JustPressed(pixelgl.KeyEnter) {
		err := g.context.SetState(StateNameGame)
		if err != nil {
			return err
		}
	}

	return nil
}
