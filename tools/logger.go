package tools

import (
	"fmt"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			fmt.Printf("groups: %v, key: %s, value: %v\n", groups, a.Key, a.Value)
			if a.Key == "session_key" {
				a.Value = slog.StringValue("REDACTED")
			}
			return a
		},
	})

	Logger = slog.New(handler)
}
