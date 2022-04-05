package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	e "github.com/PatrikOlin/gordle/pkg/errors"
	s "github.com/PatrikOlin/gordle/pkg/session"
	us "github.com/PatrikOlin/gordle/pkg/user-session"
)

func GetSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	var userSession us.UserSession

	c, err := r.Cookie("user_session_token")
	if err == http.ErrNoCookie {
		fmt.Println("no cookie!")
		userSession = us.Create()
		c = &http.Cookie{
			Name:     "user_session_token",
			Value:    userSession.Token.String(),
			MaxAge:   int((365 * 5 * 24 * time.Hour).Seconds()),
			HttpOnly: true,
		}
	} else if err != nil {
		error := e.E("GetSession", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(error)
		json.NewEncoder(w).Encode(error)
		return
	} else {
		userSession, err = us.GetUserSession(c.Value)
		if err == sql.ErrNoRows {
			userSession = us.Create()
		} else if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
		}
	}

	session, err := getSession(userSession)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	session.SetWordVisibility()

	http.SetCookie(w, c)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func getSession(userSession us.UserSession) (s.Session, error) {
	var session s.Session

	session, err := s.Get(userSession)

	if err == sql.ErrNoRows || session.IsAlive() == false {
		return s.Create(userSession), nil
	} else if err != nil {
		return session, err
	}

	return session, nil
}
