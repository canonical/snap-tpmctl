package snapd

import (
	"context"
	"net/http"
)

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
		KDFTime:    nil,
		KDFType:    "",
		KeySlots:   nil,
	}

	resp, err := c.doAsyncRequest(ctx, http.MethodPost,
		"/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
