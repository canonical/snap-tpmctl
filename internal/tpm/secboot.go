package tpm

import (
	"os"
	"path/filepath"
	"syscall"

	_ "unsafe" // Needed for go:linkname.

	"github.com/snapcore/secboot"
)

// We are executing through a snap, thus we need to link to the actual package location.
//
//go:linkname systemdCryptsetupPath github.com/snapcore/secboot/internal/luks2.systemdCryptsetupPath
var systemdCryptsetupPath string

type Mounter interface {
	Mount(path, target string) error
	Unmount(target string) error
}

// Mount provides methods to interact with secboot features.
type Mount struct {
	authRequestor secboot.AuthRequestor
	mounter       Mounter
}

type mOptions struct {
	authRequestor secboot.AuthRequestor
	mounter       Mounter
}

// MountOption is a functional option for configuring the SecTPM.
type MountOption func(*mOptions)

// WithAuthRequestor configures secboot to use the authRequestor.
func WithAuthRequestor(a secboot.AuthRequestor) MountOption {
	return func(o *mOptions) {
		o.authRequestor = a
	}
}

// New creates a new SecTPM instance with the provided options.
func NewMount(args ...MountOption) Mount {
	if snapPath := os.Getenv("SNAP"); snapPath != "" {
		systemdCryptsetupPath = filepath.Join(snapPath, "usr/bin/systemd-cryptsetup")
	}

	o := mOptions{
		mounter: &defaultMounter{},
	}
	for _, f := range args {
		f(&o)
	}

	return Mount{
		authRequestor: o.authRequestor,
		mounter:       o.mounter,
	}
}

func (m Mount) ActivateVolume(volumeName, device string) error {
	return secboot.ActivateVolumeWithRecoveryKey(
		volumeName,
		device,
		m.authRequestor,
		&secboot.ActivateVolumeOptions{
			RecoveryKeyTries: 3,
		},
	)
}

func (m Mount) DeactivateVolume(volumeName string) error {
	return secboot.DeactivateVolume(volumeName)
}

type defaultMounter struct{}

func (m defaultMounter) Mount(path, target string) error {
	return syscall.Mount(path, target, "ext4", syscall.MS_RELATIME, "rw")
}

func (m defaultMounter) Unmount(target string) error {
	return syscall.Unmount(target, 0)
}
