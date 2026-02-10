package log_test

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"snap-tpmctl/internal/log"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestWithLoggerInContext(t *testing.T) {
	t.Parallel()
	is := is.NewRelaxed(t)

	ctx := context.Background()
	msg := "test message"

	var writer bytes.Buffer
	ctx = log.WithLoggerInContext(ctx, &writer)

	logger := log.LoggerFromContext(ctx)
	is.True(logger != &log.DefaultLogger) // a logger is embedded into the context.

	logger.Log(context.Background(), slog.LevelWarn, msg)
	got := writer.String()
	is.True(strings.Contains(got, msg)) // the message is logged to the provided writer.
}

func TestSetLoggerLevelInContext(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	var got bytes.Buffer
	ctx := log.WithLoggerInContext(context.Background(), &got)
	logger := log.LoggerFromContext(ctx)

	logger.Log(context.Background(), slog.LevelInfo, "info message")
	is.True(!strings.Contains(got.String(), "info message")) // info message is not logged by default.

	logger.Log(context.Background(), slog.LevelWarn, "warning message")
	is.True(strings.Contains(got.String(), "warning message")) // warning messages are logged by default.

	log.SetLoggerLevelInContext(ctx, slog.LevelInfo)

	logger.Log(context.Background(), slog.LevelInfo, "info message")
	is.True(strings.Contains(got.String(), "info message")) // info message is logged after changing the log level.
}

func TestAllLogLevels(t *testing.T) {
	t.Parallel()

	type logLevelFunc = func(ctx context.Context, msg string, args ...any)

	tests := map[string]struct {
		fn logLevelFunc

		want string
	}{
		"debug": {
			fn: log.Debug,

			want: "DEBUG msg with args: 42\n",
		},
		"info": {
			fn: log.Info,

			want: "INFO msg with args: 42\n",
		},
		"warning": {
			fn: log.Warn,

			want: "WARN msg with args: 42\n",
		},
		"error": {
			fn: log.Error,

			want: "ERROR msg with args: 42\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			var logs bytes.Buffer
			out := io.MultiWriter(&logs, t.Output())
			ctx := log.WithLoggerInContext(context.Background(), out)
			log.SetLoggerLevelInContext(ctx, slog.LevelDebug)

			tc.fn(ctx, "msg with args: %d", 42)
			is.Equal(logs.String(), tc.want)
		})
	}
}
