// Package cmd implements the cli for exposing the cli commands snap-tpmctl supports
package cmd

import (
	"context"
	"fmt"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

func newCheckCmd() *cli.Command {
	return &cli.Command{
		Name:    "check-recovery-key",
		Usage:   "Check recovery key",
		Suggest: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.New()

			key, err := tui.ReadUserSecret("Enter recovery key: ")
			if err != nil {
				return err
			}

			if err := tpm.ValidateRecoveryKey(key); err != nil {
				return err
			}

			ok, err := tui.WithSpinnerResult("Checking recovery key...", func() (bool, error) {
				return tpm.CheckKey(ctx, c, key)
			})
			if err != nil {
				return err
			}

			// TODO: print better messages
			msg := "Recovery key does not work"
			if ok {
				msg = "Recovery key works"
			}

			fmt.Println(msg)

			return nil
		},
	}
}
