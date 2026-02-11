package snapd

import (
	"context"
	"encoding/json"
	"net/http"
)

// FdeStatus retrieves the current FDE status of the system.
func (c *Client) FdeStatus(ctx context.Context) (string, error) {
	resp, err := c.doSyncRequest(ctx, http.MethodGet, "/v2/system-info/storage-encrypted", nil, nil, nil)
	if err != nil {
		return "", err
	}

	// fdeStatusResult represents the FDE status info.
	type fdeStatusResult struct {
		Status string `json:"status"`
	}

	var status fdeStatusResult
	if err := json.Unmarshal(resp.Result, &status); err != nil {
		return "", err
	}

	return status.Status, nil
}
