package zelduh

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

type GameStatePause struct {
	context  *GameStateContext
	uiSystem *UISystem
}

func NewGameStatePause(context *GameStateContext, uiSystem *UISystem) GameState {
	return GameStatePause{
		context:  context,
		uiSystem: uiSystem,
	}
}

func (g GameStatePause) Update() error {

	g.uiSystem.DrawPauseScreen()

	if g.uiSystem.Window.JustPressed(pixelgl.KeyP) {
		fmt.Println("state: pause => game")
		g.context.SetState(GameStateNameGame)
	}
	if g.uiSystem.Window.JustPressed(pixelgl.KeyEscape) {
		g.context.SetState(GameStateNameStart)
	}

	return nil
}
