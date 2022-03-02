package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/PatrikOlin/gordle/handlers"
)

func GetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/word", handlers.GuessWord)
	r.Get("/rules", handlers.GetRules)
	// r.Get("/word", handlers.GetWord)
	// r.Get("/state", handlers.GetGameState)

	return r
}
