package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/urfave/cli/v3"
)

func newStatusCmd() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "Show current TPM/FDE status",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.New()

			status, err := tpm.FdeStatus(ctx, c)
			if err != nil {
				return err
			}

			fmt.Printf("The FDE system is %s\n", strings.ToUpper(status))

			return nil
		},
	}
}
