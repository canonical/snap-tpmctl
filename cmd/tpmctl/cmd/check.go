// Package cmd implements the cli for exposing the cli commands snap-tpmctl supports
package cmd

import (
	"context"
	"fmt"

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
			s := tpm.New()

			key, err := tui.ReadUserSecret("Enter recovery key: ")
			if err != nil {
				return err
			}

			if err := tpm.ValidateRecoveryKey(key); err != nil {
				return err
			}

			stop := tui.Spin("Checking recovery key...")
			defer stop()

			ok, err := s.CheckKey(ctx, key)
			if err != nil {
				return err
			}
			stop()

			msg := "Recovery key does not work"
			if ok {
				msg = "Recovery key works"
			}

			fmt.Println(msg)

			return nil
		},
	}
}
