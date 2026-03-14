package tpm

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// MountVolume activates the specified encrypted volume using the provided device path.
func (m Mount) MountVolume(ctx context.Context, device, target string) error {
	volumeName := luksVolumeName(device)
	mapperPath := filepath.Join("dev/mapper/", volumeName)

	if err := m.fs.MkdirAll(target); err != nil {
		return fmt.Errorf("unable to create directory: %v", err)
	}

	// Check if volume is already active
	if _, err := m.fs.Stat(mapperPath); errors.Is(err, fs.ErrNotExist) {
		if err := m.vol.Activate(volumeName, device, m.authRequestor); err != nil {
			return fmt.Errorf("unable to activate volume: %v", err)
		}
	}

	if err := m.vol.Mount(mapperPath, target); err != nil {
		return fmt.Errorf("unable to mount volume: %v", err)
	}

	return nil
}

// UnmountVolume deactivates the specified volume.
func (m Mount) UnmountVolume(ctx context.Context, target string) error {
	device, err := m.getDeviceFromMount(target)
	if err != nil {
		return fmt.Errorf("unable to determine device path: %v", err)
	}

	if err := m.vol.Unmount(target); err != nil {
		return fmt.Errorf("unable to unmount volume: %v", err)
	}

	if err := m.fs.RemoveAll(target); err != nil {
		return fmt.Errorf("unable to remove mount point: %v", err)
	}

	volumeName := filepath.Base(device)
	if err := m.vol.Deactivate(volumeName); err != nil {
		return fmt.Errorf("unable to deactivate volume: %v", err)
	}

	return nil
}

// luksVolumeName converts a directory path into a valid LUKS volume name.
func luksVolumeName(p string) string {
	return strings.TrimLeft(strings.ReplaceAll(p, "/", "-"), "-")
}
