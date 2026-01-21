package snapd

import (
	"context"
	"encoding/json"
	"net/http"
)

// FdeStatusResult represents the FDE status info.
type FdeStatusResult struct {
	Status string `json:"status"`
}

// FdeStatus retrieves the current FDE status of the system.
func (c *Client) FdeStatus(ctx context.Context) (*FdeStatusResult, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/v2/system-info/storage-encrypted", nil, nil)
	if err != nil {
		return nil, err
	}

	var status FdeStatusResult
	if err := json.Unmarshal(resp.Result, &status); err != nil {
		return nil, err
	}

	return &status, nil
}
