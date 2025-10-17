package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"test-http/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticsHandler struct {
	logger *slog.Logger
	cfg    *config.Config
	pool   *pgxpool.Pool
}

type ReadyzResponse struct {
	Status string `json:"status"`
}

func NewStatisticsHandler(pool *pgxpool.Pool, cfg *config.Config, logger *slog.Logger) *StatisticsHandler {
	return &StatisticsHandler{
		logger: logger,
		cfg:    cfg,
		pool:   pool,
	}
}

func (s *StatisticsHandler) ReadyzHandler(w http.ResponseWriter, r *http.Request) {
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
}
