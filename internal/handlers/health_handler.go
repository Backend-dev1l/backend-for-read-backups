package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"test-http/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	logger *slog.Logger
	cfg    *config.Config
	pool   *pgxpool.Pool
}

type HealthResponse struct {
	Status string `json:"status"`
}

func NewHealthHandler(pool *pgxpool.Pool, cfg *config.Config, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		cfg:    cfg,
		pool:   pool,
	}
}

func (h *HealthHandler) LivezHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Livez handler called (health)")
	response := HealthResponse{
		Status: "ok",
	}
	h.logger.Info("Livez handler preparing response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
