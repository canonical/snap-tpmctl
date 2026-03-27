package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

//nolint:dupl // newReplacePassphraseCmd and newReplacePINCmd have similar behaviour
func (a App) newReplacePassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-passphrase",
		Usage: "Replace encryption passphrase",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			oldPassphrase, err := a.tui.ReadUserSecret("Enter current passphrase: ")
			if err != nil {
				return err
			}

			newPassphrase, err := a.tui.ReadUserSecret("Enter new passphrase: ")
			if err != nil {
				return err
			}

			confirmPassphrase, err := a.tui.ReadUserSecret("Confirm new passphrase: ")
			if err != nil {
				return err
			}

			if newPassphrase != confirmPassphrase {
				return fmt.Errorf("passphrase confirmation does not match")
			}

			if err := a.tpm.IsValidPassphrase(ctx, newPassphrase); err != nil {
				return err
			}

			stop := a.tui.Spin("Replacing passphrase...")
			defer stop()

			if err := a.tpm.ReplacePassphrase(ctx, oldPassphrase, newPassphrase); err != nil {
				return err
			}
			stop()

			fmt.Println("Passphrase replaced successfully")
			return nil
		},
	}
}

//nolint:dupl // newReplacePassphraseCmd and newReplacePINCmd have similar behaviour
func (a App) newReplacePINCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-pin",
		Usage: "Replace encryption PIN",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			oldPIN, err := a.tui.ReadUserSecret("Enter current PIN: ")
			if err != nil {
				return err
			}

			newPIN, err := a.tui.ReadUserSecret("Enter new PIN: ")
			if err != nil {
				return err
			}

			confirmPIN, err := a.tui.ReadUserSecret("Confirm new PIN: ")
			if err != nil {
				return err
			}

			if newPIN != confirmPIN {
				return fmt.Errorf("PIN confirmation does not match")
			}

			if err := a.tpm.IsValidPIN(ctx, newPIN); err != nil {
				return err
			}

			stop := a.tui.Spin("Replacing PIN...")
			defer stop()

			if err := a.tpm.ReplacePIN(ctx, oldPIN, newPIN); err != nil {
				return err
			}
			stop()

			fmt.Println("PIN replaced successfully")
			return nil
		},
	}
}
