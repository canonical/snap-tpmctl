package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
)

func newRegenerateKeyCmd() *cli.Command {
	var recoveryKey string

	return &cli.Command{
		Name:    "regenerate-key",
		Usage:   "Regenerate an existing local recovery key",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKey,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return regenerateKey(ctx, recoveryKey)
		},
	}
}

func regenerateKey(ctx context.Context, _ string) error {
	// TODO: decide if we want to match exactly the security center
	// behaviour showing the key, waiting for user confirmation and then
	// replace the key and removing it from the screen

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

	res, err := c.ReplaceRecoveryKey(ctx, key.KeyID, nil)
	if err != nil {
		return fmt.Errorf("failed to replace recovery key: %w", err)
	}

	fmt.Println(res.Status)
	fmt.Println(res.Change)

	return nil
}

func newRegenerateEnterpriseKeyCmd() *cli.Command {
	var recoveryKey string

	return &cli.Command{
		Name:    "regenerate-enterprise-key",
		Usage:   "Regenerate an existing enterprise recovery key",
		Suggest: true,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "key-id",
				UsageText:   "<key-id>",
				Destination: &recoveryKey,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return regenerateEnterpriseKey(ctx, recoveryKey)
		},
	}
}

func regenerateEnterpriseKey(_ context.Context, key string) error {
	fmt.Println("Regenerated enterprise key with id", key)
	return nil
}
