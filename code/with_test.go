package code

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func Test_with(t *testing.T) {
	ctx := context.Background()

	h := slog.NewJSONHandler(
		os.Stdout,
		nil,
	)

	{
		// range:with
		logger := slog.New(h)
		logger = logger.With(
			slog.String("key1", "value1"),
		)
		logger.InfoCtx(
			ctx, "emit with With",
			slog.String("key2", "value2"),
		)
		// range.end
	}
	{
		// range:withGroup
		logger := slog.New(h)
		logger = logger.WithGroup("child")
		logger.InfoCtx(
			ctx, "emit with WithGroup",
			slog.String("key3", "value3"),
		)
		// range.end
	}
	{
		// range:duplicatedAttrs
		logger := slog.New(h)
		logger = logger.With(
			slog.String("key", "value1"),
		)
		logger.InfoCtx(
			ctx, "emit... but duplicated keys",
			slog.String("key", "value2"),
			slog.String("key", "value3"),
		)
		// range.end
	}
}
