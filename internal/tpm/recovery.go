package tpm

import (
	"context"
	"fmt"

	"snap-tpmctl/internal/snapd"
)

// keyCreator defines the interface for snapd operations needed for key management.
type keyCreator interface {
	GenerateRecoveryKey(ctx context.Context) (*snapd.GenerateRecoveryKeyResult, error)
	AddRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) error
	ReplaceRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) error
}

// CreateKeyResult contains the result of creating a recovery key.
type CreateKeyResult struct {
	RecoveryKey string
	KeyID       string
}

// CreateKey creates a new recovery key with the given name. Input should be validated using ValidateRecoveryKeyNameUnique first.
func CreateKey(ctx context.Context, client keyCreator, recoveryKeyName string) (result *CreateKeyResult, err error) {
	key, err := client.GenerateRecoveryKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery key: %w", err)
	}

	keySlots := []snapd.KeySlot{{Name: recoveryKeyName}}

	if err := client.AddRecoveryKey(ctx, key.KeyID, keySlots); err != nil {
		return nil, fmt.Errorf("failed to add recovery key: %w", err)
	}

	return &CreateKeyResult{
		RecoveryKey: key.RecoveryKey,
		KeyID:       key.KeyID,
	}, nil
}

// RegenerateKey replaces an existing recovery key with a new one with the given name. Input should be validated using ValidateRecoveryKeyName first.
func RegenerateKey(ctx context.Context, client keyCreator, recoveryKeyName string) (result *CreateKeyResult, err error) {
	key, err := client.GenerateRecoveryKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery key: %w", err)
	}

	keySlots := []snapd.KeySlot{{Name: recoveryKeyName}}

	if err := client.ReplaceRecoveryKey(ctx, key.KeyID, keySlots); err != nil {
		return nil, fmt.Errorf("failed to replace recovery key: %w", err)
	}

	return &CreateKeyResult{
		RecoveryKey: key.RecoveryKey,
		KeyID:       key.KeyID,
	}, nil
}

type keyChecker interface {
	CheckRecoveryKey(ctx context.Context, recoveryKey string, containerRoles []string) (bool, error)
}

// CheckKey verifies if a recovery key is valid by checking it against the system.
func CheckKey(ctx context.Context, client keyChecker, recoveryKey string) (bool, error) {
	ok, err := client.CheckRecoveryKey(ctx, recoveryKey, nil)
	if err != nil {
		return false, fmt.Errorf("failed to check recovery key: %w", err)
	}

	return ok, nil
}
