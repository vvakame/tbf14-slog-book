package code

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func Test_group(t *testing.T) {
	defaultLogger := slog.Default()
	defer func() {
		slog.SetDefault(defaultLogger)
	}()
	h := slog.NewJSONHandler(os.Stdout, nil)

	logger := slog.New(h)
	slog.SetDefault(logger)

	ctx := context.Background()

	{
		// range:json
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger.InfoCtx(
			ctx, "print group value",
			slog.Group(
				"group",
				slog.String("foo", "bar"),
				slog.Int("fizz", 4),
			),
		)
		// range.end
	}
	{
		// range:text
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		logger.InfoCtx(
			ctx, "print group value",
			slog.Group(
				"group",
				slog.String("foo", "bar"),
				slog.Int("fizz", 4),
			),
		)
		// range.end
	}
}
