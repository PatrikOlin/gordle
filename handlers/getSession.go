package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	e "github.com/PatrikOlin/gordle/errors"
	s "github.com/PatrikOlin/gordle/session"
	us "github.com/PatrikOlin/gordle/userSession"
)

func GetSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	c, err := r.Cookie("user_session_token")
	if err == http.ErrNoCookie {
		fmt.Println("no cookie!")
		c = &http.Cookie{
			Name:     "user_session_token",
			Value:    us.Create().Token.String(),
			MaxAge:   int((365 * 5 * 24 * time.Hour).Seconds()),
			HttpOnly: true,
		}
	} else if err != nil {
		error := e.E("GetSession", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(error)
		json.NewEncoder(w).Encode(error)
		return
	}

	session, err := getSession(c.Value)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	session.SetWordVisibility()

	http.SetCookie(w, c)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func getSession(userToken string) (s.Session, error) {
	var session s.Session

	session, err := s.Get(userToken)

	if err == sql.ErrNoRows || !session.IsAlive() {
		return s.Create(userToken), nil
	} else if err != nil {
		return session, err
	}

	return session, nil
}
