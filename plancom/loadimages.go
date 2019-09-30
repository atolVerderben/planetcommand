package plancom

import (
	"image/color"
	"log"

	assets "github.com/atolVerderben/planetcommand/plancom/internal"
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

	imageManager.AddImageFromBytes("planet", assets.PLANET_PNG)
	imageManager.AddImageFromBytes("cannon", assets.CANNON_PNG)
	imageManager.AddImageFromBytes("side-panel", assets.SIDE_PANEL_PNG)
	imageManager.AddImageFromBytes("explosion", assets.EXPLOSION_PNG)
	imageManager.AddImageFromBytes("title", assets.GAMETITLESCREEN_PNG)

	/*imageManager.LoadImageFromFile("planet", "assets/planet.png")
	imageManager.LoadImageFromFile("cannon", "assets/cannon.png")
	imageManager.LoadImageFromFile("side-panel", "assets/side-panel.png")
	imageManager.LoadImageFromFile("explosion", "assets/explosion.png")
	imageManager.LoadImageFromFile("title", "assets/GameTitleScreen.png")*/

	return imageManager
}
