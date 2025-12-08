package tpm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"snap-tpmctl/internal/snapd"
)

// authReplacer defines the interface for snapd operations needed for key management.
type authReplacer interface {
	CheckPassphrase(ctx context.Context, passphrase string) (*snapd.Response, error)
	CheckPIN(ctx context.Context, pin string) (*snapd.Response, error)
	ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []snapd.KeySlot) (*snapd.AsyncResponse, error)
	ReplacePIN(ctx context.Context, oldPin string, newPin string, keySlots []snapd.KeySlot) (*snapd.AsyncResponse, error)
}

// resultValue represents the value field in validation error responses from snapd.
type resultValue struct {
	Reasons            []string `json:"reasons"`
	EntropyBits        uint     `json:"entropy-bits"`
	MinEntropyBits     uint     `json:"min-entropy-bits"`
	OptimalEntropyBits uint     `json:"optimal-entropy-bits"`
}

// handleValidationError processes snapd validation errors and returns appropriate error messages.
func handleValidationError(err error, authMode string) error {
	var snapdErr *snapd.Error
	if !errors.As(err, &snapdErr) {
		return fmt.Errorf("failed to check %s: %w", authMode, err)
	}

	switch snapdErr.Kind {
	case "invalid-passphrase", "invalid-pin":
		// Try to unmarshal the value to check for specific reasons
		var resValue resultValue
		if err := json.Unmarshal(snapdErr.Value, &resValue); err != nil {
			if snapdErr.Message != "" {
				return fmt.Errorf("%s is invalid: %s", authMode, snapdErr.Message)
			}
			return fmt.Errorf("%s is invalid", authMode)
		}

		if slices.Contains(resValue.Reasons, "low-entropy") {
			return fmt.Errorf("%s is too weak, make it longer or more complex", authMode)
		}

		// Fallback to generic message
		if snapdErr.Message != "" {
			return fmt.Errorf("%s is invalid: %s", authMode, snapdErr.Message)
		}
		return fmt.Errorf("%s is invalid", authMode)
	case "unsupported":
		if snapdErr.Message != "" {
			return fmt.Errorf("%s validation not supported: %s", authMode, snapdErr.Message)
		}
		return fmt.Errorf("%s validation not supported", authMode)
	default:
		if snapdErr.Message != "" {
			return fmt.Errorf("%s failed validation: %s", authMode, snapdErr.Message)
		}
		return fmt.Errorf("%s failed validation", authMode)
	}
}

// IsValidPassphrase validates that the passphrase and confirmation match and are not empty.
func IsValidPassphrase(ctx context.Context, client authReplacer, passphrase, confirm string) error {
	if passphrase == "" || confirm == "" {
		return fmt.Errorf("passphrase cannot be empty, try again")
	}

	if passphrase != confirm {
		return fmt.Errorf("passphrases do not match, try again")
	}

	res, err := client.CheckPassphrase(ctx, passphrase)
	if err != nil {
		return handleValidationError(err, "passphrase")
	}

	if !res.IsOK() {
		return fmt.Errorf("weak passphrase, make it longer or more complex")
	}

	return nil
}

// IsValidPIN validates that the PIN and confirmation match and are not empty.
func IsValidPIN(ctx context.Context, client authReplacer, pin, confirm string) error {
	if pin == "" || confirm == "" {
		return fmt.Errorf("PIN cannot be empty, try again")
	}

	// Check only digits in PIN
	for _, ch := range pin {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("PIN must contain only digits, try again")
		}
	}

	if pin != confirm {
		return fmt.Errorf("PINs do not match, try again")
	}

	res, err := client.CheckPIN(ctx, pin)
	if err != nil {
		return handleValidationError(err, "PIN")
	}

	if !res.IsOK() {
		return fmt.Errorf("weak PIN, make it longer or more complex")
	}

	return nil
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
