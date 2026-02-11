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

func main() {
	a := cmd.New(os.Args)
	os.Exit(run(context.Background(), a))
}

func run(ctx context.Context, a app) int {
	if err := a.Run(ctx); err != nil {
		log.Error(ctx, "the error is: %v", err)
		return 1
	}

	return 0
}
