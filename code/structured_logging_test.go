package code

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func Test_structuredLogging(t *testing.T) {
	defaultLogger := slog.Default()
	defer func() {
		slog.SetDefault(defaultLogger)
	}()
	h := slog.NewJSONHandler(os.Stdout, nil)

	// range:newLogger
	logger := slog.New(h)
	slog.SetDefault(logger)
	// range.end

	ctx := context.Background()

	userID := "a1b2c3"

	// range:example1
	slog.InfoContext(
		ctx, "start processing",
		slog.String("userID", userID),
	)
	// range.end

	// range:example2
	slog.InfoContext(
		ctx, "start processing",
		"userID", userID,
	)
	// range.end

	// range:example3
	slog.LogAttrs(
		ctx,
		slog.LevelInfo,
		"start processing",
		slog.String("userID", userID),
	)
	// range.end
}
