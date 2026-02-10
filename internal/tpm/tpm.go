// Package tpm manages TPM/FDE features
package tpm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	// Needed for go:linkname.
	_ "unsafe"

	"github.com/canonical/snap-tpmctl/internal/tui"
	"github.com/snapcore/secboot"
)

// We are executing through a snap, thus we need to link to the actual package location.
//
//go:linkname systemdCryptsetupPath github.com/snapcore/secboot/internal/luks2.systemdCryptsetupPath
var systemdCryptsetupPath string

// MountVolume activates the specified encrypted volume using the provided device path.
func MountVolume(device, target string) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	volumeName := luksVolumeName(target)
	mapperPath := filepath.Join("/dev/mapper/", volumeName)

	// Check if volume is already active
	if _, err := os.Stat(mapperPath); os.IsNotExist(err) {
		if err := secboot.ActivateVolumeWithRecoveryKey(
			volumeName,
			device,
			&authRequestor{},
			&secboot.ActivateVolumeOptions{
				RecoveryKeyTries: 3,
			},
		); err != nil {
			return fmt.Errorf("unable to activate volume: %w", err)
		}
	}

	if err := os.MkdirAll(target, 0750); err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	if err := syscall.Mount(mapperPath, target, "ext4", 0, ""); err != nil {
		return fmt.Errorf("unable to mount volume: %w", err)
	}

	return nil
}

// UnmountVolume deactivates the specified volume.
func UnmountVolume(target string) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	if err := syscall.Unmount(target, 0); err != nil {
		return fmt.Errorf("unable to unmount volume: %w", err)
	}

	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("unable to remove mount point: %w", err)
	}

	volumeName := luksVolumeName(target)
	if err := secboot.DeactivateVolume(volumeName); err != nil {
		return fmt.Errorf("unable to deactivate volume: %w", err)
	}

	return nil
}

// luksVolumeName converts a directory path into a valid LUKS volume name.
func luksVolumeName(p string) string {
	return strings.TrimLeft(strings.ReplaceAll(p, "/", "-"), "-")
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
