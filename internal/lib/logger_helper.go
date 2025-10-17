package lib

import (
	"context"
	"log/slog"

	"test-http/pkg/logger"
)

func LogDebug(ctx context.Context, log *slog.Logger, operation, message string, attrs ...any) {
	traceID := getTraceID(ctx)
	logAttrs := []any{
		slog.String("operation", operation),
		slog.String("trace_id", traceID),
	}
	logAttrs = append(logAttrs, attrs...)
	log.DebugContext(ctx, message, logAttrs...)
}

func LogInfo(ctx context.Context, log *slog.Logger, operation, message string, attrs ...any) {
	traceID := getTraceID(ctx)
	logAttrs := []any{
		slog.String("operation", operation),
		slog.String("trace_id", traceID),
	}
	logAttrs = append(logAttrs, attrs...)
	log.InfoContext(ctx, message, logAttrs...)
}

func LogError(ctx context.Context, log *slog.Logger, operation, query, message string, err error, attrs ...any) {
	traceID := getTraceID(ctx)
	logAttrs := []any{
		slog.String("operation", operation),
		slog.String("query", query),
		slog.String("error", err.Error()),
		slog.String("trace_id", traceID),
	}
	logAttrs = append(logAttrs, attrs...)
	log.ErrorContext(ctx, message, logAttrs...)
}

func getTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(logger.TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
