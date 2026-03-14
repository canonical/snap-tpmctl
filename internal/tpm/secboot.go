package tpm

import (
	"bufio"
	"fmt"
	"io/fs"
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

// FileSystem provides filesystem operations required by Mount.
type FileSystem interface {
	fs.StatFS

	MkdirAll(path string) error
	RemoveAll(path string) error
}

// Volume provides methods to activate and mount volumes.
type Volume interface {
	Activate(volumeName, device string, authRequestor secboot.AuthRequestor) error
	Deactivate(volumeName string) error
	Mount(path, target string) error
	Unmount(target string) error
}

// Mount provides methods to interact with secboot features.
type Mount struct {
	authRequestor secboot.AuthRequestor
	fs            FileSystem
	vol           Volume
}

type mOptions struct {
	authRequestor secboot.AuthRequestor
	fs            FileSystem
	vol           Volume
}

// MountOption is a functional option for configuring the SecTPM.
type MountOption func(*mOptions)

// WithAuthRequestor configures secboot to use the authRequestor.
func WithAuthRequestor(a secboot.AuthRequestor) MountOption {
	return func(o *mOptions) {
		o.authRequestor = a
	}
}

// NewMount creates a new Mount instance with the provided options.
func NewMount(args ...MountOption) Mount {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	base := os.DirFS("/")
	statfs := base.(fs.StatFS) //nolint:forcetypeassert // fs.FS is documented to implement fs.StatFS

	o := mOptions{
		fs:  &defaultFileSystem{statfs},
		vol: &defaultVolume{},
	}
	for _, f := range args {
		f(&o)
	}

	return Mount(o)
}

// getDeviceFromMount parses /proc/mounts and returns the device path for the given mount point.
func (m Mount) getDeviceFromMount(mountPoint string) (string, error) {
	file, err := m.fs.Open("proc/mounts")
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

type defaultFileSystem struct {
	fs.StatFS
}

func (fs *defaultFileSystem) MkdirAll(path string) error {
	return os.MkdirAll(hostPath(path), 0750)
}

func (fs defaultFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(hostPath(path))
}

type defaultVolume struct{}

func (v defaultVolume) Activate(volumeName, device string, authRequestor secboot.AuthRequestor) error {
	return secboot.ActivateVolumeWithRecoveryKey(
		volumeName,
		device,
		authRequestor,
		&secboot.ActivateVolumeOptions{
			RecoveryKeyTries: 3,
		},
	)
}

func (v defaultVolume) Deactivate(volumeName string) error {
	return secboot.DeactivateVolume(volumeName)
}

func (v defaultVolume) Mount(path, target string) error {
	return syscall.Mount(hostPath(path), target, "ext4", syscall.MS_RELATIME, "rw")
}

func (v defaultVolume) Unmount(target string) error {
	return syscall.Unmount(target, 0)
}

func hostPath(p string) string {
	return filepath.Join("/", p)
}
