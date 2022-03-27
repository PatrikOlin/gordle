package stats

import (
	"math"
)

type Stats struct {
	MaxStreak       int             `json:"maxStreak"`
	CurrentStreak   int             `json:"currentStreak"`
	WinPercentage   float64         `json:"winPercentage"`
	WinDistribution WinDistribution `json:"winDistribution"`
}

type WinDistribution struct {
	Unsolved    int `json:"unsolved"`
	FirstGuess  int `json:"firstGuess"`
	SecondGuess int `json:"secondGuess"`
	ThirdGuess  int `json:"thirdGuess"`
	FourthGuess int `json:"fourthGuess"`
	FifthGuess  int `json:"fifthGuess"`
	SixthGuess  int `json:"sixthGuess"`
}

func (wd *WinDistribution) Add(numOfGuesses int) {
	switch numOfGuesses {
	case 1:
		wd.FirstGuess += 1
	case 2:
		wd.SecondGuess += 1
	case 3:
		wd.ThirdGuess += 1
	case 4:
		wd.FourthGuess += 1
	case 5:
		wd.FifthGuess += 1
	case 6:
		wd.SixthGuess += 1
	}
}

func (st *Stats) CalculateWinPercentage() {
	wins :=
		st.WinDistribution.FirstGuess +
			st.WinDistribution.SecondGuess +
			st.WinDistribution.ThirdGuess +
			st.WinDistribution.FourthGuess +
			st.WinDistribution.FifthGuess +
			st.WinDistribution.SixthGuess

	total := wins + st.WinDistribution.Unsolved
	delta := (float64(wins) / float64(total)) * 100

	st.WinPercentage = math.Round(delta)
}
