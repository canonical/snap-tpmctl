package log

import (
	"context"
	"log/slog"
)

// Debug outputs messages with the level [slog.LevelDebug] using the
// configured logging handler in the context.
func Debug(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelDebug, msg, args...)
}

// Info outputs messages with the level [slog.LevelInfo] using the
// configured logging handler in the context.
func Info(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, args...)
}

// Warn outputs messages with the level [slog.LevelWarn] using the
// configured logging handler in the context.
func Warn(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelWarn, msg, args...)
}

// Error outputs messages with the level [slog.LevelError] using the
// configured logging handler in the context.
func Error(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelError, msg, args...)
}
