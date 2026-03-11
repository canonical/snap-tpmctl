package testutils

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
)

// TestLoggerWithBuffer is a helper function that creates a logger that writes to a buffer and the test output.
// It returns a context with the logger embedded, and capture the logs in buffer.
func TestLoggerWithBuffer(t *testing.T) (context.Context, *bytes.Buffer) {
	t.Helper()

	var logs bytes.Buffer
	w := io.MultiWriter(&logs, t.Output())
	return log.WithLoggerInContext(t.Context(), w), &logs
}

// ContextLoggerWithDebug is a helper function that creates a logger that writes to the test output with debug level.
// It returns a context with the logger embedded.
func ContextLoggerWithDebug(t *testing.T) context.Context {
	t.Helper()

	ctx := log.WithLoggerInContext(t.Context(), t.Output())
	log.SetLoggerLevelInContext(ctx, slog.LevelDebug)

	return ctx
}
