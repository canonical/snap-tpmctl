package cmd

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/log"
)

func newStatusCmd() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "Show TPM status",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return status(ctx)
		},
	}
}

func status(ctx context.Context) error {
	log.Debug(ctx, "Retrieve status")

	return errors.New("TODO: implement the status API when lands on snapd")
}
