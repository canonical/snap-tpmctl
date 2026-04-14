package tpm

import (
	"context"
	"fmt"

	"github.com/canonical/snap-tpmctl/internal/snapd"
)

// AddPassphrase adds passphrase authentication to the platform key.
func (s SnapTPM) AddPassphrase(ctx context.Context, passphrase string) error {
	if err := s.snapdClient.CheckPassphrase(ctx, passphrase); err != nil {
		return fmt.Errorf("failed to check passphrase: %v", err)
	}

	if err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModePassphrase, passphrase); err != nil {
		return fmt.Errorf("failed to add passphrase: %v", err)
	}

	return nil
}

// ReplacePassphrase replaces the passphrase.
func (s SnapTPM) ReplacePassphrase(ctx context.Context, oldPassphrase, newPassphrase string) error {
	if err := s.snapdClient.CheckPassphrase(ctx, newPassphrase); err != nil {
		return fmt.Errorf("failed to check passphrase: %v", err)
	}

	if err := s.snapdClient.ReplacePassphrase(ctx, oldPassphrase, newPassphrase, nil); err != nil {
		return fmt.Errorf("failed to change passphrase: %v", err)
	}

	return nil
}

// RemovePassphrase removes passphrase authentication from the platform key.
func (s SnapTPM) RemovePassphrase(ctx context.Context) error {
	if err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModeNone, ""); err != nil {
		return fmt.Errorf("failed to remove passphrase: %v", err)
	}

	return nil
}

// AddPIN adds PIN authentication to the platform key.
func (s SnapTPM) AddPIN(ctx context.Context, pin string) error {
	if err := s.snapdClient.CheckPIN(ctx, pin); err != nil {
		return fmt.Errorf("failed to validate PIN: %v", err)
	}

	if err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModePIN, pin); err != nil {
		return fmt.Errorf("failed to add PIN: %v", err)
	}

	return nil
}

// ReplacePIN replaces the PIN using the provided client.
func (s SnapTPM) ReplacePIN(ctx context.Context, oldPIN, newPIN string) error {
	if err := s.snapdClient.CheckPIN(ctx, newPIN); err != nil {
		return fmt.Errorf("failed to validate PIN: %v", err)
	}

	if err := s.snapdClient.ReplacePIN(ctx, oldPIN, newPIN, nil); err != nil {
		return fmt.Errorf("failed to change PIN: %v", err)
	}

	return nil
}

// RemovePIN removes PIN authentication from the platform key.
func (s SnapTPM) RemovePIN(ctx context.Context) error {
	if err := s.snapdClient.ReplacePlatformKey(ctx, snapd.AuthModeNone, ""); err != nil {
		return fmt.Errorf("failed to remove PIN: %v", err)
	}

	return nil
}
