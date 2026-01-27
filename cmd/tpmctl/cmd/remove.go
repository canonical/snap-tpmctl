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

			c := snapd.NewClient()

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

func newRemoveRecoveryKey() *cli.Command {
	var recoveryKeyName string

	return &cli.Command{
		Name:    "remove-recovery-key",
		Usage:   "Remove the recovery key form the system",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKeyName,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Ensure that the user's effective ID is root
			if os.Geteuid() != 0 {
				return fmt.Errorf("this command requires elevated privileges. Please run with sudo")
			}

			c := snapd.NewClient()

			// Validate the recovery key name
			if err := tpm.ValidateRecoveryKeyName(ctx, c, recoveryKeyName); err != nil {
				return err
			}

			if err := tui.WithSpinner("Removing recovery key...", func() error {
				return tpm.RemoveRecoveryKey(ctx, c, recoveryKeyName)
			}); err != nil {
				return err
			}
			fmt.Println("PIN removed successfully")
			return nil
		},
	}
}
