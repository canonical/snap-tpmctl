package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tui"
)

func newReplacePassphraseCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-passphrase",
		Usage: "Replace encryption passphrase",
		Action: func(ctx context.Context, cmd *cli.Command) error {
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

			if err := IsValidPassphrase(newPassphrase, confirmPassphrase); err != nil {
				return err
			}

			return replacePassphrase(ctx, oldPassphrase, newPassphrase)
		},
	}
}

func replacePassphrase(ctx context.Context, oldPassphrase, newPassphrase string) error {
	c := snapd.NewClient()
	defer c.Close()

	if err := c.LoadAuthFromHome(); err != nil {
		return fmt.Errorf("failed to load auth: %w", err)
	}

	res, err := c.CheckPassphrase(ctx, newPassphrase)
	if err != nil {
		return fmt.Errorf("failed to check passphrase: %w", err)
	}

	if !res.IsOK() {
		return fmt.Errorf("weak passphrase, make it longer or more complex")
	}

	res, err = c.ReplacePassphrase(ctx, oldPassphrase, newPassphrase, nil)
	if err != nil {
		return fmt.Errorf("failed to change passphrase: %w", err)
	}

	msg := "Unable to replace passphrase"
	if res.IsOK() {
		msg = "Passphrase replaced successfully"
	}

	fmt.Println(msg)

	return nil
}

func newReplacePinCmd() *cli.Command {
	return &cli.Command{
		Name:  "replace-pin",
		Usage: "Repalce encryption pin",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return replacePin(ctx)
		},
	}
}

func replacePin(_ context.Context) error {
	fmt.Println("Repalce encryption pin")
	return nil
}

// IsValidPassphrase validates that the passphrase and confirmation match and are not empty.
func IsValidPassphrase(passphrase, confirm string) error {
	if passphrase == "" || confirm == "" {
		return fmt.Errorf("passphrase cannot be empty, try again")
	}

	if passphrase != confirm {
		return fmt.Errorf("passphrases do not match, try again")
	}

	// TODO: do we need to add a regex for valid char?

	return nil
}
