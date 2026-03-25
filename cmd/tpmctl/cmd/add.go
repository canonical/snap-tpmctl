package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
)

//nolint:dupl // PIN and passphrase commands have intentionally similar structure
func (a App) newAddPassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "add-passphrase",
		Usage: "Add passphrase authentication",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if os.Geteuid() != 0 {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
			}

			// Validate auth mode is currently none
			if err := a.tpm.ValidateAuthMode(ctx, snapd.AuthModeNone); err != nil {
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

			if err := a.tpm.IsValidPassphrase(ctx, newPassphrase); err != nil {
				return err
			}

			stop := tui.Spin("Adding passphrase...")
			defer stop()

			if err := a.tpm.AddPassphrase(ctx, newPassphrase); err != nil {
				return err
			}
			stop()

			fmt.Println("Passphrase added successfully")
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
			if os.Geteuid() != 0 {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
			}

			// Validate auth mode is currently none
			if err := a.tpm.ValidateAuthMode(ctx, snapd.AuthModeNone); err != nil {
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

			if err := a.tpm.IsValidPIN(ctx, newPIN); err != nil {
				return err
			}
			stop := tui.Spin("Adding PIN...")
			defer stop()

			if err := a.tpm.AddPIN(ctx, newPIN); err != nil {
				return err
			}
			stop()

			fmt.Println("PIN added successfully")
			return nil
		},
	}
}
