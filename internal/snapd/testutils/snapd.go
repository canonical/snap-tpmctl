// Package snapdtestutils exports testing functionalities used by other packages.
//
//nolint:gosec,revive // this package is used only in tests
package snapdtestutils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
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

	Requests *[]RecordedRequest
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

	var done []string
	recordedRequests := new([]RecordedRequest)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug(ctx, "Received request: %v %v", r.Method, r.URL.Path)

		b, err := io.ReadAll(r.Body)
		is.NoErr(err) // Server: could not read request body

		*recordedRequests = append(*recordedRequests, RecordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Body:   string(b),
		})

		// Search for response from the last root and look for the test response file in each root until found.
		var resp []byte
		for _, p := range roots {
			request := filepath.Join(p, r.Method, r.URL.Path)

			if slices.Contains(done, request) {
				continue
			}

			resp, err = os.ReadFile(request)
			if os.IsNotExist(err) {
				continue
			}
			is.NoErr(err) // Setup: read the test response from test file asset

			done = append(done, request)
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

	c := snapd.New(withBaseURL(ts.URL))
	return &MockSnapdServer{
		Client:   c,
		Requests: recordedRequests,
	}
}
