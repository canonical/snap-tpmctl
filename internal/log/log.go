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

func init() {
	levelvar := &slog.LevelVar{}
	th := newSimpleHandler(os.Stderr, levelvar)
	defaultLogger = loggerWithDynamicLevel{
		Logger:   slog.New(th),
		levelvar: levelvar,
	}
}

// WithContextLogger returns a context that contains a logger that writes to the provided io.Writer.
func WithContextLogger(ctx context.Context, w io.Writer) context.Context {
	levelvar := &slog.LevelVar{}
	th := newSimpleHandler(w, levelvar)
	logger := loggerWithDynamicLevel{
		Logger:   slog.New(th),
		levelvar: levelvar,
	}
	return context.WithValue(ctx, loggerKey, &logger)
}

// SetLoggerLevelInContext sets the logger level embedded into the context.
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
