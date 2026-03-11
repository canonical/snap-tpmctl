package tpm

import (
	"context"
	"fmt"

	"github.com/canonical/snap-tpmctl/internal/snapd"
)

// AddPassphrase adds passphrase authentication to the platform key.
func (s SnapTPM) AddPassphrase(ctx context.Context, passphrase string) error {
	err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModePassphrase, passphrase)
	if err != nil {
		return fmt.Errorf("failed to add passphrase: %v", err)
	}

	return nil
}

// ReplacePassphrase replaces the passphrase.
func (s SnapTPM) ReplacePassphrase(ctx context.Context, oldPassphrase, newPassphrase string) error {
	err := s.snapdClient.ReplacePassphrase(ctx, oldPassphrase, newPassphrase, nil)
	if err != nil {
		return fmt.Errorf("failed to change passphrase: %v", err)
	}

	return nil
}

// RemovePassphrase removes passphrase authentication from the platform key.
func (s SnapTPM) RemovePassphrase(ctx context.Context) error {
	err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModeNone, "")
	if err != nil {
		return fmt.Errorf("failed to remove passphrase: %v", err)
	}

	return nil
}

// AddPIN adds PIN authentication to the platform key.
func (s SnapTPM) AddPIN(ctx context.Context, pin string) error {
	err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModePIN, pin)
	if err != nil {
		return fmt.Errorf("failed to add PIN: %v", err)
	}

	return nil
}

// ReplacePIN replaces the PIN using the provided client.
func (s SnapTPM) ReplacePIN(ctx context.Context, oldPIN, newPIN string) error {
	err := s.snapdClient.ReplacePIN(ctx, oldPIN, newPIN, nil)
	if err != nil {
		return fmt.Errorf("failed to change PIN: %v", err)
	}

	return nil
}

// RemovePIN removes PIN authentication from the platform key.
func (s SnapTPM) RemovePIN(ctx context.Context) error {
	err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModeNone, "")
	if err != nil {
		return fmt.Errorf("failed to remove PIN: %v", err)
	}

	return nil
}
