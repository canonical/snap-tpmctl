package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/urfave/cli/v3"
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

			s := tpm.New()

			// Validate auth mode is currently none
			if err := s.ValidateAuthMode(ctx, snapd.AuthModeNone); err != nil {
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

			if err := tui.WithSpinner("Adding passphrase...", func() error {
				return s.AddPassphrase(ctx, newPassphrase)
			}); err != nil {
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

			s := tpm.New()

			// Validate auth mode is currently none
			if err := s.ValidateAuthMode(ctx,  snapd.AuthModeNone); err != nil {
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

			if err := tui.WithSpinner("Adding PIN...", func() error {
				return s.AddPIN(ctx, newPin)
			}); err != nil {
				return err
			}
			fmt.Println("PIN added successfully")
			return nil
		},
	}
}
