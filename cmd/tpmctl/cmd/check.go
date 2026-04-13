// Package cmd implements the cli for exposing the cli commands snap-tpmctl supports
package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (a App) newCheckCmd() *cli.Command {
	return &cli.Command{
		Name:    "check-recovery-key",
		Usage:   "Check recovery key",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			key, err := a.tui.ReadRecoveryKey()
			if err != nil {
				return err
			}

			stop := a.tui.Spin("Checking recovery key...")
			defer stop()

			ok, err := a.tpm.CheckKey(ctx, key)
			if err != nil {
				return err
			}
			stop()

			msg := "Recovery key does not work"
			if ok {
				msg = "Recovery key works"
			}

			fmt.Fprintln(a.tui.Writer(), msg)

			return nil
		},
	}
}
