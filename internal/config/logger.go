package config

import (
	"log/slog"
	"os"
)

func SetupLogger(env string) *slog.Logger {
	var handler slog.Handler
	w := os.Stdout

	switch env {
	case "dev":
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case "prod":
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
