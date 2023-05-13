package code

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func Test_levelAndLeveler(t *testing.T) {
	ctx := context.Background()

	removeTimeKey := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	}

	// range:levelVar
	var levelVar slog.LevelVar

	h := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: &levelVar,
			// 出力の見やすさのため time を削除する関数（実装略）
			ReplaceAttr: removeTimeKey,
		},
	)
	logger := slog.New(h)

	ls := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	}
	for _, l := range ls {
		levelVar.Set(l)
		logger.WarnCtx(ctx, "warning!", "l", l)
	}
	// range.end
}
