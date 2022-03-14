package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	s "github.com/PatrikOlin/gordle/session"
)

func GetSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	sessionID := chi.URLParam(r, "id")

	session, err := getSession(sessionID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	session.SetWordVisibility()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}

func getSession(sessionID string) (s.Session, error) {
	fmt.Println(sessionID)

	if sessionID != "" {
		return s.Get(sessionID)
	}

	return s.Create(), nil
}
