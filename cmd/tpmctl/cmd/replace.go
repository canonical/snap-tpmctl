package cmd

import (
	"context"
	"fmt"

	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

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

			if err := tui.WithSpinner("Replacing passphrase...", func() error {
				return s.ReplacePassphrase(ctx, oldPassphrase, newPassphrase)
			}); err != nil {
				return err
			}
			fmt.Println("Passphrase replaced successfully")
			return nil
		},
	}
}

func newReplacePinCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-pin",
		Usage: "Replace encryption PIN",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			s := tpm.New()

			oldPin, err := tui.ReadUserSecret("Enter current PIN: ")
			if err != nil {
				return err
			}

			newPin, err := tui.ReadUserSecret("Enter new PIN: ")
			if err != nil {
				return err
			}

			confirmPin, err := tui.ReadUserSecret("Confirm new PIN: ")
			if err != nil {
				return err
			}

			if newPin != confirmPin {
				return fmt.Errorf("PIN confirmation does not match")
			}

			if err := s.IsValidPIN(ctx, newPin); err != nil {
				return err
			}

			if err := tui.WithSpinner("Replacing PIN...", func() error {
				return s.ReplacePIN(ctx, oldPin, newPin)
			}); err != nil {
				return err
			}
			fmt.Println("PIN replaced successfully")
			return nil
		},
	}
}
