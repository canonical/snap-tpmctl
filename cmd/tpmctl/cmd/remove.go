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

func newRemovePassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "remove-passphrase",
		Usage: "Remove passphrase authentication",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if os.Geteuid() != 0 {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
			}

			s := tpm.New()

			// Validate auth mode is currently passphrase
			if err := s.ValidateAuthMode(ctx, snapd.AuthModePassphrase); err != nil {
				return err
			}

			stop := tui.Spin("Removing passphrase...")
			defer stop()

			if err := s.RemovePassphrase(ctx); err != nil {
				return err
			}
			stop()

			fmt.Println("Passphrase removed successfully")
			return nil
		},
	}
}

func newRemovePINCmd() *cli.Command {
	return &cli.Command{
		Name:  "remove-pin",
		Usage: "Remove PIN authentication",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if os.Geteuid() != 0 {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
			}

			s := tpm.New()

			// Validate auth mode is currently PIN
			if err := s.ValidateAuthMode(ctx, snapd.AuthModePIN); err != nil {
				return err
			}

			stop := tui.Spin("Removing PIN...")
			defer stop()

			if err := s.RemovePIN(ctx); err != nil {
				return err
			}
			stop()

			fmt.Println("PIN removed successfully")
			return nil
		},
	}
}
