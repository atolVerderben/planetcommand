package plancom

import (
	"image/color"
	"math"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type line struct {
	X1, Y1, X2, Y2, Angle float64
}

func (l *line) angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

//Command is an integer that tells the cannons what to do
type Command int

//List of available Commands
const (
	CommandNone Command = iota
	CommandLeft
	CommandRight
)

//Cannon is the player's way of attacking incoming threats
type Cannon struct {
	*tentsuyu.BasicObject
	Name                                                        string
	commandQueue                                                []Command
	leftLine, rightLine                                         *line
	viewAngleOffset, viewLineLength                             float64
	canFire                                                     bool
	missiles                                                    []*Missile
	missileSpeed                                                float64
	Health                                                      int
	imgParts                                                    *tentsuyu.BasicImageParts
	missileColor                                                string
	heat, heatMax, heatIncrease, heatDecrease, overHeatDecrease int
	cooldownMax, cooldownCount                                  int
	overheated                                                  bool
	hitCooldown                                                 int
	maxMissiles                                                 int
	color                                                       color.Color
}

//Names of the Cannons for UI purposes
var (
	Cannon1 = "CANNON A" //"ACAPULCO"
	Cannon2 = "CANNON B" //"BEIJING"
	Cannon3 = "CANNON C" //"COLOGNE"
	Cannon4 = "CANNON D" //"DALLAS"
)

//CreateCannon returns a Cannon on the given planet at the starting angle (position on planet)
func CreateCannon(planet *Planet, angle float64, name string) *Cannon {
	c := &Cannon{
		BasicObject:      tentsuyu.NewBasicObject(0, 0, 32, 32),
		leftLine:         &line{},
		rightLine:        &line{},
		viewAngleOffset:  0.436332,
		viewLineLength:   1000,
		Name:             name,
		canFire:          true,
		missiles:         []*Missile{},
		missileSpeed:     4,
		Health:           2,
		heatMax:          100,
		heatDecrease:     1,
		overHeatDecrease: 3,
		cooldownCount:    0,
		maxMissiles:      20, //I want to make this upgradeable, but for now make it an unreachable number
		heatIncrease:     30, // Use this to increase the amount of heat generated on shooting a missile
		cooldownMax:      5,
	}

	switch name {
	case Cannon1:
		c.imgParts = tentsuyu.NewBasicImageParts(0, 0, 32, 32)
		c.missileColor = "orange"
		c.color = color.RGBA{255, 108, 17, 255}
	case Cannon2:
		c.imgParts = tentsuyu.NewBasicImageParts(32, 0, 32, 32)
		c.missileColor = "blue"
		c.color = color.RGBA{45, 226, 230, 255}
	case Cannon3:
		c.imgParts = tentsuyu.NewBasicImageParts(64, 0, 32, 32)
		c.missileColor = "yellow"
		c.color = color.RGBA{249, 200, 14, 255}
	case Cannon4:
		c.imgParts = tentsuyu.NewBasicImageParts(96, 0, 32, 32)
		c.missileColor = "pink"
		c.color = color.RGBA{212, 0, 120, 255}
	}

	c.Angle = angle
	c.Position.X, c.Position.Y = planet.GetX()+(float64(c.Height/2)+planet.GetRadius())*math.Cos(c.Angle), planet.GetY()+(float64(c.Height/2)+planet.GetRadius())*math.Sin(c.Angle)

	return c
}

//Update the cannon
func (c *Cannon) Update(planet *Planet) {

	c.Angle = math.Atan2(c.GetY()-planet.GetY(), c.GetX()-planet.GetX())

	for _, comm := range c.commandQueue {
		switch comm {
		case CommandLeft:
			c.AddAngle(-(2 * math.Pi / 180))
		case CommandRight:
			c.AddAngle(2 * math.Pi / 180)
		}
	}
	if len(c.commandQueue) > 0 {
		c.commandQueue = []Command{}
	}

	c.Position.X, c.Position.Y = planet.GetX()+(float64(c.Height/2)+planet.GetRadius())*math.Cos(c.Angle), planet.GetY()+(float64(c.Height/2)+planet.GetRadius())*math.Sin(c.Angle)

	c.setLines()

	if c.heat > 0 {
		c.cooldownCount++
		if c.cooldownCount > c.cooldownMax {
			c.cooldownCount = 0
			if !c.overheated {
				c.heat -= c.heatDecrease
			} else {
				c.heat -= c.overHeatDecrease
			}
			if c.heat <= 0 {
				c.heat = 0
				c.overheated = false
			}
		}
	}

	c.determineCanShoot()

	if c.hitCooldown > 0 {
		c.hitCooldown--
		if c.hitCooldown < 0 {
			c.hitCooldown = 0 //failsafe to not go under 0
		}
	}

}

func (c *Cannon) determineCanShoot() {

	if c.canFire == false {
		if !c.IsAlive() {
			return
		}
		if len(c.missiles) > c.maxMissiles {
			return
		}
		if c.overheated {
			return
		}
	}

	c.canFire = true

}

//InFireRadius returns true if the given coords are within the cannon's "firing cone"
func (c *Cannon) InFireRadius(x, y float64) bool {
	pt := &tentsuyu.Vector2d{X: x, Y: y}
	v1 := &tentsuyu.Vector2d{X: c.leftLine.X1, Y: c.leftLine.Y1}
	v2 := &tentsuyu.Vector2d{X: c.rightLine.X2, Y: c.rightLine.Y2}
	v3 := &tentsuyu.Vector2d{X: c.leftLine.X2, Y: c.leftLine.Y2}

	d1 := sign(pt, v1, v2)
	d2 := sign(pt, v2, v3)
	d3 := sign(pt, v3, v1)

	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNeg && hasPos)
}

func sign(p1, p2, p3 *tentsuyu.Vector2d) float64 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

func (c *Cannon) setLines() {
	offset := c.viewAngleOffset
	x := c.GetX() - (float64(c.Width/2) * math.Cos(c.Angle))
	y := c.GetY() - (float64(c.Height/2) * math.Sin(c.Angle))
	tx := x + (c.viewLineLength * math.Cos(c.Angle+offset))
	ty := y + (c.viewLineLength * math.Sin(c.Angle+offset))
	//ebitenutil.DrawLine(screen, x, y, tx, ty, color.RGBA{255, 255, 0, 150})
	c.leftLine.X1 = x
	c.leftLine.X2 = tx
	c.leftLine.Y1 = y
	c.leftLine.Y2 = ty
	c.leftLine.Angle = c.Angle + offset

	tx = x + (c.viewLineLength * math.Cos(c.Angle-offset))
	ty = y + (c.viewLineLength * math.Sin(c.Angle-offset))

	c.rightLine.X1 = x
	c.rightLine.X2 = tx
	c.rightLine.Y1 = y
	c.rightLine.Y2 = ty
	c.rightLine.Angle = c.Angle - offset
	//ebitenutil.DrawLine(screen, x, y, tx, ty, color.RGBA{255, 255, 0, 150})
}

//AddCommand to command queue
func (c *Cannon) AddCommand(command Command) {
	c.commandQueue = append(c.commandQueue, command)
}

//ShiftCommand removes the first element from the command queue and returns it.
func (c *Cannon) ShiftCommand() Command {
	x, a := c.commandQueue[0], c.commandQueue[1:]
	c.commandQueue = a
	return x
}

//Draw the Cannon and the view cone
func (c Cannon) Draw(screen *ebiten.Image) error {
	if c.Health <= 0 {
		return nil
	}
	c.DrawRays(screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(c.Width/2), -float64(c.Height/2))
	op.GeoM.Rotate(c.Angle + AngleDrawOffset)
	//op.GeoM.Rotate(s.head.drawAngleAdj)
	op.GeoM.Translate(c.GetX(), c.GetY())

	if c.Health == 1 {
		c.imgParts.Sy = 32
	}

	screen.DrawImage(c.imgParts.SubImage(Game.ImageManager.ReturnImage("cannon")), op)

	return nil
}

//DrawRays draws the "view cone" for the cannon
func (c Cannon) DrawRays(screen *ebiten.Image) { //orig color: color.RGBA{247, 6, 207, 255}
	ebitenutil.DrawLine(screen, c.leftLine.X1, c.leftLine.Y1, c.leftLine.X2, c.leftLine.Y2, c.color)
	ebitenutil.DrawLine(screen, c.rightLine.X1, c.rightLine.Y1, c.rightLine.X2, c.rightLine.Y2, c.color)
}

//FireMissile returns true if the cannon is able to fire and creates a new missile
func (c *Cannon) FireMissile(tx, ty float64) bool {
	if c.canFire {
		c.missiles = append(c.missiles,
			CreateMissile(c.GetX(), c.GetY(), tx, ty, c.missileSpeed, c.missileColor))
		Game.AudioPlayer.PlaySE("blaster")
		c.heat += c.heatIncrease
		if c.heat >= c.heatMax {
			c.overheated = true
			c.canFire = false
			Game.AudioPlayer.PlaySE("overheat")
		}
		if len(c.missiles) > c.maxMissiles {
			c.canFire = false
		}
		return true
	}

	return false
}

//RemoveMissile from the cannon's slice
func (c *Cannon) RemoveMissile(missileID string) {
	delete := -1
	for index, e := range c.missiles {
		if e.GetIDasString() == missileID {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.missiles = append(c.missiles[:delete], c.missiles[delete+1:]...)
	}
}

//IsAlive returns true if cannon health is greater than 0
func (c Cannon) IsAlive() bool {
	if c.Health > 0 {
		return true
	}
	return false
}

//Hit reduces the cannon's hp if it's not cooling down. The cooldown period prevents instant death
func (c *Cannon) Hit() {
	if c.hitCooldown == 0 {
		c.Health--
		c.hitCooldown = 30
		if c.Health <= 0 {
			c.canFire = false
		}
	}
}
