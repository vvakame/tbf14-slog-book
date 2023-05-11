package code

import (
	"context"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func Test_attr(t *testing.T) {
	defaultLogger := slog.Default()
	defer func() {
		slog.SetDefault(defaultLogger)
	}()

	ctx := context.Background()

	var h slog.Handler
	var logger *slog.Logger

	// range:textHandler
	h = slog.NewTextHandler(os.Stdout, nil)
	logger = slog.New(h)
	logger.InfoCtx(
		ctx, "start processing",
		slog.Bool("verbose", true),
	)
	// range.end

	// range:jsonHandler
	h = slog.NewJSONHandler(os.Stdout, nil)
	logger = slog.New(h)
	logger.InfoCtx(
		ctx, "start processing",
		slog.Bool("verbose", true),
	)
	// range.end

	// range:textHandlerWithHandlerOptions
	h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// 呼び出し元コードの出力
		AddSource: true,
		// 出力するログレベル
		Level: slog.LevelDebug,
		// 属性の置き換え・削除など
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			a.Key += "!"
			a.Value = slog.StringValue(a.Value.String() + "?")
			return a
		},
	})
	logger = slog.New(h)
	logger.DebugCtx(
		ctx, "start processing",
		slog.Bool("verbose", true),
	)
	// range.end
}
