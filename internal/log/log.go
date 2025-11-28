// Package log is a temporary package until we forge our log structure.
package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"maps"
)

type (
	// Level is the log level for the logs.
	Level = slog.Level

	// Handler is the log handler function.
	Handler = func(_ context.Context, _ Level, format string, args ...any)
)

var logLevel slog.Level

const (
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = slog.LevelError
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = slog.LevelWarn
	// NoticeLevel level. Normal but significant conditions. Conditions that are not error conditions, but that may
	// require special handling. slog doesn't have a Notice level, so we use the average between Info and Warn.
	NoticeLevel = (slog.LevelInfo + slog.LevelWarn) / 2
	// InfoLevel level. General operational entries about what's going on inside the application.
	InfoLevel = slog.LevelInfo
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = slog.LevelDebug
)

func logFuncAdapter(slogFunc func(ctx context.Context, msg string, args ...any)) Handler {
	return func(ctx context.Context, _ Level, format string, args ...any) {
		slogFunc(ctx, fmt.Sprintf(format, args...))
	}
}

var allLevels = []slog.Level{
	DebugLevel,
	InfoLevel,
	NoticeLevel,
	WarnLevel,
	ErrorLevel,
}

var defaultHandlers = map[Level]Handler{
	DebugLevel: logFuncAdapter(slog.DebugContext),
	InfoLevel:  logFuncAdapter(slog.InfoContext),
	// slog doesn't have a Notice level, so in the default handler, we use Warn instead.
	NoticeLevel: logFuncAdapter(slog.WarnContext),
	WarnLevel:   logFuncAdapter(slog.WarnContext),
	ErrorLevel:  logFuncAdapter(slog.ErrorContext),
}
var handlers = maps.Clone(defaultHandlers)

// GetLevel gets the standard logger level.
func GetLevel() Level {
	return logLevel
}

// IsLevelEnabled checks if the log level is greater than the level param.
func IsLevelEnabled(level Level) bool {
	return isLevelEnabled(context.Background(), level)
}

func isLevelEnabled(context context.Context, level Level) bool {
	return slog.Default().Enabled(context, level)
}

// SetLevel sets the standard logger level.
func SetLevel(level Level) {
	logLevel = level
	slog.SetLogLoggerLevel(level)
}

// SetOutput sets the log output.
func SetOutput(out io.Writer) {
	slog.SetDefault(slog.New(NewSimpleHandler(out, GetLevel())))
}

func log(ctx context.Context, level Level, args ...any) {
	if !isLevelEnabled(ctx, level) {
		return
	}

	logf(ctx, level, fmt.Sprint(args...))
}

func logf(ctx context.Context, level Level, format string, args ...any) {
	if !isLevelEnabled(ctx, level) {
		return
	}

	handler := handlers[level]
	handler(ctx, level, format, args...)
}

// Debug outputs messages with the level [DebugLevel] (when that is enabled) using the
// configured logging handler.
func Debug(ctx context.Context, args ...any) {
	log(ctx, DebugLevel, args...)
}

// Debugf outputs messages with the level [DebugLevel] (when that is enabled) using the
// configured logging handler.
func Debugf(ctx context.Context, format string, args ...any) {
	logf(ctx, DebugLevel, format, args...)
}

// Info outputs messages with the level [InfoLevel] (when that is enabled) using the
// configured logging handler.
func Info(ctx context.Context, args ...any) {
	log(ctx, InfoLevel, args...)
}

// Infof outputs messages with the level [InfoLevel] (when that is enabled) using the
// configured logging handler.
func Infof(ctx context.Context, format string, args ...any) {
	logf(ctx, InfoLevel, format, args...)
}

// Notice outputs messages with the level [NoticeLevel] (when that is enabled) using the
// configured logging handler.
func Notice(ctx context.Context, args ...any) {
	log(ctx, NoticeLevel, args...)
}

// Noticef outputs messages with the level [NoticeLevel] (when that is enabled) using the
// configured logging handler.
func Noticef(ctx context.Context, format string, args ...any) {
	logf(ctx, NoticeLevel, format, args...)
}

// Warning outputs messages with the level [WarningLevel] (when that is enabled) using the
// configured logging handler.
func Warning(ctx context.Context, args ...any) {
	log(ctx, WarnLevel, args...)
}

// Warningf outputs messages with the level [WarningLevel] (when that is enabled) using the
// configured logging handler.
func Warningf(ctx context.Context, format string, args ...any) {
	logf(ctx, WarnLevel, format, args...)
}

// Error outputs messages with the level [ErrorLevel] (when that is enabled) using the
// configured logging handler.
func Error(ctx context.Context, args ...any) {
	log(ctx, ErrorLevel, args...)
}

// Errorf outputs messages with the level [ErrorLevel] (when that is enabled) using the
// configured logging handler.
func Errorf(ctx context.Context, format string, args ...any) {
	logf(ctx, ErrorLevel, format, args...)
}
