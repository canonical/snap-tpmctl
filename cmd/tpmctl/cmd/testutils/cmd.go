// Package cmdtestutils provides helpers for command-related tests.
package cmdtestutils

import (
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/canonical/snap-tpmctl/internal/tpm"
)

func init() {
	testsdetection.MustBeTesting()
}

// WithArgs is an option that configures the app to use the provided arguments.
//
//go:linkname WithArgs github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd.withArgs
func WithArgs(args ...string) cmd.Option

// WithSnapTPM is an option that configures the app to use the provided snap TPM.
//
//go:linkname WithSnapTPM github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd.withSnapTPM
func WithSnapTPM(t tpm.SnapTPM) cmd.Option
