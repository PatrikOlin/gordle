package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/PatrikOlin/gordle/handlers"
)

func GetRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.AllowAll().Handler)

	r.Get("/rules", handlers.GetRules)
	r.Post("/word/{id}", handlers.GuessWord)
	r.Get("/session", handlers.GetSession)
	r.Get("/session/{id}", handlers.GetSession)
	// r.Get("/word", handlers.GetWord)
	// r.Get("/state", handlers.GetGameState)

	return r
}
