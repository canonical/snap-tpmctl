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

// withFileSystem allows you to specify a custom filesystem for testing purposes.
func withFileSystem(f FileSystem) MountOption {
	testsdetection.MustBeTesting()
	return func(mo *mOptions) {
		mo.fs = f
	}
}

// withVolume allows you to specify a custom system mounter for testing purposes.
func withVolume(m Volume) MountOption {
	testsdetection.MustBeTesting()
	return func(mo *mOptions) {
		mo.vol = m
	}
}
