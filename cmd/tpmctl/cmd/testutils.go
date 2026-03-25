package cmd

import (
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/canonical/snap-tpmctl/internal/tpm"
)

// withArgs allows you to specify a custom args for testing purposes.
func withArgs(args ...string) Option {
	testsdetection.MustBeTesting()
	return func(o *option) {
		o.args = append([]string{""}, args...)
	}
}

// withSnapTPM allows you to specify a custom tpm for testing purposes.
func withSnapTPM(t tpm.SnapTPM) Option {
	testsdetection.MustBeTesting()
	return func(o *option) {
		o.tpm = t
	}
}
