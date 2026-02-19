package tpm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/snapcore/secboot"
)

// We are executing through a snap, thus we need to link to the actual package location.
//
//go:linkname systemdCryptsetupPath github.com/snapcore/secboot/internal/luks2.systemdCryptsetupPath
var systemdCryptsetupPath string

// MountVolume activates the specified encrypted volume using the provided device path.
func MountVolume(device, target string, authRequestor secboot.AuthRequestor) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	volumeName := luksVolumeName(device)
	mapperPath := filepath.Join("/dev/mapper/", volumeName)

	if err := os.MkdirAll(target, 0750); err != nil {
		// TODO: change all %w to %v
		return fmt.Errorf("unable to create directory: %v", err)
	}

	// Check if volume is already active
	if _, err := os.Stat(mapperPath); os.IsNotExist(err) {
		if err := secboot.ActivateVolumeWithRecoveryKey(
			volumeName,
			device,
			authRequestor,
			&secboot.ActivateVolumeOptions{
				RecoveryKeyTries: 3,
			},
		); err != nil {
			return fmt.Errorf("unable to activate volume: %v", err)
		}
	}

	if err := syscall.Mount(mapperPath, target, "ext4", syscall.MS_RELATIME, "rw"); err != nil {
		return fmt.Errorf("unable to mount volume: %w", err)
	}

	return nil
}

// UnmountVolume deactivates the specified volume.
func UnmountVolume(target string) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	// TODO: parse /proc/mounts to get the actual device path instead of relying on the target name
	device := target

	if err := syscall.Unmount(target, 0); err != nil {
		return fmt.Errorf("unable to unmount volume: %w", err)
	}

	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("unable to remove mount point: %w", err)
	}

	volumeName := luksVolumeName(device)
	if err := secboot.DeactivateVolume(volumeName); err != nil {
		return fmt.Errorf("unable to deactivate volume: %w", err)
	}

	return nil
}

// luksVolumeName converts a directory path into a valid LUKS volume name.
func luksVolumeName(p string) string {
	return strings.TrimLeft(strings.ReplaceAll(p, "/", "-"), "-")
}
