package tools

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})

	Logger = slog.New(handler)
}
