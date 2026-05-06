// Package main is the entry point for snap-tpmctl.
package main

import (
	"context"
	"os"

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	"github.com/canonical/snap-tpmctl/internal/log"
)

type app interface {
	Run(ctx context.Context) error
}

var mainApp = cmd.New()

func main() {
	os.Exit(run(context.Background(), mainApp))
}

func run(ctx context.Context, a app) int {
	if err := a.Run(ctx); err != nil {
		log.Error(ctx, "%v", err)
		return 1
	}

	return 0
}
