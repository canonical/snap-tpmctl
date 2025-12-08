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

// PINRequest represents a request to manage PINs in snapd.
type PINRequest struct {
	Action   string    `json:"action"`
	KeySlots []KeySlot `json:"keyslots,omitempty"`
	NewPin   string    `json:"new-pin,omitempty"`
	OldPin   string    `json:"old-pin,omitempty"`
	Pin      string    `json:"pin,omitempty"`
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

	resp, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// EntropyResponse contains entropy calculation results.
// type EntropyResponse struct {
// 	Entropy            uint `json:"entropy-bits"`
// 	MinEntropyBits     uint `json:"min-entropy-bits"`
// 	OptimalEntropyBits uint `json:"optimal-entropy-bits"`
// }

// CheckPassphrase checks if the provided passphrase is valid.
func (c *Client) CheckPassphrase(ctx context.Context, passphrase string) (*Response, error) {
	body := PassphraseRequest{
		Action:     "check-passphrase",
		Passphrase: passphrase,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CheckPIN checks if the provided PIN is valid.
func (c *Client) CheckPIN(ctx context.Context, pin string) (*Response, error) {
	body := PINRequest{
		Action: "check-pin",
		Pin:    pin,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
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

	resp, err := c.doAsyncRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
