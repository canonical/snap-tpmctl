package tpm

import (
	"context"
	"fmt"

	"snap-tpmctl/internal/snapd"
)

// keyCreator defines the interface for snapd operations needed for key management.
type keyCreator interface {
	GenerateRecoveryKey(ctx context.Context) (*snapd.GenerateRecoveryKeyResult, error)
	AddRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) (*snapd.AsyncResponse, error)
	ReplaceRecoveryKey(ctx context.Context, keyID string, slots []snapd.KeySlot) (*snapd.AsyncResponse, error)
}

// CreateKeyResult contains the result of creating a recovery key.
type CreateKeyResult struct {
	RecoveryKey string
	KeyID       string
	Status      string
}

// CreateKey creates a new recovery key with the given name. Input should be validated using ValidateRecoveryKeyName first.
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

// RegenerateKey replaces an existing recovery key with a new one with the given name. Input should be validated using ValidateRecoveryKeyName first.
func RegenerateKey(ctx context.Context, client keyCreator, recoveryKeyName string) (result *CreateKeyResult, err error) {
	key, err := client.GenerateRecoveryKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery key: %w", err)
	}

	keySlots := []snapd.KeySlot{{Name: recoveryKeyName}}

	resp, err := client.ReplaceRecoveryKey(ctx, key.KeyID, keySlots)
	if err != nil {
		return nil, fmt.Errorf("failed to replace recovery key: %w", err)
	}

	return &CreateKeyResult{
		RecoveryKey: key.RecoveryKey,
		KeyID:       key.KeyID,
		Status:      resp.Status,
	}, nil
}

type keyChecker interface {
	CheckRecoveryKey(ctx context.Context, recoveryKey string, containerRoles []string) (*snapd.Response, error)
}

// CheckKey verifies if a recovery key is valid by checking it against the system.
func CheckKey(ctx context.Context, client keyChecker, recoveryKey string) (bool, error) {
	res, err := client.CheckRecoveryKey(ctx, recoveryKey, nil)
	if err != nil {
		return false, fmt.Errorf("failed to check recovery key: %w", err)
	}

	return res.IsOK(), nil
}
