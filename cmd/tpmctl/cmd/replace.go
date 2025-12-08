package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tpm"
	"snap-tpmctl/internal/tui"
)

//nolint:dupl // PIN and passphrase commands have intentionally similar structure
func newReplacePassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-passphrase",
		Usage: "Replace encryption passphrase",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

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

			if err := tpm.IsValidPassphrase(ctx, c, newPassphrase, confirmPassphrase); err != nil {
				return err
			}

			if err := tpm.ReplacePassphrase(ctx, c, oldPassphrase, newPassphrase); err != nil {
				return err
			}
			fmt.Println("Passphrase replaced successfully")
			return nil
		},
	}
}

//nolint:dupl // PIN and passphrase commands have intentionally similar structure
func newReplacePinCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-pin",
		Usage: "Replace encryption PIN",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

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

			if err := tpm.IsValidPIN(ctx, c, newPin, confirmPin); err != nil {
				return err
			}

			if err := tpm.ReplacePIN(ctx, c, oldPin, newPin); err != nil {
				return err
			}
			fmt.Println("PIN replaced successfully")
			return nil
		},
	}
}
