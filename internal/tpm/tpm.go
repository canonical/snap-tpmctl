// Package tpm manages TPM/FDE features
package tpm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	// Needed for go:linkname.
	_ "unsafe"

	"github.com/snapcore/secboot"
	"snap-tpmctl/internal/tui"
)

// We are executing through a snap, thus we need to link to the actual package location.
//
//go:linkname systemdCryptsetupPath github.com/snapcore/secboot/internal/luks2.systemdCryptsetupPath
var systemdCryptsetupPath string

// MountVolume activates the specified encrypted volume using the provided device path.
func MountVolume(volumeName string, devicePath string) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}
	return secboot.ActivateVolumeWithRecoveryKey(
		volumeName,
		devicePath,
		&authRequestor{},
		&secboot.ActivateVolumeOptions{
			RecoveryKeyTries: 3,
		},
	)
}

// UnmountVolume deactivates the specified volume.
func UnmountVolume(volumeName string) error {
	return secboot.DeactivateVolume(volumeName)
}

type authRequestor struct{}

func (r *authRequestor) RequestUserCredential(ctx context.Context, name, path string, authTypes secboot.UserAuthType) (string, error) {
	if authTypes != secboot.UserAuthTypeRecoveryKey {
		return "", fmt.Errorf("authentication type not supported")
	}

	key, err := tui.ReadUserSecret("Enter recovery key: ")
	if err != nil {
		return "", err
	}

	return key, nil
}

// GetLuksKey validates and converts the recovery key to a binary key format by parsing and formatting it.
func GetLuksKey(recoveryKey string) (secboot.DiskUnlockKey, error) {
	binKey, err := secboot.ParseRecoveryKey(recoveryKey)
	if err != nil {
		return nil, err
	}

	return binKey[:], nil
}
