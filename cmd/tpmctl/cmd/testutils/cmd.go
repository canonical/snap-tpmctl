// Package cmdtestutils provides helpers for command-related tests.
package cmdtestutils

import (
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/canonical/snap-tpmctl/internal/tui"
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

// WithTui is an option that configures the app to use the provided Tui.
//
//go:linkname WithTui github.com/canonical/snap-tpmctl/cmd/tpmctl/cmd.withTui
func WithTui(t tui.Tui) cmd.Option
