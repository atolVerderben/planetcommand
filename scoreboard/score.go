package scoreboard

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

//Record contains the name and score
type Record struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

//LeaderBoard is a collection of records
type LeaderBoard struct {
	Records []*Record `json:"records"`
}

//SetPosition of the new Record at the given position in the leaderboard
func (l *LeaderBoard) SetPosition(index int, name string, score int) bool {
	if index < 0 || index >= len(l.Records) {
		return false
	}
	l.Records = append(l.Records, &Record{Name: name, Score: score})
	sort.Slice(l.Records, func(i int, j int) bool {
		return l.Records[i].Score > l.Records[j].Score
	})
	l.Records = l.Records[:len(l.Records)-1] // remove last record
	//l.Records[index] = &Record{Name: name, Score: score}
	return true
}

//DeterminePosition takes the score and determines what position on the leaderboard it is.
//Returns -1 if the score does not fall on the leaderboard
func (l LeaderBoard) DeterminePosition(score int) int {
	found := false
	spot := -1
	for i, r := range l.Records {
		if r.Score < score && !found {
			spot = i
			found = true
			continue
		}
	}
	return spot
}

//SaveHighScores of the leaderboard to a json file
func (lb *LeaderBoard) SaveHighScores() bool {

	leader, err := json.Marshal(lb)
	if err != nil {
		return false
	}
	err = ioutil.WriteFile("gamedata/scores.json", leader, 0644)
	if err != nil {
		return false
	}
	return true
}

//LoadHighScores from the scores.json file into a leaderboard
func LoadHighScores() *LeaderBoard {
	lb := &LeaderBoard{
		Records: []*Record{},
	}

	raw, err := ioutil.ReadFile("gamedata/scores.json")
	if err != nil {
		log.Println(err.Error())
		//os.Exit(1)
	}

	json.Unmarshal(raw, lb)
	return lb
}
