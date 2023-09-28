package main

import (
	"pong-inverso-pixel/scenes"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	mainScene := scenes.NewMainScene()
	pixelgl.Run(mainScene.Draw)
}
