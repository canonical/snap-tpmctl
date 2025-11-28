package snapd

import (
	"context"
	"encoding/json"
	"net/http"
)

// RecoveryKeySlot describes a recovery keyslot target.
// If ContainerRole is omitted, the keyslot will be implicitly expanded
// into two target keyslots for both "system-data" and "system-save".
type RecoveryKeySlot struct {
	ContainerRole string `json:"container-role,omitempty"`
	Name          string `json:"name"`
}

// GenerateRecoveryKeyResult describes the response from `generate-recovery-key` API.
type GenerateRecoveryKeyResult struct {
	RecoveryKey string `json:"recovery-key"`
	KeyID       string `json:"key-id"`
}

// RecoveryKeyRequest represents a request to manage recovery keys in snapd.
type RecoveryKeyRequest struct {
	Action         string            `json:"action"`
	KeyID          string            `json:"key-id,omitempty"`
	KeySlots       []RecoveryKeySlot `json:"keyslots,omitempty"`
	RecoveryKey    string            `json:"recovery-key,omitempty"`
	ContainerRoles []string          `json:"container-role,omitempty"`
}

// GenerateRecoveryKey creates a new recovery key and returns the key and its ID.
func (c *Client) GenerateRecoveryKey(ctx context.Context) (*GenerateRecoveryKeyResult, error) {
	body := RecoveryKeyRequest{
		Action: "generate-recovery-key",
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	var result GenerateRecoveryKeyResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// AddRecoveryKey adds a recovery key to the specified keyslots.
func (c *Client) AddRecoveryKey(ctx context.Context, keyID string, keySlots []RecoveryKeySlot) (*Response, error) {
	body := RecoveryKeyRequest{
		Action:   "add-recovery-key",
		KeyID:    keyID,
		KeySlots: keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ReplaceRecoveryKey replaces a recovery key to the specified keyslots.
func (c *Client) ReplaceRecoveryKey(ctx context.Context, keyID string, keySlots []RecoveryKeySlot) (*Response, error) {
	body := RecoveryKeyRequest{
		Action:   "replace-recovery-key",
		KeyID:    keyID,
		KeySlots: keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CheckRecoveryKey check a recovery key to the specified keyslots.
func (c *Client) CheckRecoveryKey(ctx context.Context, recoveryKey string, containerRoles []string) (*Response, error) {
	body := RecoveryKeyRequest{
		Action:         "check-recovery-key",
		RecoveryKey:    recoveryKey,
		ContainerRoles: containerRoles,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
