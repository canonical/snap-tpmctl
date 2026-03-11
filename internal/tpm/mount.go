package tpm

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	_ "unsafe" // Needed for go:linkname.

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
		return fmt.Errorf("unable to mount volume: %v", err)
	}

	return nil
}

// UnmountVolume deactivates the specified volume.
func UnmountVolume(target string) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	device, err := getDeviceFromMount(target)
	if err != nil {
		return fmt.Errorf("unable to determine device path: %v", err)
	}

	if err := syscall.Unmount(target, 0); err != nil {
		return fmt.Errorf("unable to unmount volume: %v", err)
	}

	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("unable to remove mount point: %v", err)
	}

	volumeName := filepath.Base(device)
	if err := secboot.DeactivateVolume(volumeName); err != nil {
		return fmt.Errorf("unable to deactivate volume: %v", err)
	}

	return nil
}

// luksVolumeName converts a directory path into a valid LUKS volume name.
func luksVolumeName(p string) string {
	return strings.TrimLeft(strings.ReplaceAll(p, "/", "-"), "-")
}

// getDeviceFromMount parses /proc/mounts and returns the device path for the given mount point.
func getDeviceFromMount(mountPoint string) (string, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return "", fmt.Errorf("unable to open /proc/mounts: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		// Each line format: device mount_point fstype options dummy dummy
		if len(fields) >= 2 && fields[1] == mountPoint {
			return fields[0], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading /proc/mounts: %v", err)
	}

	return "", fmt.Errorf("mount point %q doesn't exist", mountPoint)
}
