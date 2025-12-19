package tpm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"snap-tpmctl/internal/snapd"
)

// authValidator defines the interface for snapd operations needed for validation.
type authValidator interface {
	CheckPassphrase(ctx context.Context, passphrase string) (*snapd.Response, error)
	CheckPIN(ctx context.Context, pin string) (*snapd.Response, error)
	EnumerateKeySlots(ctx context.Context) (*snapd.SystemVolumesResult, error)
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
func IsValidPassphrase(ctx context.Context, client authValidator, passphrase, confirm string) error {
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
func IsValidPIN(ctx context.Context, client authValidator, pin, confirm string) error {
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

// ValidateAuthMode checks if the current authentication mode matches the expected mode.
func ValidateAuthMode(ctx context.Context, client authValidator, expectedAuthMode snapd.AuthMode) error {
	result, err := client.EnumerateKeySlots(ctx)
	if err != nil {
		return fmt.Errorf("failed to enumerate key slots: %w", err)
	}

	systemData, ok := result.ByContainerRole["system-data"]
	if !ok {
		return fmt.Errorf("system-data container role not found")
	}

	defaultKeyslot, ok := systemData.KeySlots["default"]
	if !ok {
		return fmt.Errorf("default key slot not found in system-data")
	}

	defaultFallbackKeyslot, ok := systemData.KeySlots["default-fallback"]
	if !ok {
		return fmt.Errorf("default-fallback key slot not found in system-data")
	}

	if defaultKeyslot.AuthMode != string(expectedAuthMode) || defaultFallbackKeyslot.AuthMode != string(expectedAuthMode) {
		return fmt.Errorf("authentication mode mismatch: expected %s, got default=%s, default-fallback=%s",
			expectedAuthMode,
			defaultKeyslot.AuthMode,
			defaultFallbackKeyslot.AuthMode,
		)
	}

	return nil
}

// ValidateRecoveryKeyName validates that a recovery key name is valid.
func ValidateRecoveryKeyName(ctx context.Context, client authValidator, recoveryKeyName string) error {
	// Recovery key name cannot be empty.
	if recoveryKeyName == "" {
		return fmt.Errorf("recovery key name cannot be empty")
	}

	// Recovery key name cannot start with 'snap' or 'default'.
	if strings.HasPrefix(recoveryKeyName, "snap") || strings.HasPrefix(recoveryKeyName, "default") {
		return fmt.Errorf("recovery key name cannot start with 'snap' or 'default'")
	}

	return nil
}

// ValidateRecoveryKeyNameUnique validates that a recovery key name is valid and not in use.
func ValidateRecoveryKeyNameUnique(ctx context.Context, client authValidator, recoveryKeyName string) error {
	if err := ValidateRecoveryKeyName(ctx, client, recoveryKeyName); err != nil {
		return err
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

var isValidRecoveryKey = regexp.MustCompile("^([0-9]{5}-){7}[0-9]{5}$").MatchString

// ValidateRecoveryKey validates that recovery key matches expected formatting.
func ValidateRecoveryKey(key string) error {
	if key == "" {
		return fmt.Errorf("recovery key cannot be empty")
	}

	if !isValidRecoveryKey(key) {
		return fmt.Errorf("invalid recovery key format: must contain only alphanumeric characters and hyphens")
	}

	return nil
}
