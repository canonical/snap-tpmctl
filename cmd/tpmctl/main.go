// Package main is the entry point for snap-tpmctl.
package main

import (
	"context"
	"log/slog"
	"os"

	"snap-tpmctl/cmd/tpmctl/cmd"
)

type app interface {
	Run() error
}

func main() {
	a := cmd.New(os.Args)
	os.Exit(run(context.Background(), a))
}

func run(ctx context.Context, a app) int {
	if err := a.Run(); err != nil {
		logError(ctx, err.Error())
		return 1
	}

	return 0
}

type loggerKeyType string

const loggerKey loggerKeyType = "logger"

func logError(ctx context.Context, msg string, args ...any) {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		// If no logger is set, fallback to the default logger.
		logger = slog.Default()
	}
	logger.Error(msg, args...)
}
