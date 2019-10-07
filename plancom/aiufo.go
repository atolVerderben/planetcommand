package plancom

import (
	"math"
	"math/rand"

	"github.com/atolVerderben/tentsuyu"
	"github.com/atolVerderben/tentsuyu/tentsuyutils"
	"github.com/hajimehoshi/ebiten"
)

//AIUfo controls the ufo that attacks the player planet
type AIUfo struct {
	*tentsuyu.BasicObject
	Health             int
	Destination        *tentsuyu.Vector2d
	anim               *tentsuyu.Animation
	isIdle             bool
	idleCount, idleMax int
	maxSpeed           float64
	acc                float64
	hitCooldown        int
}

//SpawnUFO at the provided coords
func SpawnUFO(x, y float64) *AIUfo {
	u := &AIUfo{
		BasicObject: tentsuyu.NewBasicObject(x, y, 32, 32),
		Destination: &tentsuyu.Vector2d{X: -1, Y: -1},
		maxSpeed:    3.0,
		anim:        tentsuyu.NewAnimation(SpriteSheets["UFO"], []int{0, 1}, 5),
		acc:         0.1,
		Health:      2,
	}
	u.Speed = 3.0
	return u
}

//Update the AIUfo
func (u *AIUfo) Update(p *Planet, playArea *tentsuyu.Rectangle) {
	u.anim.Update()
	if u.hitCooldown > 0 {
		if u.isIdle {
			u.setIdle(false)
		}
		u.hitCooldown++
		if u.hitCooldown > 120 {
			u.hitCooldown = 0
		}
	}
	if u.isIdle {
		u.anim.SetFrameSpeed(9)
		u.idleCount++
		if u.idleCount > u.idleMax {
			u.setIdle(false)
		}
		return
	}
	if !u.hasDestination() {
		//Get a new destination to move to
		x := tentsuyutils.RandomBetweenf(0, float64(playArea.W)) //playArea.X, playArea.X+float64(playArea.W))
		y := tentsuyutils.RandomBetweenf(0, float64(playArea.H)) //playArea.Y, playArea.Y+float64(playArea.H))
		u.setDestination(x, y)
		return
	}
	if tentsuyutils.Distance(u.GetX(), u.GetY(), u.Destination.X, u.Destination.Y) < 200 {
		u.Speed = u.maxSpeed / 2
	} else {
		u.Speed = u.maxSpeed
	}
	if tentsuyutils.Distance(u.GetX(), u.GetY(), u.Destination.X, u.Destination.Y) < 5 {
		u.reachedDestination()
		return
	}
	if u.hasDestination() {
		if u.acc < u.Speed {
			u.acc += 0.1
		}
		u.SetAngle(math.Atan2(u.Destination.Y-u.GetY(), u.Destination.X-u.GetX()))
		u.Velocity.X = u.acc * math.Cos(u.Angle)
		u.Velocity.Y = u.acc * math.Sin(u.Angle)
		u.Velocity.Limit(u.Speed)
		u.Position.Add(*u.Velocity)
		u.anim.SetFrameSpeed(5)

	}
	u.X = u.Position.X
	u.Y = u.Position.Y

}

//Draw the ufo
func (u *AIUfo) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(u.Width)/2, -float64(u.Height)/2)
	//op.GeoM.Rotate(u.Angle)
	op.GeoM.Translate(u.GetX(), u.GetY())
	if u.hitCooldown > 0 {
		op.ColorM.Scale(1, 1, 1, 0.5)
	}
	screen.DrawImage(u.anim.ReturnImageParts().SubImage(Game.ImageManager.ReturnImage("ufo")), op)
	return nil
}

func (u *AIUfo) setIdle(idle bool) {
	if idle {
		u.isIdle = idle
		u.idleCount = 0
	} else {
		u.isIdle = idle
		u.idleCount = 0
	}
}

func (u *AIUfo) setDestination(x, y float64) {
	u.Destination.X = x
	u.Destination.Y = y
}

func (u *AIUfo) hasDestination() bool {

	if u.Destination.X == -1 && u.Destination.Y == -1 {
		return false
	}
	return true
}

func (u *AIUfo) reachedDestination() {
	u.acc = 0.1
	u.Velocity.Mul(0)
	u.setDestination(-1, -1)
	if rand.Intn(2) == 1 {
		u.setIdle(true)
		u.idleMax = tentsuyutils.RandomBetween(60, 240)
	}
}

//IsDead returns true if the UFO has no more health
func (u AIUfo) IsDead() bool {
	if u.Health <= 0 {
		return true
	}
	return false
}

//IsAlive returns true if UFO health is greater than 0
func (u AIUfo) IsAlive() bool {
	return !u.IsDead()
}

//Hit updates the UFO to a hit state
func (u *AIUfo) Hit() {
	if u.hitCooldown > 0 {
		u.hitCooldown++
		return
	}
	u.Health--
	u.hitCooldown = 1
	if u.Health < 0 {
		u.Health = 0
	}
}
