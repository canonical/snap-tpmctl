package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

//nolint:dupl // PIN and passphrase commands have intentionally similar structure
func (a App) newAddPassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "add-passphrase",
		Usage: "Add passphrase authentication",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if !a.isUserRoot() {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
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

			stop := a.tui.Spin("Adding passphrase...")
			defer stop()

			if err := a.tpm.AddPassphrase(ctx, newPassphrase); err != nil {
				return err
			}
			stop()

			fmt.Fprintln(a.tui.Writer(), "Passphrase added successfully")
			return nil
		},
	}
}

//nolint:dupl // PIN and passphrase commands have intentionally similar structure
func (a App) newAddPINCmd() *cli.Command {
	return &cli.Command{
		Name:  "add-pin",
		Usage: "Add PIN authentication",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if !a.isUserRoot() {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
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

			stop := a.tui.Spin("Adding PIN...")
			defer stop()

			if err := a.tpm.AddPIN(ctx, newPIN); err != nil {
				return err
			}
			stop()

			fmt.Fprintln(a.tui.Writer(), "PIN added successfully")
			return nil
		},
	}
}
