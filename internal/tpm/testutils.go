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

// withRoot allows you to specify a custom root for testing purposes.
func withRoot(root string) Option {
	testsdetection.MustBeTesting()
	return func(o *options) {
		o.root = root
	}
}

// withSyscall allows you to specify a custom root for testing purposes.
func withSyscall(s syscaller) Option {
	testsdetection.MustBeTesting()
	return func(o *options) {
		o.syscall = s
	}
}
