package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
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
	c := snapd.NewClient()
	defer c.Close()

	if err := c.LoadAuthFromHome(); err != nil {
		return fmt.Errorf("failed to load auth: %w", err)
	}

	key, err := c.GenerateRecoveryKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate recovery key: %w", err)
	}

	fmt.Printf("Recovery Key: %s\n", key.RecoveryKey)
	fmt.Printf("Key ID: %s\n", key.KeyID)

	resp, err := c.AddRecoveryKey(ctx, key.KeyID, nil)
	if err != nil {
		return fmt.Errorf("failed to add recovery key: %w", err)
	}

	fmt.Println(resp.Status)

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
