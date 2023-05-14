package code

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func Test_defaultKeys(t *testing.T) {
	ctx := context.Background()

	// range:keys
	h := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			ReplaceAttr: func(
				groups []string, a slog.Attr,
			) slog.Attr {
				switch a.Key {
				case slog.TimeKey:
					a.Key = "t"
				case slog.LevelKey:
					a.Key = "l"
				case slog.SourceKey:
					a.Key = "s"
				case slog.MessageKey:
					a.Key = "m"
				}
				return a
			},
		},
	)

	logger := slog.New(h)
	logger.InfoCtx(ctx, "rename keys")
	// range.end
}
