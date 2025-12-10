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
