package zelduh

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

type StatePause struct {
	context  *StateContext
	uiSystem *UISystem
}

func NewStatePause(context *StateContext, uiSystem *UISystem) State {
	return StatePause{
		context:  context,
		uiSystem: uiSystem,
	}
}

func (g StatePause) Update() error {

	g.uiSystem.DrawPauseScreen()

	if g.uiSystem.Window.JustPressed(pixelgl.KeyP) {
		fmt.Println("state: pause => game")
		err := g.context.SetState(StateNameGame)
		if err != nil {
			return err
		}
	}
	if g.uiSystem.Window.JustPressed(pixelgl.KeyEscape) {
		err := g.context.SetState(StateNameStart)
		if err != nil {
			return err
		}
	}

	return nil
}
