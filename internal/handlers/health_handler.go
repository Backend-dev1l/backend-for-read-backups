package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
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
