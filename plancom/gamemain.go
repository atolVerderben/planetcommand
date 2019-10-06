package plancom

import (
	"image/color"
	"math/rand"
	"strconv"

	"github.com/atolVerderben/planetcommand/scoreboard"

	"github.com/golang/freetype/truetype"

	"github.com/atolVerderben/tentsuyu/tentsuyutils"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

//GameMain represents the main GameState of our game
type GameMain struct {
	gameStateMsg               tentsuyu.GameStateMsg
	mainCamera                 *tentsuyu.Camera
	width, height              float64
	planet                     *Planet
	cannons                    []*Cannon
	shadowImage, triangleImage *ebiten.Image
	explosions                 []*Explosion
	currentLevel               int

	mainGameScreen *ebiten.Image
	mainGameArea   *tentsuyu.Rectangle

	paused bool

	playerDead                  bool
	deathDelayCount, deathDelay int

	scoreDisplay, highScoreDisplay *tentsuyu.MenuElement
	score, highScore               int

	textDisplays []tentsuyu.UIElement

	overlayText            *tentsuyu.TextElement
	overlayTimer           int
	overlayTimeEnd         int
	overlayClickToContinue bool
	overlayX, overlayY     float64
	overlayDelay           int

	restartMenu *tentsuyu.Menu

	cannonDisp *cannonDisplay

	aiManager *AIController

	leaderboard   *scoreboard.LeaderBoard
	scoreEntry    *ScoreEntry
	enteringScore bool
	scoreSaveSpot int
	scoreSaveName string
	highscores    *highscoreDisplay

	ufos                          []*AIUfo
	ufoSpawnRate, ufoSpawnCounter int
}

//NewGameMain returns our main gamestate
func NewGameMain(game *tentsuyu.Game) *GameMain {
	textColor := color.RGBA{247, 6, 207, 255}
	g := &GameMain{
		mainCamera: tentsuyu.CreateCamera(720, 720),
		width:      720,
		height:     720,
		cannons:    []*Cannon{},
		mainGameArea: &tentsuyu.Rectangle{
			X: 280,
			Y: 0,
			W: 720,
			H: 720,
		},
		explosions:    []*Explosion{},
		aiManager:     CreateAIController(),
		score:         0,
		textDisplays:  []tentsuyu.UIElement{},
		currentLevel:  1,
		overlayX:      300,
		overlayY:      40,
		overlayText:   tentsuyu.NewTextElementStationary(720, 40, 720, 720, game.UIController.ReturnFont(FntMain), []string{}, color.RGBA{249, 200, 14, 255}, 24),
		deathDelay:    180,
		scoreEntry:    NewScoreEntry(0, 80),
		scoreSaveSpot: -1,
		ufos:          []*AIUfo{},
		ufoSpawnRate:  360,
	}

	g.scoreDisplay = &tentsuyu.MenuElement{
		UIElement:  tentsuyu.NewUINumberDisplayInt(&g.score, 1015, 35, 100, 100, game.UIController.ReturnFont(FntMain), 24, textColor),
		Action:     nil,
		Selectable: false,
	}

	g.textDisplays = append(g.textDisplays,
		tentsuyu.NewTextElement(1012, 10, 1000, 75, game.UIController.ReturnFont(FntMain), []string{"SCORE:"}, textColor, 24),
		tentsuyu.NewTextElement(1012, 60, 1000, 75, game.UIController.ReturnFont(FntMain), []string{"HIGH SCORE:"}, textColor, 24),
	)

	g.mainCamera.SetBounds(0, 1600, 0, 1600)
	g.planet = CreatePlanet(g.width/2, g.height/2, 64, 64)
	g.cannons = append(g.cannons,
		CreateCannon(g.planet, 0, Cannon1),
		CreateCannon(g.planet, 1.5708, Cannon2),  //90 Degrees
		CreateCannon(g.planet, 3.14159, Cannon3), // 180
		CreateCannon(g.planet, 4.71239, Cannon4)) // 270

	g.restartMenu = tentsuyu.NewMenu(720, 720)
	g.restartMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 40, game.UIController.ReturnFont(FntMain), []string{"RESUME"}, textColor, 32)},
		[]func(){func() { g.paused = false }})
	g.restartMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 40, game.UIController.ReturnFont(FntMain), []string{"RESTART"}, textColor, 32)},
		[]func(){func() { g.SetMsg(GameStateMsgMain) }})
	g.restartMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 40, game.UIController.ReturnFont(FntMain), []string{"QUIT GAME"}, textColor, 32)},
		[]func(){func() {
			g.SetMsg(GameStateMsgTitle) //os.Exit(0)
		}})

	for _, m := range g.restartMenu.Elements {
		m[0].SetHighlightColor(color.RGBA{R: 249, G: 200, B: 14, A: 255})
	}
	g.shadowImage, _ = ebiten.NewImage(int(ScreenWidth), int(ScreenHeight), ebiten.FilterDefault)
	g.triangleImage, _ = ebiten.NewImage(int(ScreenWidth), int(ScreenHeight), ebiten.FilterDefault)
	g.mainGameScreen, _ = ebiten.NewImage(int(g.width), int(g.height), ebiten.FilterDefault)
	g.triangleImage.Fill(color.RGBA{212, 0, 120, 255})
	/*
		In the year 20XX, after most of the world has been destroyed by mysterious little green missiles from space, only 4 cities remain able to defend the planet!

		Command the cannons to destroy the lasers and SAVE THE WORLD!
	*/
	game.AdditionalCameras["MainCamera"] = g.mainCamera
	g.SetOverlay([]string{"In the year 20XX, after most of the world has been",
		"destroyed by mysterious little green missiles from space,",
		"only 4 cities remain able to defend the planet!",
		"Command the cannons to destroy the lasers and...", "", "                             SAVE THE WORLD!",
		"", "", "", "", "", "", "", "", "", "", "", "",
		"                                   CONTROLS:",
		"      LEFT AND RIGHT ARROWS: ROTATE CANNONS",
		"    CLICK MOUSE WITHIN THE CONES TO FIRE MISSILES",
		"", "", "                                CLICK TO BEGIN"}, 300, 0, 0, true)

	g.cannonDisp = newcannonDisplay(g.cannons, game.UIController.ReturnFont(FntMain))

	//Load High Scores
	g.leaderboard = scoreboard.LoadHighScores()
	game.GameData.Settings["Highscore"].ValueInt = g.leaderboard.Records[0].Score // Set the highscore to highest in leaderboard
	g.highscores = newHighScoreDisplay(g.leaderboard, game.UIController.ReturnFont(FntMain))
	g.highScoreDisplay = &tentsuyu.MenuElement{
		UIElement:  tentsuyu.NewUINumberDisplayInt(&game.GameData.Settings["Highscore"].ValueInt, 1015, 85, 100, 100, game.UIController.ReturnFont(FntMain), 24, textColor),
		Action:     nil,
		Selectable: false,
	}

	g.ufos = append(g.ufos, SpawnUFO(400, 150))

	/*for _, r := range g.leaderboard.Records {
		fmt.Printf("%v : %v\n", r.Name, r.Score)
	}*/

	return g
}

//StepZoom increments set zoom levels. Helps keep the game looking good at chosen levels
func (g *GameMain) StepZoom(camera *tentsuyu.Camera, zoomOut bool) {
	switch camera.Zoom {
	case 4:
		if zoomOut {
			camera.SetZoom(2)
		}
	case 2:
		if zoomOut {
			camera.SetZoom(1)
		} else {
			camera.SetZoom(4)
		}
	case 1:
		if !zoomOut {
			camera.SetZoom(2)
		}
	default:
		if zoomOut {
			camera.ZoomOut()
		} else {
			camera.ZoomIn()
		}
	}
}

//Update the gamestate
func (g *GameMain) Update(game *tentsuyu.Game) error {

	if !Game.AudioPlayer.ReturnSongPlayer("BGM").IsPlaying() && !Game.AudioPlayer.IsMusicMuted() {
		Game.AudioPlayer.ReturnSongPlayer("BGM").Rewind()
		Game.AudioPlayer.ReturnSongPlayer("BGM").Play()
		Game.AudioPlayer.ReturnSongPlayer("BGM").SetVolume(0.3)
	}

	if g.overlayTimer > 0 {
		g.overlayTimer++
		g.overlayDelay++
		if g.overlayClickToContinue {
			if game.Input.LeftClick().JustPressed() && g.overlayDelay > 30 {
				g.overlayTimer = g.overlayTimeEnd
				g.overlayClickToContinue = false
			}
		}
		if g.overlayTimer > g.overlayTimeEnd && !g.overlayClickToContinue {
			g.overlayTimer = 0
			g.paused = false
		}
	}

	if g.deathDelayCount > 0 {
		g.deathDelayCount++
		if g.deathDelayCount > g.deathDelay {
			g.playerDead = true
			g.scoreSaveSpot = g.leaderboard.DeterminePosition(g.score)
			if g.scoreSaveSpot > -1 {
				g.enteringScore = true
			}

			g.deathDelayCount = 0
			g.restartMenu.Elements[0][0].Hide(true)
		}

	}

	if g.playerDead {
		if g.enteringScore {
			g.scoreEntry.Update(game.Input)
			if g.scoreEntry.Confirmed {
				if g.leaderboard.SetPosition(g.scoreSaveSpot, g.scoreEntry.ToString(), g.score) {
					g.leaderboard.SaveHighScores()
					g.highscores = newHighScoreDisplay(g.leaderboard, game.UIController.ReturnFont(FntMain))
				}
				g.enteringScore = false
			}
		} else {
			g.restartMenu.Update(game.Input, -280, 0)
		}
		return nil
	}

	if game.Input.Button("Escape").JustPressed() {
		if g.overlayTimer <= 0 {
			g.paused = !g.paused
		}
	}

	if game.Input.Button("ToggleMute").JustPressed() {
		Game.AudioPlayer.MuteAll(!Game.AudioPlayer.IsMusicMuted())
		Game.AudioPlayer.ReturnSongPlayer("BGM").Pause()
		Game.AudioPlayer.ReturnSongPlayer("BGM").Rewind()
	}

	if g.paused {
		if g.overlayTimer <= 0 { // Don't draw the main menu if the overlay is active
			g.restartMenu.Update(game.Input, -280, 0)
		}
		return nil
	}

	if game.Input.LeftClick().JustPressed() {
		x, y := game.Input.GetMouseCoords()
		if g.mainGameArea.Contains(x, y) {

			for _, c := range g.cannons {
				fx, fy := game.Input.GetGameMouseCoords(g.mainCamera)
				fx -= 280 //We need to adjust for the sake of the screen
				//log.Printf("X:%v, Y: %v\n", fx, fy)
				if c.InFireRadius(fx, fy) {
					if c.FireMissile(fx, fy) {

					}

				}
			}

		}
	}

	if game.Input.Button("RotateLeft").Down() {
		for _, c := range g.cannons {
			c.AddCommand(CommandLeft)
		}
	}
	if game.Input.Button("RotateRight").Down() {
		for _, c := range g.cannons {
			c.AddCommand(CommandRight)
		}
	}

	//Move the Planet, for a later iteration!
	//g.updatePlanetMovement(game)

	if g.ufoSpawnCounter > 0 {
		g.ufoSpawnCounter++
		if g.ufoSpawnCounter > g.ufoSpawnRate {
			g.ufoSpawnCounter = 0
			g.ufos = append(g.ufos, SpawnUFO(-20, 200))
		}
	}

	for _, c := range g.cannons {
		c.Update(g.planet)
		for _, m := range c.missiles {
			m.Update()
			if tentsuyutils.Distance(m.X, m.Y, m.TargetX, m.TargetY) <= 2 {
				defer c.RemoveMissile(m.GetIDasString())
				//TODO: EXPLOSION
				g.explosions = append(g.explosions, CreateExplosion(m.X, m.Y, true))
				n := rand.Intn(5) + 1
				game.AudioPlayer.PlaySE("explosion" + strconv.Itoa(n))
				continue
			}
			if m.X < 0 || m.X > g.width {
				defer c.RemoveMissile(m.GetIDasString())
				continue
			}
			if m.Y < 0 || m.Y > g.height {
				defer c.RemoveMissile(m.GetIDasString())
				continue
			}
		}
	}
	g.aiManager.Update(g.planet)
	for _, ai := range g.aiManager.launchers {
		//ai.Update(g.planet)
		for _, m := range ai.missiles {
			m.Update()
			for _, c := range g.cannons {
				if !c.IsAlive() {
					continue
				}
				if tentsuyu.Collision(c.BasicObject, m.BasicObject) {
					c.Hit()
					defer ai.RemoveMissile(m.GetIDasString())
					g.mainCamera.StartShaking(false)
					game.AudioPlayer.PlaySE("hit-planet")
					g.explosions = append(g.explosions, CreateExplosion(m.X, m.Y, false))
					continue
				}
			}
			if g.planet.IsAlive() && tentsuyu.Collision(g.planet.BasicObject, m.BasicObject) {
				defer ai.RemoveMissile(m.GetIDasString())
				g.mainCamera.StartShaking(true)
				game.AudioPlayer.PlaySE("hit-planet")
				g.explosions = append(g.explosions, CreateExplosion(m.X, m.Y, false))
				g.planet.Health--
				if g.planet.IsAlive() {
					g.planet.Width = g.planet.Width / 2
					g.planet.Height = g.planet.Height / 2
				} else {
					//Planet just died
					g.explosions = append(g.explosions, CreateExplosion(g.planet.X, g.planet.Y, false))
					g.explosions = append(g.explosions, CreateExplosion(g.planet.X+10, g.planet.Y, false))
					g.explosions = append(g.explosions, CreateExplosion(g.planet.X-10, g.planet.Y, false))
					g.explosions = append(g.explosions, CreateExplosion(g.planet.X, g.planet.Y+10, false))
					g.explosions = append(g.explosions, CreateExplosion(g.planet.X, g.planet.Y-10, false))
					g.deathDelayCount = 1
					for _, c := range g.cannons {
						c.Health = 0
					}
				}
				continue
			}

			for _, e := range g.explosions {
				if tentsuyu.Collision(e.BasicObject, m.BasicObject) {
					defer ai.RemoveMissile(m.GetIDasString())
					if e.givePoints {

						switch distance := tentsuyutils.Distance(g.planet.X, g.planet.Y, m.X, m.Y); {
						case distance < 100:
							g.score += 5
						case distance < 250:
							g.score += 3
						case distance < 400:
							g.score += 2
						default:
							g.score++

						}

					}
					g.explosions = append(g.explosions, CreateExplosion(m.X, m.Y, e.givePoints)) //give points if the previous explosion should give points
					n := rand.Intn(5) + 1
					game.AudioPlayer.PlaySE("explosion" + strconv.Itoa(n))
				}
			}
		}
	}

	for _, u := range g.ufos {
		for _, e := range g.explosions {
			if e.givePoints { //Fired from Player
				if tentsuyu.Collision(e.BasicObject, u.BasicObject) {
					u.Hit()
					if u.IsDead() {
						defer g.RemoveUFO(u.GetIDasString())
						g.explosions = append(g.explosions, CreateExplosion(u.GetX(), u.GetY(), e.givePoints))
						n := rand.Intn(5) + 1
						game.AudioPlayer.PlaySE("explosion" + strconv.Itoa(n))
						g.ufoSpawnCounter = 1
					}
				}
			}
		}
		u.Update(g.planet, g.mainGameArea)
	}

	for _, e := range g.explosions {
		e.Update()
		if e.done {
			defer g.RemoveExplosion(e.GetIDasString())
		}
	}
	g.mainCamera.FollowObjectInBounds(g.planet)
	g.mainCamera.Update()
	if g.score > game.GameData.Settings["Highscore"].ValueInt { //g.highScore {
		g.highScore = g.score
		game.GameData.Settings["Highscore"].ValueInt = g.score
	}
	g.scoreDisplay.Update(Game.Input, 0, 0)
	g.highScoreDisplay.Update(Game.Input, 0, 0)
	for _, u := range g.textDisplays {
		u.Update()
	}
	return nil
}

func (g *GameMain) updateBullets() {

}

func (g *GameMain) updatePlanetMovement(game *tentsuyu.Game) {
	if game.Input.Button("Left").Down() {
		g.planet.X -= 2
	}
	if game.Input.Button("Right").Down() {
		g.planet.X += 2
	}
}

//Draw the gamestate scene
func (g *GameMain) Draw(game *tentsuyu.Game) error {

	g.drawMain(g.mainGameScreen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(280, 0)
	g.mainCamera.ApplyCameraTransform(op, true)
	game.Screen.DrawImage(g.mainGameScreen, op)
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(1000, 0)
	game.Screen.DrawImage(game.ImageManager.ReturnImage("side-panel"), op)
	op = &ebiten.DrawImageOptions{}

	game.Screen.DrawImage(game.ImageManager.ReturnImage("side-panel"), op)
	g.scoreDisplay.Draw(game.Screen)
	g.highScoreDisplay.Draw(game.Screen)
	for _, u := range g.textDisplays {
		u.Draw(game.Screen)
	}

	for _, t := range g.cannonDisp.names {
		t.Draw(game.Screen)
	}

	for _, c := range g.cannonDisp.bars {
		c.draw(game.Screen)
	}

	for _, t := range g.highscores.scores {
		t.Draw(game.Screen)
	}

	if g.overlayTimer > 0 {
		g.overlayText.Draw(game.Screen)
	}

	return nil
}

func (g *GameMain) drawMain(screen *ebiten.Image) error {
	screen.Clear()
	screen.Fill(color.RGBA{R: 13, G: 2, B: 33, A: 255})
	g.shadowImage.Clear()
	g.shadowImage.Fill(color.Black)
	// Subtract ray triangles from shadow
	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.CompositeMode = ebiten.CompositeModeSourceOut

	for _, c := range g.cannons {
		if c.IsAlive() {
			v := rayVertices(float64(c.leftLine.X1), float64(c.leftLine.Y1), c.rightLine.X2, c.rightLine.Y2, c.leftLine.X2, c.leftLine.Y2)
			g.shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, g.triangleImage, opt)
			if c.overheated {
				//if c.cooldown%5 == 0 {
				screen.DrawTriangles(v, []uint16{0, 1, 2}, g.triangleImage, &ebiten.DrawTrianglesOptions{})
				//}
			}
		}
	}

	drawGrid(g.width, g.height, screen)

	for _, ai := range g.aiManager.launchers {

		for _, m := range ai.missiles {
			m.DrawTrail(screen)

		}
	}
	//DRAW UFO
	for _, u := range g.ufos {
		u.Draw(screen)
	}
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.8)
	screen.DrawImage(g.shadowImage, op)
	//AFTER THIS POINT IS DRAWN ABOVE THE "SHADOW" IMAGE ==================================================================

	if g.planet.IsAlive() {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(g.planet.Width)/128, float64(g.planet.Height)/128)
		op.GeoM.Translate(-float64(g.planet.Width/2), -float64(g.planet.Height/2))
		//op.GeoM.Rotate(g.planet.Angle)
		op.GeoM.Translate(g.planet.X, g.planet.Y)
		//g.mainCamera.ApplyCameraTransform(op, true)
		screen.DrawImage(Game.ImageManager.ReturnImage("planet"), op)
	}
	for _, c := range g.cannons {
		c.Draw(screen)
		for _, m := range c.missiles {
			m.Draw(screen)
			m.DrawTrail(screen)
		}
	}

	for _, ai := range g.aiManager.launchers {
		for _, m := range ai.missiles {
			m.Draw(screen)
		}
	}

	for _, e := range g.explosions {
		e.Draw(screen)
	}

	if g.playerDead || g.paused && g.overlayTimer <= 0 {
		if g.enteringScore {
			g.scoreEntry.Draw(screen)
		} else {
			g.restartMenu.Draw(screen)
		}
	}
	return nil
}

//Msg returns the gamestatemsg and achieves the GameState interface
func (g GameMain) Msg() tentsuyu.GameStateMsg {
	return g.gameStateMsg
}

//SetMsg sets the gamestatemsg value
func (g *GameMain) SetMsg(gs tentsuyu.GameStateMsg) {
	g.gameStateMsg = gs
}

func drawGrid(width, height float64, screen *ebiten.Image) {
	x, y := 0.0, 0.0
	//Draw Horizontal
	for i := 0.0; i < height; i += 40 {
		ebitenutil.DrawLine(screen, x, y+i, x+width, y+i, color.RGBA{146, 0, 117, 255})
	}
	//Draw Vertical
	for i := 0.0; i < width; i += 40 {
		ebitenutil.DrawLine(screen, x+i, y, x+i, y+height, color.RGBA{146, 0, 117, 255})
	}

}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 0.5},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 0.5},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 0.5},
	}
}

//RemoveExplosion from the game
func (g *GameMain) RemoveExplosion(exID string) {
	delete := -1
	for index, e := range g.explosions {
		if e.GetIDasString() == exID {
			delete = index
			break
		}
	}
	if delete >= 0 {
		g.explosions = append(g.explosions[:delete], g.explosions[delete+1:]...)
	}
}

//RemoveUFO from the game
func (g *GameMain) RemoveUFO(exID string) {
	delete := -1
	for index, e := range g.ufos {
		if e.GetIDasString() == exID {
			delete = index
			break
		}
	}
	if delete >= 0 {
		g.ufos = append(g.ufos[:delete], g.ufos[delete+1:]...)
	}
}

//SetOverlay sets the text to display, the duration, the offeset of X,Y and if the text can be skipped by clicking
func (g *GameMain) SetOverlay(text []string, displayTime int, offsetX, offsetY float64, clickToContinue bool) {
	g.overlayText.SetText(text)
	g.overlayText.SetX(g.overlayX + offsetX)
	g.overlayText.SetY(g.overlayY + offsetY)
	g.overlayTimer = 1
	g.overlayTimeEnd = displayTime
	g.overlayClickToContinue = clickToContinue
	g.paused = true
	g.overlayDelay = 0
}

type cannonCooldownBar struct {
	x, y             float64
	associatedCannon *Cannon
	colorPixel       *ebiten.Image
}

func newcannonCooldownBar(x, y float64, barcolor color.Color, cannon *Cannon) *cannonCooldownBar {
	c := &cannonCooldownBar{
		x:                x,
		y:                y,
		associatedCannon: cannon,
	}
	c.colorPixel, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	c.colorPixel.Fill(cannon.color)
	return c
}

func (c *cannonCooldownBar) draw(screen *ebiten.Image) error {
	var r, g, b float64
	if c.associatedCannon.overheated {
		r = rand.Float64()
		g = rand.Float64()
		b = rand.Float64()
	}
	for i := 0; i < c.associatedCannon.heat; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2.25, 20)
		op.GeoM.Translate(c.x+float64(i)*2.25, c.y)
		if c.associatedCannon.overheated {
			op.ColorM.Scale(r, g, b, 1)
		}
		if err := screen.DrawImage(c.colorPixel, op); err != nil {
			return err
		}
	}
	return nil
}

type cannonDisplay struct {
	names map[string]*tentsuyu.MenuElement
	bars  map[string]*cannonCooldownBar
}

func newcannonDisplay(cannons []*Cannon, font *truetype.Font) *cannonDisplay {
	c := &cannonDisplay{
		names: map[string]*tentsuyu.MenuElement{},
		bars:  map[string]*cannonCooldownBar{},
	}
	textColor := color.RGBA{247, 6, 207, 255}
	y := 10.0
	c.names["Title"] = &tentsuyu.MenuElement{
		UIElement:  tentsuyu.NewTextElement(15, y, 250, 100, font, []string{"HEAT LEVELS:"}, textColor, 26),
		Action:     nil,
		Selectable: false,
	}
	y += 55
	for _, cannon := range cannons {

		c.names[cannon.Name] = &tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElement(15, y, 250, 100, font, []string{cannon.Name}, cannon.color, 26),
			Action:     nil,
			Selectable: false,
		}
		y += 30
		c.bars[cannon.Name] = newcannonCooldownBar(15, y, textColor, cannon)
		y += 33
	}

	return c
}

type highscoreDisplay struct {
	scores []*tentsuyu.TextElement
}

func newHighScoreDisplay(lb *scoreboard.LeaderBoard, font *truetype.Font) *highscoreDisplay {
	textColor := color.RGBA{45, 226, 230, 255}
	h := &highscoreDisplay{
		scores: []*tentsuyu.TextElement{},
	}
	y := 200.0
	h.scores = append(h.scores, tentsuyu.NewTextElement(1015, y, 250, 100, font, []string{"TOP SCORES:"}, color.RGBA{249, 200, 14, 255}, 26))
	y += 35
	for i, r := range lb.Records {
		h.scores = append(h.scores, tentsuyu.NewTextElement(1015, y, 250, 100, font, []string{strconv.Itoa(i+1) + ". " + r.Name + " : " + strconv.Itoa(r.Score)}, textColor, 26))
		y += 35
	}
	return h
}
