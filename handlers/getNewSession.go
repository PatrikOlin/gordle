package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Session struct {
	ID uuid.UUID
}

func GetNewSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	session := Session{ID: uuid.New()}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(session)
}
