// Package tpmtestutils provides helpers for TPM-related tests.
package tpmtestutils

import (
	"slices"
	"strings"
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/internal/snapd"
	snapdtestutils "github.com/canonical/snap-tpmctl/internal/snapd/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/canonical/snap-tpmctl/internal/tpm"
	"github.com/matryer/is"
)

func init() {
	testsdetection.MustBeTesting()
}

// WithSnapdClient is an option that configures the TPM to use the provided snapd client.
//
//go:linkname WithSnapdClient github.com/canonical/snap-tpmctl/internal/tpm.withSnapdClient
func WithSnapdClient(snapdClient *snapd.Client) tpm.Option

// WithRoot is an option that configures the TPM to use the provided system mounter.
//
//go:linkname WithRoot github.com/canonical/snap-tpmctl/internal/tpm.withRoot
func WithRoot(r string) tpm.Option

// WithSyscall is an option that configures the TPM to use the provided system mounter.
//
//go:linkname WithSyscall github.com/canonical/snap-tpmctl/internal/tpm.withSyscall
func WithSyscall(s tpm.Syscall) tpm.Option

// OneRequestBodyContains checks that at least one request contains all the expected wanted contents.
func OneRequestBodyContains(is *is.I, requests []snapdtestutils.RecordedRequest, wants ...string) {
	is.Helper()

	if len(wants) == 0 {
		panic("Programmer error: RequestsBodyContains checked for nothings as it doesn’t have any wants")
	}

	is.True(slices.ContainsFunc(requests, func(r snapdtestutils.RecordedRequest) bool {
		for _, want := range wants {
			if !strings.Contains(r.Body, want) {
				return false
			}
		}
		return true
	})) // Didn't find all wants in any of the requests.
}
