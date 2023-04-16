package code

import (
	"context"
	"github.com/vvakame/tbf14-slog-book/code/customhandler"
	"github.com/vvakame/tbf14-slog-book/code/wrapslog"
	"testing"

	"github.com/galecore/xslog/xtesting"
	"golang.org/x/exp/slog"
)

func Test_xslog_xtesting(t *testing.T) {
	var h slog.Handler = xtesting.NewTestingHandler(t)
	h = h.WithAttrs([]slog.Attr{slog.String("id", "b")})
	h = h.WithGroup("test")
	logger := slog.NewLogLogger(h, slog.LevelDebug)
	logger.Println("id", "a")
}

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

func TestJSONHandler_sameKey(t *testing.T) {
	var h slog.Handler = slog.
		HandlerOptions{Level: slog.LevelDebug}.
		NewJSONHandler(writerFunc(func(p []byte) (n int, err error) {
			t.Log(string(p))
			return len(p), nil
		}))
	h = h.WithAttrs([]slog.Attr{slog.String("id", "a")})

	logger := slog.New(h)

	logger.Debug("test", "id", "b", "id", "c")
}

func Test_wrapped(t *testing.T) {
	ctx := context.Background()

	var h slog.Handler = slog.
		HandlerOptions{AddSource: true, Level: slog.LevelDebug}.
		NewJSONHandler(writerFunc(func(p []byte) (n int, err error) {
			t.Log(string(p))
			return len(p), nil
		}))
	h = customhandler.New(h)

	logger := slog.New(h)

	wrapslog.LogAttrs(ctx, logger, slog.LevelInfo, "test")
}
