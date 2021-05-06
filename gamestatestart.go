package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
)

type GameStateStart struct {
	context  *GameStateContext
	uiSystem *UISystem
}

func NewGameStateStart(context *GameStateContext, uiSystem *UISystem) GameState {
	return GameStateStart{
		context:  context,
		uiSystem: uiSystem,
	}
}

func (g GameStateStart) Update() error {
	g.uiSystem.DrawScreenStart()

	if g.uiSystem.Window.JustPressed(pixelgl.KeyEnter) {
		err := g.context.SetState(GameStateNameGame)
		if err != nil {
			return err
		}
	}

	return nil
}
