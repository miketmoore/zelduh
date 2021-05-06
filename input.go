package zelduh

import "github.com/faiface/pixel/pixelgl"

type InputImpl struct {
	window *pixelgl.Window
}

func (input InputImpl) Up() bool              { return input.window.Pressed(pixelgl.KeyUp) }
func (input InputImpl) Right() bool           { return input.window.Pressed(pixelgl.KeyRight) }
func (input InputImpl) Down() bool            { return input.window.Pressed(pixelgl.KeyDown) }
func (input InputImpl) Left() bool            { return input.window.Pressed(pixelgl.KeyLeft) }
func (input InputImpl) PrimaryAttack() bool   { return input.window.Pressed(pixelgl.KeyF) }
func (input InputImpl) SecondaryAttack() bool { return input.window.Pressed(pixelgl.KeyG) }
func (input InputImpl) Combo() bool           { return input.window.Pressed(pixelgl.KeySpace) }
