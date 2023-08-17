package code

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/vvakame/tbf14-slog-book/code/customhandler"
	"github.com/vvakame/tbf14-slog-book/code/wrapslog"
)

func Test_slog_testing(t *testing.T) {
	var h slog.Handler = slog.NewTextHandler(os.Stdout, nil)
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
		NewJSONHandler(writerFunc(func(p []byte) (n int, err error) {
			t.Log(string(p))
			return len(p), nil
		}), &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	h = h.WithAttrs([]slog.Attr{slog.String("id", "a")})

	logger := slog.New(h)

	logger.Debug("test", "id", "b", "id", "c")
}

func Test_wrapped(t *testing.T) {
	ctx := context.Background()

	var h slog.Handler = slog.
		NewJSONHandler(writerFunc(func(p []byte) (n int, err error) {
			t.Log(string(p))
			return len(p), nil
		}), &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	h = customhandler.New(h)

	logger := slog.New(h)

	wrapslog.LogAttrs(ctx, logger, slog.LevelInfo, "test")
}
