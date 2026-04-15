package tpm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/snapcore/snapd/client"
)

// resultValue represents the value field in validation error responses from snapd.
type resultValue struct {
	Reasons            []string `json:"reasons"`
	EntropyBits        uint     `json:"entropy-bits"`
	MinEntropyBits     uint     `json:"min-entropy-bits"`
	OptimalEntropyBits uint     `json:"optimal-entropy-bits"`
}

// handleValidationError processes snapd validation errors and returns appropriate error messages.
func handleValidationError(err error, authMode string) error {
	snapdErr, ok := errors.AsType[*snapd.Error](err)
	if !ok {
		return fmt.Errorf("failed to check %s: %v", authMode, err)
	}

	switch snapdErr.Kind {
	case client.ErrorKindInvalidPassphrase, client.ErrorKindInvalidPIN:
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
	case client.ErrorKindUnsupportedByTargetSystem:
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

// IsValidPassphrase checks the entropy of the passphrase.
/* TODO:
 What is returned by snapd on invalid passphrase/PIN?
-> ask for forgivness, not permission

If it makes sense: drop all the validation here.

If the error message is "suboptimal", then, check in AddPassphrase/AddPIN… directly the validity of the passphrase and
turns those into private functions.

*/
func (s SnapTPM) IsValidPassphrase(ctx context.Context, passphrase string) error {
	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty, try again")
	}

	if err := s.snapdClient.CheckPassphrase(ctx, passphrase); err != nil {
		return handleValidationError(err, "passphrase")
	}

	return nil
}

// IsValidPIN checks that the PIN is only made of digits.
func (s SnapTPM) IsValidPIN(ctx context.Context, pin string) error {
	if pin == "" {
		return fmt.Errorf("PIN cannot be empty, try again")
	}

	// Check only digits in PIN
	for _, ch := range pin {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("PIN must contain only digits, try again")
		}
	}

	if err := s.snapdClient.CheckPIN(ctx, pin); err != nil {
		return handleValidationError(err, "PIN")
	}

	return nil
}

// ValidateAuthMode checks if the current authentication mode matches the expected mode.
func (s SnapTPM) ValidateAuthMode(ctx context.Context, expectedAuthMode snapd.AuthMode) error {
	result, err := s.snapdClient.ListVolumeInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to enumerate key slots: %v", err)
	}

	systemData, ok := result.ByContainerRole["system-data"]
	if !ok {
		return fmt.Errorf("system-data container role not found")
	}

	defaultKeyslot, ok := systemData.Keyslots["default"]
	if !ok {
		return fmt.Errorf("default key slot not found in system-data")
	}

	defaultFallbackKeyslot, ok := systemData.Keyslots["default-fallback"]
	if !ok {
		return fmt.Errorf("default-fallback key slot not found in system-data")
	}

	if defaultKeyslot.AuthMode != expectedAuthMode || defaultFallbackKeyslot.AuthMode != expectedAuthMode {
		return fmt.Errorf("authentication mode mismatch: expected %s, got default=%s, default-fallback=%s",
			expectedAuthMode,
			defaultKeyslot.AuthMode,
			defaultFallbackKeyslot.AuthMode,
		)
	}

	return nil
}

// ValidateRecoveryKeyNameUnique validates that a recovery key name is valid and not in use.
func (s SnapTPM) ValidateRecoveryKeyNameUnique(ctx context.Context, recoveryKeyName string) error {
	if err := ValidateRecoveryKeyName(ctx, recoveryKeyName); err != nil {
		return err
	}

	// Recovery key name cannot already be in use.
	result, err := s.snapdClient.ListVolumeInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to enumerate key slots: %v", err)
	}

	for _, volumeInfo := range result.ByContainerRole {
		for slotName := range volumeInfo.Keyslots {
			if slotName == recoveryKeyName {
				return fmt.Errorf("recovery key name %q is already in use", recoveryKeyName)
			}
		}
	}

	return nil
}

var isValidRecoveryKey = regexp.MustCompile("^([0-9]{5}-?){7}[0-9]{5}$").MatchString

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

// ValidateRecoveryKeyName validates that a recovery key name is valid.
func ValidateRecoveryKeyName(ctx context.Context, recoveryKeyName string) error {
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
