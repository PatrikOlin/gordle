package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	g "github.com/PatrikOlin/gordle/guess"
	r "github.com/PatrikOlin/gordle/rules"
	s "github.com/PatrikOlin/gordle/session"
)

type IncomingGuess struct {
	Word string `json:"word"`
}

func GuessWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	sessionID := chi.URLParam(r, "id")
	session, err := s.Get(sessionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	var guess IncomingGuess
	json.NewDecoder(r.Body).Decode(&guess)

	session, err = guessWord(session, guess)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func guessWord(session s.Session, guess IncomingGuess) (sess s.Session, err error) {
	if utf8.RuneCountInString(guess.Word) != 5 {
		return session, errors.New("Guess a FIVE LETTER word")
	}
	if session.NumOfGuesses >= r.Get().MaxGuesses {
		return session, errors.New("Out of guesses")
	}

	// session.WordState = session.TestWord(guess.Word)
	answer := g.MakeGuess(guess.Word, session.Word, session.ID.String())
	session.NumOfGuesses++

	if answer.WordState == "GGGGG" {
		session.Status = "solved"
	}

	session.Update()
	session.GetGuesses()
	return session, nil
}
