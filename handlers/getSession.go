package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	s "github.com/PatrikOlin/gordle/session"
	"github.com/go-chi/chi/v5"
)

func GetSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	sessionID := chi.URLParam(r, "id")

	fmt.Println(sessionID)
	var session s.Session
	if sessionID != "" {
		session = s.Get(sessionID)
	} else {
		session = s.Create()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}
