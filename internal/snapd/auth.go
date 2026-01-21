package snapd

import (
	"context"
	"net/http"
)

// ReplacePassphrase replaces a passphrase to the specified keyslots.
// This is an async operation that waits for completion.
func (c *Client) ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []KeySlot) error {
	body := struct {
		Action        string    `json:"action"`
		KeySlots      []KeySlot `json:"keyslots"`
		NewPassphrase string    `json:"new-passphrase"`
		OldPassphrase string    `json:"old-passphrase"`
	}{
		Action:        "change-passphrase",
		NewPassphrase: newPassphrase,
		OldPassphrase: oldPassphrase,
		KeySlots:      keySlots,
	}

	_, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return err
	}

	return nil
}

// CheckPassphrase checks if the provided passphrase is valid.
func (c *Client) CheckPassphrase(ctx context.Context, passphrase string) error {
	body := struct {
		Action     string `json:"action"`
		Passphrase string `json:"passphrase"`
	}{
		Action:     "check-passphrase",
		Passphrase: passphrase,
	}

	if _, err := c.doSyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
		return err
	}

	return nil
}

// CheckPIN checks if the provided PIN is valid.
func (c *Client) CheckPIN(ctx context.Context, pin string) error {
	body := struct {
		Action string `json:"action"`
		Pin    string `json:"pin"`
	}{
		Action: "check-pin",
		Pin:    pin,
	}

	if _, err := c.doSyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
		return err
	}

	return nil
}

// ReplacePIN replaces a PIN to the specified keyslots.
// This is an async operation that waits for completion.
func (c *Client) ReplacePIN(ctx context.Context, oldPin string, newPin string, keySlots []KeySlot) error {
	body := struct {
		Action   string    `json:"action"`
		KeySlots []KeySlot `json:"keyslots"`
		NewPin   string    `json:"new-pin"`
		OldPin   string    `json:"old-pin"`
	}{
		Action:   "change-pin",
		NewPin:   newPin,
		OldPin:   oldPin,
		KeySlots: keySlots,
	}

	_, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return err
	}

	return nil
}

// AuthMode represents the authentication mode for platform keys.
type AuthMode string

// Supported authentication modes for platform keys.
const (
	AuthModePin        AuthMode = "pin"
	AuthModePassphrase AuthMode = "passphrase"
	AuthModeNone       AuthMode = "none"
)

// ReplacePlatformKey replaces the platform key with the specified authentication.
func (c *Client) ReplacePlatformKey(ctx context.Context, authMode AuthMode, pin, passphrase string) error {
	body := struct {
		Action     string   `json:"action"`
		AuthMode   AuthMode `json:"auth-mode"`
		Passphrase string   `json:"passphrase"`
		Pin        string   `json:"pin"`
	}{
		Action:     "replace-platform-key",
		AuthMode:   authMode,
		Pin:        pin,
		Passphrase: passphrase,
	}

	_, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return err
	}

	return nil
}
