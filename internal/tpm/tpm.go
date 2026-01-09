// Package tpm manages TPM/FDE features
package tpm

import (
	"context"
	"fmt"

	"github.com/snapcore/secboot"
	"snap-tpmctl/internal/tui"
)

// MountVolume activates the specified encrypted volume using the provided device path.
func MountVolume(volumeName string, devicePath string) error {
	return secboot.ActivateVolumeWithKeyData(
		volumeName,
		devicePath,
		&authRequestor{},
		&secboot.ActivateVolumeOptions{
			PassphraseTries:  3,
			RecoveryKeyTries: 3,
		})
}

// UnmountVolume deactivates the specified volume.
func UnmountVolume(volumeName string) error {
	return secboot.DeactivateVolume(volumeName)
}

type authRequestor struct{}

func (r *authRequestor) RequestUserCredential(ctx context.Context, name, path string, authTypes secboot.UserAuthType) (string, error) {
	var prompt string

	switch authTypes {
	case secboot.UserAuthTypePassphrase:
		prompt = "Enter passphrase: "
	case secboot.UserAuthTypePIN:
		prompt = "Enter pin: "
	case secboot.UserAuthTypeRecoveryKey:
		prompt = "Enter recovery key: "
	default:
		return "", fmt.Errorf("authentication type not supported")
	}

	key, err := tui.ReadUserSecret(prompt)
	if err != nil {
		return "", err
	}

	return key, nil
}
