package plancom

import "github.com/atolVerderben/tentsuyu"

//TODO: Implement a UFO for some variety

//AIUfo controls the ufo that attacks the player planet
type AIUfo struct {
	*tentsuyu.BasicObject
	Health int
}

//Update the AIUfo
func (u *AIUfo) Update() {

}
