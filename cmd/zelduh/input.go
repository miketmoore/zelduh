package main

import "github.com/faiface/pixel/pixelgl"

type Input struct {
	window *pixelgl.Window
}

func (input Input) Up() bool              { return input.window.Pressed(pixelgl.KeyUp) }
func (input Input) Right() bool           { return input.window.Pressed(pixelgl.KeyRight) }
func (input Input) Down() bool            { return input.window.Pressed(pixelgl.KeyDown) }
func (input Input) Left() bool            { return input.window.Pressed(pixelgl.KeyLeft) }
func (input Input) PrimaryAttack() bool   { return input.window.Pressed(pixelgl.KeyF) }
func (input Input) SecondaryAttack() bool { return input.window.Pressed(pixelgl.KeyG) }
func (input Input) Combo() bool           { return input.window.Pressed(pixelgl.KeySpace) }
