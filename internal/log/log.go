// package log provides logging utilities for the application.
// It supports redirecting logs via a context to a specific logger, like test ones.
package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type loggerWithDynamicLevel struct {
	*slog.Logger
	levelvar *slog.LevelVar
}

// defaultLogger is the default logger instance used when no logger is found in the context.
var defaultLogger loggerWithDynamicLevel

func new(w io.Writer) loggerWithDynamicLevel {
	levelvar := &slog.LevelVar{}
	levelvar.Set(slog.LevelWarn) // Set default log level to Warning.
	th := newSimpleHandler(w, levelvar)
	return loggerWithDynamicLevel{
		Logger:   slog.New(th),
		levelvar: levelvar,
	}
}

// stderrWriter is a writer that always writes to the current os.Stderr.
// This allows tests to redirect os.Stderr and have the logger follow the redirection.
type stderrWriter struct{}

func (w stderrWriter) Write(p []byte) (n int, err error) {
	return os.Stderr.Write(p)
}

func init() {
	defaultLogger = new(stderrWriter{})
}

// WithLoggerInContext returns a context that embeds a logger that writes to the io.Writer.
func WithLoggerInContext(ctx context.Context, w io.Writer) context.Context {
	logger := new(w)
	return context.WithValue(ctx, loggerKey, &logger)
}

// SetLoggerLevelInContext sets the log level for the logger embedded into the context.
func SetLoggerLevelInContext(ctx context.Context, level slog.Level) {
	logger := loggerFromContext(ctx)
	logger.levelvar.Set(level)
}

type loggerKeyType string

const loggerKey loggerKeyType = "logger"

// log is a helper function that retrieves the logger from the context and logs the message.
func log(ctx context.Context, level slog.Level, msg string, args ...any) {
	msg = fmt.Sprintf(msg, args...)
	logger := loggerFromContext(ctx)
	logger.Log(ctx, level, msg)
}

func loggerFromContext(ctx context.Context) *loggerWithDynamicLevel {
	logger, ok := ctx.Value(loggerKey).(*loggerWithDynamicLevel)
	if !ok {
		// If no logger is set, fallback to the default logger.
		return &defaultLogger
	}
	return logger
}
