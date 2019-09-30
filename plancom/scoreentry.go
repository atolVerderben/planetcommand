package plancom

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//ScoreEntry updates and displays the entry screen for entering the initials for the high score screen
type ScoreEntry struct {
	choices                        []string
	Entry1, Entry2, Entry3         int
	Text1, Text2, Text3, TitleText *tentsuyu.TextElement
	CurrSelection                  int
	StartX, StartY                 float64
	Confirmed                      bool
}

//NewScoreEntry returns a new ScoreEntry at the given x,y coords
func NewScoreEntry(x, y float64) *ScoreEntry {
	secondRow := 145.0 //Distance for the initials below the title
	textColor := color.RGBA{247, 6, 207, 255}
	s := &ScoreEntry{
		StartX: x,
		StartY: y,
		choices: []string{
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
			"!", "@", "$", "%", "^", "&", "*",
		},
		TitleText: tentsuyu.NewTextElement(x+200, y, 1000, 100, Game.UIController.ReturnFont(FntMain), []string{"ENTER YOUR INITALS.", "USE THE ARROW KEYS.", "PRESS ENTER TO CONFIRM"}, textColor, 32),
		Text1:     tentsuyu.NewTextElement(x+280, y+secondRow, 200, 40, Game.UIController.ReturnFont(FntMain), []string{"A"}, textColor, 32),
		Text2:     tentsuyu.NewTextElement(x+280+50, y+secondRow, 200, 40, Game.UIController.ReturnFont(FntMain), []string{"A"}, textColor, 32),
		Text3:     tentsuyu.NewTextElement(x+280+100, y+secondRow, 200, 40, Game.UIController.ReturnFont(FntMain), []string{"A"}, textColor, 32),
	}
	return s
}

//Update the ScoreEntry.
//Up and Down change the selected letter space either up or down the alphabet.
//Left and Right change which letter space is selected.
//Enter confirms the entry to a high score.
func (s *ScoreEntry) Update(input *tentsuyu.InputController) {
	prev1 := s.Entry1
	prev2 := s.Entry2
	prev3 := s.Entry3
	if input.Button("Enter").JustPressed() {
		switch s.CurrSelection {
		case 0:
			s.CurrSelection++
		case 1:
			s.CurrSelection++
		case 2:
			s.Confirmed = true
		}
	}
	if input.Button("Left").JustPressed() {
		s.CurrSelection--
		if s.CurrSelection < 0 {
			s.CurrSelection = 2
		}
	}
	if input.Button("Right").JustPressed() {
		s.CurrSelection++
		if s.CurrSelection > 2 {
			s.CurrSelection = 0
		}
	}

	if input.Button("Up").JustPressed() {
		switch s.CurrSelection {
		case 0:
			s.Entry1--
			if s.Entry1 < 0 {
				s.Entry1 = len(s.choices) - 1
			}
		case 1:
			s.Entry2--
			if s.Entry2 < 0 {
				s.Entry2 = len(s.choices) - 1
			}
		case 2:
			s.Entry3--
			if s.Entry3 < 0 {
				s.Entry3 = len(s.choices) - 1
			}

		}

	}

	if input.Button("Down").JustPressed() {

		switch s.CurrSelection {
		case 0:
			s.Entry1++
			if s.Entry1 > len(s.choices)-1 {
				s.Entry1 = 0
			}
		case 1:
			s.Entry2++
			if s.Entry2 > len(s.choices)-1 {
				s.Entry2 = 0
			}
		case 2:
			s.Entry3++
			if s.Entry3 > len(s.choices)-1 {
				s.Entry3 = 0
			}

		}
	}

	if s.Entry1 != prev1 {
		s.Text1.SetText([]string{s.choices[s.Entry1]})
	}
	if s.Entry2 != prev2 {
		s.Text2.SetText([]string{s.choices[s.Entry2]})
	}
	if s.Entry3 != prev3 {
		s.Text3.SetText([]string{s.choices[s.Entry3]})
	}

	switch s.CurrSelection {
	case 0:
		s.Text1.Highlighted()
		s.Text2.UnHighlighted()
		s.Text3.UnHighlighted()
	case 1:
		s.Text1.UnHighlighted()
		s.Text2.Highlighted()
		s.Text3.UnHighlighted()
	case 2:
		s.Text1.UnHighlighted()
		s.Text2.UnHighlighted()
		s.Text3.Highlighted()
	default:
		s.Text1.UnHighlighted()
		s.Text2.UnHighlighted()
		s.Text3.UnHighlighted()
	}

}

//Draw the score Entry
func (s *ScoreEntry) Draw(screen *ebiten.Image) error {
	s.TitleText.Draw(screen)
	s.Text1.Draw(screen)
	s.Text2.Draw(screen)
	s.Text3.Draw(screen)

	return nil
}

//ToString returns the score entry 3 values as a single string
func (s ScoreEntry) ToString() string {
	return s.choices[s.Entry1] + s.choices[s.Entry2] + s.choices[s.Entry3]
}
