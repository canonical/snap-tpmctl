package cmd

import (
	"context"
	"fmt"

	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

//nolint:dupl // newReplacePassphraseCmd and newReplacePINCmd have similar behaviour
func newReplacePassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-passphrase",
		Usage: "Replace encryption passphrase",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			s := tpm.New()

			oldPassphrase, err := tui.ReadUserSecret("Enter current passphrase: ")
			if err != nil {
				return err
			}

			newPassphrase, err := tui.ReadUserSecret("Enter new passphrase: ")
			if err != nil {
				return err
			}

			confirmPassphrase, err := tui.ReadUserSecret("Confirm new passphrase: ")
			if err != nil {
				return err
			}

			if newPassphrase != confirmPassphrase {
				return fmt.Errorf("passphrase confirmation does not match")
			}

			if err := s.IsValidPassphrase(ctx, newPassphrase); err != nil {
				return err
			}

			stop := tui.Spin("Replacing passphrase...")
			defer stop()

			if err := s.ReplacePassphrase(ctx, oldPassphrase, newPassphrase); err != nil {
				return err
			}
			stop()

			fmt.Println("Passphrase replaced successfully")
			return nil
		},
	}
}

//nolint:dupl // newReplacePassphraseCmd and newReplacePINCmd have similar behaviour
func newReplacePINCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-pin",
		Usage: "Replace encryption PIN",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			s := tpm.New()

			oldPIN, err := tui.ReadUserSecret("Enter current PIN: ")
			if err != nil {
				return err
			}

			newPIN, err := tui.ReadUserSecret("Enter new PIN: ")
			if err != nil {
				return err
			}

			confirmPIN, err := tui.ReadUserSecret("Confirm new PIN: ")
			if err != nil {
				return err
			}

			if newPIN != confirmPIN {
				return fmt.Errorf("PIN confirmation does not match")
			}

			if err := s.IsValidPIN(ctx, newPIN); err != nil {
				return err
			}

			stop := tui.Spin("Replacing PIN...")
			defer stop()

			if err := s.ReplacePIN(ctx, oldPIN, newPIN); err != nil {
				return err
			}
			stop()

			fmt.Println("PIN replaced successfully")
			return nil
		},
	}
}
