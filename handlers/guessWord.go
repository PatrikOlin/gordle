package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	e "github.com/PatrikOlin/gordle/errors"
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
		error := e.E("GuessWord", err, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(error)
		json.NewEncoder(w).Encode(error)
		return
	}

	var guess IncomingGuess
	json.NewDecoder(r.Body).Decode(&guess)

	_, err = isGuessValid(session, guess)
	if err != nil {
		error := e.E("GuessWord", err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(error)
		json.NewEncoder(w).Encode(error)
		return
	}

	session = guessWord(session, guess)

	session.SetWordVisibility()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func guessWord(session s.Session, guess IncomingGuess) s.Session {
	answer := g.MakeGuess(strings.ToLower(guess.Word), strings.ToLower(session.Word), session.ID.String())
	session.NumOfGuesses++

	if answer.WordState == "GGGGG" {
		session.Status = "solved"
	}

	session.Update()
	session.GetGuesses()
	return session
}

func isGuessValid(session s.Session, guess IncomingGuess) (ok bool, err error) {
	if utf8.RuneCountInString(guess.Word) != 5 {
		return false, errors.New("Guess a FIVE LETTER word")
	}
	if session.NumOfGuesses >= r.Get().MaxGuesses {
		return false, errors.New("Out of guesses")
	}
	if !g.IsWordInList(guess.Word) {
		return false, errors.New("That is not a word")
	}

	return true, nil
}
