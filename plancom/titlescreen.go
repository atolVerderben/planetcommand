package plancom

import (
	"image/color"
	"os"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//TitleScreen represents the main GameState of our game
type TitleScreen struct {
	gameStateMsg  tentsuyu.GameStateMsg
	mainCamera    *tentsuyu.Camera
	width, height float64
	message       *tentsuyu.TextElement
}

//NewTitleScreen returns our main gamestate
func NewTitleScreen(game *tentsuyu.Game) *TitleScreen {

	g := &TitleScreen{
		message: tentsuyu.NewTextElement(915, 500, 250, 100, game.UIController.ReturnFont(FntMain), []string{"PRESS ENTER"}, color.RGBA{249, 200, 14, 255}, 26),
	}

	return g
}

//Update the gamestate
func (g *TitleScreen) Update(game *tentsuyu.Game) error {
	if !Game.AudioPlayer.ReturnSongPlayer("BGM").IsPlaying() && !Game.AudioPlayer.IsMusicMuted() {
		Game.AudioPlayer.ReturnSongPlayer("BGM").Rewind()
		Game.AudioPlayer.ReturnSongPlayer("BGM").Play()
		Game.AudioPlayer.ReturnSongPlayer("BGM").SetVolume(0.3)
	}
	if game.Input.Button("Escape").JustPressed() {
		os.Exit(0)
	}
	if game.Input.Button("Enter").JustPressed() {
		g.SetMsg(GameStateMsgMain)
	}
	return nil
}

//Draw the gamestate scene
func (g *TitleScreen) Draw(game *tentsuyu.Game) error {
	op := &ebiten.DrawImageOptions{}

	game.Screen.DrawImage(game.ImageManager.ReturnImage("title"), op)

	g.message.Draw(game.Screen)
	return nil
}

//Msg returns the gamestatemsg and achieves the GameState interface
func (g TitleScreen) Msg() tentsuyu.GameStateMsg {
	return g.gameStateMsg
}

//SetMsg sets the gamestatemsg value
func (g *TitleScreen) SetMsg(gs tentsuyu.GameStateMsg) {
	g.gameStateMsg = gs
}
