package cmd

import (
	"context"
	"fmt"
	"snap-tpmctl/internal/snapd"

	"github.com/urfave/cli/v3"
)

func newCreateKeyCmd() *cli.Command {
	return &cli.Command{
		Name:  "create-key",
		Usage: "Create a new local recovery key",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// TODO: let req fail if not root

			return createKey(ctx)
		},
	}
}

func createKey(ctx context.Context) error {
	c := snapd.NewClient() // should default be true?
	defer c.Close()

	// Load authorization for API request
	if err := c.LoadAuthFromHome(); err != nil {
		return err
	}

	// Generate the recovery key
	key, err := c.GenerateRecoveryKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate recovery key: %w", err)
	}

	fmt.Printf("Recovery Key: %s\n", key.RecoveryKey)
	fmt.Printf("Key ID: %s\n", key.KeyID)

	// TODO: when the API is ready add `add-recovery-key` call
	// resp, err := c.AddRecoveryKey(ctx, key.KeyID, nil)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(resp)

	return nil
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
