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

			c := snapd.New()

			// Validate auth mode is currently passphrase
			if err := tpm.ValidateAuthMode(ctx, c, snapd.AuthModePassphrase); err != nil {
				return err
			}

			if err := tui.WithSpinner("Removing passphrase...", func() error {
				return tpm.RemovePassphrase(ctx, c)
			}); err != nil {
				return err
			}

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

			c := snapd.New()

			// Validate auth mode is currently PIN
			if err := tpm.ValidateAuthMode(ctx, c, snapd.AuthModePin); err != nil {
				return err
			}

			if err := tui.WithSpinner("Removing PIN...", func() error {
				return tpm.RemovePIN(ctx, c)
			}); err != nil {
				return err
			}
			fmt.Println("PIN removed successfully")
			return nil
		},
	}
}
