package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

func (a App) newStatusCmd() *cli.Command {
	return &cli.Command{
		Name:    "status",
		Usage:   "Show current TPM/FDE status",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			status, err := a.tpm.FdeStatus(ctx)
			if err != nil {
				return err
			}

			fmt.Printf("The FDE system is %s\n", strings.ToUpper(status))

			return nil
		},
	}
}
