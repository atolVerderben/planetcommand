package plancom

import "github.com/atolVerderben/tentsuyu"

//Orbiter is an object designed to orbit a planet
type Orbiter struct {
	*tentsuyu.BasicObject
	target   tentsuyu.GameObject
	distance float64
}

//NewOrbiter creates a new orbiter object
func NewOrbiter(x, y float64, w, h int, orbitObject tentsuyu.GameObject, distance float64) *Orbiter {
	o := &Orbiter{
		BasicObject: tentsuyu.NewBasicObject(x, y, w, h),
		target:      orbitObject,
		distance:    distance,
	}

	return o
}

//Update the orbiter
func (o *Orbiter) Update() {

}
