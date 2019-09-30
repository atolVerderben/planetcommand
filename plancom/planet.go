package plancom

import "github.com/atolVerderben/tentsuyu"

//Planet is what the player must protect
type Planet struct {
	*tentsuyu.BasicObject
	Health int
}

//CreatePlanet at given coords with given width and height
func CreatePlanet(x, y float64, w, h int) *Planet {
	p := &Planet{
		BasicObject: tentsuyu.NewBasicObject(x, y, w, h),
		Health:      3,
	}

	return p
}

//GetRadius returns the radius of the planet (i.e. width/2)
func (p Planet) GetRadius() float64 {
	return float64(p.Width / 2)
}

//IsAlive returns true if the planet has health above 0
func (p Planet) IsAlive() bool {
	if p.Health > 0 {
		return true
	}
	return false
}
