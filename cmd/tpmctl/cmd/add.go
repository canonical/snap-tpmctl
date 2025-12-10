package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tpm"
	"snap-tpmctl/internal/tui"
)

//nolint:dupl // PIN and passphrase commands have intentionally similar structure
func newAddPassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "add-passphrase",
		Usage: "Add passphrase authentication",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if os.Geteuid() != 0 {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
			}

			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			// Validate auth mode is currently none
			if err := tpm.ValidateAuthMode(ctx, c, snapd.AuthModeNone); err != nil {
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

			if err := tpm.AddPassphrase(ctx, c, newPassphrase); err != nil {
				return err
			}
			fmt.Println("Passphrase added successfully")
			return nil
		},
	}
}

//nolint:dupl // PIN and passphrase commands have intentionally similar structure
func newAddPINCmd() *cli.Command {
	return &cli.Command{
		Name:  "add-pin",
		Usage: "Add PIN authentication",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if os.Geteuid() != 0 {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
			}

			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			// Validate auth mode is currently none
			if err := tpm.ValidateAuthMode(ctx, c, snapd.AuthModeNone); err != nil {
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
			if err := tpm.AddPIN(ctx, c, newPin); err != nil {
				return err
			}
			fmt.Println("PIN added successfully")
			return nil
		},
	}
}
