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

// Error represents an error from snapd.
type Error struct {
	Message string
	Kind    snapdClient.ErrorKind
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

// setGenericHeaders sets the common HTTP headers for snapd API requests.
func (c *Client) setGenericHeaders(headers map[string]string) map[string]string {
	if headers == nil {
		headers = map[string]string{}
	}

	headers["Content-Type"] = "application/json"

	return headers
}

// response is the base response structure from snapd.
type response struct {
	Result json.RawMessage
}

//nolint:unparam // path parameter kept for future extensibility
func (c *Client) doSyncRequest(_ context.Context, method, path string, query url.Values, headers map[string]string, body any) (*response, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&body); err != nil {
		return nil, err
	}

	var result json.RawMessage
	_, err := doSync(c.snapd, method, path, query, c.setGenericHeaders(headers), &b, &result)
	var snapdErr *snapdClient.Error
	if errors.As(err, &snapdErr) {
		value, err := json.Marshal(snapdErr.Value)
		if err != nil {
			return nil, err
		}

		return nil, &Error{
			Kind:    snapdErr.Kind,
			Message: snapdErr.Message,
			Value:   value,
		}
	}
	if err != nil {
		return nil, err
	}

	return &response{Result: result}, nil
}

// asyncResponse represents the status of a change.
type asyncResponse struct {
	ID string
}

//nolint:unparam // asyncResponse parameter kept for future extensibility
func (c *Client) doAsyncRequest(ctx context.Context, method, path string, query url.Values, headers map[string]string, body any) (*asyncResponse, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&body); err != nil {
		return nil, err
	}

	changeID, err := doAsync(c.snapd, method, path, query, c.setGenericHeaders(headers), &b)
	var snapdErr *snapdClient.Error
	if errors.As(err, &snapdErr) {
		value, err := json.Marshal(snapdErr.Value)
		if err != nil {
			return nil, err
		}

		return nil, &Error{
			Kind:    snapdErr.Kind,
			Message: snapdErr.Message,
			Value:   value,
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
				return &asyncResponse{
					ID: change.ID,
				}, nil
			}
		}
	}
}

//go:linkname doSync github.com/snapcore/snapd/client.(*Client).doSync
func doSync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader, v any) (*snapdClient.ResultInfo, error)

//go:linkname doAsync github.com/snapcore/snapd/client.(*Client).doAsync
func doAsync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader) (changeID string, err error)
