package zelduh

import (
	"fmt"

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
		fmt.Println("state: start => game")
		g.context.SetState(GameStateNameGame)
	}

	return nil
}
