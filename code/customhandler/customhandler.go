package customhandler

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
)

func New(base slog.Handler) slog.Handler {
	return &handler{base}
}

type handler struct {
	base slog.Handler
}

func (h *handler) clone() *handler {
	return &handler{
		base: h.base,
	}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	if record.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{record.PC})
		f, more := fs.Next()
		fmt.Println(f, more)

		record.AddAttrs(
			slog.Group(
				"sourceLocation",
				slog.String("file", f.File),
				slog.String("line", strconv.Itoa(f.Line)),
				slog.String("function", f.Function),
			),
		)
	}

	return h.base.Handle(ctx, record)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h = h.clone()
	h.base = h.base.WithAttrs(attrs)

	return h
}

func (h *handler) WithGroup(name string) slog.Handler {
	h = h.clone()
	h.base = h.base.WithGroup(name)

	return h
}
