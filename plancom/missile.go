package plancom

import (
	"image/color"
	"math"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

//Missile is the main projectile of the game
type Missile struct {
	*tentsuyu.BasicObject
	OriginX, OriginY float64
	TargetX, TargetY float64
	Velocity         *tentsuyu.Vector2d
	Active           bool
	cr, cg, cb       float64
	ctr, ctg, ctb    uint8
}

//CreateMissile returns a new Missile at the specified x,y coords with the target tx,ty
func CreateMissile(x, y, tx, ty, speed float64, color string) *Missile {
	m := &Missile{
		BasicObject: tentsuyu.NewBasicObject(x, y, 2, 2),
		OriginX:     x,
		OriginY:     y,
		TargetX:     tx,
		TargetY:     ty,
		Active:      true,
		Velocity:    &tentsuyu.Vector2d{},
	}
	m.Speed = speed
	m.SetAngle(math.Atan2(ty-y, tx-x))
	m.Velocity.X = m.Speed * math.Cos(m.Angle)
	m.Velocity.Y = m.Speed * math.Sin(m.Angle)

	//Set the colors
	switch color {
	case "red":
		m.cr = 1
		m.cg = 0
		m.cb = 0
		m.ctr = 253
		m.ctg = 29
		m.ctb = 83
	case "blue":
		m.cr = 0
		m.cg = 0
		m.cb = 1
		m.ctr = 45
		m.ctg = 226
		m.ctb = 230
	case "green":
		m.cr = 0
		m.cg = 1
		m.cb = 0
		m.ctr = 0
		m.ctg = 255
		m.ctb = 0
	case "orange":
		m.cr = 1
		m.cg = 0.7
		m.cb = 0
		m.ctr = 255
		m.ctg = 108
		m.ctb = 17
	case "yellow":
		m.cr = 1
		m.cg = 1
		m.cb = 0
		m.ctr = 249
		m.ctg = 200
		m.ctb = 14
	case "pink":
		m.cr = 0.8
		m.cg = 0
		m.cb = 0.5
		m.ctr = 212
		m.ctg = 0
		m.ctb = 120
	default:
		m.cr = 1
		m.cg = 1
		m.cb = 1
		m.ctr = 255
		m.ctg = 255
		m.ctb = 255
	}
	return m
}

//Update the missile
func (m *Missile) Update() {
	m.X += m.Velocity.X
	m.Y += m.Velocity.Y
}

//Draw the missile
func (m *Missile) Draw(screen *ebiten.Image) error {

	//m.DrawTrail(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(-float64(m.Width/2), -float64(m.Height/2))
	//op.GeoM.Rotate(g.planet.Angle)
	op.GeoM.Translate(m.X, m.Y)
	op.ColorM.Scale(m.cr, m.cg, m.cb, 1)

	screen.DrawImage(Game.ImageManager.ReturnImage("pixel"), op)

	return nil
}

//DrawTrail of the Missile. This is separate so that it can be drawn at a different level
func (m *Missile) DrawTrail(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, m.X, m.Y, m.OriginX, m.OriginY, color.RGBA{m.ctr, m.ctg, m.ctb, 255})
}

//Projectile is the generic
type Projectile interface {
	Draw(*ebiten.Image) error
	DrawTrail(*ebiten.Image)
	Update()
}
