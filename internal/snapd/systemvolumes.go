package snapd

import (
	"context"
	"encoding/json"
	"net/http"
)

// VolumeInfo describes a system volume
type VolumeInfo struct {
	Name       string                 `json:"name"`
	VolumeName string                 `json:"volume-name"`
	Encrypted  bool                   `json:"encrypted"`
	KeySlots   map[string]KeySlotInfo `json:"keyslots,omitempty"`
}

// SystemVolumesResult describes the system volumes response
type SystemVolumesResult struct {
	ByContainerRole map[string]VolumeInfo `json:"by-container-role"`
}

// GetSystemVolumes gets information about system volumes
func (c *Client) GetSystemVolumes(ctx context.Context) (*SystemVolumesResult, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/v2/system-volumes", nil, nil)
	if err != nil {
		return nil, err
	}

	var volumes SystemVolumesResult
	if err := json.Unmarshal(resp.Result, &volumes); err != nil {
		return nil, err
	}

	return &volumes, nil
}

// AddSystemVolumeKeySlots adds keyslots to a system volume
func (c *Client) AddSystemVolumeKeySlots(ctx context.Context, volume string, keySlots []VolumeKeySlot) (string, error) {
	body := map[string]any{
		"action":   "add-key",
		"volume":   volume,
		"keyslots": keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return "", err
	}

	return resp.Change, nil
}

// RemoveSystemVolumeKeySlots removes keyslots from a system volume
func (c *Client) RemoveSystemVolumeKeySlots(ctx context.Context, volume string, keySlots []VolumeKeySlot) (string, error) {
	body := map[string]any{
		"action":   "remove-key",
		"volume":   volume,
		"keyslots": keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return "", err
	}

	return resp.Change, nil
}
