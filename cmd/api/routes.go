package api

import (
	"test-http/internal/config"
	"test-http/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(pool *pgxpool.Pool, cfg *config.Config) {

	r := chi.NewRouter()

	r.Route("api/v1", func(g chi.Router) {
		g.Get("/livez", handlers.LivezHandler())
		g.Get("/readyz", handlers.ReadyzHandler(pool, cfg))
	})
}
