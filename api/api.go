package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/PatrikOlin/gordle/handlers"
	m "github.com/PatrikOlin/gordle/middleware"
)

func GetRouter(log *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
		m.SetLogger(log),
		// cors.AllowAll().Handler,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://*", "https://*"},
			AllowCredentials: true,
		}),
	)

	r.Get("/rules", handlers.GetRules)
	r.Post("/word", handlers.GuessWord)
	r.Get("/session", handlers.GetSession)
	// r.Get("/session/{id}", handlers.GetSession)
	// r.Get("/word", handlers.GetWord)
	// r.Get("/state", handlers.GetGameState)

	return r
}
