package plancom

import "github.com/atolVerderben/tentsuyu"

//TODO: Implement a UFO for some variety

//AIUfo controls the ufo that attacks the player planet
type AIUfo struct {
	*tentsuyu.BasicObject
	Health      int
	Destination *tentsuyu.Vector2d
	anim        *tentsuyu.Animation
}

//SpawnUFO at the provided coords
func SpawnUFO(x, y float64) *AIUfo {
	u := &AIUfo{
		BasicObject: tentsuyu.NewBasicObject(x, y, 32, 32),
		Destination: &tentsuyu.Vector2d{X: -1, Y: -1},
	}
	return u
}

//Update the AIUfo
func (u *AIUfo) Update(p *Planet) {
	if !u.hasDestination() {
		//Get a new destination to move to
	}
}

func (u *AIUfo) hasDestination() bool {
	if u.Destination.X == -1 && u.Destination.Y == -1 {
		return false
	}
	return true
}
