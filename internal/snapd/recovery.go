package snapd

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/snapcore/snapd/client"
)

// KeySlot describes a recovery keyslot target.
// If ContainerRole is omitted, the keyslot will be implicitly expanded
// into two target keyslots for both "system-data" and "system-save".
type KeySlot struct {
	ContainerRole string `json:"container-role,omitempty"`
	Name          string `json:"name"`
}

// GenerateRecoveryKeyResult describes the response from `generate-recovery-key` API.
type GenerateRecoveryKeyResult struct {
	RecoveryKey string `json:"recovery-key"`
	KeyID       string `json:"key-id"`
}

// GenerateRecoveryKey creates a new recovery key and returns the key and its ID.
func (c *Client) GenerateRecoveryKey(ctx context.Context) (*GenerateRecoveryKeyResult, error) {
	body := struct {
		Action string `json:"action"`
	}{
		Action: "generate-recovery-key",
	}

	resp, err := c.doSyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, &body)
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
// This is an async operation that waits for completion.
func (c *Client) AddRecoveryKey(ctx context.Context, keyID string, keySlots []KeySlot) error {
	body := struct {
		Action   string    `json:"action"`
		KeyID    string    `json:"key-id"`
		KeySlots []KeySlot `json:"keyslots"`
	}{
		Action:   "add-recovery-key",
		KeyID:    keyID,
		KeySlots: keySlots,
	}

	if err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
		return err
	}

	return nil
}

// ReplaceRecoveryKey replaces a recovery key to the specified keyslots.
// This is an async operation that waits for completion.
func (c *Client) ReplaceRecoveryKey(ctx context.Context, keyID string, keySlots []KeySlot) error {
	body := struct {
		Action   string    `json:"action"`
		KeyID    string    `json:"key-id"`
		KeySlots []KeySlot `json:"keyslots"`
	}{
		Action:   "replace-recovery-key",
		KeyID:    keyID,
		KeySlots: keySlots,
	}

	if err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
		return err
	}

	return nil
}

// CheckRecoveryKey check a recovery key to the specified keyslots.
func (c *Client) CheckRecoveryKey(ctx context.Context, recoveryKey string, containerRoles []string) (bool, error) {
	body := struct {
		Action         string   `json:"action"`
		RecoveryKey    string   `json:"recovery-key"`
		ContainerRoles []string `json:"container-role"`
	}{
		Action:         "check-recovery-key",
		RecoveryKey:    recoveryKey,
		ContainerRoles: containerRoles,
	}

	_, err := c.doSyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	var e *Error
	if errors.As(err, &e) && e.Kind == client.ErrorKindInvalidRecoveryKey {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
