package session

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/PatrikOlin/gordle/db"
	g "github.com/PatrikOlin/gordle/guess"
	"github.com/PatrikOlin/gordle/word"
)

type Session struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Word string    `json:"-" db:"word"`
	// WordState    string    `json:"wordState" db:"word_state"`
	Status       string    `json:"status" db:"status"`
	Guesses      []g.Guess `json:"guesses" db:"-"`
	NumOfGuesses int       `json:"numberOfGuesses" db:"number_of_guesses"`
}

func Create() Session {
	s := Session{
		ID:   uuid.New(),
		Word: word.New(),
		// WordState: ".....",
		Status:  "unsolved",
		Guesses: make([]g.Guess, 0, 6),
		// GuessesString: "",
		NumOfGuesses: 0,
	}

	persistSession(s)

	return s
}

func Get(sessionID string) (Session, error) {
	var session Session
	stmt := "SELECT * FROM sessions WHERE id = $1"
	err := db.DBClient.Get(&session, stmt, sessionID)

	if err != nil {
		fmt.Println(err)
		return session, err
	}

	// if session.GuessesString != "" {
	// 	session.Guesses = strings.Split(session.GuessesString, ",")
	// }

	return session, nil
}

func persistSession(s Session) {
	stmt := "INSERT INTO sessions (id, word, status, number_of_guesses) VALUES (?, ?, ?, ?)"

	_, err := db.DBClient.Exec(stmt, s.ID, s.Word, s.Status, s.NumOfGuesses)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Session) Update() {
	stmt := "UPDATE sessions SET status=:status, number_of_guesses=:number_of_guesses WHERE id=:id"

	_, err := db.DBClient.NamedExec(stmt, s)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Session) GetGuesses() {
	stmt := "SELECT word, word_state FROM guesses WHERE session_id = ?"
	db.DBClient.Select(&s.Guesses, stmt, s.ID)
}
