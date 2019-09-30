package plancom

import "github.com/atolVerderben/tentsuyu"

//MissileLauncher launches missiles
type MissileLauncher struct {
	*tentsuyu.BasicObject
	Name       string
	isVertical bool
	missiles   []*Missile
}

//RemoveMissile from the launcher's slice
func (m *MissileLauncher) RemoveMissile(missileID string) {
	delete := -1
	for index, e := range m.missiles {
		if e.GetIDasString() == missileID {
			delete = index
			break
		}
	}
	if delete >= 0 {
		m.missiles = append(m.missiles[:delete], m.missiles[delete+1:]...)
	}
}

//CreateMissileLauncher creates a MissileLauncher, which launches Missiles
func CreateMissileLauncher(x, y float64, width, height int, name string) *MissileLauncher {
	m := &MissileLauncher{
		BasicObject: tentsuyu.NewBasicObject(x, y, width, height),
		Name:        name,
	}
	if width < height {
		m.isVertical = true
	}
	return m
}

//CreateConstantMissileLauncher creates one of the 4 just off screen Missile Launchers that constantly
//launch missiles at our player
func CreateConstantMissileLauncher(name string) *MissileLauncher {
	switch name {
	case MissileLauncherLeft:
		return CreateMissileLauncher(-10, 0, 10, 720, name)
	case MissileLauncherRight:
		return CreateMissileLauncher(730, 0, 10, 720, name)
	case MissileLauncherTop:
		return CreateMissileLauncher(0, -10, 720, 10, name)
	case MissileLauncherBottom:
		return CreateMissileLauncher(0, 730, 720, 10, name)

	}
	return nil
}

//These specifiy the edge of the screen MissileLauncher names/positions
const (
	MissileLauncherLeft   = "LeftLauncher"
	MissileLauncherRight  = "RightLauncher"
	MissileLauncherTop    = "TopLauncher"
	MissileLauncherBottom = "BottomLauncher"
)
