package guess

import (
	"fmt"
	"strings"

	"github.com/PatrikOlin/gordle/db"
)

type Guess struct {
	Word      string `json:"word" db:"word"`
	WordState string `json:"wordState" db:"word_state"`
}

func MakeGuess(guess string, correctWord string, sessionID string) Guess {
	var (
		correct    string = "G"
		correctish string = "Y"
		wrong      string = "."
		allCorrect string = "GGGGG"
	)

	if correctWord == guess {
		g := Guess{Word: guess, WordState: allCorrect}
		err := persistGuess(g, sessionID)
		if err != nil {
			fmt.Println(err)
		}
		return g
	}

	var wordState []string
	wordState = make([]string, 5)

	runesGuess := []rune(guess)
	runesCorrect := []rune(correctWord)

	for i := 0; i < len(runesGuess); i++ {
		if runesGuess[i] == runesCorrect[i] {
			wordState[i] = correct
		} else if contains(runesCorrect, runesGuess[i]) {
			wordState[i] = correctish
		} else {
			wordState[i] = wrong
		}
	}

	g := Guess{Word: guess, WordState: strings.Join(wordState, "")}
	err := persistGuess(g, sessionID)
	if err != nil {
		fmt.Println(err)
	}

	return g
}

func persistGuess(g Guess, sessionID string) error {
	stmt := "INSERT INTO guesses (session_id, word, word_state) VALUES (?, ?, ?)"
	_, err := db.DBClient.DB.Exec(stmt, sessionID, g.Word, g.WordState)
	if err != nil {
		return err
	}
	return nil
}

func contains(r []rune, ru rune) bool {
	for _, v := range r {
		if v == ru {
			return true
		}
	}
	return false
}
