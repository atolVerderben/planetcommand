package main

import (
	"log"

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
	/*
		//Now uses ebiten's Game interface to run, this is for legacy ebiten < 1.11
		if err := ebiten.Run(game.Loop, ScreenWidth, ScreenHeight, 1, "Planet Command"); err != nil {
			panic(err)
		}
	*/
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Planet Command")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
