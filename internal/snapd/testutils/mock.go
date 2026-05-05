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

type MockSnapdIntegrationServer struct {
	*snapd.Client

	Requests        []RecordedRequest
	currentRequests map[string]int
}

// NewMockSnapdServer creates a new snapd client with a mock server that responds with the contents of the test file asset.
// We are looking first at root/<method>/<url-path>:<currentRequest> and fallacbk to
//
//	root/<method>/<url-path> for the test response file asset, where <method> is the HTTP method of the request,
//
// <url-path> is the URL path of the request and <currentRequest> is the number of times that a request with
// the same method and URL path has been received by the server.
// If no match is found, a 404 response is returned.
func NewMockSnapdIntegrationServer(root string) *MockSnapdIntegrationServer {
	m := MockSnapdIntegrationServer{
		currentRequests: map[string]int{},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic("err")
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
				panic("err")
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
			panic("err")
		}

		w.WriteHeader(response.StatusCode)
		_, err = w.Write(resp)
		if err != nil {
			panic("err")
		}

	}))

	m.Client = snapd.New(withBaseURL(ts.URL))
	return &m
}
