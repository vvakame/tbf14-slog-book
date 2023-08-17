package code

import (
	"context"
	"log/slog"
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
	logger.InfoContext(ctx, "emit source")
	// range.end
}
