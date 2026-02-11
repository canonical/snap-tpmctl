// Package snapdtestutils exports testing functionalities used by other packages.
package snapdtestutils

import (
	"testing"
	_ "unsafe"

	"github.com/canonical/snap-tpmctl/internal/snapd"
)

func init() {
	if !testing.Testing() {
		panic("snapdtestutils should only be used in tests")
	}
}

// WithSocketPath configures the snapd socket path for the client.
//
//go:linkname WithSocketPath github.com/canonical/snap-tpmctl/internal/snapd.withSocketPath
func WithSocketPath(p string) snapd.Option
