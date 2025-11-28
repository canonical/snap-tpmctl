package snapd

import (
	"context"
	"net/http"
)

// PassphraseRequest represents a request to manage passphrase in snapd.
type PassphraseRequest struct {
	Action        string    `json:"action"`
	KeySlots      []KeySlot `json:"keyslots,omitempty"`
	NewPassphrase string    `json:"new-passphrase,omitempty"`
	OldPassphrase string    `json:"old-passphrase,omitempty"`
	Passphrase    string    `json:"passphrase,omitempty"`
}

// ReplacePassphrase replaces a passphrase to the specified keyslots.
func (c *Client) ReplacePassphrase(ctx context.Context, oldPassphrase string, newPassphrase string, keySlots []KeySlot) (*Response, error) {
	body := PassphraseRequest{
		Action:        "change-passphrase",
		NewPassphrase: newPassphrase,
		OldPassphrase: oldPassphrase,
		KeySlots:      keySlots,
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/v2/system-volumes", nil, body)
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
