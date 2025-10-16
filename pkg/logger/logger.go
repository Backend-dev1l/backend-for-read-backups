package logger

import (
	"log/slog"
	"os"
	"test-http/internal/config"
)

type ctxKey string

const TraceIDKey ctxKey = "trace_id"

type Options struct {
	Level   slog.Leveler
	Service string
	Env     string
	Version string
}

func New(opts Options) *slog.Logger {
	handler := NewCustomHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
func InitLogger(cfg *config.Config) (*slog.Logger, error) {
	level := parseLogLevel(cfg.Logger.Level)

	opts := Options{
		Level:   level,
		Service: cfg.App.Service,
		Env:     cfg.App.Env,
		Version: cfg.App.Version,
	}

	handler := NewCustomHandler(os.Stdout, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger, nil
}

func parseLogLevel(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
