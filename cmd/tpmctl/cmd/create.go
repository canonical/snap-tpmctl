package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	"snap-tpmctl/internal/snapd"
)

type keyCreator interface {
	LoadAuthFromHome() error
	GenerateRecoveryKey(ctx context.Context) (*snapd.GenerateRecoveryKeyResult, error)
	EnumerateKeySlots(ctx context.Context) (*snapd.SystemVolumesResult, error)
	AddRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) (*snapd.AsyncResponse, error)
	Close() error
}

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
			return createKey(ctx, c, recoveryKeyName)
		},
	}
}

// FIXME: Keep io here
func createKey(ctx context.Context, client keyCreator, recoveryKeyName string) error {
	if err := client.LoadAuthFromHome(); err != nil {
		return fmt.Errorf("failed to load auth: %w", err)
	}

	// FIXME: move this out such that only the printf remains, move to internal to do the snapd interactions, returning err, key, etc.
	key, err := client.GenerateRecoveryKey(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate recovery key: %w", err)
	}

	if err := isValidRecoveryKeyName(ctx, client, recoveryKeyName); err != nil {
		return err
	}

	fmt.Printf("Recovery Key: %s\n", key.RecoveryKey)
	fmt.Printf("Key ID: %s\n", key.KeyID)

	keySlots := []snapd.KeySlot{{Name: recoveryKeyName}}

	resp, err := client.AddRecoveryKey(ctx, key.KeyID, keySlots)
	if err != nil {
		return fmt.Errorf("failed to add recovery key: %w", err)
	}

	fmt.Println(resp.Status)

	return nil
}

func isValidRecoveryKeyName(ctx context.Context, client keyCreator, recoveryKeyName string) error {
	// Recovery key name cannot be empty.
	if recoveryKeyName == "" {
		return fmt.Errorf("recovery key name cannot be empty")
	}

	// Recovery key name cannot start with 'snap' or 'default'.
	if strings.HasPrefix(recoveryKeyName, "snap") || strings.HasPrefix(recoveryKeyName, "default") {
		return fmt.Errorf("recovery key name cannot start with 'snap' or 'default'")
	}

	// Recovery key name cannot already be in use.
	result, err := client.EnumerateKeySlots(ctx)
	if err != nil {
		return fmt.Errorf("failed to enumerate key slots: %w", err)
	}

	for _, volumeInfo := range result.ByContainerRole {
		for slotName := range volumeInfo.KeySlots {
			if slotName == recoveryKeyName {
				return fmt.Errorf("recovery key name %q is already in use", recoveryKeyName)
			}
		}
	}

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
