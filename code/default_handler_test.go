package code

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func Test_defaultHandler(t *testing.T) {
	defaultLogger := slog.Default()
	defer func() {
		slog.SetDefault(defaultLogger)
	}()

	{
		// range:textHandler
		ctx := context.Background()
		h := slog.NewTextHandler(os.Stdout, nil)
		logger := slog.New(h)
		logger.InfoContext(
			ctx, "start processing",
			slog.Bool("verbose", true),
		)
		// range.end
	}

	{
		// range:jsonHandler
		ctx := context.Background()
		h := slog.NewJSONHandler(os.Stdout, nil)
		logger := slog.New(h)
		logger.InfoContext(
			ctx, "start processing",
			slog.Bool("verbose", true),
		)
		// range.end
	}

	{
		// range:textHandlerWithHandlerOptions
		ctx := context.Background()
		h := slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				// 呼び出し元コードの出力
				AddSource: true,
				// 出力するログレベル
				Level: slog.LevelDebug,
				// 属性の置き換え・削除など
				ReplaceAttr: func(
					groups []string, a slog.Attr,
				) slog.Attr {
					a.Key += "!"
					a.Value =
						slog.StringValue(a.Value.String() + "?")
					return a
				},
			},
		)
		logger := slog.New(h)
		logger.DebugContext(
			ctx, "start processing",
			slog.Bool("verbose", true),
		)
		// range.end
	}
}
