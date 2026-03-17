// Package snapdtestutils exports testing functionalities used by other packages.
//
//nolint:gosec,revive // this package is used only in tests
package snapdtestutils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/internal/log"
	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/testutils"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
	"github.com/matryer/is"
)

func init() {
	testsdetection.MustBeTesting()
}

// withBaseURL configures the snapd socket path for the client.
//
//go:linkname withBaseURL github.com/canonical/snap-tpmctl/internal/snapd.withBaseURL
func withBaseURL(p string) snapd.Option

type RecordedRequest struct {
	Method string
	Path   string
	Body   string
}

type MockSnapdServer struct {
	*snapd.Client

	Requests       []RecordedRequest
	currentRequest int
}

// NewMockSnapdServer creates a new snapd client with a mock server that responds with the contents of the test file asset.
// The test file asset is searched in the provided roots in order, the first match is used as the response and than marked as done.
// If no match is found, a 404 response is returned.
func NewMockSnapdServer(t *testing.T, ctx context.Context, roots ...string) *MockSnapdServer {
	t.Helper()
	is := is.New(t)

	if roots == nil {
		roots = []string{testutils.TestPath(t)}
	}

	m := MockSnapdServer{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug(ctx, "Received request: %v %v", r.Method, r.URL.Path)
		m.currentRequest++

		b, err := io.ReadAll(r.Body)
		is.NoErr(err) // Server: could not read request body

		m.Requests = append(m.Requests, RecordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Body:   string(b),
		})

		// Look for file with path: <root>/<method>/<url-path>:<currentRequest> in each root, the first match is used as the response and than marked as done.
		var requests []string
		for _, p := range roots {
			requests = append(requests, fmt.Sprintf("%s:%d", filepath.Join(p, r.Method, r.URL.Path), m.currentRequest))
		}
		for _, p := range roots {
			requests = append(requests, filepath.Join(p, r.Method, r.URL.Path))
		}

		// Search for response and look for the test response file in each root until found.
		var resp []byte
		for _, r := range requests {
			resp, err = os.ReadFile(r)
			if os.IsNotExist(err) {
				continue
			}
			is.NoErr(err) // Setup: read the test response from test file asset

			break
		}

		if resp == nil {
			log.Debug(ctx, "Test response file not found for request: %v %v", r.Method, r.URL.Path)
			http.NotFound(w, r)
			return
		}

		log.Debug(ctx, "Returning back: %s", resp)

		var response struct {
			StatusCode int `json:"status-code"`
		}
		err = json.Unmarshal(resp, &response)
		is.NoErr(err) // Server: could not unmarshal test response

		w.WriteHeader(response.StatusCode)
		_, err = w.Write(resp)
		is.NoErr(err) // Server: could not write response to client
	}))
	t.Cleanup(ts.Close)

	m.Client = snapd.New(withBaseURL(ts.URL))
	return &m
}
