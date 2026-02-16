package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

// version is set at build time via ldflags.
var version = "dev"

func newVersionCmd() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Usage:   "Print version",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cli.DefaultPrintVersion(cmd.Root())

			return nil
		},
	}
}
