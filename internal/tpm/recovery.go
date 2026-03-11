package tpm

import (
	"context"
	"fmt"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/snapcore/secboot"
)

// CreateKey creates a new recovery key with the given name. Input should be validated using ValidateRecoveryKeyNameUnique first.
func (s SnapTPM) CreateKey(ctx context.Context, recoveryKeyName string) (recoveryKey string, err error) {
	key, err := s.snapdClient.GenerateRecoveryKey(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate recovery key: %v", err)
	}

	keySlots := []snapd.Keyslot{{Name: recoveryKeyName}}

	if err := s.snapdClient.AddRecoveryKey(ctx, key.KeyID, keySlots); err != nil {
		return "", fmt.Errorf("failed to add recovery key: %v", err)
	}

	return key.RecoveryKey, nil
}

// RegenerateKey replaces an existing recovery key with a new one with the given name. Input should be validated using ValidateRecoveryKeyName first.
func (s SnapTPM) RegenerateKey(ctx context.Context, recoveryKeyName string) (recoveryKey string, err error) {
	key, err := s.snapdClient.GenerateRecoveryKey(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate recovery key: %v", err)
	}

	keySlots := []snapd.Keyslot{{Name: recoveryKeyName}}

	if err := s.snapdClient.ReplaceRecoveryKey(ctx, key.KeyID, keySlots); err != nil {
		return "", fmt.Errorf("failed to replace recovery key: %v", err)
	}

	return key.RecoveryKey, nil
}

// CheckKey verifies if a recovery key is valid by checking it against the system.
func (s SnapTPM) CheckKey(ctx context.Context, recoveryKey string) (bool, error) {
	ok, err := s.snapdClient.CheckRecoveryKey(ctx, recoveryKey, nil)
	if err != nil {
		return false, fmt.Errorf("failed to check recovery key: %v", err)
	}

	return ok, nil
}

// GetLuksKey validates and converts the recovery key to a binary key format by parsing and formatting it.
func GetLuksKey(recoveryKey string) (secboot.DiskUnlockKey, error) {
	binKey, err := secboot.ParseRecoveryKey(recoveryKey)
	if err != nil {
		return nil, err
	}

	return binKey[:], nil
}
