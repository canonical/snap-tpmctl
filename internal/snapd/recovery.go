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

// RecoveryKeyResult describes the response from generate-recovery-key.
type GenerateRecoveryKeyResult struct {
	RecoveryKey string `json:"recovery-key"`
	KeyID       string `json:"key-id"`
}

// GenerateRecoveryKey creates a new recovery key and returns the key and its ID.
func (c *Client) GenerateRecoveryKey(ctx context.Context) (*GenerateRecoveryKeyResult, error) {
	body := map[string]any{
		"action": "generate-recovery-key",
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
func (c *Client) AddRecoveryKey(ctx context.Context, keyID string, keySlots []RecoveryKeySlot) (string, error) {
	body := map[string]any{
		"action":   "add-recovery-key",
		"key-id":   keyID,
		"keyslots": keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return "", err
	}

	return resp.Change, nil
}

// ReplaceRecoveryKey replace a recovery key to the specified keyslots.
func (c *Client) ReplaceRecoveryKey(ctx context.Context, keyID string, keySlots []RecoveryKeySlot) (string, error) {
	body := map[string]any{
		"action":   "replace-recovery-key",
		"key-id":   keyID,
		"keyslots": keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return "", err
	}

	return resp.Change, nil
}
