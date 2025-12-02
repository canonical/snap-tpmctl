package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
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
			return createKey(ctx, recoveryKeyName)
		},
	}
}

func createKey(ctx context.Context, recoveryKeyName string) error {
	c := snapd.NewClient()
	defer c.Close()

	if err := c.LoadAuthFromHome(); err != nil {
		return fmt.Errorf("failed to load auth: %w", err)
	}

	if recoveryKeyName == "" {
		return fmt.Errorf("recovery key name cannot be empty")
	}

	if strings.HasPrefix(recoveryKeyName, "snap") || strings.HasPrefix(recoveryKeyName, "default") {
		return fmt.Errorf("recovery key name cannot start with 'snap' or 'default'")
	}

	key, err := c.GenerateRecoveryKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate recovery key: %w", err)
	}

	fmt.Printf("Recovery Key: %s\n", key.RecoveryKey)
	fmt.Printf("Key ID: %s\n", key.KeyID)

	keySlots := []snapd.KeySlot{{Name: recoveryKeyName}}

	resp, err := c.AddRecoveryKey(ctx, key.KeyID, keySlots)
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
