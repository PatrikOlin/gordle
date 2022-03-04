package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	s "github.com/PatrikOlin/gordle/session"
)

type Guess struct {
	Word string `json:"word"`
}

func GuessWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	sessionID := chi.URLParam(r, "id")
	session := s.Get(sessionID)

	var guess Guess
	json.NewDecoder(r.Body).Decode(&guess)

	session, err := guessWord(session, guess)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func guessWord(session s.Session, guess Guess) (sess s.Session, err error) {
	if utf8.RuneCountInString(guess.Word) != 5 {
		return session, errors.New("Guess a FIVE LETTER word")
	}

	session.WordState = session.TestWord(guess.Word)
	if len(session.Guesses) > 0 {
		session.Guesses = append(session.Guesses, guess.Word)
	} else {
		session.Guesses = []string{guess.Word}
	}
	session.GuessesString = strings.Join(session.Guesses, ",")
	session.NumOfGuesses++

	if session.WordState == "GGGGG" {
		session.Status = "solved"
	}

	session.Update()
	return session, nil
}
