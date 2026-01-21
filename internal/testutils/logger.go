package testutils

import (
	"bytes"
	"io"
	"log/slog"
)

// TestLogger writes logs to both a buffer and the test output.
func TestLogger(w io.Writer) {
	_, _ = TestLoggerWithBuffer(w)
}

// TestLoggerWithBuffer writes logs to both a buffer and the test output,
// and returns both the logger and the buffer for log assertions.
func TestLoggerWithBuffer(w io.Writer) (*slog.Logger, *bytes.Buffer) {
	var logs bytes.Buffer

	out := io.MultiWriter(&logs, w)
	h := slog.NewTextHandler(out, nil)

	return slog.New(h), &logs
}
