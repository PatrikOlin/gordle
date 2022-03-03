package session

import (
	"log"

	"github.com/google/uuid"

	"github.com/PatrikOlin/gordle/db"
	"github.com/PatrikOlin/gordle/word"
)

type Session struct {
	ID     uuid.UUID
	Word   string
	Status string
}

func Create() Session {
	s := Session{
		ID:     uuid.New(),
		Word:   word.New(),
		Status: "unsolved",
	}

	persistSession(s)

	return s
}

func persistSession(s Session) {
	stmt := "INSERT INTO sessions (id, status, word) VALUES (?, ?, ?)"

	_, err := db.DBClient.Exec(stmt, s.ID, s.Status, s.Word)
	if err != nil {
		log.Fatalln(err)
	}
}
