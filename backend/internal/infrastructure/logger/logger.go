package logger

import (
	"log/slog"
	"os"

	"github.com/personal-excalidraw/backend/internal/infrastructure/config"
)

// New creates a new structured logger based on configuration
func New(cfg *config.LoggerConfig) *slog.Logger {
	level := parseLevel(cfg.Level)

	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if cfg.Format == "json" {
		// JSON format for production (machine-readable)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Text format for development (human-readable)
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// parseLevel parses log level string to slog.Level
func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
