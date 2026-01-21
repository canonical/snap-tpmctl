package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tpm"
)

func newStatusCmd() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "Show current TPM/FDE status",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()

			status, err := tpm.FdeStatus(ctx, c)
			if err != nil {
				return err
			}

			fmt.Printf("The FDE system is %s\n", strings.ToUpper(status))

			return nil
		},
	}
}
