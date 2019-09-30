package plancom

import (
	"math/rand"
	"time"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
	//assets "github.com/atolVerderben/planetcommand/plancom/internal"
)

//Game is the tentsuyu.Game which contains all the relevant library attributes
var Game *tentsuyu.Game

//Kanit Font
//https://fonts.google.com/specimen/Kanit?selection.family=Kanit
//https://github.com/cadsondemak/kanit

var (
	//ScreenWidth of the game
	ScreenWidth float64
	//ScreenHeight of the game
	ScreenHeight float64
	//FntMain is Kanit font
	FntMain = "FntKanit"
	//AngleDrawOffset is 90 degrees (All images in spritesheet face UP... Right is 0 degrees/radians)
	AngleDrawOffset = 1.5708
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

//NewGame returns a new Game while setting the width and height of the screen
func NewGame(w, h float64) (game *tentsuyu.Game, err error) {
	ScreenWidth = w
	ScreenHeight = h
	game, err = tentsuyu.NewGame(w, h)
	game.Input.RegisterButton("ToggleZoom", ebiten.KeyTab)

	game.Input.RegisterButton("RotateLeft", ebiten.KeyLeft, ebiten.KeyA)
	game.Input.RegisterButton("RotateRight", ebiten.KeyRight, ebiten.KeyD)

	game.Input.RegisterButton("Left", ebiten.KeyLeft, ebiten.KeyA)
	game.Input.RegisterButton("Right", ebiten.KeyRight, ebiten.KeyD)
	game.Input.RegisterButton("Enter", ebiten.KeyEnter)

	game.Input.RegisterButton("ZoomIn", ebiten.KeyQ)
	game.Input.RegisterButton("ZoomOut", ebiten.KeyE)

	game.Input.RegisterButton("ToggleMute", ebiten.KeyM)

	game.UIController.AddFontFile(FntMain, "assets/font/Kanit-Bold.ttf")

	game.LoadImages(func() *tentsuyu.ImageManager {
		return loadImages()
	})
	game.LoadAudio(func() *tentsuyu.AudioPlayer {
		return loadAudio()
	})

	//game.SetGameState(NewGameMain(game))
	game.SetGameStateLoop(func() error {
		switch game.GetGameState().Msg() {
		case GameStateMsgMain:
			game.SetGameState(NewGameMain(game))
		case GameStateMsgTitle:
			game.SetGameState(NewTitleScreen(game))
		case tentsuyu.GameStateMsgNotStarted:
			game.SetGameState(NewTitleScreen(game))
		default:

		}
		return nil
	})
	game.GameData.Settings["Highscore"] = &tentsuyu.GameValuePair{ValueType: tentsuyu.GameValueInt, ValueInt: 0}
	Game = game
	return
}

//GameState Messages used for this game
var (
	GameStateMsgMain  tentsuyu.GameStateMsg = "MainGame"
	GameStateMsgTitle tentsuyu.GameStateMsg = "Game Title Screen"
)
