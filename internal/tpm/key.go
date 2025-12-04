package tpm

import (
	"context"
	"fmt"
	"strings"

	"snap-tpmctl/internal/snapd"
)

// keyCreator defines the interface for snapd operations needed for key management.
type keyCreator interface {
	GenerateRecoveryKey(ctx context.Context) (*snapd.GenerateRecoveryKeyResult, error)
	EnumerateKeySlots(ctx context.Context) (*snapd.SystemVolumesResult, error)
	AddRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) (*snapd.AsyncResponse, error)
}

// CreateKeyResult contains the result of creating a recovery key.
type CreateKeyResult struct {
	RecoveryKey string
	KeyID       string
	Status      string
}

// ValidateRecoveryKeyName validates that a recovery key name is valid and not in use.
func ValidateRecoveryKeyName(ctx context.Context, client keyCreator, recoveryKeyName string) error {
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

// CreateKey creates a new recovery key with the given name. Input should be validated using [ValidateRecoveryKeyName] first.
func CreateKey(ctx context.Context, client keyCreator, recoveryKeyName string) (result *CreateKeyResult, err error) {
	key, err := client.GenerateRecoveryKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery key: %w", err)
	}

	keySlots := []snapd.KeySlot{{Name: recoveryKeyName}}

	resp, err := client.AddRecoveryKey(ctx, key.KeyID, keySlots)
	if err != nil {
		return nil, fmt.Errorf("failed to add recovery key: %w", err)
	}

	return &CreateKeyResult{
		RecoveryKey: key.RecoveryKey,
		KeyID:       key.KeyID,
		Status:      resp.Status,
	}, nil
}
