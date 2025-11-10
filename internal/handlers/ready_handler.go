package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"test-http/internal/config"
	"test-http/pkg/helper"

	"test-http/pkg/fault"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadyHandler struct {
	pool   *pgxpool.Pool
	cfg    *config.Config
	logger *slog.Logger
}

func NewReadyHandler(pool *pgxpool.Pool, cfg *config.Config, logger *slog.Logger) *ReadyHandler {
	return &ReadyHandler{
		pool:   pool,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *ReadyHandler) ReadyzHandler(w http.ResponseWriter, r *http.Request) error {
	s.logger.Info("Readyz handler called")

	ctx, cancel := context.WithTimeout(r.Context(), s.cfg.TimeOuts.PerRequestTimeout)
	defer cancel()

	if err := s.pool.Ping(ctx); err != nil {
		s.logger.Error("Database ping failed", "err", err)
		render.Status(r, http.StatusInternalServerError)
		return helper.HTTPError(w, fault.UnhandledError.Err())

	}
	render.Status(r, http.StatusOK)

	return nil
}
