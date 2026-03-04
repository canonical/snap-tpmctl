// Package snapdtestutils exports testing functionalities used by other packages.
package snapdtestutils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/matryer/is"
)

func init() {
	if !testing.Testing() {
		panic("snapdtestutils should only be used in tests")
	}
}

// withBaseURL configures the snapd socket path for the client.
//
//go:linkname withBaseURL github.com/canonical/snap-tpmctl/internal/snapd.withBaseURL
func withBaseURL(p string) snapd.Option

// NewMockSnapdServer creates a new snapd client with a mock server that responds with the contents of the test file asset.
func NewMockSnapdServer(t *testing.T, url string) *snapd.Client {
	t.Helper()
	is := is.New(t)

	resp, err := os.ReadFile(testutils.TestPath(t))
	is.NoErr(err) // Setup: read the test response from test file asset

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != url {
			http.NotFound(w, r)
			return
		}
		_, err := w.Write(resp)
		is.NoErr(err) // Server: could not write response to client
	}))
	t.Cleanup(ts.Close)

	c := snapd.New(withBaseURL(ts.URL))
	return c
}
