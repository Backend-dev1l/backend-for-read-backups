package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"test-http/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReadyzResponse struct {
	Status string `json:"status"`
}

func ReadyzHandler(pool *pgxpool.Pool, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), cfg.TimeOuts.PerRequestTimeout)
		defer cancel()

		w.Header().Set("Content-Type", "application/json")

		if err := pool.Ping(ctx); err != nil {
			response := ReadyzResponse{
				Status: "error",
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(response)
			return
		}

		response := ReadyzResponse{
			Status: "ok",
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
