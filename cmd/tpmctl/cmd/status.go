package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
)

func newStatusCmd() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "Show current TPM/FDE status",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			result, err := c.FdeStatus(ctx)
			if err != nil {
				return err
			}

			fmt.Printf("FDE status: %s\n", result.State)

			return nil
		},
	}
}
