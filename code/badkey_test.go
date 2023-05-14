package code

import (
	"context"
	"os"
	"testing"

	"golang.org/x/exp/slog"
)

func Test_badkey(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// range:badkey
	logger.InfoCtx(ctx, "test", "value(not-key)")
	// range.end
}
