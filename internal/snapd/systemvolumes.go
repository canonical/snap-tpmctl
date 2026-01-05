package snapd

import (
	"context"
	"encoding/json"
	"net/http"
)

// KeySlotInfo describes a keyslot in a volume.
type KeySlotInfo struct {
	Type         string   `json:"type"`
	AuthMode     string   `json:"auth-mode,omitempty"`
	PlatformName string   `json:"platform-name,omitempty"`
	Roles        []string `json:"roles,omitempty"`
}

// VolumeInfo describes a system volume.
type VolumeInfo struct {
	Name       string                 `json:"name"`
	VolumeName string                 `json:"volume-name"`
	Encrypted  bool                   `json:"encrypted"`
	KeySlots   map[string]KeySlotInfo `json:"keyslots,omitempty"`
}

// SystemVolumesResult describes the system volumes response.
type SystemVolumesResult struct {
	ByContainerRole map[string]VolumeInfo `json:"by-container-role"`
}

// EnumerateKeySlots gets information about system volumes.
func (c *Client) EnumerateKeySlots(ctx context.Context) (*SystemVolumesResult, error) {
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
