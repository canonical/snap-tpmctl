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
