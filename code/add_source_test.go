package code

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func Test_addSource(t *testing.T) {
	ctx := context.Background()

	// range:source
	h := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
		},
	)

	logger := slog.New(h)
	logger.InfoCtx(ctx, "emit source")
	// range.end
}
