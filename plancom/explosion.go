package plancom

import (
	"math/rand"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//Explosion is generated when a missile hits its target or reaches its destination
type Explosion struct {
	*tentsuyu.BasicObject
	scaleCount int
	eScale     float64
	w, h       float64
	r, g, b    float64
	done       bool
	givePoints bool
}

//CreateExplosion at given coords
func CreateExplosion(x, y float64, givePoints bool) *Explosion {
	e := &Explosion{
		BasicObject: tentsuyu.NewBasicObject(x, y, 32, 32),
		eScale:      0.1,
		w:           32 * 0.1,
		h:           32 * 0.1,
		r:           1,
		g:           1,
		b:           1,
		givePoints:  givePoints,
	}

	return e
}

//Update explosion
func (e *Explosion) Update() {
	e.scaleCount++
	if e.scaleCount == 5 {
		e.eScale += 0.1
		e.w = 32 * e.eScale
		e.h = 32 * e.eScale
		e.Width = int(e.w)
		e.Height = int(e.h)
		e.scaleCount = 0
		e.r = rand.Float64()
		e.g = rand.Float64()
		e.b = rand.Float64()
	}
	if e.w >= 48 {
		e.done = true
	}
}

//Draw the explosion
func (e *Explosion) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(e.eScale, e.eScale)
	op.GeoM.Translate(-float64(e.w/2), -float64(e.h/2))
	//op.GeoM.Rotate(g.planet.Angle)
	op.GeoM.Translate(e.GetX(), e.GetY())
	op.ColorM.Scale(e.r, e.g, e.b, 1)
	screen.DrawImage(Game.ImageManager.ReturnImage("explosion"), op)

	return nil
}
