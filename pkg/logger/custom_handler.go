package logger

import (
	"context"

	"io"
	"log/slog"
)

type CustomHandler struct {
	handler slog.Handler
	opts    Options
}

func NewCustomHandler(w io.Writer, opts Options) slog.Handler {
	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: true,
	})

	return &CustomHandler{
		handler: jsonHandler,
		opts:    opts,
	}
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := []slog.Attr{
		slog.String("service", h.opts.Service),
		slog.String("env", h.opts.Env),
		slog.String("version", h.opts.Version),
	}

	// Добавляем trace_id если есть
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}

	// Маскируем PII в атрибутах записи
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, maskPII(a))
		return true
	})

	// Создаём новую запись с замаскированными атрибутами
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	newRecord.AddAttrs(attrs...)

	return h.handler.Handle(ctx, newRecord)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	maskedAttrs := make([]slog.Attr, 0, len(attrs))
	for _, attr := range attrs {
		maskedAttrs = append(maskedAttrs, maskPII(attr))
	}

	return &CustomHandler{
		handler: h.handler.WithAttrs(maskedAttrs),
		opts:    h.opts,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithGroup(name),
		opts:    h.opts,
	}
}
