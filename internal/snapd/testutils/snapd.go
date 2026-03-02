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

// WithBaseURL configures the snapd socket path for the client.
//
//go:linkname WithBaseURL github.com/canonical/snap-tpmctl/internal/snapd.withBaseURL
func WithBaseURL(p string) snapd.Option
