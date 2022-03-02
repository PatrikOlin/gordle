package handlers

import (
	"encoding/json"
	"net/http"
)

func GuessWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("RÅGAX")
}
