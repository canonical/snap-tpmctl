package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

// simpleHandler writes logs in the format: <level> <message>.
type simpleHandler struct {
	level slog.Leveler
	w     io.Writer
}

// newSimpleHandler creates a new simpleHandler that writes to the provided io.Writer.
func newSimpleHandler(w io.Writer, leveler slog.Leveler) slog.Handler {
	return &simpleHandler{
		level: leveler,
		w:     w,
	}
}

// Enabled checks if the handler is enabled for the given log level.
func (h *simpleHandler) Enabled(ctx context.Context, l slog.Level) bool {
	minLevel := slog.LevelWarn
	if h.level != nil {
		minLevel = h.level.Level()
	}
	return l >= minLevel
}

// Handle formats the record to include the level and then the message.
func (h *simpleHandler) Handle(ctx context.Context, r slog.Record) error {
	_, err := fmt.Fprintf(h.w, "%s %s\n", r.Level.String(), r.Message)
	return err
}

// WithAttrs is not implemented in our simpleHandler.
func (h *simpleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("WithAttrs is not implemented for simpleHandler")
}

// WithGroup is not implemented in our simpleHandler.
func (h *simpleHandler) WithGroup(name string) slog.Handler {
	panic("WithGroup is not implemented for simpleHandler")
}
