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


	s.logger.Info("Readyz handler called")
	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeOuts.PerRequestTimeout)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	if err := s.pool.Ping(ctx); err != nil {
		response := ReadyzResponse{
			Status: "error",
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(response)
		s.logger.Error("Failed to ping database", "error", err)
		return
	}

	response := ReadyzResponse{
		Status: "ok",
	}
	s.logger.Info("Database ping successful")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
