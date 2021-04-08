package zelduh

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func GameStateOver(ui UI, currLocaleMsgs map[string]string, gameModel *GameModel) {
	ui.Window.Clear(colornames.Darkgray)
	DrawMapBackground(ui.Window, MapX, MapY, MapW, MapH, colornames.Black)
	DrawCenterText(ui.Window, ui.Text, currLocaleMsgs["gameOverScreenMessage"], colornames.White)

	if ui.Window.JustPressed(pixelgl.KeyEnter) {
		gameModel.CurrentState = StateStart
	}
}
