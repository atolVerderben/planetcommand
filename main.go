package main

import (
	"github.com/atolVerderben/planetcommand/plancom"

	"github.com/hajimehoshi/ebiten"
)

//Size of the game
const (
	ScreenWidth  = 1280
	ScreenHeight = 720
)

func main() {
	game, err := plancom.NewGame(ScreenWidth, ScreenHeight)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(game.Loop, ScreenWidth, ScreenHeight, 1, "Planet Command"); err != nil {
		panic(err)
	}
}
