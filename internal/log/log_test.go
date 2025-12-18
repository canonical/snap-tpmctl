package log_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nalgeon/be"
	"snap-tpmctl/internal/log"
)

var supportedLevels = []log.Level{
	log.DebugLevel,
	log.InfoLevel,
	log.NoticeLevel,
	log.WarnLevel,
	log.ErrorLevel,
}

func TestLevelEnabled(t *testing.T) {
	// This can't be parallel.
	defaultLevel := log.GetLevel()
	t.Cleanup(func() {
		log.SetLevel(defaultLevel)
	})

	for _, level := range supportedLevels {
		t.Run(fmt.Sprintf("Set log level to %s", level), func(t *testing.T) {
			log.SetLevel(level)
			be.Equal(t, level, log.GetLevel())    // Mismatched log level
			be.True(t, log.IsLevelEnabled(level)) // Log level should be enabled
		})
	}
}

func callLogHandler(ctx context.Context, level log.Level, args ...any) {
	switch level {
	case log.ErrorLevel:
		log.Error(ctx, args...)
	case log.WarnLevel:
		log.Warning(ctx, args...)
	case log.NoticeLevel:
		log.Notice(ctx, args...)
	case log.InfoLevel:
		log.Info(ctx, args...)
	case log.DebugLevel:
		log.Debug(ctx, args...)
	}
}

func callLogHandlerf(ctx context.Context, level log.Level, format string, args ...any) {
	switch level {
	case log.ErrorLevel:
		log.Errorf(ctx, format, args...)
	case log.WarnLevel:
		log.Warningf(ctx, format, args...)
	case log.NoticeLevel:
		log.Noticef(ctx, format, args...)
	case log.InfoLevel:
		log.Infof(ctx, format, args...)
	case log.DebugLevel:
		log.Debugf(ctx, format, args...)
	}
}

func TestSetLevelHandler(t *testing.T) {
	defaultLevel := log.GetLevel()
	t.Cleanup(func() {
		log.SetLevel(defaultLevel)
		for _, level := range supportedLevels {
			log.SetLevelHandler(level, nil)
		}
	})

	for _, level := range supportedLevels {
		t.Run(fmt.Sprintf("Set log handler for %s", level), func(t *testing.T) {
			handlerCalled := false
			wantArgs := []any{true, 5.5, []string{"bar"}}
			wantCtx := context.TODO()
			log.SetLevelHandler(level, func(ctx context.Context, l log.Level, format string, args ...any) {
				handlerCalled = true
				be.Equal(t, wantCtx, ctx)                    // Mismatched log context
				be.Equal(t, level, l)                        // Mismatched log level
				be.Equal(t, fmt.Sprint(wantArgs...), format) // Mismatched log format
			})

			log.SetLevel(level)

			callLogHandler(wantCtx, level, wantArgs...)
			be.True(t, handlerCalled) // Expected handler to be called

			handlerCalled = false
			callLogHandler(wantCtx, level+1, wantArgs...)
			be.True(t, !handlerCalled) // Expected handler not to be called

			handlerCalled = false
			log.SetLevelHandler(level, nil)

			callLogHandler(wantCtx, level, wantArgs...)
			be.True(t, !handlerCalled) // Expected handler not to be called
		})
	}

	for _, level := range supportedLevels {
		t.Run(fmt.Sprintf("Set log handler for %s, using formatting", level), func(t *testing.T) {
			handlerCalled := false
			wantArgs := []any{true, 5.5, []string{"bar"}}
			wantFormat := "Bool is %v, float is %f, array is %v"
			wantCtx := context.TODO()
			log.SetLevelHandler(level, func(ctx context.Context, l log.Level, format string, args ...any) {
				handlerCalled = true
				be.Equal(t, wantCtx, ctx)       // Mismatched log context
				be.Equal(t, level, l)           // Mismatched log level
				be.Equal(t, wantFormat, format) // Mismatched log format
				be.Equal(t, wantArgs, args)     // Mismatched log args
			})

			handlerCalled = false
			callLogHandlerf(wantCtx, level, wantFormat, wantArgs...)

			handlerCalled = false
			callLogHandlerf(wantCtx, level+1, wantFormat, wantArgs...)
			be.True(t, !handlerCalled) // Expected handler not to be called

			handlerCalled = false
			log.SetLevelHandler(level, nil)

			callLogHandlerf(wantCtx, level, wantFormat, wantArgs...)
			be.True(t, !handlerCalled) // Expected handler not to be called
		})
	}

	log.SetLevelHandler(log.Level(99999999), nil)
}

func TestSetHandler(t *testing.T) {
	defaultLevel := log.GetLevel()
	t.Cleanup(func() {
		log.SetLevel(defaultLevel)
		log.SetHandler(nil)
	})

	handlerCalled := false
	wantLevel := log.Level(0)
	wantArgs := []any{true, 5.5, []string{"bar"}}
	wantCtx := context.TODO()

	log.SetHandler(func(ctx context.Context, l log.Level, format string, args ...any) {
		handlerCalled = true
		be.Equal(t, wantCtx, ctx)                    // Mismatched log context
		be.Equal(t, wantLevel, l)                    // Mismatched log level
		be.Equal(t, fmt.Sprint(wantArgs...), format) // Mismatched log format
	})
	for idx, level := range supportedLevels {
		t.Run(fmt.Sprintf("Set log handler, testing level %s", level), func(t *testing.T) {})

		wantLevel = level
		handlerCalled = false
		log.SetLevel(level)
		fmt.Println("Logging at level", level)

		callLogHandler(wantCtx, level, wantArgs...)
		be.True(t, handlerCalled) // Expected handler to be called

		handlerCalled = false
		nextLevel := level + 1
		if idx < len(supportedLevels)-1 {
			nextLevel = supportedLevels[idx+1]
		}
		log.SetLevel(nextLevel)
		callLogHandler(wantCtx, level, wantArgs...)
		be.True(t, !handlerCalled) // Expected handler not to be called
	}

	log.SetHandler(nil)
	for _, level := range supportedLevels {
		t.Run(fmt.Sprintf("Set log handler, ignoring level %s", level), func(t *testing.T) {})

		wantLevel = level
		handlerCalled = false
		log.SetLevel(level)
		callLogHandler(wantCtx, level, wantArgs...)
		be.True(t, !handlerCalled) // Expected handler not to be called
	}

	wantFormat := "Bool is %v, float is %f, array is %v"
	log.SetHandler(func(ctx context.Context, l log.Level, format string, args ...any) {
		handlerCalled = true
		be.Equal(t, wantCtx, ctx)       // Mismatched log context
		be.Equal(t, wantLevel, l)       // Mismatched log level
		be.Equal(t, wantFormat, format) // Mismatched log format
		be.Equal(t, wantArgs, args)     // Mismatched log args
	})
	for idx, level := range supportedLevels {
		t.Run(fmt.Sprintf("Set log handler, testing level %s", level), func(t *testing.T) {})

		wantLevel = level
		handlerCalled = false
		log.SetLevel(level)
		callLogHandlerf(wantCtx, level, wantFormat, wantArgs...)
		be.True(t, handlerCalled) // Expected handler to be called

		handlerCalled = false
		nextLevel := level + 1
		if idx < len(supportedLevels)-1 {
			nextLevel = supportedLevels[idx+1]
		}
		log.SetLevel(nextLevel)
		callLogHandlerf(wantCtx, level, wantFormat, wantArgs...)
		be.True(t, !handlerCalled) // Expected handler not to be called
	}

	log.SetHandler(nil)
	for _, level := range supportedLevels {
		t.Run(fmt.Sprintf("Set log handler, ignoring level %s", level), func(t *testing.T) {})

		wantLevel = level
		handlerCalled = false
		log.SetLevel(level)
		callLogHandlerf(wantCtx, level, wantFormat, wantArgs...)
		be.True(t, !handlerCalled) // Expected handler not to be called
	}
}
