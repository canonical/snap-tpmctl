package tpm

import (
	"context"
	"fmt"

	"snap-tpmctl/internal/snapd"
)

// authReplacer defines the interface for snapd operations needed for changing authentication.
type authReplacer interface {
	ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []snapd.KeySlot) (*snapd.AsyncResponse, error)
	ReplacePIN(ctx context.Context, oldPin string, newPin string, keySlots []snapd.KeySlot) (*snapd.AsyncResponse, error)
	ReplacePlatformKey(ctx context.Context, authMode snapd.AuthMode, pin, passphrase string) (*snapd.AsyncResponse, error)
}

// ReplacePassphrase replaces the passphrase using the provided client.
func ReplacePassphrase(ctx context.Context, client authReplacer, oldPassphrase, newPassphrase string) error {
	ares, err := client.ReplacePassphrase(ctx, oldPassphrase, newPassphrase, nil)
	if err != nil {
		return fmt.Errorf("failed to change passphrase: %w", err)
	}

	if !ares.IsOK() {
		return fmt.Errorf("unable to replace passphrase: %s", ares.Err)
	}

	return nil
}

// ReplacePIN replaces the PIN using the provided client.
func ReplacePIN(ctx context.Context, client authReplacer, oldPin, newPin string) error {
	ares, err := client.ReplacePIN(ctx, oldPin, newPin, nil)
	if err != nil {
		return fmt.Errorf("failed to change PIN: %w", err)
	}

	if !ares.IsOK() {
		return fmt.Errorf("unable to replace PIN: %s", ares.Err)
	}

	return nil
}

// AddPassphrase adds passphrase authentication to the platform key.
func AddPassphrase(ctx context.Context, client authReplacer, passphrase string) error {
	ares, err := client.ReplacePlatformKey(ctx, snapd.AuthModePassphrase, "", passphrase)
	if err != nil {
		return fmt.Errorf("failed to add passphrase: %w", err)
	}

	if !ares.IsOK() {
		return fmt.Errorf("unable to add passphrase: %s", ares.Err)
	}

	return nil
}

// AddPIN adds PIN authentication to the platform key.
func AddPIN(ctx context.Context, client authReplacer, pin string) error {
	ares, err := client.ReplacePlatformKey(ctx, snapd.AuthModePin, pin, "")
	if err != nil {
		return fmt.Errorf("failed to add PIN: %w", err)
	}

	if !ares.IsOK() {
		return fmt.Errorf("unable to add PIN: %s", ares.Err)
	}

	return nil
}

// RemovePassphrase removes passphrase authentication from the platform key.
func RemovePassphrase(ctx context.Context, client authReplacer) error {
	ares, err := client.ReplacePlatformKey(ctx, snapd.AuthModeNone, "", "")
	if err != nil {
		return fmt.Errorf("failed to remove passphrase: %w", err)
	}

	if !ares.IsOK() {
		return fmt.Errorf("unable to remove passphrase: %s", ares.Err)
	}

	return nil
}

// RemovePIN removes PIN authentication from the platform key.
func RemovePIN(ctx context.Context, client authReplacer) error {
	ares, err := client.ReplacePlatformKey(ctx, snapd.AuthModeNone, "", "")
	if err != nil {
		return fmt.Errorf("failed to remove PIN: %w", err)
	}

	if !ares.IsOK() {
		return fmt.Errorf("unable to remove PIN: %s", ares.Err)
	}

	return nil
}
