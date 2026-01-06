// Package snapd provides a client for making calls to the systems local snapd service
package snapd

import (
	"bytes"
	"encoding/json"
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

func (c *Client) Foo() error {
	var body bytes.Buffer
	req := struct {
		Action string `json:"action"`
	}{
		Action: "foo",
	}

	if err := json.NewEncoder(&body).Encode(&req); err != nil {
		return err
	}
	if _, err := doSync(c.snapd, "POST", "/v2/foo", nil, nil, &body, nil); err != nil {
		return fmt.Errorf("cannot request system reboot: %v", err)
	}
	return nil

}

func (c *Client) doSync(method, path string, query url.Values, headers map[string]string, body io.Reader) (response, error) {
	var resp response

	_, err := doSync(c.snapd, method, path, query, headers, body, &resp)
	return resp, err
}

// go:linkname doSync github.com/snapcore/snap/client.(*Client).doSync
func doSync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader, v any) (*snapdClient.ResultInfo, error)

// go:linkname doAsync github.com/snapcore/snap/client.(*Client).doAsync
func doAsync(c *snapdClient.Client, method, path string, query url.Values, headers map[string]string, body io.Reader) (changeID string, err error)

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
