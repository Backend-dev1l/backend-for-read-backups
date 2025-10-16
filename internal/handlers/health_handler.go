package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func RegisterHealthRoutes(r chi.Router) {
	r.Get("/livez", LivezHandler())
}

func LivezHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		response := HealthResponse{
			Status: "ok",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}
}
