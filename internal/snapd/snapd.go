// Package snapd provides a client for making calls to the systems local snapd service
package snapd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"
	_ "unsafe" // Required for go:linkname directives

	snapdClient "github.com/snapcore/snapd/client"
)

const (
	defaultSocketPath = "/var/run/snapd.socket"
	defaultUserAgent  = "snapd.go"
)

// Response is the base response structure from snapd.
// TODO: make this struct private, but export for test mocks.
type Response struct {
	ErrorMessage string
	Result       json.RawMessage
	StatusCode   int
}

// IsOK checks if a commonly know snapd accepted status was returned.
func (r *Response) IsOK() bool {
	return r.StatusCode == 200 || r.StatusCode == 202
}

// AsyncResponse represents the status of a change.
type AsyncResponse struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
	Ready   bool   `json:"ready"`
	Err     string `json:"err,omitempty"`
}

// IsOK checks if the asynchronous operation completed successfully.
func (r *AsyncResponse) IsOK() bool {
	return r.Ready && r.Status == "Done"
}

// Error represents an error from snapd.
type Error struct {
	Message string
	Kind    string
	Value   json.RawMessage
}

func (e *Error) Error() string {
	if e.Kind != "" {
		return fmt.Sprintf("snapd error: %s (%s)", e.Message, e.Kind)
	}
	return fmt.Sprintf("snapd error: %s", e.Message)
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
		snapd: snapdClient.New(&snapdClient.Config{
			Interactive: true,
			Socket:      defaultSocketPath,
			UserAgent:   defaultUserAgent,
		}),
	}
}

// newRequestBody marshals the given body into JSON format and returns it as an io.Reader.
func (c *Client) newRequestBody(body any) (io.Reader, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewReader(data)
	}
	return reqBody, nil
}

// setGenericHeaders sets the common HTTP headers for snapd API requests.
func (c *Client) setGenericHeaders(headers map[string]string) map[string]string {
	if headers == nil {
		headers = map[string]string{}
	}

	headers["Content-Type"] = "application/json"

	return headers
}

//nolint:unparam // path parameter kept for future extensibility
func (c *Client) doSyncRequest(_ context.Context, method, path string, query url.Values, headers map[string]string, body any) (*Response, error) {
	b, err := c.newRequestBody(body)
	if err != nil {
		return nil, err
	}

	var result json.RawMessage
	_, err = doSync(c.snapd, method, path, query, c.setGenericHeaders(headers), b, &result)
	var snapdErr *snapdClient.Error
	if errors.As(err, &snapdErr) {
		return nil, &Error{
			Kind:    string(snapdErr.Kind),
			Message: snapdErr.Message,
			Value:   json.RawMessage(snapdErr.Error()),
		}
	}
	if err != nil {
		return nil, err
	}

	return &Response{Result: result, StatusCode: 200}, nil
}

func (c *Client) doAsyncRequest(ctx context.Context, method, path string, query url.Values, headers map[string]string, body any) (*AsyncResponse, error) {
	b, err := c.newRequestBody(body)
	if err != nil {
		return nil, err
	}

	changeID, err := doAsync(c.snapd, method, path, query, c.setGenericHeaders(headers), b)
	var snapdErr *snapdClient.Error
	if errors.As(err, &snapdErr) {
		return nil, &Error{
			Kind:    string(snapdErr.Kind),
			Message: snapdErr.Message,
			Value:   json.RawMessage(snapdErr.Error()),
		}
	}
	if err != nil {
		return nil, err
	}

	// TODO: use notices api
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			change, err := c.snapd.Change(changeID)
			if err != nil {
				return nil, err
			}

			if change.Ready {
				return &AsyncResponse{
					ID:      change.ID,
					Kind:    change.Kind,
					Summary: change.Summary,
					Status:  change.Status,
					Ready:   change.Ready,
					Err:     change.Err,
				}, nil
			}
		}
	}
}

//go:linkname doSync github.com/snapcore/snapd/client.(*Client).doSync
func doSync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader, v any) (*snapdClient.ResultInfo, error)

//go:linkname doAsync github.com/snapcore/snapd/client.(*Client).doAsync
func doAsync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader) (changeID string, err error)
