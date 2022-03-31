package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	rules := fmt.Sprint(`
	Rules:
		- You have 5 guesses
		- You can only guess using 5 letters words
		- You can only guess using lowercase letters
		- Green means correct at the right position
		- Yellow means correct at the wrong position
	Give a filename in argument to have more words

	Good luck!
	`)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rules)
}
