package api

import (
	"log/slog"
	"test-http/internal/config"
	"test-http/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r chi.Router, pool *pgxpool.Pool, cfg *config.Config, logger *slog.Logger) {
	statisticsHandler := handlers.NewStatisticsHandler(pool, cfg, logger)
	healthHandler := handlers.NewHealthHandler(pool, cfg, logger)

	r.Route("/api/v1", func(g chi.Router) {
		g.Get("/livez", healthHandler.LivezHandler)
		g.Get("/readyz", statisticsHandler.ReadyzHandler)
	})
}
