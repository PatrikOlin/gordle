package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	e "github.com/PatrikOlin/gordle/pkg/errors"
	g "github.com/PatrikOlin/gordle/pkg/guess"
	"github.com/PatrikOlin/gordle/pkg/rules"
	r "github.com/PatrikOlin/gordle/pkg/rules"
	s "github.com/PatrikOlin/gordle/pkg/session"
	us "github.com/PatrikOlin/gordle/pkg/user-session"
)

type IncomingGuess struct {
	Word string `json:"word"`
}

func GuessWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	c, err := r.Cookie("user_session_token")
	if err != nil {
		error := e.E("GetSession", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(error)
		json.NewEncoder(w).Encode(error)

		return
	}

	userSession, err := us.GetUserSession(c.Value)
	if err != nil {
		error := e.E("GuessWord", err, http.StatusInternalServerError, "User session not found, try clearing your cookies.")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(error)
		json.NewEncoder(w).Encode(error)

		return
	}

	session, err := s.Get(userSession)
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

	if session.Status == "solved" || session.NumOfGuesses >= rules.Get().MaxGuesses {
		fs := session.GetStats(userSession.Token.String())
		if session.IsDaily {
			userSession.SetDailyFinished()
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fs)

		return
	}

	session.SetWordVisibility()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func guessWord(session s.Session, guess IncomingGuess) s.Session {
	newGuess := g.MakeGuess(strings.ToLower(guess.Word), strings.ToLower(session.Word), session.ID.String())
	session.Guesses = append(session.Guesses, newGuess)

	session.NumOfGuesses++

	if newGuess.WordState == "GGGGG" {
		session.Status = "solved"
		session.FinishedAt = int(time.Now().Unix())
	}

	session.Update()
	return session
}

func isGuessValid(session s.Session, guess IncomingGuess) (ok bool, err error) {
	if utf8.RuneCountInString(guess.Word) != 5 {
		return false, errors.New("Gissa på ett ord med FEM bokstäver")
	}
	if session.NumOfGuesses >= r.Get().MaxGuesses {
		return false, errors.New("Slut på gissningar")
	}
	if !g.IsWordInList(guess.Word) {
		return false, errors.New("Finns inte i ordlistan")
	}
	if session.Status == "solved" {
		return false, errors.New("Du har redan löst det här ordet, ladda om sidan för att få ett nytt")
	}

	return true, nil
}
