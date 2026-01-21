package tpm

import (
	"context"
	"fmt"

	"snap-tpmctl/internal/snapd"
)

// authReplacer defines the interface for snapd operations needed for changing authentication.
type authReplacer interface {
	ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []snapd.KeySlot) error
	ReplacePIN(ctx context.Context, oldPin string, newPin string, keySlots []snapd.KeySlot) error
	ReplacePlatformKey(ctx context.Context, authMode snapd.AuthMode, pin, passphrase string) error
}

// ReplacePassphrase replaces the passphrase using the provided client.
func ReplacePassphrase(ctx context.Context, client authReplacer, oldPassphrase, newPassphrase string) error {
	err := client.ReplacePassphrase(ctx, oldPassphrase, newPassphrase, nil)
	if err != nil {
		return fmt.Errorf("failed to change passphrase: %w", err)
	}

	return nil
}

// ReplacePIN replaces the PIN using the provided client.
func ReplacePIN(ctx context.Context, client authReplacer, oldPin, newPin string) error {
	err := client.ReplacePIN(ctx, oldPin, newPin, nil)
	if err != nil {
		return fmt.Errorf("failed to change PIN: %w", err)
	}

	return nil
}

// AddPassphrase adds passphrase authentication to the platform key.
func AddPassphrase(ctx context.Context, client authReplacer, passphrase string) error {
	err := client.ReplacePlatformKey(ctx, snapd.AuthModePassphrase, "", passphrase)
	if err != nil {
		return fmt.Errorf("failed to add passphrase: %w", err)
	}

	return nil
}

// AddPIN adds PIN authentication to the platform key.
func AddPIN(ctx context.Context, client authReplacer, pin string) error {
	err := client.ReplacePlatformKey(ctx, snapd.AuthModePin, pin, "")
	if err != nil {
		return fmt.Errorf("failed to add PIN: %w", err)
	}

	return nil
}

// RemovePassphrase removes passphrase authentication from the platform key.
func RemovePassphrase(ctx context.Context, client authReplacer) error {
	err := client.ReplacePlatformKey(ctx, snapd.AuthModeNone, "", "")
	if err != nil {
		return fmt.Errorf("failed to remove passphrase: %w", err)
	}

	return nil
}

// RemovePIN removes PIN authentication from the platform key.
func RemovePIN(ctx context.Context, client authReplacer) error {
	err := client.ReplacePlatformKey(ctx, snapd.AuthModeNone, "", "")
	if err != nil {
		return fmt.Errorf("failed to remove PIN: %w", err)
	}

	return nil
}
