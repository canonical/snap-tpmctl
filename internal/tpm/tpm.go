// Package tpm manages TPM/FDE features
package tpm

import (
	"context"
	"fmt"
	"syscall"
	_ "unsafe" // Needed for go:linkname.

	"github.com/canonical/snap-tpmctl/internal/snapd"
)

// SnapTPM provides methods to interact with TPM/FDE features via snapd.
type SnapTPM struct {
	options
}

type options struct {
	snapdClient *snapd.Client

	root    string
	syscall syscaller
}

// Option is a functional option for configuring the SnapTPM.
type Option func(*options)

// syscaller abstracts mount and unmount system calls used by SnapTPM.
type syscaller interface {
	Mount(path, target string) error
	Unmount(target string) error
}

// New creates a new SnapTPM instance with the provided options.
func New(args ...Option) SnapTPM {
	o := options{
		snapdClient: snapd.New(),

		root:    "/",
		syscall: defaultSyscall{},
	}
	for _, f := range args {
		f(&o)
	}

	return SnapTPM{
		options: o,
	}
}

// FdeStatus retrieves the Full Disk Encryption status from snapd.
func (s SnapTPM) FdeStatus(ctx context.Context) (string, error) {
	status, err := s.snapdClient.FdeStatus(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve the FDE status: %v", err)
	}

	return status, nil
}

// ListVolumeInfo retrieves information about system volumes.
func (s SnapTPM) ListVolumeInfo(ctx context.Context) (result snapd.SystemVolumesResult, err error) {
	result, err = s.snapdClient.ListVolumeInfo(ctx)
	if err != nil {
		return result, fmt.Errorf("failed to retrieve the volume info: %v", err)
	}

	return result, nil
}

type defaultSyscall struct{}

func (defaultSyscall) Mount(path, target string) error {
	return syscall.Mount(path, target, "ext4", syscall.MS_RELATIME, "rw")
}
func (defaultSyscall) Unmount(target string) error {
	return syscall.Unmount(target, 0)
}
