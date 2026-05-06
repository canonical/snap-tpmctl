package tpm

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	_ "unsafe" // Needed for go:linkname.

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/snapcore/secboot"
)

// We are executing through a snap, thus we need to link to the actual package location.
//
//go:linkname systemdCryptsetupPath github.com/snapcore/secboot/internal/luks2.systemdCryptsetupPath
var systemdCryptsetupPath string

// Mount activates and mounts the TPM-protected volume at the given path to the target mount point.
func (s SnapTPM) Mount(ctx context.Context, device, target string, authRequestor secboot.AuthRequestor) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	// Check if the volume is active and mapped by other tools
	p, err := s.getMapperFromDevice(device)
	if err != nil {
		return fmt.Errorf("unable to locate device: %v", err)
	}
	if p != "" {
		return fmt.Errorf("unable to activate device: resource is already mapped as %q", p)
	}

	volumeName := luksVolumeName(device)
	mapperPath := filepath.Join(s.root, "dev", "mapper", volumeName)

	if err := os.MkdirAll(target, 0750); err != nil {
		return fmt.Errorf("unable to create directory: %v", err)
	}

	// Check if the volume is already mounted by the tool
	p, err = s.getMountFromMapper(mapperPath)
	if err != nil {
		return fmt.Errorf("unable to locate volume: %v", err)
	}
	if p != "" {
		return fmt.Errorf("unable to activate volume: resource is already mounted as %q", p)
	}

	// Check if volume is already active
	if _, err := os.Stat(mapperPath); os.IsNotExist(err) {
		if err := secboot.ActivateVolumeWithRecoveryKey(
			volumeName,
			device,
			authRequestor,
			&secboot.ActivateVolumeOptions{
				RecoveryKeyTries: 3,
			}); err != nil {
			return fmt.Errorf("unable to activate volume: %v", err)
		}
	}

	log.Debug(ctx, "Mounting %q to %q", mapperPath, target)
	if err := s.syscall.Mount(mapperPath, target); err != nil {
		return fmt.Errorf("unable to mount volume: %v", err)
	}

	return nil
}

// Unmount unmounts and deactivate the TPM-protected volume from the target mount point.
func (s SnapTPM) Unmount(ctx context.Context, target string) error {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	mapperPath, err := s.getMapperFromMount(target)
	if err != nil {
		return fmt.Errorf("unable to determine device path: %v", err)
	}
	if mapperPath == "" {
		return errors.New("path not found in /proc/mounts")
	}

	if err := s.syscall.Unmount(target); err != nil {
		return fmt.Errorf("unable to unmount volume: %v", err)
	}

	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("unable to remove mount point: %v", err)
	}

	volumeName := filepath.Base(mapperPath)
	if err := secboot.DeactivateVolume(volumeName); err != nil {
		return fmt.Errorf("unable to deactivate volume: %v", err)
	}

	return nil
}

type mountsFieldType int

const (
	// deviceField is the field index for the device in /proc/mounts.
	deviceField mountsFieldType = iota
	// mountPointField is the field index for the mount point in /proc/mounts.
	mountPointField
)

// getMapperFromMount parses /proc/mounts and returns the mapper path for the given mount point.
func (s SnapTPM) getMapperFromMount(mountPoint string) (string, error) {
	return s.searchInProcMounts(mountPoint, mountPointField, deviceField)
}

// mountFromMapper parses /proc/mounts and returns the mapper path for the given mount point.
func (s SnapTPM) getMountFromMapper(mapperPath string) (string, error) {
	return s.searchInProcMounts(mapperPath, deviceField, mountPointField)
}

// searchInProcMounts searches /proc/mounts for a specific path in the specified field and returns the corresponding field result.
func (s SnapTPM) searchInProcMounts(path string, fieldPath, fieldResult mountsFieldType) (string, error) {
	file, err := os.Open(filepath.Join(s.root, "proc", "mounts"))
	if err != nil {
		return "", fmt.Errorf("unable to open /proc/mounts: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		// Each line format: device mount_point fstype options dummy dummy
		if len(fields) < 2 {
			continue
		}

		if fields[fieldPath] == path {
			return fields[fieldResult], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading /proc/mounts: %v", err)
	}

	return "", nil
}

// getMapperFromDevice checks if a device is already mapped by other tools by reading /sys/class/block/holders.
// It returns the mapper path if the device is in use, or an empty string if not.
func (s SnapTPM) getMapperFromDevice(device string) (string, error) {
	holdersPath := filepath.Join(s.root, "sys", "class", "block", filepath.Base(device), "holders")
	holders, err := os.ReadDir(holdersPath)
	if err != nil {
		return "", fmt.Errorf("unable to read holders: %v", err)
	}

	if len(holders) == 0 {
		return "", nil
	}

	// /sys/class/block/<holder>/dm/name contains the device mapper name.
	dmNamePath := filepath.Join(s.root, "sys", "class", "block", holders[0].Name(), "dm", "name")
	mapperName, err := os.ReadFile(dmNamePath)
	if err != nil {
		return "", fmt.Errorf("unable to read mapper: %v", err)
	}

	return filepath.Join(s.root, "dev", "mapper", strings.TrimSpace(string(mapperName))), nil
}

// luksVolumeName converts a directory path into a valid LUKS volume name.
func luksVolumeName(p string) string {
	return strings.TrimLeft(strings.ReplaceAll(p, "/", "-"), "-")
}
