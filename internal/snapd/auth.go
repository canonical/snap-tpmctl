package snapd

import (
	"context"
	"net/http"
)

// PassphraseRequest represents a request to manage passphrases in snapd.
type PassphraseRequest struct {
	Action        string    `json:"action"`
	KeySlots      []KeySlot `json:"keyslots,omitempty"`
	NewPassphrase string    `json:"new-passphrase,omitempty"`
	OldPassphrase string    `json:"old-passphrase,omitempty"`
	Passphrase    string    `json:"passphrase,omitempty"`
}

// ReplacePassphrase replaces a passphrase to the specified keyslots.
// This is an async operation that waits for completion.
func (c *Client) ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []KeySlot) (*AsyncResponse, error) {
	body := PassphraseRequest{
		Action:        "change-passphrase",
		NewPassphrase: newPassphrase,
		OldPassphrase: oldPassphrase,
		KeySlots:      keySlots,
	}

	resp, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CheckPassphrase checks if the provided passphrase is valid.
func (c *Client) CheckPassphrase(ctx context.Context, passphrase string) (*Response, error) {
	body := PassphraseRequest{
		Action:     "check-passphrase",
		Passphrase: passphrase,
	}

	resp, err := c.doSyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// PINRequest represents a request to manage PINs in snapd.
type PINRequest struct {
	Action   string    `json:"action"`
	KeySlots []KeySlot `json:"keyslots,omitempty"`
	NewPin   string    `json:"new-pin,omitempty"`
	OldPin   string    `json:"old-pin,omitempty"`
	Pin      string    `json:"pin,omitempty"`
}

// CheckPIN checks if the provided PIN is valid.
func (c *Client) CheckPIN(ctx context.Context, pin string) (*Response, error) {
	body := PINRequest{
		Action: "check-pin",
		Pin:    pin,
	}

	resp, err := c.doSyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ReplacePIN replaces a PIN to the specified keyslots.
// This is an async operation that waits for completion.
func (c *Client) ReplacePIN(ctx context.Context, oldPin string, newPin string, keySlots []KeySlot) (*AsyncResponse, error) {
	body := PINRequest{
		Action:   "change-pin",
		NewPin:   newPin,
		OldPin:   oldPin,
		KeySlots: keySlots,
	}

	resp, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AuthMode represents the authentication mode for platform keys.
type AuthMode string

// Supported authentication modes for platform keys.
const (
	AuthModePin        AuthMode = "pin"
	AuthModePassphrase AuthMode = "passphrase"
	AuthModeNone       AuthMode = "none"
)

// KDFType represents the key derivation function type.
type KDFType string

// KDF (Key Derivation Function) types supported for password-based key derivation.
const (
	KDFTypeArgon2id KDFType = "argon2id"
	KDFTypeArgon2i  KDFType = "argon2i"
	KDFTypePBKDF2   KDFType = "pbkdf2"
)

// PlatformKeyRequest represents the request body for replacing a platform key.
type PlatformKeyRequest struct {
	Action     string    `json:"action"`
	AuthMode   AuthMode  `json:"auth-mode"`
	Passphrase string    `json:"passphrase,omitempty"`
	Pin        string    `json:"pin,omitempty"`
	KDFTime    *int      `json:"kdf-time,omitempty"`
	KDFType    KDFType   `json:"kdf-type,omitempty"`
	KeySlots   []KeySlot `json:"keyslots,omitempty"`
}

// ReplacePlatformKey replaces the platform key with the specified authentication.
func (c *Client) ReplacePlatformKey(ctx context.Context, authMode AuthMode, pin, passphrase string) (*AsyncResponse, error) {
	body := PlatformKeyRequest{
		Action:     "replace-platform-key",
		AuthMode:   authMode,
		Pin:        pin,
		Passphrase: passphrase,
	}

	resp, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
