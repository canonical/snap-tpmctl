package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tpm"
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

			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			// Validate auth mode is currently passphrase
			if err := tpm.ValidateAuthMode(ctx, c, snapd.AuthModePassphrase); err != nil {
				return err
			}

			if err := tpm.RemovePassphrase(ctx, c); err != nil {
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

			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			// Validate auth mode is currently PIN
			if err := tpm.ValidateAuthMode(ctx, c, snapd.AuthModePin); err != nil {
				return err
			}

			if err := tpm.RemovePIN(ctx, c); err != nil {
				return err
			}
			fmt.Println("PIN removed successfully")
			return nil
		},
	}
}
