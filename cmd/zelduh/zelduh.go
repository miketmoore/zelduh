package main

import (
	"fmt"
	_ "image/png"
	"os"

	"github.com/miketmoore/zelduh"

	"github.com/faiface/pixel/pixelgl"
)

func run() {

	debugMode := false

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) >= 1 && argsWithoutProg[0] == "debug" {
		debugMode = true
	}

	// TileSize defines the width and height of a tile
	const tileSize float64 = 48

	// FrameRate is used to determine which sprite to use for animations
	const frameRate int = 5

	main, err := zelduh.NewMain(
		debugMode,
		tileSize,
		frameRate,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = main.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	os.Exit(1)
}

func main() {
	pixelgl.Run(run)
}
