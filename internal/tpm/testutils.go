//nolint:unused // helper functions used only in tests
package tpm

import (
	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
)

// withSnapdClient allows you to specify a custom snapd client for testing purposes.
func withSnapdClient(snapdClient *snapd.Client) Option {
	testsdetection.MustBeTesting()
	return func(opts *options) {
		opts.snapdClient = snapdClient
	}
}

// withActivator allows you to specify a custom TPM volume activator for testing purposes.
func withActivator(a Activator) MountOption {
	testsdetection.MustBeTesting()
	return func(mo *mOptions) {
		mo.activator = a
	}
}

// withFileSystem allows you to specify a custom filesystem for testing purposes.
func withFileSystem(f FileSystem) MountOption {
	testsdetection.MustBeTesting()
	return func(mo *mOptions) {
		mo.filesystem = f
	}
}

// withMounter allows you to specify a custom system mounter for testing purposes.
func withMounter(m Mounter) MountOption {
	testsdetection.MustBeTesting()
	return func(mo *mOptions) {
		mo.mounter = m
	}
}
