// Package snapd provides a client for making calls to the systems local snapd service
package snapd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	snapdClient "github.com/snapcore/snapd/client"
)

const (
	defaultSocketPath = "/var/run/snapd.socket"
	defaultUserAgent  = "snapd.go"
)

// Response is the base response structure from snapd.
// TODO: make this struct private, but export for test mocks
type Response struct {
	Type       string          `json:"type"`
	StatusCode int             `json:"status-code"`
	Status     string          `json:"status"`
	Result     json.RawMessage `json:"result,omitempty"`
	Change     string          `json:"change,omitempty"`
}

// Result is the result structure returned from snapd in a response.
type Result struct {
	Kind    string          `json:"kind"`
	Message string          `json:"message"`
	Value   json.RawMessage `json:"value,omitempty"`
}

// IsOK checks if a commonly know snapd accepted status was returned.
func (r *Response) IsOK() bool {
	return r.Status == "Accepted" || r.Status == "OK" || r.StatusCode == 200 || r.StatusCode == 202
}

// TODO: better fields parsing with status-code and type

// AsyncResponse represents the status of a change.
type AsyncResponse struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
	Ready   bool   `json:"ready"`
	Err     string `json:"err,omitempty"`
	// Tasks   json.RawMessage `json:"tasks,omitempty"`

}

// IsOK checks if the asynchronous operation completed successfully.
func (r *AsyncResponse) IsOK() bool {
	return r.Ready && r.Status == "Done"
}

// Error represents an error from snapd.
type Error struct {
	Message    string
	Kind       string
	StatusCode int
	Status     string
	Value      json.RawMessage
}

func (e *Error) Error() string {
	if e.Kind != "" {
		return fmt.Sprintf("snapd error: %s (%s)", e.Message, e.Kind)
	}
	return fmt.Sprintf("snapd error: %s", e.Message)
}

// NewResponseBody parses a JSON response body from snapd and returns a Response.
// If the response type is "error", it extracts error details from the Result field and returns an Error.
func (c *Client) NewResponseBody(body []byte) (*Response, error) {
	var snapdResp Response
	if err := json.Unmarshal(body, &snapdResp); err != nil {
		return nil, err
	}

	if snapdResp.Type == "error" {
		var errResp struct {
			Message string          `json:"message"`
			Kind    string          `json:"kind,omitempty"`
			Value   json.RawMessage `json:"value,omitempty"`
		}

		if err := json.Unmarshal(snapdResp.Result, &errResp); err != nil {
			return nil, err
		}

		return nil, &Error{
			Message:    errResp.Message,
			Kind:       errResp.Kind,
			StatusCode: snapdResp.StatusCode,
			Status:     snapdResp.Status,
			Value:      errResp.Value,
		}
	}

	return &snapdResp, nil
}

// Client is a snapd client.
type Client struct {
	snapd *snapdClient.Client
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// NewClient creates a new snapd client.
func NewClient(opts ...ClientOption) *Client {
	return &Client{
		snapd: snapdClient.New(nil),
	}
}

func (c *Client) doSyncRequest(ctx context.Context, method, path string, query url.Values, headers map[string]string, body io.Reader) (*Response, error) {
	var resp response
	_, err := doSync(c.snapd, method, path, query, headers, body, &resp)
	var snapdErr snapdClient.Error
	if errors.As(err, &snapdErr) {
		return nil, fmt.Errorf("%w %d", snapdErr, snapdErr.StatusCode)
	}
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	r, err := c.NewResponseBody(bodyBytes)
	if err != nil {
		return nil, err
	}

	return r, err
}

func (c *Client) doAsyncRequest(ctx context.Context, method, path string, query url.Values, headers map[string]string, body io.Reader) (*AsyncResponse, error) {
	var resp response
	_, err := do(c.snapd, method, path, query, headers, body, &resp, nil)
	if err != nil {
		return nil, err
	}

	if resp.Type != "async" {
		return nil, fmt.Errorf("expected async response for %q on %q, got %q", method, path, resp.Type)
	}
	if resp.Change == "" {
		return nil, fmt.Errorf("async response without change reference")
	}

	// changeId, err := doAsync(c.snapd, method, path, query, headers, body)
	// if err != nil {
	// 	return nil, err
	// }

	// TODO: find a way to do it without polling (?)
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:

			change, err := c.snapd.Change(resp.Change)

			if err != nil {
				return nil, err
			}

			if change.Ready {
				return &AsyncResponse{}, nil
			}
		}
	}
}

// go:linkname doSync github.com/snapcore/snap/client.(*Client).doSync
func doSync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader, v any) (*snapdClient.ResultInfo, error)

// // go:linkname doAsync github.com/snapcore/snap/client.(*Client).doAsync
// func doAsync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader) (changeID string, err error)

// go:linkname do github.com/snapcore/snap/client.(*Client).do
func do(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader, v any, opts *doOptions) (statusCode int, err error)

type doOptions struct {
	// Timeout is the overall request timeout
	Timeout time.Duration
	// Retry interval.
	// Note for a request with a Timeout but without a retry, Retry should just
	// be set to something larger than the Timeout.
	Retry time.Duration
}

// A response produced by the REST API will usually fit in this
// (exceptions are the icons/ endpoints obvs)
type response struct {
	Result json.RawMessage `json:"result"`
	Type   string          `json:"type"`
	Change string          `json:"change"`

	WarningCount     int       `json:"warning-count"`
	WarningTimestamp time.Time `json:"warning-timestamp"`

	snapdClient.ResultInfo

	Maintenance *snapdClient.Error `json:"maintenance"`
}
