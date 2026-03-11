package tpmtestutils

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"
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

// GetTestPath returns the test path based on the service.
func GetTestPath(t *testing.T, wantErr bool, service string) string {
	t.Helper()

	path := filepath.Join("testdata/snapdservice", service)
	if wantErr {
		path = "testdata/snapdservicefail"
	}

	return path
}

// HasBodyContent checks that at least one request contains all the expected body content
func HasBodyContent(is *is.I, requests []snapdtestutils.RecordedRequest, content ...string) bool {
	is.Helper()

	if content == nil {
		is.Fail()
	}

	return slices.ContainsFunc(requests, func(r snapdtestutils.RecordedRequest) bool {
		for _, c := range content {
			if !strings.Contains(r.Body, c) {
				return false
			}
		}
		return true
	})
}
