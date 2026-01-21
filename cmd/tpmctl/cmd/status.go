package cmd

import (
	"context"
	"fmt"
	"strings"

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

			result, err := c.FdeStatus(ctx)
			if err != nil {
				return err
			}

			fmt.Printf("The FDE system is %s\n", strings.ToUpper(result.Status))

			return nil
		},
	}
}
