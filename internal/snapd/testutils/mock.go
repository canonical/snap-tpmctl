//go:build integrationtests

package snapdtestutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	_ "unsafe" // Required for go:linkname directives

	"github.com/canonical/snap-tpmctl/internal/snapd"
	"github.com/canonical/snap-tpmctl/internal/testutils/testsdetection"
)

func init() {
	testsdetection.MustBeTesting()
}

// NewMockSnapdServerWithPath creates a new snapd client with a mock server that responds with the contents of the test file asset from the specified path.
// This is almost identical to NewMockSnapdServer, but it's meant to be used only for integration tests.
func NewMockSnapdServerWithPath(root string) *MockSnapdServer {
	m := MockSnapdServer{
		currentRequests: map[string]int{},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(fmt.Sprintf("read request body: %v", err))
		}

		m.Requests = append(m.Requests, RecordedRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Body:   string(b),
		})

		// Search for response in <root>/<method>/<url-path>:<currentRequest> and fallback to <root>/<method>/<url-path>.
		var resp []byte
		uri := filepath.Join(root, r.Method, r.URL.Path)
		m.currentRequests[uri]++
		for _, r := range []string{fmt.Sprintf("%s:%d", uri, m.currentRequests[uri]), uri} {
			resp, err = os.ReadFile(r)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				panic(fmt.Sprintf("read mock response file %q: %v", r, err))
			}

			break
		}

		if resp == nil {
			http.NotFound(w, r)
			return
		}

		var response struct {
			StatusCode int `json:"status-code"`
		}
		err = json.Unmarshal(resp, &response)
		if err != nil {
			panic(fmt.Sprintf("decode mock response JSON: %v", err))
		}

		w.WriteHeader(response.StatusCode)
		_, err = w.Write(resp)
		if err != nil {
			panic(fmt.Sprintf("write mock response: %v", err))
		}

	}))

	m.Client = snapd.New(withBaseURL(ts.URL))
	return &m
}
