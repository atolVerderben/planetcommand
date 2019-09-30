package plancom

import (
	"image/color"
	"log"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

func loadImages() *tentsuyu.ImageManager {
	imageManager := tentsuyu.NewImageManager()

	whiteImg, err := ebiten.NewImage(1, 1, ebiten.FilterNearest)
	if err != nil {
		log.Fatal(err)
	}
	whiteImg.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	imageManager.AddImage("pixel", whiteImg)

	imageManager.LoadImageFromFile("planet", "assets/planet.png")
	imageManager.LoadImageFromFile("cannon", "assets/cannon.png")
	imageManager.LoadImageFromFile("side-panel", "assets/side-panel.png")
	imageManager.LoadImageFromFile("explosion", "assets/explosion.png")

	return imageManager
}
