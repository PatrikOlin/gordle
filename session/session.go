package session

import (
	"log"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/PatrikOlin/gordle/db"
	"github.com/PatrikOlin/gordle/word"
)

type Session struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Word          string    `json:"-" db:"word"`
	WordState     string    `json:"wordState" db:"word_state"`
	Status        string    `json:"status" db:"status"`
	Guesses       []string  `json:"guesses" db:"-"`
	GuessesString string    `json:"-" db:"guesses"'`
	NumOfGuesses  int       `json:"numberOfGuesses" db:"number_of_guesses"`
}

func Create() Session {
	s := Session{
		ID:            uuid.New(),
		Word:          word.New(),
		WordState:     ".....",
		Status:        "unsolved",
		Guesses:       make([]string, 0, 6),
		GuessesString: "",
		NumOfGuesses:  0,
	}

	persistSession(s)

	return s
}

func Get(sessionID string) Session {
	var session Session
	stmt := "SELECT * FROM sessions WHERE id = $1"
	err := db.DBClient.Get(&session, stmt, sessionID)

	if err != nil {
		log.Fatalln(err)
	}

	if session.GuessesString != "" {
		session.Guesses = strings.Split(session.GuessesString, ",")
	}

	return session
}

func persistSession(s Session) {
	stmt := "INSERT INTO sessions (id, status, word, word_state, guesses, number_of_guesses) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := db.DBClient.Exec(stmt, s.ID, s.Status, s.Word, s.WordState, s.GuessesString, s.NumOfGuesses)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Session) Update() {
	stmt := "UPDATE sessions SET status=:status, word_state=:word_state, guesses=:guesses, number_of_guesses=:number_of_guesses WHERE id=:id"

	_, err := db.DBClient.NamedExec(stmt, s)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Session) TestWord(guess string) string {
	var (
		correct    string = "G"
		correctish string = "Y"
		wrong      string = "."
		allCorrect string = "GGGGG"
	)

	if s.Word == guess {
		return allCorrect
	}

	var wordState []string
	wordState = make([]string, 5)

	guessSlice := strings.Split(guess, "")
	wordSlice := strings.Split(s.Word, "")

	for i := 0; i < utf8.RuneCountInString(s.Word); i++ {
		if wordSlice[i] == guessSlice[i] {
			wordState[i] = correct
		} else if strings.ContainsAny(s.Word, string(guessSlice[i])) {
			wordState[i] = correctish
		} else {
			wordState[i] = wrong
		}
	}

	return strings.Join(wordState, "")
}
