package snapd

import (
	"context"
	"fmt"
	"net/http"

	snapdClient "github.com/snapcore/snapd/client"
	"github.com/snapcore/snapd/gadget/device"
)

// ReplacePassphrase replaces a passphrase to the specified keyslots.
// This is an async operation that waits for completion.
func (c *Client) ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []Keyslot) error {
	body := struct {
		Action   string    `json:"action"`
		KeySlots []Keyslot `json:"keyslots"`

		snapdClient.ChangePassphraseOptions
	}{
		Action:   "change-passphrase",
		KeySlots: keySlots,
		ChangePassphraseOptions: snapdClient.ChangePassphraseOptions{
			NewPassphrase: newPassphrase,
			OldPassphrase: oldPassphrase,
		},
	}

	if err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
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
		PIN    string `json:"pin"`
	}{
		Action: "check-pin",
		PIN:    pin,
	}

	if _, err := c.doSyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
		return err
	}

	return nil
}

// ReplacePIN replaces a PIN to the specified keyslots.
// This is an async operation that waits for completion.
func (c *Client) ReplacePIN(ctx context.Context, oldPIN string, newPIN string, keySlots []Keyslot) error {
	body := struct {
		Action   string    `json:"action"`
		KeySlots []Keyslot `json:"keyslots"`

		snapdClient.ChangePINOptions
	}{
		Action:   "change-pin",
		KeySlots: keySlots,
		ChangePINOptions: snapdClient.ChangePINOptions{
			NewPIN: newPIN,
			OldPIN: oldPIN,
		},
	}

	if err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
		return err
	}

	return nil
}

// AuthMode represents the authentication mode for platform keys.
type AuthMode = device.AuthMode

// Supported authentication modes for platform keys.
const (
	AuthModePIN        = device.AuthModePIN
	AuthModePassphrase = device.AuthModePassphrase
	AuthModeNone       = device.AuthModeNone
)

// ReplacePlatformKey replaces the platform key with the specified authentication.
func (c *Client) ReplacePlatformKey(ctx context.Context, authMode AuthMode, secret string) error {
	if authMode == AuthModeNone && secret != "" {
		return fmt.Errorf("expected no secret when auth mode is none, got: %q", secret)
	}

	var passphrase, pin string
	switch authMode {
	case AuthModePIN:
		pin = secret
	case AuthModePassphrase:
		passphrase = secret
	}

	body := struct {
		Action string `json:"action"`

		snapdClient.PlatformKeyOptions
	}{
		Action: "replace-platform-key",
		PlatformKeyOptions: snapdClient.PlatformKeyOptions{
			AuthMode:   authMode,
			PIN:        pin,
			Passphrase: passphrase,
		},
	}

	if err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body); err != nil {
		return err
	}

	return nil
}
