// Package snapd provides a client for making calls to the systems local snapd service
package snapd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const (
	defaultSocketPath = "/var/run/snapd.socket"
	defaultUserAgent  = "snapd.go"
)

// Client is a snapd client.
type Client struct {
	httpClient       *http.Client
	socketPath       string
	userAgent        string
	macaroon         string
	discharges       []string
	allowInteraction bool
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// NewClient creates a new snapd client.
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		socketPath:       defaultSocketPath,
		userAgent:        defaultUserAgent,
		allowInteraction: true, // TODO: should default be true?
	}

	for _, opt := range opts {
		opt(client)
	}

	client.httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", client.socketPath)
			},
		},
	}

	return client
}

// Close closes the client.
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// SetAuthorization sets the authorization credentials.
func (c *Client) SetAuthorization(req *http.Request) {
	// Enable interactive authentication (Polkit)
	if c.allowInteraction {
		req.Header.Set("X-Allow-Interaction", "true")
	}

	// Add user key
	if c.macaroon != "" {
		auth := fmt.Sprintf("Macaroon root=%q", c.macaroon)
		for _, discharge := range c.discharges {
			auth += fmt.Sprintf(",discharge=%q", discharge)
		}
		req.Header.Set("Authorization", auth)
	}
}

// SetGenericHeaders sets the common HTTP headers for snapd API requests.
func (c *Client) SetGenericHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
}

// NewRequestBody marshals the given body into JSON format and returns it as an io.Reader.
func (c *Client) NewRequestBody(body any) (io.Reader, error) {
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

// NewURL constructs a new URL for the snapd REST API.
func (c *Client) NewURL(path string, query url.Values) url.URL {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   path,
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u
}

// Response is the base response structure from snapd.
type Response struct {
	Type       string          `json:"type"`
	StatusCode int             `json:"status-code"`
	Status     string          `json:"status"`
	Result     json.RawMessage `json:"result,omitempty"`
	Change     string          `json:"change,omitempty"`
}

// IsOK checks if a commonly know snapd accepted status was returned.
func (r *Response) IsOK() bool {
	return r.Status == "Accepted" || r.Status == "OK"
}

// snapdError represents an error from snapd.
type snapdError struct {
	Message    string
	Kind       string
	StatusCode int
	Status     string
	Value      any
}

func (e *snapdError) Error() string {
	if e.Kind != "" {
		return fmt.Sprintf("snapd error: %s (%s)", e.Message, e.Kind)
	}
	return fmt.Sprintf("snapd error: %s", e.Message)
}

// NewResponseBody parses a JSON response body from snapd and returns a snapdResponse.
// If the response type is "error", it extracts error details from the Result field and returns a snapdError.
func (c *Client) NewResponseBody(body []byte) (*Response, error) {
	var snapdResp Response
	if err := json.Unmarshal(body, &snapdResp); err != nil {
		return nil, err
	}

	if snapdResp.Type == "error" {
		var errResp struct {
			Message string         `json:"message"`
			Kind    string         `json:"kind,omitempty"`
			Value   map[string]any `json:"value,omitempty"`
		}

		if err := json.Unmarshal(snapdResp.Result, &errResp); err != nil {
			return nil, err
		}

		return nil, &snapdError{
			Message:    errResp.Message,
			Kind:       errResp.Kind,
			StatusCode: snapdResp.StatusCode,
			Status:     snapdResp.Status,
			Value:      errResp.Value,
		}
	}

	return &snapdResp, nil
}

// LoadAuthFromFile loads authentication credentials from file.
func (c *Client) LoadAuthFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read auth file: %w", err)
	}

	// auth represents the authentication data from snapd's auth.json
	var auth struct {
		Macaroon   string   `json:"macaroon"`
		Discharges []string `json:"discharges,omitempty"`
	}

	if err := json.Unmarshal(data, &auth); err != nil {
		return fmt.Errorf("failed to parse auth file: %w", err)
	}

	c.macaroon = auth.Macaroon
	c.discharges = auth.Discharges

	return nil
}

// LoadAuthFromHome loads authentication credentials from user's home directory.
func (c *Client) LoadAuthFromHome() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	authPath := filepath.Join(homeDir, ".snap", "auth.json")
	return c.LoadAuthFromFile(authPath)
}

// doRequest performs an HTTP request to snapd.
//
//nolint:unparam // path parameter may vary in future
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body any) (*Response, error) {
	reqBody, err := c.NewRequestBody(body)
	if err != nil {
		return nil, err
	}

	u := c.NewURL(path, query)

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	c.SetGenericHeaders(req)
	c.SetAuthorization(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	snapdResp, err := c.NewResponseBody(bodyBytes)
	if err != nil {
		return nil, err
	}

	return snapdResp, nil
}
