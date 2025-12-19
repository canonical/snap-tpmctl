package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tpm"
	"snap-tpmctl/internal/tui"
)

func newRegenerateKeyCmd() *cli.Command {
	var recoveryKeyName string

	return &cli.Command{
		Name:    "regenerate-recovery-key",
		Usage:   "Regenerate an existing recovery key",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKeyName,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// TODO: decide if we want to match exactly the security center
			// behaviour showing the key, waiting for user confirmation and then
			// replace the key and removing it from the screen

			c := snapd.NewClient()
			defer c.Close()

			// Load auth before validation
			if err := c.LoadAuthFromHome(); err != nil {
				return fmt.Errorf("failed to load auth: %w", err)
			}

			// Validate the recovery key name
			if err := tpm.ValidateRecoveryKeyName(ctx, c, recoveryKeyName); err != nil {
				return err
			}

			result, err := tui.WithSpinnerResult("Regenerating recovery key...", func() (*tpm.CreateKeyResult, error) {
				return tpm.RegenerateKey(ctx, c, recoveryKeyName)
			})
			if err != nil {
				return err
			}

			fmt.Printf("Recovery Key: %s\n", result.RecoveryKey)
			fmt.Printf("Key ID: %s\n", result.KeyID)
			fmt.Println(result.Status)

			return nil
		},
	}
}
