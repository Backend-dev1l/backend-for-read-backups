package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"test-http/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testService = "test-service"
	testEnv     = "test"
	testVersion = "1.0.0"
)

type fieldExpectation struct {
	key  string
	want any
}

// TestSecretsNotLeakedInLogs проверяет, что секретные данные не попадают в логи
func TestSecretsNotLeakedInLogs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		level     slog.Level
		logFn     func(*slog.Logger)
		expected  []fieldExpectation
		forbidden []string
	}{
		{
			name:  "authorization header should be masked",
			level: slog.LevelInfo,
			logFn: func(log *slog.Logger) {
				log.Info("request received",
					slog.String("authorization", "Bearer secret-token-12345"),
					slog.String("method", "GET"),
				)
			},
			expected: []fieldExpectation{
				{key: "msg", want: "request received"},
				{key: "authorization", want: "***REDACTED***"},
				{key: "method", want: "GET"},
			},
			forbidden: []string{"Bearer", "secret-token-12345"},
		},
		{
			name:  "api key should be masked",
			level: slog.LevelInfo,
			logFn: func(log *slog.Logger) {
				log.Info("api call",
					slog.String("x-api-key", "sk-1234567890abcdef"),
					slog.String("endpoint", "/api/users"),
				)
			},
			expected: []fieldExpectation{
				{key: "msg", want: "api call"},
				{key: "x-api-key", want: "***REDACTED***"},
				{key: "endpoint", want: "/api/users"},
			},
			forbidden: []string{"sk-1234567890abcdef"},
		},
		{
			name:  "set-cookie should be masked",
			level: slog.LevelInfo,
			logFn: func(log *slog.Logger) {
				log.Info("response headers",
					slog.String("set-cookie", "session=abc123; HttpOnly; Secure"),
					slog.Int("status", 200),
				)
			},
			expected: []fieldExpectation{
				{key: "msg", want: "response headers"},
				{key: "set-cookie", want: "***REDACTED***"},
				{key: "status", want: float64(200)},
			},
			forbidden: []string{"session=abc123", "HttpOnly"},
		},
		{
			name:  "email should be partially masked",
			level: slog.LevelInfo,
			logFn: func(log *slog.Logger) {
				log.Info("user registered",
					slog.String("email", "user@example.com"),
					slog.String("user_id", "123"),
				)
			},
			expected: []fieldExpectation{
				{key: "msg", want: "user registered"},
				{key: "email", want: "u**r@example.com"},
				{key: "user_id", want: "123"},
			},
			forbidden: []string{"user@example.com"},
		},
		{
			name:  "phone should be partially masked",
			level: slog.LevelInfo,
			logFn: func(log *slog.Logger) {
				log.Info("phone verification",
					slog.String("phone", "+79991234567"),
					slog.String("status", "verified"),
				)
			},
			expected: []fieldExpectation{
				{key: "msg", want: "phone verification"},
				{key: "phone", want: "********4567"},
				{key: "status", want: "verified"},
			},
			forbidden: []string{"+79991234567"},
		},
		{
			name:  "multiple sensitive fields",
			level: slog.LevelInfo,
			logFn: func(log *slog.Logger) {
				log.Info("authentication attempt",
					slog.String("authorization", "Bearer token123"),
					slog.String("email", "admin@test.com"),
					slog.String("x-api-key", "key-secret"),
					slog.String("ip", "192.168.1.1"),
				)
			},
			expected: []fieldExpectation{
				{key: "msg", want: "authentication attempt"},
				{key: "authorization", want: "***REDACTED***"},
				{key: "email", want: "a***n@test.com"},
				{key: "x-api-key", want: "***REDACTED***"},
				{key: "ip", want: "192.168.1.1"},
			},
			forbidden: []string{"Bearer token123", "admin@test.com", "key-secret"},
		},
		{
			name:  "case insensitive masking",
			level: slog.LevelInfo,
			logFn: func(log *slog.Logger) {
				log.Info("headers check",
					slog.String("Authorization", "Bearer UPPERCASE-TOKEN"),
					slog.String("X-API-KEY", "UPPERCASE-KEY"),
				)
			},
			expected: []fieldExpectation{
				{key: "msg", want: "headers check"},
				{key: "Authorization", want: "***REDACTED***"},
				{key: "X-API-KEY", want: "***REDACTED***"},
			},
			forbidden: []string{"UPPERCASE-TOKEN", "UPPERCASE-KEY"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			raw, entry := captureLogEntry(t, tt.level, tt.logFn)

			assertCommonMetadata(t, entry, tt.name)

			for _, exp := range tt.expected {
				assertFieldValue(t, entry, exp.key, exp.want, tt.name)
			}

			for _, secret := range tt.forbidden {
				assertSecretNotPresent(t, raw, secret, tt.name)
			}
		})
	}
}

// TestSecretsNotLeakedWithContext проверяет маскирование с использованием контекста
func TestSecretsNotLeakedWithContext(t *testing.T) {
	t.Parallel()

	raw, entry := captureLogEntry(t, slog.LevelInfo, func(log *slog.Logger) {
		ctx := context.WithValue(context.Background(), logger.TraceIDKey, "trace-123")
		log.InfoContext(ctx, "request with secrets",
			slog.String("authorization", "Bearer secret-token"),
			slog.String("email", "test@example.com"),
		)
	})

	assertCommonMetadata(t, entry, "context log")
	assertFieldValue(t, entry, "trace_id", "trace-123", "context log")
	assertFieldValue(t, entry, "authorization", "***REDACTED***", "context log")
	assertFieldValue(t, entry, "email", "t**t@example.com", "context log")
	assertSecretNotPresent(t, raw, "secret-token", "context log")
}

// TestSecretsNotLeakedWithAttrs проверяет маскирование при использовании WithAttrs
func TestSecretsNotLeakedWithAttrs(t *testing.T) {
	t.Parallel()

	raw, entry := captureLogEntry(t, slog.LevelInfo, func(log *slog.Logger) {
		log.With(
			slog.String("authorization", "Bearer preset-token"),
			slog.String("user_id", "user-123"),
		).Info("operation completed",
			slog.String("status", "success"),
		)
	})

	assertCommonMetadata(t, entry, "with attrs")
	assertFieldValue(t, entry, "authorization", "***REDACTED***", "with attrs")
	assertFieldValue(t, entry, "user_id", "user-123", "with attrs")
	assertFieldValue(t, entry, "status", "success", "with attrs")
	assertSecretNotPresent(t, raw, "preset-token", "with attrs")
}

// TestSecretsNotLeakedInErrorLogs проверяет маскирование в логах ошибок
func TestSecretsNotLeakedInErrorLogs(t *testing.T) {
	t.Parallel()

	raw, entry := captureLogEntry(t, slog.LevelError, func(log *slog.Logger) {
		log.Error("authentication failed",
			slog.String("authorization", "Bearer failed-token"),
			slog.String("email", "hacker@evil.com"),
			slog.String("error", "invalid credentials"),
		)
	})

	assertCommonMetadata(t, entry, "error log")
	assertFieldValue(t, entry, "level", "ERROR", "error log")
	assertFieldValue(t, entry, "authorization", "***REDACTED***", "error log")
	assertFieldValue(t, entry, "email", "h****r@evil.com", "error log")
	assertFieldValue(t, entry, "error", "invalid credentials", "error log")
	assertSecretNotPresent(t, raw, "failed-token", "error log")
}

func captureLogEntry(t *testing.T, level slog.Level, logFn func(*slog.Logger)) (string, map[string]any) {
	t.Helper()

	buf, log := newBufferedLogger(t, level)
	logFn(log)

	raw := strings.TrimSpace(buf.String())
	require.NotEmpty(t, raw, "expected log output")

	lines := strings.Split(raw, "\n")
	entryLine := lines[len(lines)-1]

	var decoded map[string]any
	err := json.Unmarshal([]byte(entryLine), &decoded)
	require.NoErrorf(t, err, "failed to decode log output %q", entryLine)

	return raw, decoded
}

func newBufferedLogger(t *testing.T, level slog.Level) (*bytes.Buffer, *slog.Logger) {
	t.Helper()

	var buf bytes.Buffer
	opts := logger.Options{
		Level:   level,
		Service: testService,
		Env:     testEnv,
		Version: testVersion,
	}

	handler := logger.NewCustomHandler(&buf, opts)
	return &buf, slog.New(handler)
}

func assertCommonMetadata(t *testing.T, entry map[string]any, scenario string) {
	t.Helper()

	assertFieldValue(t, entry, "service", testService, scenario)
	assertFieldValue(t, entry, "env", testEnv, scenario)
	assertFieldValue(t, entry, "version", testVersion, scenario)
}

func assertFieldValue(t *testing.T, entry map[string]any, key string, want any, scenario string) {
	t.Helper()

	got, ok := entry[key]
	require.Truef(t, ok, "%s: expected field %q to be present. entry=%v", scenario, key, entry)

	switch wantTyped := want.(type) {
	case string:
		gotStr, ok := got.(string)
		require.Truef(t, ok, "%s: field %q expected string value, got %T", scenario, key, got)
		assert.Equalf(t, wantTyped, gotStr, "%s: field %q mismatch", scenario, key)
	case float64:
		gotFloat, ok := got.(float64)
		require.Truef(t, ok, "%s: field %q expected float64 value, got %T", scenario, key, got)
		assert.Equalf(t, wantTyped, gotFloat, "%s: field %q mismatch", scenario, key)
	default:
		assert.Equalf(t, want, got, "%s: field %q mismatch", scenario, key)
	}
}

func assertSecretNotPresent(t *testing.T, raw, secret, scenario string) {
	t.Helper()

	if secret == "" {
		return
	}

	assert.NotContainsf(t, raw, secret, "%s: SECRET LEAKED! log output: %s", scenario, raw)
}
