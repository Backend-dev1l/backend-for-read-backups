package middleware

import (
	"context"
	"net/http"

	"test-http/pkg/logger"

	"github.com/google/uuid"
)

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.NewString()

		w.Header().Set("X-Trace-Id", traceID)

		ctx := context.WithValue(r.Context(), logger.TraceIDKey, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
