package testutils

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/canonical/snap-tpmctl/internal/log"
)

// TestLoggerWithBuffer is a helper function that creates a logger that writes to a buffer and the test output.
// It returns a context with the logger embedded, and capture the logs in buffer.
func TestLoggerWithBuffer(t *testing.T) (context.Context, *bytes.Buffer) {
	t.Helper()

	var logs bytes.Buffer
	w := io.MultiWriter(&logs, t.Output())
	return log.WithLoggerInContext(context.Background(), w), &logs
}
