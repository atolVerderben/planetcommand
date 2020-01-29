package plancom

import (
	"math/rand"

	"github.com/atolVerderben/tentsuyu/tentsuyutils"
)

//AI will represent the artificial intelligence of the game
type AI interface{}

//AIController controls all the AI elements for the game
type AIController struct {
	launchers  []*AILauncher
	updateTick int
}

//CreateAIController returns a new AIController
func CreateAIController() *AIController {
	a := &AIController{
		launchers: CreateEdgeAILaunchers(),
	}

	return a
}

//Update the AIController. This will coordinate all different AI Types
func (a *AIController) Update(p *Planet) {
	increaseDifficulty := false
	a.updateTick++
	if a.updateTick > 7200 { //2 minutes
		a.updateTick = 0
		increaseDifficulty = true
	}
	for _, l := range a.launchers {
		if increaseDifficulty {
			//log.Println("INCREASE!!")
			l.extraLaunches++
			l.minTime -= 60
			l.maxTime -= 60
			if l.minTime < 60 {
				l.minTime = 60
			}
			if l.maxTime < 120 {
				l.maxTime = 120
			}
		}
		l.Update(p)
	}
}

//AILauncher controls the various missile launchers
type AILauncher struct {
	*MissileLauncher
	prevLaunch       int
	nextLaunch       int
	minTime, maxTime int
	extraLaunches    int
	maxShots         int
}

//CreateEdgeAILaunchers is a convenient function to return the 4 edge launchers
func CreateEdgeAILaunchers() []*AILauncher {
	return []*AILauncher{CreateAILauncher(MissileLauncherLeft),
		CreateAILauncher(MissileLauncherRight),
		CreateAILauncher(MissileLauncherTop),
		CreateAILauncher(MissileLauncherBottom),
	}
}

//CreateAILauncher returns a new AILauncher
func CreateAILauncher(name string) *AILauncher {
	a := &AILauncher{
		MissileLauncher: CreateConstantMissileLauncher(name),
		minTime:         300,
		maxTime:         600,
		extraLaunches:   3,
	}
	a.nextLaunch = tentsuyutils.RandomBetween(a.minTime, a.maxTime)

	return a
}

//Update the AILauncher and determine when to fire new missiles
func (a *AILauncher) Update(planet *Planet) {
	a.prevLaunch++
	if a.prevLaunch >= a.nextLaunch { //TODO: Change this
		if planet.IsAlive() {
			a.LaunchMissile(planet.GetX(), planet.GetY())
			for i := 0; i < a.extraLaunches; i++ {
				if rand.Intn(8) == 1 {
					a.LaunchMissile(planet.GetX(), planet.GetY())
				}
			}
		}
		a.nextLaunch = tentsuyutils.RandomBetween(a.minTime, a.maxTime)
		a.prevLaunch = 0
	}
}

//LaunchMissile at the provided target x,y coords
func (a *AILauncher) LaunchMissile(tx, ty float64) {
	if a.isVertical {
		r := tentsuyutils.RandomBetween(0, a.Height)
		a.missiles = append(a.missiles,
			CreateMissile(a.GetX(), float64(r), tx, ty, 0.5, "green"))
	} else {
		r := tentsuyutils.RandomBetween(0, a.Width)
		a.missiles = append(a.missiles,
			CreateMissile(float64(r), a.GetY(), tx, ty, 0.5, "green"))
	}
	Game.AudioPlayer.PlaySE("blaster2")
}
