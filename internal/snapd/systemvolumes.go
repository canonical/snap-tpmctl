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

// VolumeKeySlot describes a keyslot configuration.
type VolumeKeySlot struct {
	Number   int    `json:"number"`
	Name     string `json:"name,omitempty"`
	Priority int    `json:"priority,omitempty"`
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

// SystemVolumesRequest represents a request to manage volumes in snapd.
type SystemVolumesRequest struct {
	Action   string          `json:"action"`
	Volume   string          `json:"volume,omitempty"`
	KeySlots []VolumeKeySlot `json:"keyslots,omitempty"`
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

// AddSystemVolumeKeySlots adds keyslots to a system volume.
func (c *Client) AddSystemVolumeKeySlots(ctx context.Context, volume string, keySlots []VolumeKeySlot) (*Response, error) {
	body := SystemVolumesRequest{
		Action:   "add-key",
		Volume:   volume,
		KeySlots: keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RemoveSystemVolumeKeySlots removes keyslots from a system volume.
func (c *Client) RemoveSystemVolumeKeySlots(ctx context.Context, volume string, keySlots []VolumeKeySlot) (*Response, error) {
	body := SystemVolumesRequest{
		Action:   "remove-key",
		Volume:   volume,
		KeySlots: keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
