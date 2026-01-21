package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tpm"
	"snap-tpmctl/internal/tui"
)

func newCreateKeyCmd() *cli.Command {
	var recoveryKeyName string

	return &cli.Command{
		Name:  "create-recovery-key",
		Usage: "Create a new recovery key",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKeyName,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			c := snapd.NewClient()

			// Validate the recovery key name
			if err := tpm.ValidateRecoveryKeyNameUnique(ctx, c, recoveryKeyName); err != nil {
				return err
			}

			result, err := tui.WithSpinnerResult("Generating recovery key...", func() (*tpm.CreateKeyResult, error) {
				return tpm.CreateKey(ctx, c, recoveryKeyName)
			})
			if err != nil {
				return err
			}

			fmt.Printf("Recovery Key: %s\n", result.RecoveryKey)
			fmt.Printf("Key ID: %s\n", result.KeyID)

			return nil
		},
	}
}
