package tools

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "session_key" {
				a.Value = slog.StringValue("REDACTED")
			}
			return a
		},
	})

	logger = slog.New(handler)
}
