package logger

import (
	"log/slog"
	"strings"
)

var piiKeys = map[string]bool{
	"authorization": true,
	"set-cookie":    true,
	"x-api-key":     true,
}

func maskPII(attr slog.Attr) slog.Attr {
	key := strings.ToLower(attr.Key)

	if piiKeys[key] {
		return slog.String(attr.Key, "***REDACTED***")
	}

	if key == "email" {
		if val, ok := attr.Value.Any().(string); ok && val != "" {
			return slog.String(attr.Key, maskEmail(val))
		}
	}

	if key == "phone" {
		if val, ok := attr.Value.Any().(string); ok && val != "" {
			return slog.String(attr.Key, maskPhone(val))
		}
	}

	return attr
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	user := parts[0]
	if len(user) <= 2 {
		return "***@" + parts[1]
	}
	return user[:1] + strings.Repeat("*", len(user)-2) + user[len(user)-1:] + "@" + parts[1]
}

func maskPhone(phone string) string {
	r := []rune(phone)
	if len(r) <= 4 {
		return strings.Repeat("*", len(r))
	}
	return strings.Repeat("*", len(r)-4) + string(r[len(r)-4:])
}
