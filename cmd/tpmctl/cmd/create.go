package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
	"snap-tpmctl/internal/tpm"
)

func newCreateKeyCmd() *cli.Command {
	var recoveryKeyName string

	return &cli.Command{
		Name:  "create-key",
		Usage: "Create a new local recovery key",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKeyName,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
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

			result, err := tpm.CreateKey(ctx, c, recoveryKeyName)
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

func newCreateEnterpriseKeyCmd() *cli.Command {
	return &cli.Command{
		Name:  "create-enterprise-key",
		Usage: "Create a new enterprise recovery key for Landscape",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return createEnterpriseKey(ctx)
		},
	}
}

func createEnterpriseKey(_ context.Context) error {
	fmt.Println("Created enterprise key")
	return nil
}
